package feature

import (
	"fmt"

	"github.com/turtacn/go-qemu/internal/version"
)

type Arch string

const (
	ArchX86 Arch = "x86_64"
	ArchARM Arch = "aarch64"
)

type FeatureGate struct {
	Name        string
	Description string
	Available   version.VersionRange
	Deprecated  *version.QEMUVersion
	Removed     *version.QEMUVersion
	Arches      []Arch
	SchemaKey   string
}

type Reason struct {
	VersionOK bool
	ArchOK    bool
	SchemaOK  bool
	Detail    string
}

func (r Reason) Available() bool {
	return r.VersionOK && r.ArchOK && r.SchemaOK
}

func (r Reason) String() string {
	if r.Available() {
		return "Available"
	}
	return fmt.Sprintf("Unavailable (VersionOK=%t, ArchOK=%t, SchemaOK=%t): %s",
		r.VersionOK, r.ArchOK, r.SchemaOK, r.Detail)
}
