package process

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProcessStartWithNonexistentBinary(t *testing.T) {
	_, err := Start(context.Background(), Config{
		Binary: "/nonexistent/qemu-binary",
	})
	if err == nil {
		t.Errorf("expected error when starting nonexistent binary")
	}
}

func TestProcessMock(t *testing.T) {
	tmpDir := t.TempDir()
	mockSock := filepath.Join(tmpDir, "mock.sock")

	// Create a dummy file to simulate socket ready
	if err := os.WriteFile(mockSock, []byte{}, 0600); err != nil {
		t.Fatalf("failed to create dummy socket file: %v", err)
	}

	p := &Process{QMPSocket: mockSock}
	if err := p.WaitReady(context.Background(), 1*time.Second); err != nil {
		t.Errorf("WaitReady should succeed when file exists, got %v", err)
	}
}
