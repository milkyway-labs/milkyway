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

func (u *UserInsuranceFund) Validate() error {
	if err := u.Balance.Validate(); err != nil {
		return err
	}

	return nil
}

func (u *UserInsuranceFund) CanCoverDecCoins(insurancePercentage sdkmath.LegacyDec, coins sdk.DecCoins) (bool, sdk.Coins) {
	return CanCoverDecCoins(u.Balance, insurancePercentage, coins)
}

// CoverableDecCoins returns the amount of dec coins that can be covered by the
// insurance fund.
func (u *UserInsuranceFund) CoverableDecCoins(insurancePercentage sdkmath.LegacyDec) sdk.DecCoins {
	return GetCoverableDecCoins(u.Balance, insurancePercentage)
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

// NewUserInsuranceFundEntry creates a new UserInsuranceFundState.
func NewUserInsuranceFundEntry(userAddress string, balance sdk.Coins) UserInsuranceFundEntry {
	return UserInsuranceFundEntry{
		UserAddress: userAddress,
		Balance:     balance,
	}
}

func (uif *UserInsuranceFundEntry) Validate() error {
	_, err := sdk.AccAddressFromBech32(uif.UserAddress)
	if err != nil {
		return err
	}

	return uif.Balance.Validate()
}
