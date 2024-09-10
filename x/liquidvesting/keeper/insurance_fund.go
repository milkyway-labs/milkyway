package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

// AddToUserInsuranceFund adds the provided amount to the user's insurance fund.
func (k *Keeper) AddToUserInsuranceFund(
	ctx sdk.Context,
	user sdk.AccAddress,
	amount sdk.Coins,
) error {
	insuranceFund, err := k.InsuranceFunds.Get(ctx, user)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			insuranceFund = types.NewInsuranceFund()
		} else {
			return err
		}
	}

	// Update the user's insurance fund
	insuranceFund.Add(amount)
	// Store the updated user's insurance fund
	err = k.InsuranceFunds.Set(ctx, user, insuranceFund)
	if err != nil {
		return err
	}

	// Dispatch the deposit event.
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventDepositToUserInsuranceFund,
			sdk.NewAttribute("user", user.String()),
			sdk.NewAttribute("deposited", amount.String()),
		),
	)

	return nil
}

// WithdrawFromUserInsuranceFund withdraws coins from the user's insurance fund
// and sends them to the user.
func (k *Keeper) WithdrawFromUserInsuranceFund(
	ctx sdk.Context,
	user sdk.AccAddress,
	amount sdk.Coins,
) error {
	panic("unimplemented")
}

// GetUserInsuranceFundBalance returns the amount of coins in the user's insurance fund.
func (k *Keeper) GetUserInsuranceFundBalance(
	ctx context.Context,
	user sdk.AccAddress,
) (sdk.Coins, error) {
	insuranceFund, err := k.InsuranceFunds.Get(ctx, user)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return sdk.NewCoins(), nil
		} else {
			return nil, err
		}
	}

	return insuranceFund.Balance, nil
}

// GetInsuranceFundBalance returns the amount of coins in the insurance fund.
func (k *Keeper) GetInsuranceFundBalance(ctx context.Context) (sdk.Coins, error) {
	accAddr, err := sdk.AccAddressFromBech32(k.ModuleAddress)
	if err != nil {
		return nil, err
	}

	return k.BankKeeper.GetAllBalances(ctx, accAddr), nil
}
