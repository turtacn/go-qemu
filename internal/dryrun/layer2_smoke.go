package dryrun

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/turtacn/go-qemu/internal/process"
	"github.com/turtacn/go-qemu/qmp"
)

func runLayer2(ctx context.Context, cfg Config) (*LayerResult, error) {
	lr := &LayerResult{
		Layer: Layer2Smoke,
		OK:    true,
	}

	if cfg.QEMUBinary == "" {
		lr.Checks = append(lr.Checks, CheckResult{
			Name:    "qemu-binary-check",
			OK:      false,
			Message: "QEMU binary path not provided for Layer 2 smoke check",
		})
		lr.OK = false
		return lr, nil
	}

	socketDir := cfg.QMPSocketDir
	if socketDir == "" {
		socketDir = os.TempDir()
	}

	qmpSocket := fmt.Sprintf("%s/dryrun-%d.qmp", socketDir, os.Getpid())

	smokeArgs := append([]string{}, cfg.Args...)
	smokeArgs = append(smokeArgs, "-S")

	proc, err := process.Start(ctx, process.Config{
		Binary:    cfg.QEMUBinary,
		Args:      smokeArgs,
		QMPSocket: qmpSocket,
	})

	if err != nil {
		lr.OK = false
		lr.Checks = append(lr.Checks, CheckResult{
			Name:    "process-start",
			OK:      false,
			Message: fmt.Sprintf("QEMU process start failed: %v", err),
		})
		return lr, nil
	}
	defer proc.Kill()
	defer os.Remove(qmpSocket)

	if err := proc.WaitReady(ctx, 3*time.Second); err != nil {
		lr.OK = false
		lr.Checks = append(lr.Checks, CheckResult{
			Name:    "qmp-socket-ready",
			OK:      false,
			Message: fmt.Sprintf("QMP socket ready wait failed: %v", err),
		})
		return lr, nil
	}

	client, err := qmp.Dial(ctx, qmpSocket)
	if err != nil {
		lr.OK = false
		lr.Checks = append(lr.Checks, CheckResult{
			Name:    "qmp-dial",
			OK:      false,
			Message: fmt.Sprintf("QMP dial failed: %v", err),
		})
		return lr, nil
	}
	defer client.Close()

	commands := []string{"query-version", "query-migrate-capabilities"}
	for _, cmd := range commands {
		_, err := client.Execute(ctx, cmd, nil)
		ok := (err == nil)
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		lr.Checks = append(lr.Checks, CheckResult{
			Name:    "qmp-cmd-" + cmd,
			OK:      ok,
			Message: msg,
		})
		if !ok {
			lr.OK = false
		}
	}

	return lr, nil
}
