package rules

import (
	"github.com/turtacn/go-qemu/internal/cmdline"
	"github.com/turtacn/go-qemu/internal/feature"
	"github.com/turtacn/go-qemu/internal/version"
	"github.com/turtacn/go-qemu/internal/vmconfig"
)

type BuildRule interface {
	Name() string
	Applies(cfg *vmconfig.VMConfig, ver version.QEMUVersion, arch feature.Arch) bool
	Apply(cfg *vmconfig.VMConfig, ver version.QEMUVersion, reg *feature.FeatureGateRegistry) ([]cmdline.ArgEntry, error)
	Conflicts() []string
}
