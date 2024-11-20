package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

// SetWithdrawAddress sets a new address that will receive the rewards upon withdrawal
func (k *Keeper) SetWithdrawAddress(ctx context.Context, addr, withdrawAddr sdk.AccAddress) error {
	// Check if the withdraw address is blocked
	if k.bankKeeper.BlockedAddr(withdrawAddr) {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive external funds", withdrawAddr)
	}

	// Set the withdraw address
	err := k.DelegatorWithdrawAddrs.Set(ctx, addr, withdrawAddr)
	if err != nil {
		return err
	}
	return nil
}

// WithdrawDelegationRewards withdraws the rewards from the delegation and reinitializes it
func (k *Keeper) WithdrawDelegationRewards(
	ctx context.Context, delAddr sdk.AccAddress, target restakingtypes.DelegationTarget,
) (types.Pools, error) {
	// Get the delegation
	delegator, err := k.accountKeeper.AddressCodec().BytesToString(delAddr)
	if err != nil {
		return nil, err
	}

	delegation, found, err := k.restakingKeeper.GetDelegationForTarget(ctx, target, delegator)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, sdkerrors.ErrNotFound.Wrapf("delegation not found: %d, %s", target.GetID(), delAddr.String())
	}

	// Withdraw the rewards
	rewards, err := k.withdrawDelegationRewards(ctx, target, delegation)
	if err != nil {
		return nil, err
	}

	// Reinitialize the delegation
	err = k.initializeDelegation(ctx, target, delAddr)
	if err != nil {
		return nil, err
	}

	return rewards, nil
}

// WithdrawOperatorCommission withdraws the operator's accumulated commission
func (k *Keeper) WithdrawOperatorCommission(ctx context.Context, operatorID uint32) (types.Pools, error) {
	operator, found, err := k.operatorsKeeper.GetOperator(ctx, operatorID)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, operatorstypes.ErrOperatorNotFound
	}

	// Fetch the operator accumulated commission
	accumCommission, err := k.OperatorAccumulatedCommissions.Get(ctx, operatorID)
	if err != nil {
		return nil, err
	}
	if accumCommission.Commissions.IsEmpty() {
		return nil, types.ErrNoOperatorCommission
	}

	commissions, remainder := accumCommission.Commissions.TruncateDecimal()

	// Leave the remainder to withdraw later
	err = k.OperatorAccumulatedCommissions.Set(ctx, operatorID, types.AccumulatedCommission{
		Commissions: remainder,
	})
	if err != nil {
		return nil, err
	}

	// Update the outstanding rewards
	outstanding, err := k.OperatorOutstandingRewards.Get(ctx, operatorID)
	if err != nil {
		return nil, err
	}

	err = k.OperatorOutstandingRewards.Set(ctx, operatorID, types.OutstandingRewards{
		Rewards: outstanding.Rewards.Sub(types.NewDecPoolsFromPools(commissions)),
	})
	if err != nil {
		return nil, err
	}

	// Send the commission to the operator
	commissionCoins := commissions.Sum()
	if !commissionCoins.IsZero() {
		adminAddr, err := k.accountKeeper.AddressCodec().StringToBytes(operator.Admin)
		if err != nil {
			return nil, err
		}
		withdrawAddr, err := k.GetDelegatorWithdrawAddr(ctx, adminAddr)
		if err != nil {
			return nil, err
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.RewardsPoolName, withdrawAddr, commissionCoins)
		if err != nil {
			return nil, err
		}
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawCommission,
			sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, fmt.Sprint(operatorID)),
			sdk.NewAttribute(sdk.AttributeKeyAmount, commissions.Sum().String()),
			sdk.NewAttribute(types.AttributeKeyAmountPerPool, commissions.String()),
		),
	)

	return commissions, nil
}
