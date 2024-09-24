package types

import (
	fmt "fmt"

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
		return fmt.Errorf("invalid depositor address: %s", err)
	}
	return i.Amount.Validate()
}

// MsgDepositInsurance defines a struct for depositing tokens
// into the insurance fund.
type MsgDepositInsurance struct {
	Amounts []InsuranceDeposit `json:"amounts"`
}

func (msg MsgDepositInsurance) ValidateBasic() error {
	for i, deposit := range msg.Amounts {
		// Ensure that the deposits have all the same denom
		if i > 0 && deposit.Amount.Denom != msg.Amounts[i].Amount.Denom {
			return fmt.Errorf("can't deposit multiple coins")
		}

		if err := deposit.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}

func (msg MsgDepositInsurance) GetTotalDepositAmount() (*sdk.Coin, error) {
	if len(msg.Amounts) == 0 {
		return nil, fmt.Errorf("no coins to deposit")
	}

	totalAmount := msg.Amounts[0].Amount
	for _, deposit := range msg.Amounts[1:] {
		// Ensure that the deposits have all the same denom
		if deposit.Amount.Denom != totalAmount.Denom {
			return nil, fmt.Errorf("can't deposit multiple denoms")
		} else {
			totalAmount = totalAmount.Add(deposit.Amount)
		}
	}

	return &totalAmount, nil
}
