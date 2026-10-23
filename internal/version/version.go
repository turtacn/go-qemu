package version

import (
	"fmt"
	"strconv"
	"strings"
)

// QEMUVersion represents a specific QEMU version.
type QEMUVersion struct {
	Major  int
	Minor  int
	Micro  int
	HCITag string // e.g. "hci-6.8.0", empty for community version
}

// V creates a QEMUVersion with major, minor, micro.
func V(major, minor, micro int) QEMUVersion {
	return QEMUVersion{Major: major, Minor: minor, Micro: micro}
}

// WithHCITag attaches an HCI tag.
func (v QEMUVersion) WithHCITag(tag string) QEMUVersion {
	v.HCITag = tag
	return v
}

// GTE checks if v >= other.
func (v QEMUVersion) GTE(other QEMUVersion) bool {
	if v.Major != other.Major {
		return v.Major > other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor > other.Minor
	}
	return v.Micro >= other.Micro
}

// LT checks if v < other.
func (v QEMUVersion) LT(other QEMUVersion) bool {
	return !v.GTE(other)
}

// String implements Stringer.
func (v QEMUVersion) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Micro)
	if v.HCITag != "" {
		s += "-" + v.HCITag
	}
	return s
}

// Parse parses a version string like "5.0.0" or "10.2.2-hci-6.8.0".
func Parse(s string) (QEMUVersion, error) {
	parts := strings.SplitN(s, "-", 2)
	nums := strings.Split(parts[0], ".")
	if len(nums) != 3 {
		return QEMUVersion{}, fmt.Errorf("invalid version format: %q", s)
	}
	v := QEMUVersion{}
	var err error
	if v.Major, err = strconv.Atoi(nums[0]); err != nil {
		return v, err
	}
	if v.Minor, err = strconv.Atoi(nums[1]); err != nil {
		return v, err
	}
	if v.Micro, err = strconv.Atoi(nums[2]); err != nil {
		return v, err
	}
	if len(parts) > 1 {
		v.HCITag = parts[1]
	}
	return v, nil
}
