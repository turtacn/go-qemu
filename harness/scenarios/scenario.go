package scenarios

import "context"

type Scenario interface {
	Name() string
	Description() string
	Prerequisites(ctx context.Context) error
	Setup(ctx context.Context) error
	Execute(ctx context.Context) error
	Verify(ctx context.Context) error
	Teardown(ctx context.Context) error
}
