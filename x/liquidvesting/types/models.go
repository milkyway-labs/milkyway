package types

import (
	"fmt"
	"time"

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
	u.Balance = u.Balance.Sort().Add(amount.Sort()...)
}

func (u *UserInsuranceFund) Validate() error {
	return u.Balance.Validate()
}

func NewBurnCoins(delegator string, completionTime time.Time, amount sdk.Coins) BurnCoins {
	return BurnCoins{
		DelegatorAddress: delegator,
		CompletionTime:   completionTime,
		Amount:           amount,
	}
}

func (bc BurnCoins) Validate() error {
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

func (uif UserInsuranceFundState) Validate() error {
	if _, err := sdk.AccAddressFromBech32(uif.UserAddress); err != nil {
		return err
	}

	return uif.InsuranceFund.Validate()
}
