package scenarios

import (
	"context"
	"fmt"

	"github.com/turtacn/go-qemu/internal/vmconfig"
)

type Cfg2CmdCompatScenario struct {
	cfg *vmconfig.VMConfig
}

func (s *Cfg2CmdCompatScenario) Name() string { return "cfg2cmd-compat-pve" }
func (s *Cfg2CmdCompatScenario) Description() string {
	return "Proxmox cfg2cmd historical format compatibility check"
}
func (s *Cfg2CmdCompatScenario) Prerequisites(ctx context.Context) error { return nil }
func (s *Cfg2CmdCompatScenario) Setup(ctx context.Context) error         { return nil }
func (s *Cfg2CmdCompatScenario) Execute(ctx context.Context) error {
	cfg, err := vmconfig.LoadFromHCI(resolveFixturePath("harness/env/fixtures/vm_pve_compat.yaml"))
	if err != nil {
		return err
	}
	s.cfg = cfg
	return nil
}
func (s *Cfg2CmdCompatScenario) Verify(ctx context.Context) error {
	if s.cfg == nil {
		return fmt.Errorf("loaded config is nil")
	}
	if s.cfg.ID != 101 {
		return fmt.Errorf("expected VM ID 101, got %d", s.cfg.ID)
	}
	return nil
}
func (s *Cfg2CmdCompatScenario) Teardown(ctx context.Context) error { return nil }
