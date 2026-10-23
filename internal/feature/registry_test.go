package feature

import (
	"testing"

	"github.com/turtacn/go-qemu/internal/version"
	"github.com/turtacn/go-qemu/qapi"
)

func TestPostcopyPreemptNotAvailableIn500(t *testing.T) {
	reg := New(version.V(5, 0, 0), ArchX86)
	enabled, reason, err := reg.IsEnabled("migrate.postcopy-preempt")
	if err != nil {
		t.Fatalf("IsEnabled failed: %v", err)
	}
	if enabled {
		t.Errorf("postcopy-preempt should NOT be enabled in QEMU 5.0.0")
	}
	if reason.VersionOK {
		t.Errorf("postcopy-preempt VersionOK should be false for 5.0.0")
	}
}

func TestPostcopyPreemptAvailableIn10(t *testing.T) {
	reg := New(version.V(10, 2, 2), ArchX86)
	enabled, _, err := reg.IsEnabled("migrate.postcopy-preempt")
	if err != nil {
		t.Fatalf("IsEnabled failed: %v", err)
	}
	if !enabled {
		t.Errorf("postcopy-preempt SHOULD be enabled in QEMU 10.2.2")
	}
}

func TestLiveSchemaOverridesVersion(t *testing.T) {
	reg := New(version.V(10, 2, 2), ArchX86)
	mockSchema := qapi.NewQAPISchema("10.2.2") // empty schema, no postcopy-preempt
	reg.InjectLiveSchema(mockSchema)

	enabled, reason, err := reg.IsEnabled("migrate.postcopy-preempt")
	if err != nil {
		t.Fatalf("IsEnabled failed: %v", err)
	}
	if enabled {
		t.Errorf("postcopy-preempt should be disabled when live schema lacks it")
	}
	if !reason.VersionOK {
		t.Errorf("version check should pass for 10.2.2")
	}
	if reason.SchemaOK {
		t.Errorf("schema check should fail when key is absent in live schema")
	}
}

func TestBlockdevPreferred(t *testing.T) {
	reg5 := New(version.V(5, 0, 0), ArchX86)
	enabled5, _, _ := reg5.IsEnabled("disk.blockdev-preferred")
	if enabled5 {
		t.Errorf("disk.blockdev-preferred should be false in 5.0.0")
	}

	reg10 := New(version.V(10, 2, 2), ArchX86)
	enabled10, _, _ := reg10.IsEnabled("disk.blockdev-preferred")
	if !enabled10 {
		t.Errorf("disk.blockdev-preferred should be true in 10.2.2")
	}
}
