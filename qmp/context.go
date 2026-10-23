package qmp

import (
	"context"
	"encoding/json"
	"fmt"
)

// Client is a context-aware interface for communicating with QMP.
type Client interface {
	Connect(ctx context.Context) error
	Disconnect() error
	Execute(ctx context.Context, cmd string, args interface{}) ([]byte, error)
	QueryVersion(ctx context.Context) (*Version, error)
	QueryMigrateCapabilities(ctx context.Context) ([]MigrateCapabilityStatus, error)
	SetMigrateCapabilities(ctx context.Context, caps []MigrateCapabilityStatus) error
	Migrate(ctx context.Context, uri string) error
	Resume(ctx context.Context) error
	Quit(ctx context.Context) error
	Close() error
}

// MigrateCapabilityStatus represents a migration capability name and its boolean status.
type MigrateCapabilityStatus struct {
	Capability string `json:"capability"`
	State      bool   `json:"state"`
}

// ContextMonitor wraps a QMP Monitor to provide Context-aware methods.
type ContextMonitor struct {
	mon Monitor
}

// NewContextMonitor wraps an existing Monitor in a ContextMonitor.
func NewContextMonitor(mon Monitor) *ContextMonitor {
	return &ContextMonitor{mon: mon}
}

// Execute executes a QMP command with context support.
func (cm *ContextMonitor) Execute(ctx context.Context, cmd string, args interface{}) ([]byte, error) {
	cmdStruct := Command{
		Execute: cmd,
		Args:    args,
	}
	raw, err := json.Marshal(cmdStruct)
	if err != nil {
		return nil, &TransportError{Op: "Marshal", Err: err}
	}

	type result struct {
		out []byte
		err error
	}
	ch := make(chan result, 1)
	go func() {
		out, err := cm.mon.Run(raw)
		ch <- result{out: out, err: err}
	}()

	select {
	case <-ctx.Done():
		return nil, &TransportError{Op: cmd, Err: ctx.Err()}
	case res := <-ch:
		if res.err != nil {
			return nil, res.err
		}
		var resp response
		if err := json.Unmarshal(res.out, &resp); err == nil {
			if resp.Error.Desc != "" {
				return nil, &CommandError{
					Class:   ErrorClass(resp.Error.Class),
					Desc:    resp.Error.Desc,
					Command: cmd,
				}
			}
		}
		return res.out, nil
	}
}

// Dial connects to a QMP socket at socketPath and performs the initial handshake.
func Dial(ctx context.Context, socketPath string) (*ContextMonitor, error) {
	if err := WaitForSocket(ctx, socketPath); err != nil {
		return nil, &TransportError{Op: "WaitForSocket", Err: err}
	}
	mon, err := NewSocketMonitor("unix", socketPath, 2)
	if err != nil {
		return nil, &TransportError{Op: "NewSocketMonitor", Err: err}
	}
	if err := mon.Connect(); err != nil {
		return nil, &TransportError{Op: "Connect", Err: err}
	}

	cm := NewContextMonitor(mon)
	if _, err := cm.Execute(ctx, "qmp_capabilities", nil); err != nil {
		mon.Disconnect()
		return nil, fmt.Errorf("qmp_capabilities handshake failed: %w", err)
	}

	return cm, nil
}

func (cm *ContextMonitor) Connect(ctx context.Context) error {
	return cm.mon.Connect()
}

func (cm *ContextMonitor) Disconnect() error {
	return cm.mon.Disconnect()
}

func (cm *ContextMonitor) Close() error {
	return cm.mon.Disconnect()
}

func (cm *ContextMonitor) QueryVersion(ctx context.Context) (*Version, error) {
	raw, err := cm.Execute(ctx, "query-version", nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Return struct {
			Package string `json:"package"`
			QEMU    struct {
				Major int `json:"major"`
				Minor int `json:"minor"`
				Micro int `json:"micro"`
			} `json:"qemu"`
		} `json:"return"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	v := &Version{
		Package: resp.Return.Package,
	}
	v.QEMU.Major = resp.Return.QEMU.Major
	v.QEMU.Minor = resp.Return.QEMU.Minor
	v.QEMU.Micro = resp.Return.QEMU.Micro
	return v, nil
}

func (cm *ContextMonitor) QueryMigrateCapabilities(ctx context.Context) ([]MigrateCapabilityStatus, error) {
	raw, err := cm.Execute(ctx, "query-migrate-capabilities", nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Return []MigrateCapabilityStatus `json:"return"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	return resp.Return, nil
}

func (cm *ContextMonitor) SetMigrateCapabilities(ctx context.Context, caps []MigrateCapabilityStatus) error {
	args := map[string]interface{}{
		"capabilities": caps,
	}
	_, err := cm.Execute(ctx, "migrate-set-capabilities", args)
	return err
}

func (cm *ContextMonitor) Migrate(ctx context.Context, uri string) error {
	args := map[string]interface{}{
		"uri": uri,
	}
	_, err := cm.Execute(ctx, "migrate", args)
	return err
}

func (cm *ContextMonitor) Resume(ctx context.Context) error {
	_, err := cm.Execute(ctx, "cont", nil)
	return err
}

func (cm *ContextMonitor) Quit(ctx context.Context) error {
	_, err := cm.Execute(ctx, "quit", nil)
	return err
}
