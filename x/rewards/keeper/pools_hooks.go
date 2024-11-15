package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
)

var (
	_ poolstypes.PoolsHooks = PoolsHooks{}
)

type PoolsHooks struct {
	k *Keeper
}

func (k *Keeper) PoolsHooks() PoolsHooks {
	return PoolsHooks{k}
}

// AfterPoolCreated implements poolstypes.PoolsHooks
func (h PoolsHooks) AfterPoolCreated(ctx sdk.Context, poolID uint32) error {
	return h.k.AfterDelegationTargetCreated(ctx, restakingtypes.DELEGATION_TYPE_POOL, poolID)
}
