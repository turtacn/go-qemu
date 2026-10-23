package dryrun

import (
	"github.com/turtacn/go-qemu/internal/capability"
	"github.com/turtacn/go-qemu/internal/vmconfig"
)

type LayerMask int

const (
	Layer1Static LayerMask = 1 << iota
	Layer2Smoke
	Layer3Runtime
	AllLayers = Layer1Static | Layer2Smoke | Layer3Runtime
)

type Config struct {
	Args           []string
	QEMUBinary     string
	Layers         LayerMask
	CapabilityMeta *capability.Meta
	VMConfig       *vmconfig.VMConfig
	QMPSocketDir   string
}

type Result struct {
	OK      bool
	Layers  map[LayerMask]*LayerResult
	Summary string
}

type LayerResult struct {
	Layer    LayerMask
	OK       bool
	Checks   []CheckResult
	Duration string
}

type CheckResult struct {
	Name    string
	OK      bool
	Message string
}
