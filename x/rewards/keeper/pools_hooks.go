package keeper

import (
	"context"

	poolstypes "github.com/milkyway-labs/milkyway/v3/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v3/x/restaking/types"
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
func (h PoolsHooks) AfterPoolCreated(ctx context.Context, poolID uint32) error {
	return h.k.AfterDelegationTargetCreated(ctx, restakingtypes.DELEGATION_TYPE_POOL, poolID)
}
