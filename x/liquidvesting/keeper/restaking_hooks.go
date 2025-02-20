package keeper

import (
	"context"

	"github.com/milkyway-labs/milkyway/v9/x/liquidvesting/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
)

var _ restakingtypes.RestakingHooks = RestakingHooks{}

type RestakingHooks struct {
	k *Keeper
}

func (k *Keeper) RestakingHooks() RestakingHooks {
	return RestakingHooks{k}
}

func (h RestakingHooks) BeforeDelegationSharesModified(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32, delegator string) error {
	delegation, found, err := h.k.restakingKeeper.GetDelegation(ctx, delType, targetID, delegator)
	if err != nil {
		return err
	}
	if !found {
		return restakingtypes.ErrDelegationNotFound
	}

	// If the delegation has no locked shares, exit early.
	if !types.HasLockedShares(delegation.Shares) {
		return nil
	}

	coveredLockedShares, err := h.k.GetCoveredLockedShares(ctx, delegation)
	if err != nil {
		return err
	}
	if coveredLockedShares.IsZero() {
		return nil
	}

	return h.k.DecrementTargetCoveredLockedShares(ctx, delType, targetID, coveredLockedShares)
}

func (h RestakingHooks) AfterDelegationModified(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32, delegator string) error {
	delegation, found, err := h.k.restakingKeeper.GetDelegation(ctx, delType, targetID, delegator)
	if err != nil {
		return err
	}
	if !found {
		return restakingtypes.ErrDelegationNotFound
	}

	// If the delegation has no locked shares, remove the delegator from the list and
	// exit early.
	if !types.HasLockedShares(delegation.Shares) {
		return h.k.RemoveLockedRepresentationDelegator(ctx, delegation.UserAddress)
	}

	// If the delegation has locked representation inside, mark the delegator as
	// locked representation delegator.
	err = h.k.SetLockedRepresentationDelegator(ctx, delegation.UserAddress)
	if err != nil {
		return err
	}

	coveredLockedShares, err := h.k.GetCoveredLockedShares(ctx, delegation)
	if err != nil {
		return err
	}
	if coveredLockedShares.IsZero() {
		return nil
	}

	return h.k.IncrementTargetCoveredLockedShares(ctx, delType, targetID, coveredLockedShares)
}

// BeforePoolDelegationSharesModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforePoolDelegationSharesModified(ctx context.Context, poolID uint32, delegator string) error {
	return h.BeforeDelegationSharesModified(ctx, restakingtypes.DELEGATION_TYPE_POOL, poolID, delegator)
}

// AfterPoolDelegationModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) AfterPoolDelegationModified(ctx context.Context, poolID uint32, delegator string) error {
	return h.AfterDelegationModified(ctx, restakingtypes.DELEGATION_TYPE_POOL, poolID, delegator)
}

// BeforeOperatorDelegationSharesModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforeOperatorDelegationSharesModified(ctx context.Context, operatorID uint32, delegator string) error {
	return h.BeforeDelegationSharesModified(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operatorID, delegator)
}

// AfterOperatorDelegationModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) AfterOperatorDelegationModified(ctx context.Context, operatorID uint32, delegator string) error {
	return h.AfterDelegationModified(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operatorID, delegator)
}

// BeforeServiceDelegationSharesModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforeServiceDelegationSharesModified(ctx context.Context, serviceID uint32, delegator string) error {
	return h.BeforeDelegationSharesModified(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID, delegator)
}

// AfterServiceDelegationModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) AfterServiceDelegationModified(ctx context.Context, serviceID uint32, delegator string) error {
	return h.AfterDelegationModified(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID, delegator)
}

// BeforePoolDelegationCreated implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforePoolDelegationCreated(context.Context, uint32, string) error {
	return nil
}

// BeforeOperatorDelegationCreated implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforeOperatorDelegationCreated(context.Context, uint32, string) error {
	return nil
}

// BeforeServiceDelegationCreated implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforeServiceDelegationCreated(context.Context, uint32, string) error {
	return nil
}

// BeforePoolDelegationRemoved implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforePoolDelegationRemoved(context.Context, uint32, string) error {
	return nil
}

// BeforeOperatorDelegationRemoved implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforeOperatorDelegationRemoved(context.Context, uint32, string) error {
	return nil
}

// BeforeServiceDelegationRemoved implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforeServiceDelegationRemoved(context.Context, uint32, string) error {
	return nil
}

// AfterUnbondingInitiated implements restakingtypes.RestakingHooks
func (h RestakingHooks) AfterUnbondingInitiated(context.Context, uint64) error {
	return nil
}

// AfterUserPreferencesModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) AfterUserPreferencesModified(context.Context, string, restakingtypes.UserPreferences, restakingtypes.UserPreferences) error {
	return nil
}
