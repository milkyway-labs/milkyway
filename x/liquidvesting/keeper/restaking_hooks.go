package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v9/x/liquidvesting/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
)

var _ restakingtypes.RestakingHooks = RestakingHooks{}

type RestakingHooks struct {
	*Keeper
}

func (k *Keeper) RestakingHooks() RestakingHooks {
	return RestakingHooks{k}
}

func (h RestakingHooks) BeforeDelegationSharesModified(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32, delegator string) error {
	delegation, found, err := h.restakingKeeper.GetDelegation(ctx, delType, targetID, delegator)
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

	insuranceFund, err := h.GetUserInsuranceFundBalance(ctx, delegation.UserAddress)
	if err != nil {
		return err
	}
	// Exit early if the user doesn't have insurance fund balance
	if insuranceFund.IsZero() {
		return nil
	}

	// Temporarily save the delegation so that it can be used in the
	// AfterDelegationModified hook.
	target, err := h.restakingKeeper.GetDelegationTarget(ctx, delegation.Type, delegation.TargetID)
	if err != nil {
		return err
	}
	err = h.SetPreviousDelegationTokens(ctx, delegator, delType, targetID, target.TokensFromShares(delegation.Shares))
	if err != nil {
		return err
	}

	// Calculate covered locked shares
	params, err := h.GetParams(ctx)
	if err != nil {
		return err
	}
	activeLockedTokens, err := h.GetAllUserActiveLockedRepresentations(ctx, delegation.UserAddress)
	if err != nil {
		return err
	}
	coveredLockedShares, err := types.GetCoveredLockedShares(
		target,
		delegation,
		insuranceFund,
		params.InsurancePercentage,
		activeLockedTokens,
	)
	if err != nil {
		return err
	}
	if coveredLockedShares.IsZero() {
		return nil
	}

	return h.DecrementTargetCoveredLockedShares(ctx, delType, targetID, coveredLockedShares)
}

func (h RestakingHooks) AfterDelegationModified(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32, delegator string) error {
	delegation, found, err := h.restakingKeeper.GetDelegation(ctx, delType, targetID, delegator)
	if err != nil {
		return err
	}
	if !found {
		return restakingtypes.ErrDelegationNotFound
	}

	// Depending on whether the delegation has locked shares or not, either remove
	// the delegator from the locked representation delegators list or mark them as
	// a locked representation delegator.
	if !types.HasLockedShares(delegation.Shares) {
		return h.RemoveLockedRepresentationDelegator(ctx, delegator)
	}
	err = h.SetLockedRepresentationDelegator(ctx, delegator)
	if err != nil {
		return err
	}

	insuranceFund, err := h.GetUserInsuranceFundBalance(ctx, delegation.UserAddress)
	if err != nil {
		return err
	}
	// Exit early if the user doesn't have insurance fund balance
	if insuranceFund.IsZero() {
		return nil
	}

	// Get the cached delegation from the BeforeDelegationSharesModified hook.
	prevDelegationTokens, err := h.GetPreviousDelegationTokens(ctx, delegator, delType, targetID)
	if err != nil {
		return err
	}
	err = h.RemovePreviousDelegationTokens(ctx, delegator, delType, targetID)
	if err != nil {
		return err
	}

	// Calculate the previous active locked tokens before modifying the delegation.
	target, err := h.restakingKeeper.GetDelegationTarget(ctx, delegation.Type, delegation.TargetID)
	if err != nil {
		return err
	}
	newDelegationTokens := target.TokensFromShares(delegation.Shares)
	newActiveLockedTokens, err := h.GetAllUserActiveLockedRepresentations(ctx, delegator)
	if err != nil {
		return err
	}
	prevActiveLockedTokens := newActiveLockedTokens.Sub(newDelegationTokens).Add(prevDelegationTokens...)

	params, err := h.GetParams(ctx)
	if err != nil {
		return err
	}

	// Withdraw restaking rewards from the user's all other delegations except for
	// the current one since active locked tokens amount has changed and covered
	// locked shares for the other delegations need to be updated.
	err = h.WithdrawUserLockedRestakingRewards(
		ctx,
		delegator,
		func(del restakingtypes.Delegation) bool {
			return del.Type != delType || del.TargetID != targetID
		},
		func() (sdk.Coins, sdkmath.LegacyDec, sdk.DecCoins) {
			return insuranceFund, params.InsurancePercentage, prevActiveLockedTokens
		},
		func() (sdk.Coins, sdkmath.LegacyDec, sdk.DecCoins) {
			return insuranceFund, params.InsurancePercentage, newActiveLockedTokens
		},
	)
	if err != nil {
		return err
	}

	coveredLockedShares, err := types.GetCoveredLockedShares(
		target,
		delegation,
		insuranceFund,
		params.InsurancePercentage,
		newActiveLockedTokens,
	)
	if err != nil {
		return err
	}
	if coveredLockedShares.IsZero() {
		return nil
	}

	return h.IncrementTargetCoveredLockedShares(ctx, delType, targetID, coveredLockedShares)
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

// Methods below are no-op, but are required to satisfy the interface.

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
