package rules

import (
	"fmt"

	"github.com/turtacn/go-qemu/internal/cmdline"
	"github.com/turtacn/go-qemu/internal/feature"
	"github.com/turtacn/go-qemu/internal/version"
	"github.com/turtacn/go-qemu/internal/vmconfig"
)

type NetworkRule struct{}

func (r *NetworkRule) Name() string { return "network" }

func (r *NetworkRule) Applies(cfg *vmconfig.VMConfig, _ version.QEMUVersion, _ feature.Arch) bool {
	return len(cfg.Networks) > 0
}

func (r *NetworkRule) Apply(cfg *vmconfig.VMConfig, ver version.QEMUVersion, reg *feature.FeatureGateRegistry) ([]cmdline.ArgEntry, error) {
	var args []cmdline.ArgEntry
	for i, nic := range cfg.Networks {
		id := fmt.Sprintf("net%d", i)
		model := nic.Model
		if model == "" {
			model = "virtio-net-pci"
		}
		netdevVal := fmt.Sprintf("tap,id=%s", id)
		if nic.Bridge != "" {
			netdevVal = fmt.Sprintf("tap,id=%s,br=%s", id, nic.Bridge)
		}
		deviceVal := fmt.Sprintf("%s,netdev=%s,id=nic%d", model, id, i)
		if nic.MAC != "" {
			deviceVal += fmt.Sprintf(",mac=%s", nic.MAC)
		}

		args = append(args,
			cmdline.ArgEntry{Key: "-netdev", Value: netdevVal, RuleName: r.Name()},
			cmdline.ArgEntry{Key: "-device", Value: deviceVal, RuleName: r.Name()},
		)
	}
	return args, nil
}

func (r *NetworkRule) Conflicts() []string { return nil }
