package qmp

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWaitForSocket(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "qmp-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	sockPath := filepath.Join(tmpDir, "test.sock")

	// Test timeout when socket doesn't exist
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancelTimeout()
	if err := WaitForSocket(ctxTimeout, sockPath); err == nil {
		t.Errorf("expected timeout error when socket file does not exist")
	}

	// Create dummy file
	if err := os.WriteFile(sockPath, []byte{}, 0600); err != nil {
		t.Fatalf("failed to create dummy socket file: %v", err)
	}

	ctxReady, cancelReady := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelReady()
	if err := WaitForSocket(ctxReady, sockPath); err != nil {
		t.Errorf("expected WaitForSocket to succeed, got %v", err)
	}
}

func TestErrors(t *testing.T) {
	cmdErr := &CommandError{
		Class:   ErrClassCommandNotFound,
		Desc:    "The command has not been found",
		Command: "invalid_cmd",
	}
	expectedCmdErrStr := `QMP command "invalid_cmd" failed [CommandNotFound]: The command has not been found`
	if cmdErr.Error() != expectedCmdErrStr {
		t.Errorf("got %q, expected %q", cmdErr.Error(), expectedCmdErrStr)
	}

	baseErr := errors.New("connection reset by peer")
	transErr := &TransportError{
		Op:  "Execute",
		Err: baseErr,
	}
	expectedTransErrStr := "QMP transport error during Execute: connection reset by peer"
	if transErr.Error() != expectedTransErrStr {
		t.Errorf("got %q, expected %q", transErr.Error(), expectedTransErrStr)
	}
	if !errors.Is(transErr, baseErr) {
		t.Errorf("expected errors.Is to match unwrapped error")
	}
}
