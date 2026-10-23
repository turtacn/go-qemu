package qapi

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSchema500(t *testing.T) {
	schemaDir := "../schemas/qemu-5.0.0/qapi"
	schema, err := ParseDir(schemaDir, "5.0.0")
	if err != nil {
		t.Fatalf("ParseDir failed for 5.0.0: %v", err)
	}

	if !schema.CommandExists("migrate") {
		t.Errorf("expected command 'migrate' to exist in 5.0.0 schema")
	}
	if !schema.CommandExists("migrate-set-capabilities") {
		t.Errorf("expected command 'migrate-set-capabilities' to exist in 5.0.0 schema")
	}
	if !schema.CommandExists("query-qmp-schema") {
		t.Errorf("expected command 'query-qmp-schema' to exist in 5.0.0 schema")
	}

	if schema.EnumValueExists("MigrationCapability", "postcopy-preempt") {
		t.Errorf("postcopy-preempt should NOT exist in 5.0.0 MigrationCapability enum")
	}
	if !schema.EnumValueExists("MigrationCapability", "postcopy-ram") {
		t.Errorf("postcopy-ram should exist in 5.0.0 MigrationCapability enum")
	}
}

func TestParseSchema10(t *testing.T) {
	schemaDir := "../schemas/qemu-10.2.2/qapi"
	schema, err := ParseDir(schemaDir, "10.2.2")
	if err != nil {
		t.Fatalf("ParseDir failed for 10.2.2: %v", err)
	}

	if !schema.CommandExists("migrate") {
		t.Errorf("expected command 'migrate' to exist in 10.2.2 schema")
	}
	if !schema.EnumValueExists("MigrationCapability", "postcopy-preempt") {
		t.Errorf("postcopy-preempt SHOULD exist in 10.2.2 MigrationCapability enum")
	}
	if !schema.CommandExists("calc-dirty-rate") {
		t.Errorf("calc-dirty-rate command should exist in 10.2.2 schema")
	}
}

func TestCodeGenV5(t *testing.T) {
	schemaDir := "../schemas/qemu-5.0.0/qapi"
	schema, err := ParseDir(schemaDir, "5.0.0")
	if err != nil {
		t.Fatalf("ParseDir failed: %v", err)
	}

	tmpOutDir := t.TempDir()
	gen := NewGenerator(schema, "5.0.0", "gen_v5", tmpOutDir)
	if err := gen.Run(); err != nil {
		t.Fatalf("gen.Run failed: %v", err)
	}

	expectedFiles := []string{"types_gen.go", "commands_gen.go", "schema_meta_gen.go"}
	for _, f := range expectedFiles {
		p := filepath.Join(tmpOutDir, f)
		if _, err := os.Stat(p); err != nil {
			t.Errorf("expected generated file %s to exist", f)
		}
	}
}

func TestCodeGenV10(t *testing.T) {
	schemaDir := "../schemas/qemu-10.2.2/qapi"
	schema, err := ParseDir(schemaDir, "10.2.2")
	if err != nil {
		t.Fatalf("ParseDir failed: %v", err)
	}

	tmpOutDir := t.TempDir()
	gen := NewGenerator(schema, "10.2.2", "gen_v10", tmpOutDir)
	if err := gen.Run(); err != nil {
		t.Fatalf("gen.Run failed: %v", err)
	}

	expectedFiles := []string{"types_gen.go", "commands_gen.go", "schema_meta_gen.go"}
	for _, f := range expectedFiles {
		p := filepath.Join(tmpOutDir, f)
		if _, err := os.Stat(p); err != nil {
			t.Errorf("expected generated file %s to exist", f)
		}
	}
}
