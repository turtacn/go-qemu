package rules

import (
	"testing"

	"github.com/turtacn/go-qemu/internal/cmdline"
	"github.com/turtacn/go-qemu/internal/version"
	"github.com/turtacn/go-qemu/internal/vmconfig"
)

func TestDefaultRules(t *testing.T) {
	cfg := &vmconfig.VMConfig{
		ID:      100,
		CPU:     vmconfig.CPUConfig{Cores: 2, Sockets: 1, Model: "host"},
		Memory:  vmconfig.MemoryConfig{SizeMB: 2048},
		Machine: vmconfig.MachineConfig{Type: "pc-i440fx"},
		Disks:   []vmconfig.DiskConfig{{Bus: "scsi", Index: 0, Path: "/dev/sda", Format: "raw"}},
	}

	b5, _ := cmdline.NewBuilder(cmdline.WithQEMUVersion(version.V(5, 0, 0)))
	reg := DefaultRegistry()
	var ruleSlice []cmdline.Rule
	for _, r := range reg.Rules() {
		ruleSlice = append(ruleSlice, r)
	}
	b5.SetRules(ruleSlice)

	args5, report5, err := b5.Build(cfg)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if len(args5) == 0 {
		t.Errorf("expected non-empty args")
	}

	hasDrive := false
	for _, a := range report5.Args {
		if a.Key == "-drive" {
			hasDrive = true
		}
	}
	if !hasDrive {
		t.Errorf("expected -drive in 5.0.0 output")
	}

	b10, _ := cmdline.NewBuilder(cmdline.WithQEMUVersion(version.V(10, 2, 2)))
	b10.SetRules(ruleSlice)
	_, report10, err := b10.Build(cfg)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	hasBlockdev := false
	for _, a := range report10.Args {
		if a.Key == "-blockdev" {
			hasBlockdev = true
		}
	}
	if !hasBlockdev {
		t.Errorf("expected -blockdev in 10.2.2 output")
	}
}
