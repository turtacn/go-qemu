package rules

import "github.com/turtacn/go-qemu/internal/cmdline"

type Registry struct {
	rules []BuildRule
}

func DefaultRegistry() *Registry {
	return &Registry{
		rules: []BuildRule{
			&CPURule{},
			&MemoryRule{},
			&MachineTypeRule{},
			&DiskSelectorRule{},
			&NetworkRule{},
			&QMPSocketRule{},
		},
	}
}

func (r *Registry) Rules() []BuildRule {
	return r.rules
}

func init() {
	reg := DefaultRegistry()
	cmdlineRules := make([]cmdline.Rule, len(reg.rules))
	for i, r := range reg.rules {
		cmdlineRules[i] = r
	}
	cmdline.RegisterDefaultRules(cmdlineRules)
}
