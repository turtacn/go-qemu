package qapi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// QAPISchema represents a parsed complete schema.
type QAPISchema struct {
	Version  string
	Commands map[string]*CommandDef
	Structs  map[string]*StructDef
	Enums    map[string]*EnumDef
	Unions   map[string]*UnionDef
	Events   map[string]*EventDef
}

func NewQAPISchema(version string) *QAPISchema {
	return &QAPISchema{
		Version:  version,
		Commands: make(map[string]*CommandDef),
		Structs:  make(map[string]*StructDef),
		Enums:    make(map[string]*EnumDef),
		Unions:   make(map[string]*UnionDef),
		Events:   make(map[string]*EventDef),
	}
}

func (s *QAPISchema) Merge(other *QAPISchema) {
	for k, v := range other.Commands {
		s.Commands[k] = v
	}
	for k, v := range other.Structs {
		s.Structs[k] = v
	}
	for k, v := range other.Enums {
		s.Enums[k] = v
	}
	for k, v := range other.Unions {
		s.Unions[k] = v
	}
	for k, v := range other.Events {
		s.Events[k] = v
	}
}

func safeParseFeatures(raw interface{}) []string {
	if raw == nil {
		return nil
	}
	if slice, ok := raw.([]interface{}); ok {
		return ParseFeatures(slice)
	}
	return nil
}

// ParseDir parses all qapi schema files starting from entry file (qapi-schema.json or dir/qapi-schema.json).
func ParseDir(dir string, version string) (*QAPISchema, error) {
	entryFile := filepath.Join(dir, "qapi-schema.json")
	if _, err := os.Stat(entryFile); os.IsNotExist(err) {
		entryFile = ""
	}

	visited := make(map[string]bool)
	if entryFile != "" {
		return parseWithIncludes(entryFile, version, visited)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	mainSchema := NewQAPISchema(version)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			path := filepath.Join(dir, entry.Name())
			subSchema, err := ParseFile(path, version)
			if err != nil {
				return nil, err
			}
			mainSchema.Merge(subSchema)
		}
	}
	return mainSchema, nil
}

func parseWithIncludes(entryFile string, version string, visited map[string]bool) (*QAPISchema, error) {
	absPath, err := filepath.Abs(entryFile)
	if err != nil {
		absPath = entryFile
	}
	if visited[absPath] {
		return nil, fmt.Errorf("circular include detected: %s", entryFile)
	}
	visited[absPath] = true

	data, err := os.ReadFile(entryFile)
	if err != nil {
		return nil, err
	}

	var rawEntries []map[string]interface{}
	if err := json.Unmarshal(data, &rawEntries); err != nil {
		return nil, fmt.Errorf("failed to parse JSON in %s: %w", entryFile, err)
	}

	schema := NewQAPISchema(version)
	dir := filepath.Dir(entryFile)

	for _, raw := range rawEntries {
		if inc, ok := raw["include"].(string); ok {
			incPath := filepath.Join(dir, inc)
			subSchema, err := parseWithIncludes(incPath, version, visited)
			if err != nil {
				return nil, err
			}
			schema.Merge(subSchema)
		} else {
			processEntry(schema, raw)
		}
	}

	return schema, nil
}

func ParseFile(path string, version string) (*QAPISchema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rawEntries []map[string]interface{}
	if err := json.Unmarshal(data, &rawEntries); err != nil {
		return nil, fmt.Errorf("failed to parse JSON in %s: %w", path, err)
	}

	schema := NewQAPISchema(version)
	for _, raw := range rawEntries {
		processEntry(schema, raw)
	}
	return schema, nil
}

func processEntry(schema *QAPISchema, raw map[string]interface{}) {
	if cmdName, ok := raw["command"].(string); ok {
		cmd := &CommandDef{
			Name:     cmdName,
			Returns:  raw["returns"],
			If:       ParseConditionExpr(raw["if"]),
			Features: safeParseFeatures(raw["features"]),
		}
		schema.Commands[cmdName] = cmd
	} else if enumName, ok := raw["enum"].(string); ok {
		enumDef := &EnumDef{
			Name:     enumName,
			If:       ParseConditionExpr(raw["if"]),
			Features: safeParseFeatures(raw["features"]),
		}
		if rawData, ok := raw["data"].([]interface{}); ok {
			for _, item := range rawData {
				if s, ok := item.(string); ok {
					enumDef.Data = append(enumDef.Data, EnumValue{Name: s})
				} else if m, ok := item.(map[string]interface{}); ok {
					if name, ok := m["name"].(string); ok {
						enumDef.Data = append(enumDef.Data, EnumValue{
							Name:     name,
							If:       ParseConditionExpr(m["if"]),
							Features: safeParseFeatures(m["features"]),
						})
					}
				}
			}
		}
		schema.Enums[enumName] = enumDef
	} else if structName, ok := raw["struct"].(string); ok {
		structDef := &StructDef{
			Name:     structName,
			Data:     make(map[string]FieldDef),
			If:       ParseConditionExpr(raw["if"]),
			Features: safeParseFeatures(raw["features"]),
		}
		if base, ok := raw["base"].(string); ok {
			structDef.Base = base
		}
		if rawData, ok := raw["data"].(map[string]interface{}); ok {
			for fName, fVal := range rawData {
				isOpt := strings.HasPrefix(fName, "*")
				if fStr, ok := fVal.(string); ok {
					structDef.Data[fName] = FieldDef{
						Type:     fStr,
						Optional: isOpt,
					}
				} else if fSlice, ok := fVal.([]interface{}); ok {
					if len(fSlice) > 0 {
						if elemStr, ok := fSlice[0].(string); ok {
							structDef.Data[fName] = FieldDef{
								Type:     "[]" + elemStr,
								Optional: isOpt,
							}
						}
					}
				}
			}
		}
		schema.Structs[structName] = structDef
	} else if eventName, ok := raw["event"].(string); ok {
		eventDef := &EventDef{
			Name:     eventName,
			If:       ParseConditionExpr(raw["if"]),
			Features: safeParseFeatures(raw["features"]),
		}
		schema.Events[eventName] = eventDef
	} else if unionName, ok := raw["union"].(string); ok {
		unionDef := &UnionDef{
			Name:     unionName,
			If:       ParseConditionExpr(raw["if"]),
			Features: safeParseFeatures(raw["features"]),
		}
		schema.Unions[unionName] = unionDef
	}
}

func (s *QAPISchema) CommandExists(cmdName string) bool {
	_, ok := s.Commands[cmdName]
	return ok
}

func (s *QAPISchema) EnumValueExists(enumName, valName string) bool {
	enumDef, ok := s.Enums[enumName]
	if !ok {
		return false
	}
	for _, val := range enumDef.Data {
		if val.Name == valName {
			return true
		}
	}
	return false
}

func (s *QAPISchema) HasKey(key string) bool {
	parts := strings.SplitN(key, ":", 3)
	if len(parts) < 2 {
		return false
	}
	kind := parts[0]
	switch kind {
	case "command":
		return s.CommandExists(parts[1])
	case "type", "struct":
		_, okStruct := s.Structs[parts[1]]
		_, okEnum := s.Enums[parts[1]]
		return okStruct || okEnum
	case "enum":
		_, ok := s.Enums[parts[1]]
		return ok
	case "enum-value":
		if len(parts) == 3 {
			return s.EnumValueExists(parts[1], parts[2])
		}
	}
	return false
}
