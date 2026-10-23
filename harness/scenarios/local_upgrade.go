package scenarios

import (
	"context"
	"fmt"

	"github.com/turtacn/go-qemu/internal/upgrade"
)

type LocalUpgradeScenario struct {
	vmid     int
	upgrader *upgrade.LocalUpgrader
	result   *upgrade.UpgradeResult
}

func NewLocalUpgradeScenario(vmid int) *LocalUpgradeScenario {
	return &LocalUpgradeScenario{
		vmid: vmid,
		upgrader: upgrade.NewLocalUpgrader(upgrade.LocalUpgradeConfig{
			VMID:         vmid,
			VMConfigPath: resolveFixturePath("harness/env/fixtures/vm_basic_x86.yaml"),
		}),
	}
}

func (s *LocalUpgradeScenario) Name() string { return "local-upgrade-5.0-to-10.2.2" }
func (s *LocalUpgradeScenario) Description() string {
	return "Local QEMU hot upgrade: 5.0.0 to 10.2.2 with zero VM downtime"
}
func (s *LocalUpgradeScenario) Prerequisites(ctx context.Context) error {
	return nil
}
func (s *LocalUpgradeScenario) Setup(ctx context.Context) error {
	s.vmid = 100
	s.upgrader = upgrade.NewLocalUpgrader(upgrade.LocalUpgradeConfig{
		VMID:         s.vmid,
		VMConfigPath: resolveFixturePath("harness/env/fixtures/vm_basic_x86.yaml"),
	})
	return nil
}
func (s *LocalUpgradeScenario) Execute(ctx context.Context) error {
	res, err := s.upgrader.RunWithResult(ctx)
	if err != nil {
		return fmt.Errorf("local upgrade execution failed: %w", err)
	}
	s.result = res
	return nil
}
func (s *LocalUpgradeScenario) Verify(ctx context.Context) error {
	if s.result == nil {
		return fmt.Errorf("upgrade result is nil")
	}
	if s.result.DstVersion.String() != "10.2.2" {
		return fmt.Errorf("expected target version 10.2.2, got %s", s.result.DstVersion)
	}
	return nil
}
func (s *LocalUpgradeScenario) Teardown(ctx context.Context) error { return nil }
