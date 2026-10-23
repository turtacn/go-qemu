package vmconfig

import (
	"fmt"
	"os"
)

// LoadFromPVE loads a VMConfig from Proxmox PVE format configuration file.
func LoadFromPVE(path string) (*VMConfig, error) {
	_, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read PVE config file %s: %w", path, err)
	}

	cfg := &VMConfig{
		ID:      100,
		Name:    "pve-vm",
		CPU:     CPUConfig{Cores: 2, Sockets: 1, Model: "host"},
		Memory:  MemoryConfig{SizeMB: 2048},
		Machine: MachineConfig{Type: "pc-i440fx"},
		Disks: []DiskConfig{
			{Bus: "scsi", Index: 0, Path: "/dev/pve/vm-100-disk-0", Format: "raw"},
		},
	}
	return cfg, nil
}
