package rules

import (
	"fmt"

	"github.com/turtacn/go-qemu/internal/cmdline"
	"github.com/turtacn/go-qemu/internal/feature"
	"github.com/turtacn/go-qemu/internal/version"
	"github.com/turtacn/go-qemu/internal/vmconfig"
)

type CPURule struct{}

func (r *CPURule) Name() string { return "cpu" }

func (r *CPURule) Applies(cfg *vmconfig.VMConfig, _ version.QEMUVersion, _ feature.Arch) bool {
	return true
}

func (r *CPURule) Apply(cfg *vmconfig.VMConfig, ver version.QEMUVersion, reg *feature.FeatureGateRegistry) ([]cmdline.ArgEntry, error) {
	model := cfg.CPU.Model
	if model == "" {
		model = "host"
	}
	smp := fmt.Sprintf("cpus=%d,sockets=%d,cores=%d,threads=1",
		cfg.CPU.Cores*cfg.CPU.Sockets, cfg.CPU.Sockets, cfg.CPU.Cores)

	return []cmdline.ArgEntry{
		{Key: "-cpu", Value: model, RuleName: r.Name()},
		{Key: "-smp", Value: smp, RuleName: r.Name()},
	}, nil
}

func (r *CPURule) Conflicts() []string { return nil }
