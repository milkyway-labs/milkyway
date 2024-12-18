package keeper

import (
	"context"

	restakingtypes "github.com/milkyway-labs/milkyway/v6/x/restaking/types"
)

// AfterUserTrustedServiceUpdated is called after a user's trust in a service is
// updated. AfterUserTrustedServiceUpdated updates the total delegator shares for
// the service in all pools where the user has a delegation.
func (k *Keeper) AfterUserTrustedServiceUpdated(
	ctx context.Context,
	userAddress string,
	serviceID uint32,
	trusted bool,
) error {
	delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(userAddress)
	if err != nil {
		return err
	}

	return k.restakingKeeper.IterateUserPoolDelegations(ctx, userAddress, func(del restakingtypes.Delegation) (stop bool, err error) {
		isSecured, err := k.restakingKeeper.IsServiceSecuredByPool(ctx, serviceID, del.TargetID)
		if err != nil {
			return true, err
		}

		if !isSecured {
			return false, nil
		}

		pool, err := k.GetDelegationTarget(ctx, restakingtypes.DELEGATION_TYPE_POOL, del.TargetID)
		if err != nil {
			return true, err
		}

		// Calling these two methods has same effect as calling
		// BeforePoolDelegationSharesModified and then AfterPoolDelegationModified.
		_, err = k.withdrawDelegationRewards(ctx, pool, del)
		if err != nil {
			return true, err
		}

		err = k.initializeDelegation(ctx, pool, delAddr)
		if err != nil {
			return true, err
		}

		if trusted {
			err = k.IncrementPoolServiceTotalDelegatorShares(ctx, del.TargetID, serviceID, del.Shares)
		} else {
			err = k.DecrementPoolServiceTotalDelegatorShares(ctx, del.TargetID, serviceID, del.Shares)
		}
		if err != nil {
			return true, err
		}

		return false, nil
	})
}
