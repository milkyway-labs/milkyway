package types

import (
	"fmt"
	"strconv"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// GetOperatorAddress generates an operator address from its id
func GetOperatorAddress(operatorID uint32) sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("operator-%d", operatorID))
}

// ParseOperatorID tries parsing the given value as an operator id
func ParseOperatorID(value string) (uint32, error) {
	operatorID, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid operator ID: %s", value)
	}
	return uint32(operatorID), nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewOperator creates a new Operator object
func NewOperator(
	id uint32,
	status OperatorStatus,
	moniker string,
	website string,
	pictureURL string,
	admin string,
) Operator {
	return Operator{
		ID:              id,
		Status:          status,
		Admin:           admin,
		Moniker:         moniker,
		Website:         website,
		PictureURL:      pictureURL,
		Address:         GetOperatorAddress(id).String(),
		Tokens:          sdk.NewCoins(),
		DelegatorShares: sdk.NewDecCoins(),
	}
}

// Validate checks that the Operator has valid values.
func (o *Operator) Validate() error {
	if o.ID == 0 {
		return fmt.Errorf("invalid id: %d", o.ID)
	}

	if o.Status == OPERATOR_STATUS_UNSPECIFIED {
		return fmt.Errorf("invalid status: %s", o.Status)
	}

	if strings.TrimSpace(o.Moniker) == "" {
		return fmt.Errorf("invalid moniker: %s", o.Moniker)
	}

	_, err := sdk.AccAddressFromBech32(o.Admin)
	if err != nil {
		return fmt.Errorf("invalid admin address: %s", o.Admin)
	}

	_, err = sdk.AccAddressFromBech32(o.Address)
	if err != nil {
		return fmt.Errorf("invalid address: %s", o.Address)
	}

	return nil
}

// IsActive returns whether the operator is active.
func (o Operator) IsActive() bool {
	return o.Status == OPERATOR_STATUS_ACTIVE
}

// InvalidExRate returns whether the exchange rates is invalid.
// This can happen e.g. if Operator loses all tokens due to slashing. In this case,
// make all future delegations invalid.
func (o Operator) InvalidExRate() bool {
	for _, token := range o.Tokens {
		if token.IsZero() && o.DelegatorShares.AmountOf(token.Denom).IsPositive() {
			return true
		}
	}
	return false
}

// SharesFromTokens returns the shares of a delegation given a bond amount. It
// returns an error if the operator has no tokens.
func (o Operator) SharesFromTokens(tokens sdk.Coin) (sdkmath.LegacyDec, error) {
	if o.Tokens.IsZero() {
		return sdkmath.LegacyZeroDec(), ErrInsufficientShares
	}

	delegatorTokenShares := o.DelegatorShares.AmountOf(tokens.Denom)
	operatorTokenAmount := o.Tokens.AmountOf(tokens.Denom)

	return delegatorTokenShares.MulInt(tokens.Amount).QuoInt(operatorTokenAmount), nil
}

// AddTokensFromDelegation adds the given amount of tokens to the operator's total tokens,
// also updating the operator's delegator shares.
// It returns the updated operator and the shares issued.
func (o Operator) AddTokensFromDelegation(amount sdk.Coins) (Operator, sdk.DecCoins) {
	// calculate the shares to issue
	issuedShares := sdk.NewDecCoins()
	for _, token := range amount {
		var tokenShares sdk.DecCoin
		delegatorShares := o.DelegatorShares.AmountOf(token.Denom)

		if delegatorShares.IsZero() {
			// The first delegation to an operator sets the exchange rate to one
			tokenShares = sdk.NewDecCoinFromCoin(token)
		} else {
			shares, err := o.SharesFromTokens(token)
			if err != nil {
				panic(err)
			}
			tokenShares = sdk.NewDecCoinFromDec(token.Denom, shares)
		}

		issuedShares = issuedShares.Add(tokenShares)
	}

	o.Tokens = o.Tokens.Add(amount...)
	o.DelegatorShares = o.DelegatorShares.Add(issuedShares...)

	return o, issuedShares
}

// --------------------------------------------------------------------------------------------------------------------

// OperatorUpdate defines the fields that can be updated in an Operator.
type OperatorUpdate struct {
	Moniker    string
	Website    string
	PictureURL string
}

func NewOperatorUpdate(
	moniker string,
	website string,
	pictureURL string,
) OperatorUpdate {
	return OperatorUpdate{
		Moniker:    moniker,
		Website:    website,
		PictureURL: pictureURL,
	}
}

// Update returns a new Operator with updated fields.
func (o *Operator) Update(update OperatorUpdate) Operator {
	if update.Moniker == DoNotModify {
		update.Moniker = o.Moniker
	}

	if update.Website == DoNotModify {
		update.Website = o.Website
	}

	if update.PictureURL == DoNotModify {
		update.PictureURL = o.PictureURL
	}

	return NewOperator(
		o.ID,
		o.Status,
		update.Moniker,
		update.Website,
		update.PictureURL,
		o.Admin,
	)
}
