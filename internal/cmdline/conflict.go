package cmdline

import "fmt"

// CheckConflictingArgs verifies there are no incompatible arguments in args slice.
func CheckConflictingArgs(args []string) error {
	hasDrive := false
	hasBlockdev := false

	for i := 0; i < len(args); i++ {
		if args[i] == "-drive" {
			hasDrive = true
		}
		if args[i] == "-blockdev" {
			hasBlockdev = true
		}
	}

	if hasDrive && hasBlockdev {
		return fmt.Errorf("conflict detected: both -drive and -blockdev are present")
	}

	return nil
}
