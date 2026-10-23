package scenarios

import (
	"context"
	"fmt"

	"github.com/turtacn/go-qemu/internal/migrate"
	"github.com/turtacn/go-qemu/qmp"
)

type mockClient struct {
	caps []qmp.MigrateCapabilityStatus
}

func (m *mockClient) Connect(ctx context.Context) error { return nil }
func (m *mockClient) Disconnect() error                 { return nil }
func (m *mockClient) Close() error                      { return nil }
func (m *mockClient) Execute(ctx context.Context, cmd string, args interface{}) ([]byte, error) {
	return nil, nil
}
func (m *mockClient) QueryVersion(ctx context.Context) (*qmp.Version, error) {
	return &qmp.Version{}, nil
}
func (m *mockClient) QueryMigrateCapabilities(ctx context.Context) ([]qmp.MigrateCapabilityStatus, error) {
	return m.caps, nil
}
func (m *mockClient) SetMigrateCapabilities(ctx context.Context, caps []qmp.MigrateCapabilityStatus) error {
	return nil
}
func (m *mockClient) Migrate(ctx context.Context, uri string) error { return nil }
func (m *mockClient) Resume(ctx context.Context) error              { return nil }
func (m *mockClient) Quit(ctx context.Context) error                { return nil }

type CrossVersionMigrateScenario struct {
	disabledCaps []string
}

func (s *CrossVersionMigrateScenario) Name() string {
	return "cross-version-migration-negotiation"
}
func (s *CrossVersionMigrateScenario) Description() string {
	return "Migration capability negotiation between QEMU 5.0 and 10.2.2"
}
func (s *CrossVersionMigrateScenario) Prerequisites(ctx context.Context) error { return nil }
func (s *CrossVersionMigrateScenario) Setup(ctx context.Context) error         { return nil }
func (s *CrossVersionMigrateScenario) Execute(ctx context.Context) error {
	src := &mockClient{caps: []qmp.MigrateCapabilityStatus{{Capability: "postcopy-ram", State: true}}}
	dst := &mockClient{caps: []qmp.MigrateCapabilityStatus{
		{Capability: "postcopy-ram", State: true},
		{Capability: "postcopy-preempt", State: true},
	}}

	negotiator := migrate.NewNegotiator(src, dst)
	res, err := negotiator.Negotiate(ctx)
	if err != nil {
		return err
	}
	s.disabledCaps = res.Disabled
	return nil
}
func (s *CrossVersionMigrateScenario) Verify(ctx context.Context) error {
	for _, d := range s.disabledCaps {
		if d == "postcopy-preempt" {
			return nil
		}
	}
	return fmt.Errorf("expected postcopy-preempt to be disabled during 5.0 -> 10.2.2 migration")
}
func (s *CrossVersionMigrateScenario) Teardown(ctx context.Context) error { return nil }
