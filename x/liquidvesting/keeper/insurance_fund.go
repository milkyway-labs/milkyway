package keeper

import (
	"errors"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

// AddToUserInsuranceFund adds the provided amount to the user's insurance fund.
// NOTE: We assume that the amount that will be added to the user's insurance fund
// is already present in the module account balance.
func (k *Keeper) AddToUserInsuranceFund(
	ctx sdk.Context,
	user sdk.AccAddress,
	amount sdk.Coins,
) error {
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
func (k *Keeper) WithdrawFromUserInsuranceFund(
	ctx sdk.Context,
	user sdk.AccAddress,
	amount sdk.Coins,
) error {
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
func (k *Keeper) GetUserInsuranceFund(
	ctx sdk.Context,
	user sdk.AccAddress,
) (types.UserInsuranceFund, error) {
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
func (k *Keeper) GetUserInsuranceFundBalance(
	ctx sdk.Context,
	user string,
) (sdk.Coins, error) {
	accAddr, err := sdk.AccAddressFromBech32(user)
	if err != nil {
		return nil, err
	}

	insuranceFund, err := k.GetUserInsuranceFund(ctx, accAddr)
	if err != nil {
		return nil, err
	}

	return insuranceFund.Balance, nil
}

// GetInsuranceFundBalance returns the amount of coins in the insurance fund.
func (k *Keeper) GetInsuranceFundBalance(ctx sdk.Context) (sdk.Coins, error) {
	accAddr, err := sdk.AccAddressFromBech32(k.ModuleAddress)
	if err != nil {
		return nil, err
	}

	return k.bankKeeper.GetAllBalances(ctx, accAddr), nil
}

// GetUserUsedInsuranceFund returns the amount of coins that are used
// to cover the user's vested representation tokens that have been restaked.
func (k *Keeper) GetUserUsedInsuranceFund(ctx sdk.Context, userAddress string) (sdk.Coins, error) {
	// Get vested representations that the insurance fund covers
	vestedRepresentations, err := k.GetAllUserActiveVestedRepresentations(ctx, userAddress)
	if err != nil {
		return nil, err
	}

	// No vested representation tokens were restaked, the used
	// insurance fund is zero
	if vestedRepresentations.IsZero() {
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
	for _, coin := range vestedRepresentations {
		nativeDenom, err := types.VestedDenomToNative(coin.Denom)
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
func (k *Keeper) CanWithdrawFromInsuranceFund(ctx sdk.Context, user sdk.AccAddress, amount sdk.Coins) (bool, error) {
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

	userAddress, err := k.accountKeeper.AddressCodec().BytesToString(user)
	if err != nil {
		return false, err
	}

	// Get all the vested representations that are currently being
	// covered by the user's insurance fund.
	vestedRepresentations, err := k.GetAllUserActiveVestedRepresentations(ctx, userAddress)
	if err != nil {
		return false, err
	}

	// Ensure that the user's insurance fund can cover the user's restaked
	// vested representations after the withdrawal.
	userInsuranceFund.Balance = userInsuranceFund.Balance.Sub(amount...)
	canCover, _, err := userInsuranceFund.CanCoverDecCoins(params.InsurancePercentage, vestedRepresentations)

	return canCover, err
}
