package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v9/x/liquidvesting/types"
	operatorstypes "github.com/milkyway-labs/milkyway/v9/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v9/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v9/x/services/types"
)

// AddToUserInsuranceFund adds the provided amount to the user's insurance fund.
// NOTE: We assume that the amount that will be added to the user's insurance fund
// is already present in the module account balance.
func (k *Keeper) AddToUserInsuranceFund(ctx context.Context, user string, amount sdk.Coins) error {
	insuranceFund, err := k.insuranceFunds.Get(ctx, user)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			insuranceFund = types.NewEmptyInsuranceFund()
		} else {
			return err
		}
	}

	newBalance := insuranceFund.Balance.Add(amount...)

	err = k.beforeUserInsuranceFundModified(ctx, user, insuranceFund.Balance, newBalance)
	if err != nil {
		return err
	}

	// Update the user's insurance fund
	insuranceFund.Balance = newBalance
	// Store the updated user's insurance fund
	err = k.insuranceFunds.Set(ctx, user, insuranceFund)
	if err != nil {
		return err
	}

	return nil
}

// WithdrawFromUserInsuranceFund withdraws coins from the user's insurance fund
// and sends them to the user.
func (k *Keeper) WithdrawFromUserInsuranceFund(ctx context.Context, user string, amount sdk.Coins) error {
	insuranceFund, err := k.insuranceFunds.Get(ctx, user)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.ErrInsufficientInsuranceFundBalance
		}
		return err
	}

	// Ensure that the user can withdraw that amount from their insurance fund
	if !insuranceFund.Balance.IsAllGTE(amount) {
		return types.ErrInsufficientInsuranceFundBalance
	}

	newBalance := insuranceFund.Balance.Sub(amount...)

	err = k.beforeUserInsuranceFundModified(ctx, user, insuranceFund.Balance, newBalance)
	if err != nil {
		return err
	}

	// Update the user insurance fund
	insuranceFund.Balance = newBalance
	err = k.insuranceFunds.Set(ctx, user, insuranceFund)
	if err != nil {
		return err
	}

	// Send the coins back to the user
	userAddress, err := k.accountKeeper.AddressCodec().StringToBytes(user)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, userAddress, amount)
	if err != nil {
		return err
	}

	return nil
}

// beforeUserInsuranceFundModified is called before modifying user's insurance
// fund.
func (k *Keeper) beforeUserInsuranceFundModified(ctx context.Context, user string, oldInsuranceFund, newInsuranceFund sdk.Coins) error {
	delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(user)
	if err != nil {
		return err
	}

	// Calculate old coverable coins and new coverable coins
	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	// Iterating through all the user's delegations, withdraw rewards and update
	// each delegation target's total covered locked shares.
	err = k.restakingKeeper.IterateUserDelegations(ctx, user, func(del restakingtypes.Delegation) (stop bool, err error) {
		target, err := k.restakingKeeper.GetDelegationTarget(ctx, del.Type, del.TargetID)
		if err != nil {
			return true, err
		}

		oldCoveredLockedShares, err := types.GetCoveredLockedShares(target, del, oldInsuranceFund, params.InsurancePercentage)
		if err != nil {
			return true, err
		}
		newCoveredLockedShares, err := types.GetCoveredLockedShares(target, del, newInsuranceFund, params.InsurancePercentage)
		if err != nil {
			return true, err
		}

		prevTargetCoveredLockedShares, err := k.GetTargetCoveredLockedShares(ctx, del.Type, del.TargetID)
		if err != nil {
			return true, err
		}
		newTargetCoveredLockedShares := prevTargetCoveredLockedShares.Add(newCoveredLockedShares...).Sub(oldCoveredLockedShares)

		// Withdraw rewards
		k.restakingOverrider.GetDelegationTarget = func(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32) (restakingtypes.DelegationTarget, error) {
			uncoveredLockedShares := types.UncoveredLockedShares(target.GetDelegatorShares(), newTargetCoveredLockedShares)

			switch target := target.(type) {
			case poolstypes.Pool:
				target, _, err = target.RemoveDelShares(uncoveredLockedShares)
				if err != nil {
					return nil, err
				}
				return target, nil
			case operatorstypes.Operator:
				target, _ = target.RemoveDelShares(uncoveredLockedShares)
				return target, nil
			case servicestypes.Service:
				target, _ = target.RemoveDelShares(uncoveredLockedShares)
				return target, nil
			default:
				return nil, fmt.Errorf("invalid target type %T", target)
			}
		}
		k.restakingOverrider.GetDelegation = func(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32, delegator string) (restakingtypes.Delegation, bool, error) {
			del.Shares = types.DeductUncoveredLockedShares(del.Shares, newCoveredLockedShares)
			return del, true, nil
		}
		err = k.withRestakingOverrider(func() error {
			_, err := k.rewardsKeeper.WithdrawDelegationRewards(ctx, delAddr, del.Type, del.TargetID)
			return err
		})
		if err != nil {
			return true, err
		}

		if newTargetCoveredLockedShares.IsZero() {
			err = k.TargetsCoveredLockedShares.Remove(ctx, collections.Join(int32(del.Type), del.TargetID))
			if err != nil {
				return true, err
			}
		} else {
			err = k.TargetsCoveredLockedShares.Set(
				ctx,
				collections.Join(int32(del.Type), del.TargetID),
				types.TargetCoveredLockedShares{Shares: newTargetCoveredLockedShares},
			)
			if err != nil {
				return true, err
			}
		}

		return false, nil
	})
	return err
}

// GetUserInsuranceFund returns the user's insurance fund.
func (k *Keeper) GetUserInsuranceFund(ctx context.Context, user string) (types.UserInsuranceFund, error) {
	insuranceFund, err := k.insuranceFunds.Get(ctx, user)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.NewEmptyInsuranceFund(), nil
		} else {
			return types.UserInsuranceFund{}, err
		}
	}

	return insuranceFund, nil
}

// GetUserInsuranceFundBalance returns the amount of coins in the user's insurance fund.
func (k *Keeper) GetUserInsuranceFundBalance(ctx context.Context, user string) (sdk.Coins, error) {
	insuranceFund, err := k.GetUserInsuranceFund(ctx, user)
	if err != nil {
		return nil, err
	}

	return insuranceFund.Balance, nil
}

// GetInsuranceFundBalance returns the amount of coins in the insurance fund.
func (k *Keeper) GetInsuranceFundBalance(ctx context.Context) (sdk.Coins, error) {
	accAddr, err := sdk.AccAddressFromBech32(k.ModuleAddress)
	if err != nil {
		return nil, err
	}

	return k.bankKeeper.GetAllBalances(ctx, accAddr), nil
}

// GetUserUsedInsuranceFund returns the amount of coins that are used
// to cover the user's locked representation tokens that have been restaked.
func (k *Keeper) GetUserUsedInsuranceFund(ctx context.Context, userAddress string) (sdk.Coins, error) {
	// Get locked representations that the insurance fund covers
	lockedRepresentations, err := k.GetAllUserActiveLockedRepresentations(ctx, userAddress)
	if err != nil {
		return nil, err
	}

	// No locked representation tokens were restaked, the used
	// insurance fund is zero
	if lockedRepresentations.IsZero() {
		return sdk.NewCoins(), nil
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	userInsuranceFund, err := k.GetUserInsuranceFundBalance(ctx, userAddress)
	if err != nil {
		return nil, err
	}

	// Compute the used insurance fund
	usedInsuranceFund := sdk.NewCoins()
	for _, coin := range lockedRepresentations {
		nativeDenom, err := types.LockedDenomToNative(coin.Denom)
		if err != nil {
			return nil, err
		}
		requiredAmount := params.InsurancePercentage.Mul(coin.Amount).QuoInt64(100).Ceil().TruncateInt()
		usedInsuranceFund = usedInsuranceFund.Add(sdk.NewCoin(
			nativeDenom,
			// Pick the minimum between the required amount and the amount
			// in the insurance fund to avoid incorrect values.
			math.MinInt(requiredAmount, userInsuranceFund.AmountOf(nativeDenom)),
		))
	}

	return usedInsuranceFund, nil
}

// CanWithdrawFromInsuranceFund returns true if the user can withdraw the provided amount
// from their insurance fund.
func (k *Keeper) CanWithdrawFromInsuranceFund(ctx context.Context, user string, amount sdk.Coins) (bool, error) {
	userInsuranceFund, err := k.GetUserInsuranceFund(ctx, user)
	if err != nil {
		return false, err
	}
	// Ensure that the user has enough coins in the insurance fund
	if !userInsuranceFund.Balance.IsAllGTE(amount) {
		return false, nil
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		return false, err
	}

	// Get all the locked representations that are currently being
	// covered by the user's insurance fund.
	lockedRepresentations, err := k.GetAllUserActiveLockedRepresentations(ctx, user)
	if err != nil {
		return false, err
	}

	// Ensure that the user's insurance fund can cover the user's restaked
	// locked representations after the withdrawal.
	userInsuranceFund.Balance = userInsuranceFund.Balance.Sub(amount...)
	canCover, _ := userInsuranceFund.CanCoverDecCoins(params.InsurancePercentage, lockedRepresentations)

	return canCover, nil
}

func (k *Keeper) GetCoveredLockedShares(ctx context.Context, delegation restakingtypes.Delegation) (sdk.DecCoins, error) {
	// Get coverable dec coins by the user's insurance fund
	insuranceFund, err := k.GetUserInsuranceFundBalance(ctx, delegation.UserAddress)
	if err != nil {
		return nil, err
	}
	// Exit early if the user doesn't have insurance fund balance
	if insuranceFund.IsZero() {
		return nil, nil
	}
	target, err := k.restakingKeeper.GetDelegationTarget(ctx, delegation.Type, delegation.TargetID)
	if err != nil {
		return nil, err
	}
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	return types.GetCoveredLockedShares(target, delegation, insuranceFund, params.InsurancePercentage)
}
