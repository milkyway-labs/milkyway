package keeper

import (
	"context"

	"github.com/milkyway-labs/milkyway/v7/x/restaking/types"
)

var _ types.RestakingHooks = &Keeper{}

// BeforePoolDelegationCreated implements types.RestakingHooks
func (k *Keeper) BeforePoolDelegationCreated(ctx context.Context, poolID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.BeforePoolDelegationCreated(ctx, poolID, delegator)
	}
	return nil
}

// BeforePoolDelegationSharesModified implements types.RestakingHooks
func (k *Keeper) BeforePoolDelegationSharesModified(ctx context.Context, poolID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.BeforePoolDelegationSharesModified(ctx, poolID, delegator)
	}
	return nil
}

// AfterPoolDelegationModified implements types.RestakingHooks
func (k *Keeper) AfterPoolDelegationModified(ctx context.Context, poolID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.AfterPoolDelegationModified(ctx, poolID, delegator)
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// BeforeOperatorDelegationCreated implements types.RestakingHooks
func (k *Keeper) BeforeOperatorDelegationCreated(ctx context.Context, operatorID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.BeforeOperatorDelegationCreated(ctx, operatorID, delegator)
	}
	return nil
}

// BeforeOperatorDelegationSharesModified implements types.RestakingHooks
func (k *Keeper) BeforeOperatorDelegationSharesModified(ctx context.Context, operatorID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.BeforeOperatorDelegationSharesModified(ctx, operatorID, delegator)
	}
	return nil
}

// AfterOperatorDelegationModified implements types.RestakingHooks
func (k *Keeper) AfterOperatorDelegationModified(ctx context.Context, operatorID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.AfterOperatorDelegationModified(ctx, operatorID, delegator)
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// BeforeServiceDelegationCreated implements types.RestakingHooks
func (k *Keeper) BeforeServiceDelegationCreated(ctx context.Context, serviceID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.BeforeServiceDelegationCreated(ctx, serviceID, delegator)
	}
	return nil
}

// BeforeServiceDelegationSharesModified implements types.RestakingHooks
func (k *Keeper) BeforeServiceDelegationSharesModified(ctx context.Context, serviceID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.BeforeServiceDelegationSharesModified(ctx, serviceID, delegator)
	}
	return nil
}

// AfterServiceDelegationModified implements types.RestakingHooks
func (k *Keeper) AfterServiceDelegationModified(ctx context.Context, serviceID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.AfterServiceDelegationModified(ctx, serviceID, delegator)
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// BeforePoolDelegationRemoved implements types.RestakingHooks
func (k *Keeper) BeforePoolDelegationRemoved(ctx context.Context, poolID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.BeforePoolDelegationRemoved(ctx, poolID, delegator)
	}
	return nil
}

// BeforeOperatorDelegationRemoved implements types.RestakingHooks
func (k *Keeper) BeforeOperatorDelegationRemoved(ctx context.Context, operatorID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.BeforeOperatorDelegationRemoved(ctx, operatorID, delegator)
	}
	return nil
}

// BeforeServiceDelegationRemoved implements types.RestakingHooks
func (k *Keeper) BeforeServiceDelegationRemoved(ctx context.Context, serviceID uint32, delegator string) error {
	if k.hooks != nil {
		return k.hooks.BeforeServiceDelegationRemoved(ctx, serviceID, delegator)
	}
	return nil
}

// AfterUnbondingInitiated implements types.RestakingHooks
func (k *Keeper) AfterUnbondingInitiated(ctx context.Context, unbondingDelegationID uint64) error {
	if k.hooks != nil {
		return k.hooks.AfterUnbondingInitiated(ctx, unbondingDelegationID)
	}
	return nil
}

// AfterUserPreferencesModified implements types.RestakingHooks
func (k *Keeper) AfterUserPreferencesModified(ctx context.Context, userAddress string, oldPreferences, newPreferences types.UserPreferences) error {
	if k.hooks != nil {
		return k.hooks.AfterUserPreferencesModified(ctx, userAddress, oldPreferences, newPreferences)
	}
	return nil
}
