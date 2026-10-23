package migrate

import (
	"context"
	"fmt"

	"github.com/turtacn/go-qemu/qmp"
)

type NegotiatedCapabilitySet struct {
	SrcVersion   string
	DstVersion   string
	SrcCaps      []qmp.MigrateCapabilityStatus
	DstCaps      []qmp.MigrateCapabilityStatus
	Intersection []qmp.MigrateCapabilityStatus
	Disabled     []string
}

type MigrationNegotiator struct {
	srcClient qmp.Client
	dstClient qmp.Client
}

func NewNegotiator(src, dst qmp.Client) *MigrationNegotiator {
	return &MigrationNegotiator{srcClient: src, dstClient: dst}
}

func (n *MigrationNegotiator) Negotiate(ctx context.Context) (*NegotiatedCapabilitySet, error) {
	srcCaps, err := n.srcClient.QueryMigrateCapabilities(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query source migration capabilities: %w", err)
	}

	dstCaps, err := n.dstClient.QueryMigrateCapabilities(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query destination migration capabilities: %w", err)
	}

	srcCapNames := make(map[string]bool)
	for _, c := range srcCaps {
		srcCapNames[c.Capability] = c.State
	}

	result := &NegotiatedCapabilitySet{
		SrcCaps: srcCaps,
		DstCaps: dstCaps,
	}

	finalDstCaps := make([]qmp.MigrateCapabilityStatus, 0, len(dstCaps))

	for _, dstCap := range dstCaps {
		if srcState, exists := srcCapNames[dstCap.Capability]; exists {
			finalCap := qmp.MigrateCapabilityStatus{
				Capability: dstCap.Capability,
				State:      srcState,
			}
			finalDstCaps = append(finalDstCaps, finalCap)
			if srcState {
				result.Intersection = append(result.Intersection, finalCap)
			}
		} else {
			finalDstCaps = append(finalDstCaps, qmp.MigrateCapabilityStatus{
				Capability: dstCap.Capability,
				State:      false,
			})
			result.Disabled = append(result.Disabled, dstCap.Capability)
		}
	}

	if err := n.dstClient.SetMigrateCapabilities(ctx, finalDstCaps); err != nil {
		return nil, fmt.Errorf("failed to set destination migration capabilities: %w", err)
	}

	return result, nil
}
