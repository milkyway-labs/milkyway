package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
)

var _ restakingtypes.RestakingHooks = RestakingHooks{}

type RestakingHooks struct {
	k *Keeper
}

func (k *Keeper) RestakingHooks() RestakingHooks {
	return RestakingHooks{k}
}

// BeforePoolDelegationCreated implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforePoolDelegationCreated(ctx sdk.Context, poolID uint32, _ string) error {
	return h.k.BeforeDelegationCreated(ctx, restakingtypes.DELEGATION_TYPE_POOL, poolID)
}

// BeforePoolDelegationSharesModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforePoolDelegationSharesModified(ctx sdk.Context, poolID uint32, delegator string) error {
	return h.k.BeforeDelegationSharesModified(ctx, restakingtypes.DELEGATION_TYPE_POOL, poolID, delegator)
}

// AfterPoolDelegationModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) AfterPoolDelegationModified(ctx sdk.Context, poolID uint32, delegator string) error {
	return h.k.AfterDelegationModified(ctx, restakingtypes.DELEGATION_TYPE_POOL, poolID, delegator)
}

// BeforeOperatorDelegationCreated implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforeOperatorDelegationCreated(ctx sdk.Context, operatorID uint32, _ string) error {
	return h.k.BeforeDelegationCreated(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operatorID)
}

// BeforeOperatorDelegationSharesModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforeOperatorDelegationSharesModified(ctx sdk.Context, operatorID uint32, delegator string) error {
	return h.k.BeforeDelegationSharesModified(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operatorID, delegator)
}

// AfterOperatorDelegationModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) AfterOperatorDelegationModified(ctx sdk.Context, operatorID uint32, delegator string) error {
	return h.k.AfterDelegationModified(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operatorID, delegator)
}

// BeforeServiceDelegationCreated implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforeServiceDelegationCreated(ctx sdk.Context, serviceID uint32, _ string) error {
	return h.k.BeforeDelegationCreated(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID)
}

// BeforeServiceDelegationSharesModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforeServiceDelegationSharesModified(ctx sdk.Context, serviceID uint32, delegator string) error {
	return h.k.BeforeDelegationSharesModified(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID, delegator)
}

// AfterServiceDelegationModified implements restakingtypes.RestakingHooks
func (h RestakingHooks) AfterServiceDelegationModified(ctx sdk.Context, serviceID uint32, delegator string) error {
	return h.k.AfterDelegationModified(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID, delegator)
}

// BeforePoolDelegationRemoved implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforePoolDelegationRemoved(_ sdk.Context, _ uint32, _ string) error {
	return nil
}

// BeforeOperatorDelegationRemoved implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforeOperatorDelegationRemoved(_ sdk.Context, _ uint32, _ string) error {
	return nil
}

// BeforeServiceDelegationRemoved implements restakingtypes.RestakingHooks
func (h RestakingHooks) BeforeServiceDelegationRemoved(_ sdk.Context, _ uint32, _ string) error {
	return nil
}

// AfterUnbondingInitiated implements restakingtypes.RestakingHooks
func (h RestakingHooks) AfterUnbondingInitiated(_ sdk.Context, _ uint64) error {
	return nil
}

func (h RestakingHooks) AfterUserTrustedServiceUpdated(ctx sdk.Context, userAddress string, serviceID uint32, trusted bool) error {
	delAddr, err := h.k.accountKeeper.AddressCodec().StringToBytes(userAddress)
	if err != nil {
		return err
	}

	err = h.k.restakingKeeper.IterateUserPoolDelegations(ctx, userAddress, func(del restakingtypes.Delegation) (stop bool, err error) {
		isSecured, err := h.k.restakingKeeper.IsServiceSecuredByPool(ctx, serviceID, del.TargetID)
		if err != nil {
			return true, err
		}
		if !isSecured {
			return false, nil
		}

		pool, found := h.k.poolsKeeper.GetPool(ctx, del.TargetID)
		if !found {
			return true, poolstypes.ErrPoolNotFound
		}

		// Calling these two methods has same effect as calling
		// BeforePoolDelegationSharesModified and then AfterPoolDelegationModified.
		_, err = h.k.withdrawDelegationRewards(ctx, pool, del)
		if err != nil {
			return true, err
		}
		err = h.k.initializeDelegation(ctx, pool, delAddr)
		if err != nil {
			return true, err
		}

		if trusted {
			err = h.k.IncrementPoolServiceTotalDelegatorShares(ctx, del.TargetID, serviceID, del.Shares)
		} else {
			err = h.k.DecrementPoolServiceTotalDelegatorShares(ctx, del.TargetID, serviceID, del.Shares)
		}
		if err != nil {
			return true, err
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}
