package feature

import "github.com/turtacn/go-qemu/internal/version"

var deviceGates = []FeatureGate{
	{
		Name:        "disk.blockdev-preferred",
		Description: "Use -blockdev instead of -drive (recommended in 7.0+)",
		Available:   version.Since(version.V(7, 0, 0)),
		SchemaKey:   "command:blockdev-add",
	},
	{
		Name:        "disk.drive-deprecated",
		Description: "-drive parameter marked as deprecated",
		Available:   version.Since(version.V(7, 0, 0)),
	},
	{
		Name:        "memory.explicit-backend",
		Description: "Explicit memory backend object (memory-backend-ram)",
		Available:   version.Since(version.V(5, 0, 0)),
		SchemaKey:   "command:object-add",
	},
	{
		Name:        "device.virtio-scsi",
		Description: "virtio-scsi-pci device",
		Available:   version.Since(version.V(1, 4, 0)),
	},
}
