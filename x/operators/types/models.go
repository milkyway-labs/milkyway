package types

import (
	"fmt"
	"strconv"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/milkyway-labs/milkyway/utils"
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
	params OperatorParams,
) Operator {
	return Operator{
		ID:         id,
		Status:     status,
		Admin:      admin,
		Moniker:    moniker,
		Website:    website,
		PictureURL: pictureURL,
		Address:    GetOperatorAddress(id).String(),
		Params:     params,
	}
}

func (o Operator) GetID() uint32 {
	return o.ID
}

func (o Operator) GetAddress() string {
	return o.Address
}

func (o Operator) GetTokens() sdk.Coins {
	return o.Tokens
}

func (o Operator) GetDelegatorShares() sdk.DecCoins {
	return o.DelegatorShares
}

// Validate checks that the Operator has valid values.
func (o Operator) Validate() error {
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

	err = o.Params.Validate()
	if err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	return nil
}

// GetSharesDenom returns the shares denom for an operator and token denom
func (o Operator) GetSharesDenom(tokenDenom string) string {
	return utils.GetSharesDenomFromTokenDenom("operator", o.ID, tokenDenom)
}

// IsActive returns whether the operator is active.
func (o Operator) IsActive() bool {
	return o.Status == OPERATOR_STATUS_ACTIVE
}

// InvalidExRate returns whether the exchange rates is invalid.
// This can happen e.g. if Operator loses all tokens due to slashing. In this case,
// make all future delegations invalid.
func (o Operator) InvalidExRate() bool {
	return utils.IsInvalidExRate(o.Tokens, o.DelegatorShares)
}

// TokensFromShares calculates the token worth of provided shares
func (o Operator) TokensFromShares(shares sdk.DecCoins) sdk.DecCoins {
	return utils.ComputeTokensFromShares(shares, o.Tokens, o.DelegatorShares)
}

// TokensFromSharesTruncated calculates the token worth of provided shares, truncated
func (o Operator) TokensFromSharesTruncated(shares sdk.DecCoins) sdk.DecCoins {
	return utils.ComputeTokensFromSharesTruncated(shares, o.Tokens, o.DelegatorShares)
}

// SharesFromTokens returns the shares of a delegation given a bond amount. It
// returns an error if the operator has no tokens.
func (o Operator) SharesFromTokens(tokens sdk.Coins) (sdk.DecCoins, error) {
	if o.Tokens.IsZero() {
		return sdk.NewDecCoins(), ErrInsufficientShares
	}
	return utils.SharesFromTokens(tokens, o.GetSharesDenom, o.Tokens, o.DelegatorShares)
}

// SharesFromTokensTruncated returns the truncated shares of a delegation given a bond amount.
// It returns an error if the operator has no tokens.
func (o Operator) SharesFromTokensTruncated(tokens sdk.Coins) (sdk.DecCoins, error) {
	if o.Tokens.IsZero() {
		return sdk.NewDecCoins(), ErrInsufficientShares
	}
	return utils.SharesFromTokensTruncated(tokens, o.GetSharesDenom, o.Tokens, o.DelegatorShares)
}

// SharesFromDecCoins returns the shares of a delegation given a bond amount. It
// returns an error if the operator has no tokens.
func (o Operator) SharesFromDecCoins(tokens sdk.DecCoins) (sdk.DecCoins, error) {
	if o.Tokens.IsZero() {
		return sdk.NewDecCoins(), ErrInsufficientShares
	}
	return utils.SharesFromDecCoins(tokens, o.GetSharesDenom, o.Tokens, o.DelegatorShares)
}

// AddTokensFromDelegation adds the given amount of tokens to the operator's total tokens,
// also updating the operator's delegator shares.
// It returns the updated operator and the shares issued.
func (o Operator) AddTokensFromDelegation(amount sdk.Coins) (Operator, sdk.DecCoins) {
	issuedShares := utils.IssueShares(amount, o.GetSharesDenom, o.Tokens, o.DelegatorShares)

	o.Tokens = o.Tokens.Add(amount...)
	o.DelegatorShares = o.DelegatorShares.Add(issuedShares...)

	return o, issuedShares
}

// RemoveDelShares removes delegator shares from an operator and returns
// the amount of tokens that should be issued for those shares.
// NOTE: Because token fractions are left in the operator,
// the exchange rate of future shares of this validator can increase.
func (o Operator) RemoveDelShares(delShares sdk.DecCoins) (Operator, sdk.Coins) {
	remainingShares := o.DelegatorShares.Sub(delShares)

	var issuedTokens sdk.Coins
	if remainingShares.IsZero() {
		// Last delegation share gets any trimmings
		issuedTokens = o.Tokens
		o.Tokens = sdk.NewCoins()
	} else {
		// Leave excess tokens in the operator
		// However, fully use all the delegator shares
		issuedTokens, _ = o.TokensFromShares(delShares).TruncateDecimal()
		o.Tokens = o.Tokens.Sub(issuedTokens...)

		if o.Tokens.IsAnyNegative() {
			panic("attempting to remove more tokens than available in operator")
		}
	}

	o.DelegatorShares = remainingShares

	return o, issuedTokens
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
		o.Params,
	)
}

// --------------------------------------------------------------------------------------------------------------------

// NewOperatorParams creates a new OperatorParams instance
func NewOperatorParams(commissionRate math.LegacyDec) OperatorParams {
	return OperatorParams{
		CommissionRate: commissionRate,
	}
}

// DefaultOperatorParams returns the default operator params
func DefaultOperatorParams() OperatorParams {
	return NewOperatorParams(math.LegacyZeroDec())
}

// Validate validates the operator params
func (p *OperatorParams) Validate() error {
	if p.CommissionRate.IsNegative() || p.CommissionRate.GT(math.LegacyOneDec()) {
		return fmt.Errorf("invalid commission rate: %s", p.CommissionRate.String())
	}

	return nil
}
