package scenarios

import (
	"context"
	"fmt"

	"github.com/turtacn/go-qemu/internal/feature"
	"github.com/turtacn/go-qemu/internal/version"
)

type FeatureMatrixScenario struct{}

func (s *FeatureMatrixScenario) Name() string { return "feature-gates-verification" }
func (s *FeatureMatrixScenario) Description() string {
	return "Feature Gates verification across version matrix"
}
func (s *FeatureMatrixScenario) Prerequisites(ctx context.Context) error { return nil }
func (s *FeatureMatrixScenario) Setup(ctx context.Context) error         { return nil }
func (s *FeatureMatrixScenario) Execute(ctx context.Context) error       { return nil }
func (s *FeatureMatrixScenario) Verify(ctx context.Context) error {
	reg5 := feature.New(version.V(5, 0, 0), feature.ArchX86)
	enabled5, _, _ := reg5.IsEnabled("migrate.postcopy-preempt")
	if enabled5 {
		return fmt.Errorf("postcopy-preempt should be disabled in 5.0")
	}

	reg10 := feature.New(version.V(10, 2, 2), feature.ArchX86)
	enabled10, _, _ := reg10.IsEnabled("migrate.postcopy-preempt")
	if !enabled10 {
		return fmt.Errorf("postcopy-preempt should be enabled in 10.2.2")
	}
	return nil
}
func (s *FeatureMatrixScenario) Teardown(ctx context.Context) error { return nil }
