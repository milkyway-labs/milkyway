package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/pools/types"
)

var _ types.PoolsHooks = &Keeper{}

func (k *Keeper) AfterPoolCreated(ctx sdk.Context, poolID uint32) error {
	if k.hooks != nil {
		return k.hooks.AfterPoolCreated(ctx, poolID)
	}
	return nil
}
