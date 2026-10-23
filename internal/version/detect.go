package version

import (
	"context"
	"fmt"

	"github.com/turtacn/go-qemu/qmp"
)

// DetectVersion queries the QEMU version from a live QMP client.
func DetectVersion(ctx context.Context, client qmp.Client) (QEMUVersion, error) {
	v, err := client.QueryVersion(ctx)
	if err != nil {
		return QEMUVersion{}, fmt.Errorf("failed to query qemu version: %w", err)
	}
	return V(v.QEMU.Major, v.QEMU.Minor, v.QEMU.Micro), nil
}
