package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewInsuranceFund(balance sdk.Coins, used sdk.Coins) UserInsuranceFund {
	return UserInsuranceFund{
		Balance: balance,
		Used:    used,
	}
}

func NewEmptyInsuranceFund() UserInsuranceFund {
	return NewInsuranceFund(sdk.NewCoins(), sdk.NewCoins())
}

func (u *UserInsuranceFund) Add(amount sdk.Coins) {
	u.Balance = u.Balance.Add(amount...)
}

func (u *UserInsuranceFund) AddUsed(amount ...sdk.Coin) {
	u.Used = u.Used.Add(amount...)
}

func (u *UserInsuranceFund) DecreaseUsed(amount ...sdk.Coin) {
	u.Used = u.Used.Sub(amount...)
}

func (u *UserInsuranceFund) Validate() error {
	if err := u.Balance.Validate(); err != nil {
		return err
	}
	if err := u.Used.Validate(); err != nil {
		return err
	}
	if !u.Balance.IsAllGTE(u.Used) {
		return fmt.Errorf("used balance should be lower then total insurance fund balance")
	}

	return nil
}

// Unused returns the amount of coins that are not being used to
// cover restaking positions
func (u *UserInsuranceFund) Unused() sdk.Coins {
	return u.Balance.Sub(u.Used...)
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
