package rules

import (
	"fmt"

	"github.com/turtacn/go-qemu/internal/cmdline"
	"github.com/turtacn/go-qemu/internal/feature"
	"github.com/turtacn/go-qemu/internal/version"
	"github.com/turtacn/go-qemu/internal/vmconfig"
)

type MachineTypeRule struct{}

func (r *MachineTypeRule) Name() string { return "machine.type" }

func (r *MachineTypeRule) Applies(cfg *vmconfig.VMConfig, _ version.QEMUVersion, _ feature.Arch) bool {
	return true
}

func (r *MachineTypeRule) Apply(cfg *vmconfig.VMConfig, ver version.QEMUVersion, reg *feature.FeatureGateRegistry) ([]cmdline.ArgEntry, error) {
	baseType := cfg.Machine.Type
	if baseType == "" {
		baseType = "pc-i440fx"
	}

	machineVer := fmt.Sprintf("%d.%d", ver.Major, ver.Minor)
	fullType := fmt.Sprintf("%s-%s", baseType, machineVer)

	value := fmt.Sprintf("%s,accel=kvm", fullType)
	if !cfg.Machine.USB {
		value += ",usb=off"
	}

	riskLevel := cmdline.RiskNone
	if cfg.Machine.FixedVersion != "" && cfg.Machine.FixedVersion != fullType {
		riskLevel = cmdline.RiskBehaviorChange
	}

	return []cmdline.ArgEntry{
		{
			Key:      "-machine",
			Value:    value,
			RuleName: r.Name(),
			Reason:   fmt.Sprintf("Machine type version follows QEMU %s", ver),
			Risk:     riskLevel,
		},
	}, nil
}

func (r *MachineTypeRule) Conflicts() []string { return nil }
