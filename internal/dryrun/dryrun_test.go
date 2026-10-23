package dryrun

import (
	"context"
	"testing"

	"github.com/turtacn/go-qemu/internal/vmconfig"
)

func TestLayer1ValidConfigPasses(t *testing.T) {
	cfg := Config{
		Args:   []string{"-machine", "pc-i440fx-5.0", "-m", "1024M", "-smp", "cpus=2", "-chardev", "socket,id=qmp"},
		Layers: Layer1Static,
		VMConfig: &vmconfig.VMConfig{
			Memory: vmconfig.MemoryConfig{SizeMB: 1024},
		},
	}

	res, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if !res.OK {
		t.Errorf("expected layer 1 to pass, got result: %+v", res)
	}
}

func TestLayer1MissingMachineArgFails(t *testing.T) {
	cfg := Config{
		Args:   []string{"-m", "1024M", "-smp", "cpus=2", "-chardev", "socket,id=qmp"}, // missing -machine
		Layers: Layer1Static,
	}

	res, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if res.OK {
		t.Errorf("expected layer 1 to fail when -machine is missing")
	}
}

func TestLayer2FailsWhenBinaryDoesNotExist(t *testing.T) {
	cfg := Config{
		Args:       []string{"-machine", "pc-i440fx-5.0", "-m", "1024M", "-smp", "cpus=2", "-chardev", "socket,id=qmp"},
		QEMUBinary: "/nonexistent/qemu-system-x86_64",
		Layers:     Layer2Smoke,
	}

	res, err := Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if res.OK {
		t.Errorf("expected layer 2 to fail for nonexistent QEMU binary")
	}
}
