package qapi

import (
	"fmt"
	"os"
	"path/filepath"
)

type Generator struct {
	Schema  *QAPISchema
	Version string
	Package string
	OutDir  string
}

func NewGenerator(schema *QAPISchema, version, pkg, outDir string) *Generator {
	return &Generator{
		Schema:  schema,
		Version: version,
		Package: pkg,
		OutDir:  outDir,
	}
}

func (g *Generator) Run() error {
	if err := os.MkdirAll(g.OutDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", g.OutDir, err)
	}

	if err := g.genTypes(); err != nil {
		return fmt.Errorf("genTypes failed: %w", err)
	}
	if err := g.genCommands(); err != nil {
		return fmt.Errorf("genCommands failed: %w", err)
	}
	if err := g.genMeta(); err != nil {
		return fmt.Errorf("genMeta failed: %w", err)
	}
	return nil
}

func (g *Generator) writeFile(filename, content string) error {
	path := filepath.Join(g.OutDir, filename)
	return os.WriteFile(path, []byte(content), 0644)
}
