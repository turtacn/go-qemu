package feature

import (
	"fmt"
	"strings"
	"sync"

	genV10 "github.com/turtacn/go-qemu/internal/gen/v10"
	genV5 "github.com/turtacn/go-qemu/internal/gen/v5"
	"github.com/turtacn/go-qemu/internal/version"
	"github.com/turtacn/go-qemu/qapi"
)

// FeatureGateRegistry manages feature gates and version compatibility rules.
type FeatureGateRegistry struct {
	mu               sync.RWMutex
	gates            map[string]*FeatureGate
	currentVer       version.QEMUVersion
	currentArch      Arch
	liveSchema       *qapi.QAPISchema
	staticCommandSet map[string]struct{}
}

// New creates a new FeatureGateRegistry initialized with standard feature gates.
func New(ver version.QEMUVersion, arch Arch) *FeatureGateRegistry {
	r := &FeatureGateRegistry{
		gates:       make(map[string]*FeatureGate),
		currentVer:  ver,
		currentArch: arch,
	}

	for _, g := range migrationGates {
		gCopy := g
		r.gates[gCopy.Name] = &gCopy
	}
	for _, g := range deviceGates {
		gCopy := g
		r.gates[gCopy.Name] = &gCopy
	}
	for _, g := range coreGates {
		gCopy := g
		r.gates[gCopy.Name] = &gCopy
	}

	r.staticCommandSet = r.loadStaticCommandSet(ver)
	return r
}

func (r *FeatureGateRegistry) loadStaticCommandSet(ver version.QEMUVersion) map[string]struct{} {
	switch {
	case ver.GTE(version.V(10, 0, 0)):
		return genV10.CommandSet
	case ver.GTE(version.V(5, 0, 0)):
		return genV5.CommandSet
	default:
		return map[string]struct{}{}
	}
}

// InjectLiveSchema attaches a live QAPISchema obtained from a running QEMU instance.
func (r *FeatureGateRegistry) InjectLiveSchema(schema *qapi.QAPISchema) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.liveSchema = schema
}

// IsEnabled queries if a feature is enabled under the current version, arch, and schema.
func (r *FeatureGateRegistry) IsEnabled(name string) (bool, Reason, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	gate, ok := r.gates[name]
	if !ok {
		return false, Reason{}, fmt.Errorf("unknown feature gate: %q", name)
	}

	reason := Reason{}

	// 1. Version Range Check
	reason.VersionOK = gate.Available.Contains(r.currentVer)
	if gate.Removed != nil && r.currentVer.GTE(*gate.Removed) {
		reason.VersionOK = false
		reason.Detail = fmt.Sprintf("removed in %s", gate.Removed)
	}

	// 2. Arch Check
	reason.ArchOK = true
	if len(gate.Arches) > 0 {
		reason.ArchOK = false
		for _, a := range gate.Arches {
			if a == r.currentArch {
				reason.ArchOK = true
				break
			}
		}
	}

	// 3. Schema Check
	reason.SchemaOK = true
	if gate.SchemaKey != "" {
		if r.liveSchema != nil {
			reason.SchemaOK = r.liveSchema.HasKey(gate.SchemaKey)
			if !reason.SchemaOK {
				reason.Detail = fmt.Sprintf("schema key %q not found in live schema", gate.SchemaKey)
			}
		} else {
			if cmdName := extractCommandFromSchemaKey(gate.SchemaKey); cmdName != "" {
				_, reason.SchemaOK = r.staticCommandSet[cmdName]
			}
		}
	}

	return reason.Available(), reason, nil
}

func extractCommandFromSchemaKey(key string) string {
	parts := strings.SplitN(key, ":", 2)
	if len(parts) == 2 && parts[0] == "command" {
		return parts[1]
	}
	return ""
}
