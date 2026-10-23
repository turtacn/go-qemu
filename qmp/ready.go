package qmp

import (
	"context"
	"os"
	"time"
)

// WaitForSocket waits for a Unix domain socket file to exist and be accessible.
// Since QEMU startup has a brief delay before creating the QMP socket, this encapsulates retry logic.
func WaitForSocket(ctx context.Context, socketPath string) error {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	// Initial immediate check
	if _, err := os.Stat(socketPath); err == nil {
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if _, err := os.Stat(socketPath); err == nil {
				return nil
			}
		}
	}
}
