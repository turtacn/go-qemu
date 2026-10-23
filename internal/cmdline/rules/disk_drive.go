package rules

import (
	"fmt"

	"github.com/turtacn/go-qemu/internal/cmdline"
	"github.com/turtacn/go-qemu/internal/version"
	"github.com/turtacn/go-qemu/internal/vmconfig"
)

func buildDriveArgs(disk vmconfig.DiskConfig, ver version.QEMUVersion) ([]cmdline.ArgEntry, error) {
	driveID := fmt.Sprintf("drive-%s%d", disk.Bus, disk.Index)
	fmtFormat := disk.Format
	if fmtFormat == "" {
		fmtFormat = "raw"
	}

	deprecated := ver.GTE(version.V(7, 0, 0))

	return []cmdline.ArgEntry{
		{
			Key: "-drive",
			Value: fmt.Sprintf(
				"file=%s,format=%s,if=none,id=%s,cache=none,aio=native",
				disk.Path, fmtFormat, driveID,
			),
			RuleName:    "disk.drive",
			FeatureGate: "",
			Reason:      "QEMU 5.x standard disk parameters",
			Deprecated:  deprecated,
		},
	}, nil
}
