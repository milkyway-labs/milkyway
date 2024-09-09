package types

import (
	fmt "fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InsuranceDeposit defines an individual deposit into the insurance fund.
type InsuranceDeposit struct {
	// Address of the user that deposited the tokens.
	Depositor string `json:"depositor"`
	// Amount of tokens deposited by the user in the insurance fund.
	Amount sdk.Coin `json:"amount"`
}

func (i InsuranceDeposit) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(i.Depositor); err != nil {
		return err
	}
	return i.Amount.Validate()
}

// MsgDepositInsurance defines a struct for depositing tokens
// into the insurance fund.
type MsgDepositInsurance struct {
	Amounts []InsuranceDeposit `json:"amounts"`
}

func (msg MsgDepositInsurance) ValidateBasic() error {
	denoms := make([]string, 0)
	for _, deposit := range msg.Amounts {
		// Ensure that the deposits have all the same denom
		if len(denoms) == 0 {
			denoms = append(denoms, deposit.Amount.Denom)
		} else if denoms[0] != deposit.Amount.Denom {
			return fmt.Errorf("can't deposit multiple coins")
		}

		if err := deposit.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}

func (msg MsgDepositInsurance) GetTotalDepositAmount() (sdk.Coin, error) {
	if len(msg.Amounts) == 0 {
		return sdk.NewCoin("", math.NewInt(0)), fmt.Errorf("no coins to deposit")
	}

	totalAmount := msg.Amounts[0].Amount
	for i, deposit := range msg.Amounts {
		// Skip the first deposit since we have initialized the total amount
		// with the first deposit amount.
		if i == 0 {
			continue
		}

		// Ensure that the deposits have all the same denom
		if deposit.Amount.Denom != totalAmount.Denom {
			return sdk.NewCoin("", math.NewInt(0)), fmt.Errorf("can't deposit multiple denoms")
		} else {
			totalAmount = totalAmount.Add(deposit.Amount)
		}
	}

	return totalAmount, nil
}
