package vmconfig

import "fmt"

func Validate(cfg *VMConfig) error {
	if cfg.Memory.SizeMB <= 0 {
		return fmt.Errorf("memory size must be > 0 MB")
	}
	if cfg.CPU.Cores <= 0 {
		cfg.CPU.Cores = 1
	}
	if cfg.CPU.Sockets <= 0 {
		cfg.CPU.Sockets = 1
	}
	if cfg.Machine.Type == "" {
		cfg.Machine.Type = "pc-i440fx"
	}
	return nil
}
