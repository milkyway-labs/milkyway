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
	Amount math.Int `json:"amount"`
}

func (i *InsuranceDeposit) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(i.Depositor); err != nil {
		return fmt.Errorf("invalid depositor address: %s", err)
	}
	return nil
}

// IsZero returns true if the amount is zero and false otherwise.
func (i *InsuranceDeposit) IsZero() bool {
	return i.Amount.IsZero()
}

// IsPositive returns true if the amount is positive and false otherwise.
func (i *InsuranceDeposit) IsPositive() bool {
	return i.Amount.IsPositive()
}

// IsNegative returns true if the amount is negative and false otherwise.
func (i *InsuranceDeposit) IsNegative() bool {
	return i.Amount.IsNegative()
}

// MsgDepositInsurance defines a struct for depositing tokens
// into the insurance fund.
type MsgDepositInsurance struct {
	Amounts []InsuranceDeposit `json:"amounts"`
}

// NewMsgDepositInsurance creates a new MsgDepositInsurance instance.
func NewMsgDepositInsurance(amounts []InsuranceDeposit) *MsgDepositInsurance {
	return &MsgDepositInsurance{Amounts: amounts}
}

func (msg *MsgDepositInsurance) ValidateBasic() error {
	for _, deposit := range msg.Amounts {
		if err := deposit.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}

func (msg *MsgDepositInsurance) GetTotalDepositAmount() math.Int {
	totalAmount := math.ZeroInt()
	for _, deposit := range msg.Amounts {
		totalAmount = totalAmount.Add(deposit.Amount)
	}

	return totalAmount
}
