package vmconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	cfg := &VMConfig{
		Memory: MemoryConfig{SizeMB: 1024},
	}
	if err := Validate(cfg); err != nil {
		t.Fatalf("expected valid config, got %v", err)
	}
	if cfg.CPU.Cores != 1 {
		t.Errorf("expected default 1 core")
	}
	if cfg.Machine.Type != "pc-i440fx" {
		t.Errorf("expected default pc-i440fx machine")
	}

	invalidCfg := &VMConfig{
		Memory: MemoryConfig{SizeMB: 0},
	}
	if err := Validate(invalidCfg); err == nil {
		t.Errorf("expected error for 0MB memory")
	}
}

func TestLoadFromHCI(t *testing.T) {
	tmpDir := t.TempDir()
	yamlPath := filepath.Join(tmpDir, "vm.yaml")
	content := `
id: 101
name: test-vm
memory:
  sizemb: 2048
cpu:
  cores: 4
  sockets: 1
machine:
  type: pc-i440fx
disks:
  - bus: scsi
    index: 0
    path: /var/lib/qemu/disk0.raw
    format: raw
`
	if err := os.WriteFile(yamlPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write yaml: %v", err)
	}

	cfg, err := LoadFromHCI(yamlPath)
	if err != nil {
		t.Fatalf("LoadFromHCI failed: %v", err)
	}
	if cfg.ID != 101 || cfg.Memory.SizeMB != 2048 || len(cfg.Disks) != 1 {
		t.Errorf("unexpected loaded config: %+v", cfg)
	}
}
