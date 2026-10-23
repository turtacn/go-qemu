package scenarios

import (
	"context"
	"os"
	"path/filepath"
)

type Scenario interface {
	Name() string
	Description() string
	Prerequisites(ctx context.Context) error
	Setup(ctx context.Context) error
	Execute(ctx context.Context) error
	Verify(ctx context.Context) error
	Teardown(ctx context.Context) error
}

func resolveFixturePath(rel string) string {
	if _, err := os.Stat(rel); err == nil {
		return rel
	}
	alt := filepath.Join("..", rel)
	if _, err := os.Stat(alt); err == nil {
		return alt
	}
	alt2 := filepath.Join("..", "..", rel)
	if _, err := os.Stat(alt2); err == nil {
		return alt2
	}
	return rel
}
