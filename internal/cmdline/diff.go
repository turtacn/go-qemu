package cmdline

import "github.com/turtacn/go-qemu/internal/version"

// DiffFrom compares current BuildReport against another BuildReport.
func (r *BuildReport) DiffFrom(otherVer version.QEMUVersion, other *BuildReport) *CmdlineDiff {
	diff := &CmdlineDiff{
		FromVersion: otherVer,
		ToVersion:   r.QEMUVersion,
	}

	fromMap := make(map[string]ArgEntry)
	for _, a := range other.Args {
		fromMap[a.Key] = a
	}

	toMap := make(map[string]ArgEntry)
	for _, a := range r.Args {
		toMap[a.Key] = a
	}

	for key, toArg := range toMap {
		if fromArg, exists := fromMap[key]; !exists {
			diff.Added = append(diff.Added, toArg)
		} else if fromArg.Value != toArg.Value {
			diff.Changed = append(diff.Changed, ArgChange{From: fromArg, To: toArg})
		} else {
			diff.Unchanged = append(diff.Unchanged, toArg)
		}
	}

	for key, fromArg := range fromMap {
		if _, exists := toMap[key]; !exists {
			diff.Removed = append(diff.Removed, fromArg)
		}
	}

	return diff
}
