package dryrun

import (
	"context"
	"fmt"
	"os"

	"github.com/turtacn/go-qemu/internal/cmdline"
)

func runLayer1(ctx context.Context, cfg Config) (*LayerResult, error) {
	lr := &LayerResult{
		Layer: Layer1Static,
		OK:    true,
	}

	checks := []struct {
		name string
		fn   func() (bool, string)
	}{
		{"qemu-binary-exists", func() (bool, string) {
			if cfg.QEMUBinary == "" {
				return true, ""
			}
			if _, err := os.Stat(cfg.QEMUBinary); err != nil {
				return false, fmt.Sprintf("QEMU binary does not exist at %s", cfg.QEMUBinary)
			}
			return true, ""
		}},
		{"required-args-present", func() (bool, string) {
			required := []string{"-machine", "-m", "-smp"}
			for _, r := range required {
				if !argsContainKey(cfg.Args, r) {
					return false, fmt.Sprintf("Missing required parameter: %s", r)
				}
			}
			return true, ""
		}},
		{"memory-size-alignment", func() (bool, string) {
			if cfg.VMConfig != nil {
				memMB := cfg.VMConfig.Memory.SizeMB
				if memMB%4 != 0 {
					return false, fmt.Sprintf("Memory size %dMB is not a multiple of 4MB", memMB)
				}
			}
			return true, ""
		}},
		{"no-conflicting-args", func() (bool, string) {
			if err := cmdline.CheckConflictingArgs(cfg.Args); err != nil {
				return false, err.Error()
			}
			return true, ""
		}},
		{"qmp-socket-configured", func() (bool, string) {
			if !argsContainKey(cfg.Args, "-chardev") && !argsContainKey(cfg.Args, "-mon") {
				return false, "Missing QMP socket configuration (-chardev or -mon missing)"
			}
			return true, ""
		}},
	}

	for _, c := range checks {
		ok, msg := c.fn()
		lr.Checks = append(lr.Checks, CheckResult{Name: c.name, OK: ok, Message: msg})
		if !ok {
			lr.OK = false
		}
	}

	return lr, nil
}

func argsContainKey(args []string, key string) bool {
	for _, a := range args {
		if a == key {
			return true
		}
	}
	return false
}
