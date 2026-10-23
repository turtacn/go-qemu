package rules

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
