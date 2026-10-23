package cmdline

import (
	"fmt"

	"github.com/turtacn/go-qemu/internal/feature"
	"github.com/turtacn/go-qemu/internal/version"
	"github.com/turtacn/go-qemu/internal/vmconfig"
)

type Rule interface {
	Name() string
	Applies(cfg *vmconfig.VMConfig, ver version.QEMUVersion, arch feature.Arch) bool
	Apply(cfg *vmconfig.VMConfig, ver version.QEMUVersion, reg *feature.FeatureGateRegistry) ([]ArgEntry, error)
	Conflicts() []string
}

type CmdlineBuilder struct {
	version    version.QEMUVersion
	arch       feature.Arch
	featureReg *feature.FeatureGateRegistry
	rules      []Rule
}

func NewBuilder(opts ...Option) (*CmdlineBuilder, error) {
	b := &CmdlineBuilder{
		version: version.V(5, 0, 0),
		arch:    feature.ArchX86,
	}

	for _, opt := range opts {
		opt(b)
	}

	if b.featureReg == nil {
		b.featureReg = feature.New(b.version, b.arch)
	}

	return b, nil
}

func (b *CmdlineBuilder) SetRules(rules []Rule) {
	b.rules = rules
}

func (b *CmdlineBuilder) Build(cfg *vmconfig.VMConfig) ([]string, *BuildReport, error) {
	report := &BuildReport{
		QEMUVersion: b.version,
		Arch:        b.arch,
	}

	for _, rule := range b.rules {
		if !rule.Applies(cfg, b.version, b.arch) {
			continue
		}

		entries, err := rule.Apply(cfg, b.version, b.featureReg)
		if err != nil {
			return nil, nil, fmt.Errorf("rule %s failed: %w", rule.Name(), err)
		}

		for _, entry := range entries {
			report.Args = append(report.Args, entry)
			if entry.Deprecated {
				report.Warnings = append(report.Warnings,
					fmt.Sprintf("Parameter %s is deprecated in QEMU %s", entry.Key, b.version))
			}
		}
	}

	argSlice := report.ToArgSlice()
	if err := CheckConflictingArgs(argSlice); err != nil {
		report.Errors = append(report.Errors, err.Error())
		return nil, report, err
	}

	return argSlice, report, nil
}
