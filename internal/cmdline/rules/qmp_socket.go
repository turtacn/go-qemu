package rules

import (
	"fmt"

	"github.com/turtacn/go-qemu/internal/cmdline"
	"github.com/turtacn/go-qemu/internal/feature"
	"github.com/turtacn/go-qemu/internal/version"
	"github.com/turtacn/go-qemu/internal/vmconfig"
)

type QMPSocketRule struct{}

func (r *QMPSocketRule) Name() string { return "qmp.socket" }

func (r *QMPSocketRule) Applies(cfg *vmconfig.VMConfig, _ version.QEMUVersion, _ feature.Arch) bool {
	return true
}

func (r *QMPSocketRule) Apply(cfg *vmconfig.VMConfig, ver version.QEMUVersion, reg *feature.FeatureGateRegistry) ([]cmdline.ArgEntry, error) {
	sockPath := fmt.Sprintf("/var/run/qemu/vm%d.qmp", cfg.ID)
	if cfg.ID == 0 {
		sockPath = "/var/run/qemu/vm.qmp"
	}

	return []cmdline.ArgEntry{
		{
			Key:      "-chardev",
			Value:    fmt.Sprintf("socket,id=qmpmon,path=%s,server=on,wait=off", sockPath),
			RuleName: r.Name(),
		},
		{
			Key:      "-mon",
			Value:    "chardev=qmpmon,mode=control",
			RuleName: r.Name(),
		},
	}, nil
}

func (r *QMPSocketRule) Conflicts() []string { return nil }
