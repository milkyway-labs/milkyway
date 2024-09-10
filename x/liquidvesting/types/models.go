package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewInsuranceFund() UserInsuranceFund {
	return UserInsuranceFund{
		Balance: sdk.NewCoins(),
	}
}

func (u *UserInsuranceFund) Add(amount sdk.Coins) {
	u.Balance = u.Balance.Sort().Add(amount.Sort()...)
}
