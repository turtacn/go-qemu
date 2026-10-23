package version

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	v, err := Parse("5.0.0")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if v.Major != 5 || v.Minor != 0 || v.Micro != 0 {
		t.Errorf("got %v, expected 5.0.0", v)
	}

	v2, err := Parse("10.2.2-hci-6.8.0")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if v2.Major != 10 || v2.Minor != 2 || v2.Micro != 2 || v2.HCITag != "hci-6.8.0" {
		t.Errorf("got %v, expected 10.2.2-hci-6.8.0", v2)
	}
}

func TestVersionComparison(t *testing.T) {
	v5 := V(5, 0, 0)
	v7 := V(7, 0, 0)
	v10 := V(10, 2, 2)

	if !v10.GTE(v5) {
		t.Errorf("10.2.2 should be GTE 5.0.0")
	}
	if v5.GTE(v7) {
		t.Errorf("5.0.0 should NOT be GTE 7.0.0")
	}
	if !v5.LT(v7) {
		t.Errorf("5.0.0 should be LT 7.0.0")
	}

	r := Between(v5, v10)
	if !r.Contains(v7) {
		t.Errorf("7.0.0 should be in [5.0.0, 10.2.2)")
	}
	if r.Contains(v10) {
		t.Errorf("10.2.2 should NOT be in [5.0.0, 10.2.2)")
	}

	since5 := Since(v5)
	if !since5.Contains(v10) {
		t.Errorf("10.2.2 should be in Since(5.0.0)")
	}
}
