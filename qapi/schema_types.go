package qapi

import "encoding/json"

type EntryKind string

const (
	KindCommand   EntryKind = "command"
	KindStruct    EntryKind = "struct"
	KindEnum      EntryKind = "enum"
	KindUnion     EntryKind = "union"
	KindAlternate EntryKind = "alternate"
	KindEvent     EntryKind = "event"
	KindInclude   EntryKind = "include"
)

type ConditionExpr struct {
	Simple string
	All    []string
	Any    []string
	Not    *ConditionExpr
}

type SchemaEntry struct {
	Kind     EntryKind
	Name     string
	Data     json.RawMessage
	If       *ConditionExpr
	Features []string
	SinceVer string
}

type CommandDef struct {
	Name     string
	Data     map[string]FieldDef
	Returns  interface{}
	If       *ConditionExpr
	Features []string
}

type StructDef struct {
	Name     string
	Data     map[string]FieldDef
	Base     string
	If       *ConditionExpr
	Features []string
}

type EnumValue struct {
	Name     string
	If       *ConditionExpr
	Features []string
}

type EnumDef struct {
	Name     string
	Data     []EnumValue
	If       *ConditionExpr
	Features []string
}

type UnionDef struct {
	Name     string
	Data     map[string]string
	If       *ConditionExpr
	Features []string
}

type EventDef struct {
	Name     string
	Data     map[string]FieldDef
	If       *ConditionExpr
	Features []string
}

type FieldDef struct {
	Type     string
	Optional bool
	If       *ConditionExpr
	Features []string
}
