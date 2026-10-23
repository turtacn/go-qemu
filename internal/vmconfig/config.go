package vmconfig

type VMConfig struct {
	ID       int
	Name     string
	CPU      CPUConfig
	Memory   MemoryConfig
	Machine  MachineConfig
	Disks    []DiskConfig
	Networks []NICConfig
	Display  DisplayConfig
}

type CPUConfig struct {
	Cores   int
	Sockets int
	Threads int
	Model   string
}

type MemoryConfig struct {
	SizeMB int
}

type MachineConfig struct {
	Type         string // e.g. "pc-i440fx" or "pc-q35"
	USB          bool
	FixedVersion string // e.g. "pc-i440fx-5.0"
}

type DiskConfig struct {
	Bus    string // e.g. "scsi", "virtio", "ide"
	Index  int
	Path   string
	Format string // e.g. "raw", "qcow2"
	Cache  string
	AIO    string
}

type NICConfig struct {
	Model  string
	MAC    string
	Bridge string
}

type DisplayConfig struct {
	Type string // e.g. "vnc", "none", "std"
}
