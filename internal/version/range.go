package version

// VersionRange represents [From, To).
// To being nil means no upper bound.
type VersionRange struct {
	From QEMUVersion
	To   *QEMUVersion // nil = no upper bound
}

// Since creates a VersionRange starting from v with no upper bound.
func Since(v QEMUVersion) VersionRange {
	return VersionRange{From: v}
}

// Between creates a VersionRange [from, to).
func Between(from, to QEMUVersion) VersionRange {
	return VersionRange{From: from, To: &to}
}

// Contains checks if v falls into [From, To).
func (r VersionRange) Contains(v QEMUVersion) bool {
	if !v.GTE(r.From) {
		return false
	}
	if r.To != nil && !v.LT(*r.To) {
		return false
	}
	return true
}
