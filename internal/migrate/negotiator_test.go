package migrate

import (
	"context"
	"testing"

	"github.com/turtacn/go-qemu/qmp"
)

type mockQMPClient struct {
	caps        []qmp.MigrateCapabilityStatus
	lastSetCaps []qmp.MigrateCapabilityStatus
}

func (m *mockQMPClient) Connect(ctx context.Context) error { return nil }
func (m *mockQMPClient) Disconnect() error                 { return nil }
func (m *mockQMPClient) Close() error                      { return nil }
func (m *mockQMPClient) Execute(ctx context.Context, cmd string, args interface{}) ([]byte, error) {
	return nil, nil
}
func (m *mockQMPClient) QueryVersion(ctx context.Context) (*qmp.Version, error) {
	return &qmp.Version{}, nil
}
func (m *mockQMPClient) QueryMigrateCapabilities(ctx context.Context) ([]qmp.MigrateCapabilityStatus, error) {
	return m.caps, nil
}
func (m *mockQMPClient) SetMigrateCapabilities(ctx context.Context, caps []qmp.MigrateCapabilityStatus) error {
	m.lastSetCaps = caps
	return nil
}
func (m *mockQMPClient) Migrate(ctx context.Context, uri string) error { return nil }
func (m *mockQMPClient) Resume(ctx context.Context) error              { return nil }
func (m *mockQMPClient) Quit(ctx context.Context) error                { return nil }

func TestNegotiateDisablesPostcopyPreempt(t *testing.T) {
	srcClient := &mockQMPClient{
		caps: []qmp.MigrateCapabilityStatus{
			{Capability: "xbzrle", State: true},
			{Capability: "postcopy-ram", State: true},
		},
	}

	dstClient := &mockQMPClient{
		caps: []qmp.MigrateCapabilityStatus{
			{Capability: "xbzrle", State: true},
			{Capability: "postcopy-ram", State: true},
			{Capability: "postcopy-preempt", State: true}, // 10.2.2 feature
		},
	}

	negotiator := NewNegotiator(srcClient, dstClient)
	result, err := negotiator.Negotiate(context.Background())

	if err != nil {
		t.Fatalf("Negotiate failed: %v", err)
	}

	foundDisabled := false
	for _, disabled := range result.Disabled {
		if disabled == "postcopy-preempt" {
			foundDisabled = true
		}
	}
	if !foundDisabled {
		t.Errorf("postcopy-preempt must be disabled because source QEMU 5.0 does not support it")
	}

	for _, cap := range dstClient.lastSetCaps {
		if cap.Capability == "postcopy-preempt" && cap.State {
			t.Errorf("postcopy-preempt state in target must be false")
		}
	}
}
