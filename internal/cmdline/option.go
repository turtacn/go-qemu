package cmdline

import (
	"github.com/turtacn/go-qemu/internal/feature"
	"github.com/turtacn/go-qemu/internal/version"
)

type Option func(*CmdlineBuilder)

func WithQEMUVersion(ver version.QEMUVersion) Option {
	return func(b *CmdlineBuilder) {
		b.version = ver
	}
}

func WithArch(arch feature.Arch) Option {
	return func(b *CmdlineBuilder) {
		b.arch = arch
	}
}

func WithRegistry(reg *feature.FeatureGateRegistry) Option {
	return func(b *CmdlineBuilder) {
		b.featureReg = reg
	}
}
