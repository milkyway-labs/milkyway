package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DepositToUserInsuranceFund deposits coins to the user's insurance fund.
func (k *Keeper) DepositToUserInsuranceFund(
	goCtx context.Context,
	user sdk.AccAddress,
	amount sdk.Coins,
) error {
	panic("unimplemented")
}

// WithdrawFromUserInsuranceFund withdraws coins from the user's insurance fund.
func (k *Keeper) WithdrawFromUserInsuranceFund(
	goCtx context.Context,
	user sdk.AccAddress,
	amount sdk.Coins,
) error {
	panic("unimplemented")
}

// GetUserInsuranceFundBalance returns the amount of coins in the user's insurance fund.
func (k *Keeper) GetUserInsuranceFundBalance(
	goCtx context.Context,
	user sdk.AccAddress,
) (sdk.Coins, error) {
	panic("unimplemented")
}

// GetInsuranceFundBalance returns the amount of coins in the insurance fund.
func (k *Keeper) GetInsuranceFundBalance(goCtx context.Context) (sdk.Coins, error) {
	panic("unimplemented")
}
