package types

import (
	"fmt"
	"strconv"
	"strings"

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

func (s Service) GetID() uint32 {
	return s.ID
}

func (s Service) GetAddress() string {
	return s.Address
}

func (s Service) GetTokens() sdk.Coins {
	return s.Tokens
}

func (s Service) GetDelegatorShares() sdk.DecCoins {
	return s.DelegatorShares
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

// IsActive returns whether the service is active.
func (s Service) IsActive() bool {
	return s.Status == SERVICE_STATUS_ACTIVE
}

// InvalidExRate returns whether the exchange rates is invalid.
// This can happen e.g. if Service loses all tokens due to slashing. In this case,
// make all future delegations invalid.
func (s Service) InvalidExRate() bool {
	return utils.IsInvalidExRate(s.Tokens, s.DelegatorShares)
}

// TokensFromShares calculates the token worth of provided shares
func (s Service) TokensFromShares(shares sdk.DecCoins) sdk.DecCoins {
	return utils.ComputeTokensFromShares(shares, s.Tokens, s.DelegatorShares)
}

// TokensFromSharesTruncated calculates the token worth of provided shares, truncated
func (s Service) TokensFromSharesTruncated(shares sdk.DecCoins) sdk.DecCoins {
	return utils.ComputeTokensFromSharesTruncated(shares, s.Tokens, s.DelegatorShares)
}

// SharesFromTokens returns the shares of a delegation given a bond amount. It
// returns an error if the service has no tokens.
func (s Service) SharesFromTokens(tokens sdk.Coins) (sdk.DecCoins, error) {
	if s.Tokens.IsZero() {
		return sdk.NewDecCoins(), ErrInsufficientShares
	}
	return utils.SharesFromTokens(tokens, s.GetSharesDenom, s.Tokens, s.DelegatorShares)
}

// SharesFromTokensTruncated returns the truncated shares of a delegation given a bond amount.
// It returns an error if the service has no tokens.
func (s Service) SharesFromTokensTruncated(tokens sdk.Coins) (sdk.DecCoins, error) {
	if s.Tokens.IsZero() {
		return sdk.NewDecCoins(), ErrInsufficientShares
	}
	return utils.SharesFromTokensTruncated(tokens, s.GetSharesDenom, s.Tokens, s.DelegatorShares)
}

// AddTokensFromDelegation adds the given amount of tokens to the service's total tokens,
// also updating the service's delegator shares.
// It returns the updated service and the shares issued.
func (s Service) AddTokensFromDelegation(amount sdk.Coins) (Service, sdk.DecCoins) {
	issuedShares := utils.IssueShares(amount, s.GetSharesDenom, s.Tokens, s.DelegatorShares)

	s.Tokens = s.Tokens.Add(amount...)
	s.DelegatorShares = s.DelegatorShares.Add(issuedShares...)

	return s, issuedShares
}

// RemoveDelShares removes delegator shares from a service.
// NOTE: Because token fractions are left in the service, the exchange rate of future shares
// of this validator can increase.
func (s Service) RemoveDelShares(delShares sdk.DecCoins) (Service, sdk.Coins) {
	remainingShares := s.DelegatorShares.Sub(delShares)

	var issuedTokens sdk.Coins
	if remainingShares.IsZero() {
		// Last delegation share gets any trimmings
		issuedTokens = s.Tokens
		s.Tokens = sdk.NewCoins()
	} else {
		// Leave excess tokens in the validator
		// However fully use all the delegator shares
		issuedTokens, _ = s.TokensFromShares(delShares).TruncateDecimal()
		s.Tokens = s.Tokens.Sub(issuedTokens...)

		if s.Tokens.IsAnyNegative() {
			panic("attempting to remove more tokens than available in service")
		}
	}

	s.DelegatorShares = remainingShares

	return s, issuedTokens
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
