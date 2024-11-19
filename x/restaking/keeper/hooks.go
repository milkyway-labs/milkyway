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

// --------------------------------------------------------------------------------------------------------------------

// BeforeOperatorDelegationCreated implements types.RestakingHooks
func (k *Keeper) BeforeOperatorDelegationCreated(ctx sdk.Context, operatorID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.BeforeOperatorDelegationCreated(ctx, operatorID, delegator)
	}
	return nil
}

// BeforeOperatorDelegationSharesModified implements types.RestakingHooks
func (k *Keeper) BeforeOperatorDelegationSharesModified(ctx sdk.Context, operatorID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.BeforeOperatorDelegationSharesModified(ctx, operatorID, delegator)
	}
	return nil
}

// AfterOperatorDelegationModified implements types.RestakingHooks
func (k *Keeper) AfterOperatorDelegationModified(ctx sdk.Context, operatorID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.AfterOperatorDelegationModified(ctx, operatorID, delegator)
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// BeforeServiceDelegationCreated implements types.RestakingHooks
func (k *Keeper) BeforeServiceDelegationCreated(ctx sdk.Context, serviceID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.BeforeServiceDelegationCreated(ctx, serviceID, delegator)
	}
	return nil
}

// BeforeServiceDelegationSharesModified implements types.RestakingHooks
func (k *Keeper) BeforeServiceDelegationSharesModified(ctx sdk.Context, serviceID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.BeforeServiceDelegationSharesModified(ctx, serviceID, delegator)
	}
	return nil
}

// AfterServiceDelegationModified implements types.RestakingHooks
func (k *Keeper) AfterServiceDelegationModified(ctx sdk.Context, serviceID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.AfterServiceDelegationModified(ctx, serviceID, delegator)
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// BeforePoolDelegationRemoved implements types.RestakingHooks
func (k *Keeper) BeforePoolDelegationRemoved(ctx sdk.Context, poolID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.BeforePoolDelegationRemoved(ctx, poolID, delegator)
	}
	return nil
}

// BeforeOperatorDelegationRemoved implements types.RestakingHooks
func (k *Keeper) BeforeOperatorDelegationRemoved(ctx sdk.Context, operatorID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.BeforeOperatorDelegationRemoved(ctx, operatorID, delegator)
	}
	return nil
}

// BeforeServiceDelegationRemoved implements types.RestakingHooks
func (k *Keeper) BeforeServiceDelegationRemoved(ctx sdk.Context, serviceID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.BeforeServiceDelegationRemoved(ctx, serviceID, delegator)
	}
	return nil
}

// AfterUnbondingInitiated implements types.RestakingHooks
func (k *Keeper) AfterUnbondingInitiated(ctx sdk.Context, unbondingDelegationID uint64) error {
	if k.hooks != nil {
		return k.hooks.AfterUnbondingInitiated(ctx, unbondingDelegationID)
	}
	return nil
}

// AfterUserPreferencesModified implements types.RestakingHooks
func (k *Keeper) AfterUserPreferencesModified(ctx sdk.Context, userAddress string, oldPreferences, newPreferences types.UserPreferences) error {
	if k.hooks != nil {
		return k.hooks.AfterUserPreferencesModified(ctx, userAddress, oldPreferences, newPreferences)
	}
	return nil
}
