package types

import (
	"fmt"
	"strconv"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/milkyway-labs/milkyway/utils"
)

// GetServiceAddress generates a service address from its id
func GetServiceAddress(serviceID uint32) sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("service-%d", serviceID))
}

// ParseServiceID parses a string into a uint32
func ParseServiceID(value string) (uint32, error) {
	id, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewService creates a new Service instance
func NewService(
	id uint32,
	status ServiceStatus,
	name string,
	description string,
	website string,
	pictureURL string,
	admin string,
) Service {
	return Service{
		ID:          id,
		Status:      status,
		Name:        name,
		Description: description,
		Website:     website,
		PictureURL:  pictureURL,
		Admin:       admin,
		Address:     GetServiceAddress(id).String(),
	}
}

// Validate checks that the Service has valid values.
func (s Service) Validate() error {
	if s.Status == SERVICE_STATUS_UNSPECIFIED {
		return fmt.Errorf("invalid status: %s", s.Status)
	}

	if s.ID == 0 {
		return fmt.Errorf("invalid id: %d", s.ID)
	}

	if strings.TrimSpace(s.Name) == "" {
		return fmt.Errorf("invalid name: %s", s.Name)
	}

	_, err := sdk.AccAddressFromBech32(s.Admin)
	if err != nil {
		return fmt.Errorf("invalid admin address")
	}

	_, err = sdk.AccAddressFromBech32(s.Address)
	if err != nil {
		return fmt.Errorf("invalid service address")
	}

	return nil
}

// GetSharesDenom returns the shares denom for a service and token denom
func (s Service) GetSharesDenom(tokenDenom string) string {
	return utils.GetSharesDenomFromTokenDenom("service", s.ID, tokenDenom)
}

// GetTokenDenomFromSharesDenom returns the token denom from a shares denom
func (s Service) GetTokenDenomFromSharesDenom(sharesDenom string) string {
	return utils.GetTokenDenomFromSharesDenom(sharesDenom)
}

// IsActive returns whether the service is active.
func (s Service) IsActive() bool {
	return s.Status == SERVICE_STATUS_ACTIVE
}

// InvalidExRate returns whether the exchange rates is invalid.
// This can happen e.g. if Service loses all tokens due to slashing. In this case,
// make all future delegations invalid.
func (s Service) InvalidExRate() bool {
	for _, token := range s.Tokens {
		if token.IsZero() && s.DelegatorShares.AmountOf(token.Denom).IsPositive() {
			return true
		}
	}
	return false
}

// TokensFromShares calculates the token worth of provided shares
func (s Service) TokensFromShares(shares sdk.DecCoins) sdk.DecCoins {
	tokens := sdk.NewDecCoins()
	for _, share := range shares {
		tokenDenom := s.GetTokenDenomFromSharesDenom(share.Denom)
		operatorTokenAmount := s.Tokens.AmountOf(tokenDenom)
		delegatorSharesAmount := s.DelegatorShares.AmountOf(share.Denom)

		tokenAmount := share.Amount.MulInt(operatorTokenAmount).Quo(delegatorSharesAmount)

		tokens = tokens.Add(sdk.NewDecCoinFromDec(tokenDenom, tokenAmount))
	}

	return tokens
}

// SharesFromTokens returns the shares of a delegation given a bond amount. It
// returns an error if the service has no tokens.
func (s Service) SharesFromTokens(tokens sdk.Coin) (sdkmath.LegacyDec, error) {
	if s.Tokens.IsZero() {
		return sdkmath.LegacyZeroDec(), ErrInsufficientShares
	}

	sharesDenom := s.GetSharesDenom(tokens.Denom)
	delegatorTokenShares := s.DelegatorShares.AmountOf(sharesDenom)
	operatorTokenAmount := s.Tokens.AmountOf(tokens.Denom)

	return delegatorTokenShares.MulInt(tokens.Amount).QuoInt(operatorTokenAmount), nil
}

// AddTokensFromDelegation adds the given amount of tokens to the service's total tokens,
// also updating the service's delegator shares.
// It returns the updated service and the shares issued.
func (s Service) AddTokensFromDelegation(amount sdk.Coins) (Service, sdk.DecCoins) {
	// calculate the shares to issue
	issuedShares := sdk.NewDecCoins()
	for _, token := range amount {
		var tokenShares sdk.DecCoin
		sharesDenom := s.GetSharesDenom(token.Denom)

		delegatorShares := s.DelegatorShares.AmountOf(sharesDenom)
		if delegatorShares.IsZero() {
			// The first delegation to an operator sets the exchange rate to one
			tokenShares = sdk.NewDecCoin(sharesDenom, token.Amount)
		} else {
			shares, err := s.SharesFromTokens(token)
			if err != nil {
				panic(err)
			}
			tokenShares = sdk.NewDecCoinFromDec(sharesDenom, shares)
		}

		issuedShares = issuedShares.Add(tokenShares)
	}

	s.Tokens = s.Tokens.Add(amount...)
	s.DelegatorShares = s.DelegatorShares.Add(issuedShares...)

	return s, issuedShares
}

// --------------------------------------------------------------------------------------------------------------------

// ServiceUpdate defines the fields that can be updated in a Service.
type ServiceUpdate struct {
	Name        string
	Description string
	Website     string
	PictureURL  string
}

// NewServiceUpdate returns a new ServiceUpdate instance.
func NewServiceUpdate(
	name string,
	description string,
	website string,
	pictureURL string,
) ServiceUpdate {
	return ServiceUpdate{
		Name:        name,
		Description: description,
		Website:     website,
		PictureURL:  pictureURL,
	}
}

// Update returns a new Service with updated fields.
func (a *Service) Update(update ServiceUpdate) Service {
	if update.Name == DoNotModify {
		update.Name = a.Name
	}

	if update.Description == DoNotModify {
		update.Description = a.Description
	}

	if update.Website == DoNotModify {
		update.Website = a.Website
	}

	if update.PictureURL == DoNotModify {
		update.PictureURL = a.PictureURL
	}

	return NewService(
		a.ID,
		a.Status,
		update.Name,
		update.Description,
		update.Website,
		update.PictureURL,
		a.Admin,
	)
}
