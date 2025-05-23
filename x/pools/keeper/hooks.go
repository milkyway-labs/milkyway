package keeper

import (
	"context"

	"github.com/milkyway-labs/milkyway/v12/x/pools/types"
)

var _ types.PoolsHooks = &Keeper{}

func (k *Keeper) AfterPoolCreated(ctx context.Context, poolID uint32) error {
	if k.hooks != nil {
		return k.hooks.AfterPoolCreated(ctx, poolID)
	}
	return nil
}
