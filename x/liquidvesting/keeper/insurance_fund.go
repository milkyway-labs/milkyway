package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

// AddToUserInsuranceFund adds the provided amount to the user's insurance fund.
// NOTE: We assume that the amount that will be added to the user's insurance fund
// is already present in the module account balance.
func (k *Keeper) AddToUserInsuranceFund(ctx context.Context, user sdk.AccAddress, amount sdk.Coins) error {
	insuranceFund, err := k.insuranceFunds.Get(ctx, user)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			insuranceFund = types.NewEmptyInsuranceFund()
		} else {
			return err
		}
	}

	// Update the user's insurance fund
	insuranceFund.Add(amount)
	// Store the updated user's insurance fund
	return k.insuranceFunds.Set(ctx, user, insuranceFund)
}

// WithdrawFromUserInsuranceFund withdraws coins from the user's insurance fund
// and sends them to the user.
func (k *Keeper) WithdrawFromUserInsuranceFund(ctx context.Context, user sdk.AccAddress, amount sdk.Coins) error {
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

	// Update the user insurance fund
	insuranceFund.Balance = insuranceFund.Balance.Sub(amount...)
	err = k.insuranceFunds.Set(ctx, user, insuranceFund)
	if err != nil {
		return err
	}

	// Send the coins back to the user
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, user, amount)
	if err != nil {
		return err
	}

	return nil
}

// GetUserInsuranceFund returns the user's insurance fund.
func (k *Keeper) GetUserInsuranceFund(ctx context.Context, user sdk.AccAddress) (types.UserInsuranceFund, error) {
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
func (k *Keeper) GetUserInsuranceFundBalance(ctx context.Context, user sdk.AccAddress) (sdk.Coins, error) {
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

// CanWithdrawFromInsuranceFund returns true if the user can withdraw the provided amount
// from their insurance fund.
func (k *Keeper) CanWithdrawFromInsuranceFund(ctx context.Context, user sdk.AccAddress, amount sdk.Coins) (bool, error) {
	userInsuranceFund, err := k.GetUserInsuranceFund(ctx, user)
	if err != nil {
		return false, err
	}

	return userInsuranceFund.Unused().IsAllGTE(amount), nil
}
