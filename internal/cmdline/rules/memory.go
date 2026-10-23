package rules

import (
	"fmt"

	"github.com/turtacn/go-qemu/internal/cmdline"
	"github.com/turtacn/go-qemu/internal/feature"
	"github.com/turtacn/go-qemu/internal/version"
	"github.com/turtacn/go-qemu/internal/vmconfig"
)

type MemoryRule struct{}

func (r *MemoryRule) Name() string { return "memory" }

func (r *MemoryRule) Applies(cfg *vmconfig.VMConfig, _ version.QEMUVersion, _ feature.Arch) bool {
	return cfg.Memory.SizeMB > 0
}

func (r *MemoryRule) Apply(cfg *vmconfig.VMConfig, ver version.QEMUVersion, reg *feature.FeatureGateRegistry) ([]cmdline.ArgEntry, error) {
	memValue := fmt.Sprintf("%dM", cfg.Memory.SizeMB)
	explicit, _, _ := reg.IsEnabled("memory.explicit-backend")

	args := []cmdline.ArgEntry{
		{Key: "-m", Value: memValue, RuleName: r.Name()},
	}

	if explicit {
		args = append(args, cmdline.ArgEntry{
			Key:      "-object",
			Value:    fmt.Sprintf("memory-backend-ram,id=mem,size=%dM", cfg.Memory.SizeMB),
			RuleName: r.Name(),
			Reason:   "Explicit memory backend object recommended in QEMU 5.0+",
		})
	}

	return args, nil
}

func (r *MemoryRule) Conflicts() []string { return nil }
