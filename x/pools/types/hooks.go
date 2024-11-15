package types

import (
	"context"
)

type PoolsHooks interface {
	AfterPoolCreated(ctx context.Context, poolID uint32) error
}
