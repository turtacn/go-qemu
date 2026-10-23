package rules

import (
	"github.com/turtacn/go-qemu/internal/cmdline"
	"github.com/turtacn/go-qemu/internal/feature"
	"github.com/turtacn/go-qemu/internal/version"
	"github.com/turtacn/go-qemu/internal/vmconfig"
)

type DiskSelectorRule struct{}

func (r *DiskSelectorRule) Name() string { return "disk.selector" }

func (r *DiskSelectorRule) Applies(cfg *vmconfig.VMConfig, _ version.QEMUVersion, _ feature.Arch) bool {
	return len(cfg.Disks) > 0
}

func (r *DiskSelectorRule) Apply(cfg *vmconfig.VMConfig, ver version.QEMUVersion, reg *feature.FeatureGateRegistry) ([]cmdline.ArgEntry, error) {
	blockdevPreferred, _, _ := reg.IsEnabled("disk.blockdev-preferred")

	var allArgs []cmdline.ArgEntry
	for _, disk := range cfg.Disks {
		var args []cmdline.ArgEntry
		var err error

		if blockdevPreferred {
			args, err = buildBlockdevArgs(disk, ver)
		} else {
			args, err = buildDriveArgs(disk, ver)
		}
		if err != nil {
			return nil, err
		}
		allArgs = append(allArgs, args...)
	}
	return allArgs, nil
}

func (r *DiskSelectorRule) Conflicts() []string { return nil }
