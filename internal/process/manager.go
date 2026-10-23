package process

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/turtacn/go-qemu/qmp"
)

type Config struct {
	Binary    string
	Args      []string
	QMPSocket string
}

type Process struct {
	cmd       *exec.Cmd
	QMPSocket string
}

func Start(ctx context.Context, cfg Config) (*Process, error) {
	if _, err := os.Stat(cfg.Binary); err != nil {
		return nil, fmt.Errorf("qemu binary not found at %s: %w", cfg.Binary, err)
	}

	cmd := exec.CommandContext(ctx, cfg.Binary, cfg.Args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start qemu process: %w", err)
	}

	p := &Process{
		cmd:       cmd,
		QMPSocket: cfg.QMPSocket,
	}

	return p, nil
}

func (p *Process) WaitReady(ctx context.Context, timeout time.Duration) error {
	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if p.QMPSocket != "" {
		if err := qmp.WaitForSocket(waitCtx, p.QMPSocket); err != nil {
			p.Kill()
			return fmt.Errorf("qmp socket wait timeout: %w", err)
		}
	}
	return nil
}

func (p *Process) Kill() error {
	if p.cmd != nil && p.cmd.Process != nil {
		_ = p.cmd.Process.Kill()
		_ = p.cmd.Wait()
	}
	return nil
}

func (p *Process) Pid() int {
	if p.cmd != nil && p.cmd.Process != nil {
		return p.cmd.Process.Pid
	}
	return 0
}
