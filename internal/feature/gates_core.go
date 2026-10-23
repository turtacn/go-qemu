package feature

import "github.com/turtacn/go-qemu/internal/version"

var coreGates = []FeatureGate{
	{
		Name:        "core.query-qmp-schema",
		Description: "query-qmp-schema command",
		Available:   version.Since(version.V(2, 11, 0)),
		SchemaKey:   "command:query-qmp-schema",
	},
	{
		Name:        "core.query-version",
		Description: "query-version command",
		Available:   version.Since(version.V(0, 14, 0)),
		SchemaKey:   "command:query-version",
	},
}
