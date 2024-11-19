package types

import (
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewInsuranceFund(balance sdk.Coins) UserInsuranceFund {
	return UserInsuranceFund{
		Balance: balance,
	}
}

func NewEmptyInsuranceFund() UserInsuranceFund {
	return NewInsuranceFund(sdk.NewCoins())
}

func (u *UserInsuranceFund) Add(amount sdk.Coins) {
	u.Balance = u.Balance.Add(amount...)
}

func (u *UserInsuranceFund) Validate() error {
	if err := u.Balance.Validate(); err != nil {
		return err
	}

	return nil
}

func (u *UserInsuranceFund) CanCoverDecCoins(insurancePercentage sdkmath.LegacyDec, coins sdk.DecCoins) (bool, sdk.Coins, error) {
	required := sdk.NewCoins()
	for _, coin := range coins {
		if IsVestedRepresentationDenom(coin.Denom) {
			nativeDenom, err := VestedDenomToNative(coin.Denom)
			if err != nil {
				return false, nil, err
			}
			required = required.Add(sdk.NewCoin(nativeDenom, insurancePercentage.Mul(coin.Amount).QuoInt64(100).Ceil().TruncateInt()))
		}
	}

	return u.Balance.IsAllGTE(required), required, nil
}

func NewBurnCoins(delegator string, completionTime time.Time, amount sdk.Coins) BurnCoins {
	return BurnCoins{
		DelegatorAddress: delegator,
		CompletionTime:   completionTime,
		Amount:           amount,
	}
}

func (bc *BurnCoins) Validate() error {
	if _, err := sdk.AccAddressFromBech32(bc.DelegatorAddress); err != nil {
		return err
	}
	if bc.CompletionTime.IsZero() {
		return fmt.Errorf("invalid completion time")
	}
	return bc.Amount.Validate()
}

// NewUserInsuranceFundState creates a new UserInsuranceFundState.
func NewUserInsuranceFundState(
	userAddress string,
	insuranceFund UserInsuranceFund,
) UserInsuranceFundState {
	return UserInsuranceFundState{
		UserAddress:   userAddress,
		InsuranceFund: insuranceFund,
	}
}

func (uif *UserInsuranceFundState) Validate() error {
	if _, err := sdk.AccAddressFromBech32(uif.UserAddress); err != nil {
		return err
	}

	return uif.InsuranceFund.Validate()
}
