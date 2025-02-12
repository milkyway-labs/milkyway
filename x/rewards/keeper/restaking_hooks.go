package keeper

import (
	"context"

	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
)

var _ restakingtypes.RestakingHooks = RestakingHooks{}

type RestakingHooks struct {
	k *Keeper
}

func (k *Keeper) RestakingHooks() RestakingHooks {
	return RestakingHooks{k}
}

// BeforePoolDelegationCreated implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforePoolDelegationCreated(ctx context.Context, poolID uint32, _ string) error {
	return h.k.BeforeDelegationCreated(ctx, restakingtypes.DELEGATION_TYPE_POOL, poolID)
}

// BeforePoolDelegationSharesModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforePoolDelegationSharesModified(ctx context.Context, poolID uint32, delegator string) error {
	return h.k.BeforeDelegationSharesModified(ctx, restakingtypes.DELEGATION_TYPE_POOL, poolID, delegator)
}

// AfterPoolDelegationModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) AfterPoolDelegationModified(ctx context.Context, poolID uint32, delegator string) error {
	return h.k.AfterDelegationModified(ctx, restakingtypes.DELEGATION_TYPE_POOL, poolID, delegator)
}

// BeforeOperatorDelegationCreated implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforeOperatorDelegationCreated(ctx context.Context, operatorID uint32, _ string) error {
	return h.k.BeforeDelegationCreated(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operatorID)
}

// BeforeOperatorDelegationSharesModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforeOperatorDelegationSharesModified(ctx context.Context, operatorID uint32, delegator string) error {
	return h.k.BeforeDelegationSharesModified(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operatorID, delegator)
}

// AfterOperatorDelegationModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) AfterOperatorDelegationModified(ctx context.Context, operatorID uint32, delegator string) error {
	return h.k.AfterDelegationModified(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operatorID, delegator)
}

// BeforeServiceDelegationCreated implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforeServiceDelegationCreated(ctx context.Context, serviceID uint32, _ string) error {
	return h.k.BeforeDelegationCreated(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID)
}

// BeforeServiceDelegationSharesModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforeServiceDelegationSharesModified(ctx context.Context, serviceID uint32, delegator string) error {
	return h.k.BeforeDelegationSharesModified(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID, delegator)
}

// AfterServiceDelegationModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) AfterServiceDelegationModified(ctx context.Context, serviceID uint32, delegator string) error {
	return h.k.AfterDelegationModified(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID, delegator)
}

// BeforePoolDelegationRemoved implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforePoolDelegationRemoved(_ context.Context, _ uint32, _ string) error {
	return nil
}

// BeforeOperatorDelegationRemoved implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforeOperatorDelegationRemoved(_ context.Context, _ uint32, _ string) error {
	return nil
}

// BeforeServiceDelegationRemoved implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforeServiceDelegationRemoved(_ context.Context, _ uint32, _ string) error {
	return nil
}

// AfterUnbondingInitiated implements restakingtypes.RestakingHooks
func (h RestakingHooks) AfterUnbondingInitiated(_ context.Context, _ uint64) error {
	return nil
}

// AfterUserPreferencesModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) AfterUserPreferencesModified(ctx context.Context, userAddress string, oldPreferences, newPreferences restakingtypes.UserPreferences) error {
	return h.k.AfterUserPreferencesModified(ctx, userAddress, oldPreferences, newPreferences)
}
