package feature

import "github.com/turtacn/go-qemu/internal/version"

var migrationGates = []FeatureGate{
	{
		Name:        "migrate.postcopy-ram",
		Description: "Postcopy migration mode",
		Available:   version.Since(version.V(2, 6, 0)),
		SchemaKey:   "enum-value:MigrationCapability:postcopy-ram",
	},
	{
		Name:        "migrate.postcopy-preempt",
		Description: "Postcopy preempt mode (introduced in QEMU 8.0, not supported in 5.0)",
		Available:   version.Since(version.V(8, 0, 0)),
		SchemaKey:   "enum-value:MigrationCapability:postcopy-preempt",
	},
	{
		Name:        "migrate.multifd",
		Description: "Multifd parallel migration",
		Available:   version.Since(version.V(4, 0, 0)),
		SchemaKey:   "enum-value:MigrationCapability:multifd",
	},
	{
		Name:        "migrate.dirty-bitmaps",
		Description: "Dirty bitmaps migration sync",
		Available:   version.Since(version.V(4, 2, 0)),
		SchemaKey:   "enum-value:MigrationCapability:dirty-bitmaps",
	},
	{
		Name:        "migrate.set-parameters",
		Description: "migrate-set-parameters command",
		Available:   version.Since(version.V(2, 8, 0)),
		SchemaKey:   "command:migrate-set-parameters",
	},
	{
		Name:        "migrate.calc-dirty-rate",
		Description: "calc-dirty-rate command (introduced in QEMU 7.0)",
		Available:   version.Since(version.V(7, 0, 0)),
		SchemaKey:   "command:calc-dirty-rate",
	},
}
