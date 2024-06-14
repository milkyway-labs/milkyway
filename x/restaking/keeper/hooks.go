package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

var _ types.RestakingHooks = &Keeper{}

// BeforePoolDelegationCreated implements types.RestakingHooks
func (k *Keeper) BeforePoolDelegationCreated(ctx sdk.Context, poolID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.BeforePoolDelegationCreated(ctx, poolID, delegator)
	}
	return nil
}

// BeforePoolDelegationSharesModified implements types.RestakingHooks
func (k *Keeper) BeforePoolDelegationSharesModified(ctx sdk.Context, poolID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.BeforePoolDelegationSharesModified(ctx, poolID, delegator)
	}
	return nil
}

// AfterPoolDelegationModified implements types.RestakingHooks
func (k *Keeper) AfterPoolDelegationModified(ctx sdk.Context, poolID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.AfterPoolDelegationModified(ctx, poolID, delegator)
	}
	return nil
}
