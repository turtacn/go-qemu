package scenarios

import (
	"context"
	"fmt"

	"github.com/turtacn/go-qemu/internal/cmdline"
	"github.com/turtacn/go-qemu/internal/cmdline/rules"
	"github.com/turtacn/go-qemu/internal/version"
	"github.com/turtacn/go-qemu/internal/vmconfig"
)

type CmdlineDiffScenario struct {
	diff *cmdline.CmdlineDiff
}

func (s *CmdlineDiffScenario) Name() string { return "cmdline-diff-5.0-vs-10.2.2" }
func (s *CmdlineDiffScenario) Description() string {
	return "Cmdline parameter diff and risk report between 5.0 and 10.2.2"
}
func (s *CmdlineDiffScenario) Prerequisites(ctx context.Context) error { return nil }
func (s *CmdlineDiffScenario) Setup(ctx context.Context) error         { return nil }
func (s *CmdlineDiffScenario) Execute(ctx context.Context) error {
	cfg := &vmconfig.VMConfig{
		Disks: []vmconfig.DiskConfig{{Bus: "scsi", Index: 0, Path: "/dev/sda"}},
	}
	reg := rules.DefaultRegistry()
	var ruleSlice []cmdline.Rule
	for _, r := range reg.Rules() {
		ruleSlice = append(ruleSlice, r)
	}

	b5, _ := cmdline.NewBuilder(cmdline.WithQEMUVersion(version.V(5, 0, 0)))
	b5.SetRules(ruleSlice)
	_, report5, _ := b5.Build(cfg)

	b10, _ := cmdline.NewBuilder(cmdline.WithQEMUVersion(version.V(10, 2, 2)))
	b10.SetRules(ruleSlice)
	_, report10, _ := b10.Build(cfg)

	s.diff = report10.DiffFrom(version.V(5, 0, 0), report5)
	return nil
}
func (s *CmdlineDiffScenario) Verify(ctx context.Context) error {
	if s.diff == nil {
		return fmt.Errorf("diff is nil")
	}
	hasBlockdevAdded := false
	for _, a := range s.diff.Added {
		if a.Key == "-blockdev" {
			hasBlockdevAdded = true
		}
	}
	if !hasBlockdevAdded {
		return fmt.Errorf("expected -blockdev in added diff")
	}
	return nil
}
func (s *CmdlineDiffScenario) Teardown(ctx context.Context) error { return nil }
