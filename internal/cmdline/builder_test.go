package cmdline

import (
	"testing"

	"github.com/turtacn/go-qemu/internal/feature"
	"github.com/turtacn/go-qemu/internal/version"
	"github.com/turtacn/go-qemu/internal/vmconfig"
)

type dummyDiskRule struct{}

func (r *dummyDiskRule) Name() string { return "dummy.disk" }
func (r *dummyDiskRule) Applies(cfg *vmconfig.VMConfig, _ version.QEMUVersion, _ feature.Arch) bool {
	return len(cfg.Disks) > 0
}
func (r *dummyDiskRule) Apply(cfg *vmconfig.VMConfig, ver version.QEMUVersion, reg *feature.FeatureGateRegistry) ([]ArgEntry, error) {
	blockdevPreferred, _, _ := reg.IsEnabled("disk.blockdev-preferred")
	if blockdevPreferred {
		return []ArgEntry{{Key: "-blockdev", Value: "driver=raw,node-name=drive0", RuleName: r.Name()}}, nil
	}
	return []ArgEntry{{Key: "-drive", Value: "file=/dev/sda,format=raw", RuleName: r.Name()}}, nil
}
func (r *dummyDiskRule) Conflicts() []string { return nil }

func containsArgKey(report *BuildReport, key string) bool {
	for _, a := range report.Args {
		if a.Key == key {
			return true
		}
	}
	return false
}

func TestDiskParamVersionDifference(t *testing.T) {
	cfg := &vmconfig.VMConfig{
		Disks: []vmconfig.DiskConfig{{Bus: "scsi", Index: 0, Path: "/dev/sda"}},
	}

	b5, _ := NewBuilder(WithQEMUVersion(version.V(5, 0, 0)))
	b5.SetRules([]Rule{&dummyDiskRule{}})
	_, report5, err := b5.Build(cfg)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if !containsArgKey(report5, "-drive") || containsArgKey(report5, "-blockdev") {
		t.Errorf("5.0.0 should use -drive, got %+v", report5.Args)
	}

	b10, _ := NewBuilder(WithQEMUVersion(version.V(10, 2, 2)))
	b10.SetRules([]Rule{&dummyDiskRule{}})
	_, report10, err := b10.Build(cfg)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if containsArgKey(report10, "-drive") || !containsArgKey(report10, "-blockdev") {
		t.Errorf("10.2.2 should use -blockdev, got %+v", report10.Args)
	}
}

func TestCmdlineDiffCapturesDiskChange(t *testing.T) {
	cfg := &vmconfig.VMConfig{
		Disks: []vmconfig.DiskConfig{{Bus: "scsi", Index: 0, Path: "/dev/sda"}},
	}

	b5, _ := NewBuilder(WithQEMUVersion(version.V(5, 0, 0)))
	b5.SetRules([]Rule{&dummyDiskRule{}})
	_, report5, _ := b5.Build(cfg)

	b10, _ := NewBuilder(WithQEMUVersion(version.V(10, 2, 2)))
	b10.SetRules([]Rule{&dummyDiskRule{}})
	_, report10, _ := b10.Build(cfg)

	diff := report10.DiffFrom(version.V(5, 0, 0), report5)

	foundDriveInRemoved := false
	for _, r := range diff.Removed {
		if r.Key == "-drive" {
			foundDriveInRemoved = true
		}
	}
	if !foundDriveInRemoved {
		t.Errorf("expected -drive in diff.Removed")
	}

	foundBlockdevInAdded := false
	for _, a := range diff.Added {
		if a.Key == "-blockdev" {
			foundBlockdevInAdded = true
		}
	}
	if !foundBlockdevInAdded {
		t.Errorf("expected -blockdev in diff.Added")
	}
}
