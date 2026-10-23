package rules

import (
	"fmt"

	"github.com/turtacn/go-qemu/internal/cmdline"
	"github.com/turtacn/go-qemu/internal/version"
	"github.com/turtacn/go-qemu/internal/vmconfig"
)

func buildBlockdevArgs(disk vmconfig.DiskConfig, ver version.QEMUVersion) ([]cmdline.ArgEntry, error) {
	nodeID := fmt.Sprintf("drive-%s%d", disk.Bus, disk.Index)
	fmtFormat := disk.Format
	if fmtFormat == "" {
		fmtFormat = "raw"
	}

	blockdevArg := fmt.Sprintf(
		"driver=%s,node-name=%s,file.driver=host_device,file.filename=%s,file.aio=native,cache.direct=on",
		fmtFormat, nodeID, disk.Path,
	)

	return []cmdline.ArgEntry{
		{
			Key:         "-blockdev",
			Value:       blockdevArg,
			RuleName:    "disk.blockdev",
			FeatureGate: "disk.blockdev-preferred",
			Reason:      "QEMU 7.0+ recommended -blockdev instead of -drive",
			Deprecated:  false,
		},
		{
			Key:      "-device",
			Value:    fmt.Sprintf("scsi-hd,drive=%s,id=%s%d", nodeID, disk.Bus, disk.Index),
			RuleName: "disk.blockdev",
		},
	}, nil
}
