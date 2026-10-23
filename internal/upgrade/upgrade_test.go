package upgrade

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalUpgrader(t *testing.T) {
	tmpDir := t.TempDir()
	yamlPath := filepath.Join(tmpDir, "vm.yaml")
	content := `
id: 100
name: test-vm
memory:
  sizemb: 2048
cpu:
  cores: 2
machine:
  type: pc-i440fx
`
	if err := os.WriteFile(yamlPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	upgrader := NewLocalUpgrader(LocalUpgradeConfig{
		VMID:         100,
		VMConfigPath: yamlPath,
	})

	res, err := upgrader.RunWithResult(context.Background())
	if err != nil {
		t.Fatalf("RunWithResult failed: %v", err)
	}

	if res.VMID != 100 {
		t.Errorf("got VMID %d, expected 100", res.VMID)
	}
	if res.SrcVersion.String() != "5.0.0" || res.DstVersion.String() != "10.2.2" {
		t.Errorf("unexpected versions in result: %v -> %v", res.SrcVersion, res.DstVersion)
	}
}
