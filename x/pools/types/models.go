package types

import (
	"fmt"
	"strconv"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// GetPoolAddress generates a pool address from its id
func GetPoolAddress(poolID uint32) sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("pool-%d", poolID))
}

// ParsePoolID parses a pool id from a string
func ParsePoolID(value string) (uint32, error) {
	parsed, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid pool id: %s", value)
	}
	return uint32(parsed), nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewPool creates a new Pool instance
func NewPool(id uint32, denom string) Pool {
	return Pool{
		ID:              id,
		Denom:           denom,
		Address:         GetPoolAddress(id).String(),
		Tokens:          sdkmath.ZeroInt(),
		DelegatorShares: sdkmath.LegacyZeroDec(),
	}
}

// Validate checks if the pool is valid
func (p Pool) Validate() error {
	if p.ID == 0 {
		return fmt.Errorf("invalid pool id")
	}

	if sdk.ValidateDenom(p.Denom) != nil {
		return fmt.Errorf("invalid pool denom")
	}

	_, err := sdk.AccAddressFromBech32(p.Address)
	if err != nil {
		return fmt.Errorf("invalid pool address")
	}

	return nil
}

// GetSharesDenom returns the shares denom for a pool and token denom
func (p Pool) GetSharesDenom(tokenDenom string) string {
	return fmt.Sprintf("pool/%d/%s", p.ID, tokenDenom)
}

// InvalidExRate returns whether the exchange rates is invalid.
// This can happen e.g. if Pool loses all tokens due to slashing. In this case,
// make all future delegations invalid.
func (p Pool) InvalidExRate() bool {
	return p.Tokens.IsZero() && p.DelegatorShares.IsPositive()
}

// TokensFromShares calculates the token worth of provided shares
func (p Pool) TokensFromShares(shares sdkmath.LegacyDec) sdkmath.LegacyDec {
	return (shares.MulInt(p.Tokens)).Quo(p.DelegatorShares)
}

// SharesFromTokens returns the shares of a delegation given a bond amount. It
// returns an error if the pool has no tokens.
func (p Pool) SharesFromTokens(amt sdkmath.Int) (sdkmath.LegacyDec, error) {
	if p.Tokens.IsZero() {
		return sdkmath.LegacyZeroDec(), ErrInsufficientShares
	}

	return p.DelegatorShares.MulInt(amt).QuoInt(p.Tokens), nil
}

// AddTokensFromDelegation adds the given amount of tokens to the pool's total tokens,
// also updating the pool's delegator shares.
// It returns the updated pool and the shares issued.
func (p Pool) AddTokensFromDelegation(amount sdkmath.Int) (Pool, sdkmath.LegacyDec) {
	// calculate the shares to issue
	var issuedShares sdkmath.LegacyDec
	if p.DelegatorShares.IsZero() {
		// the first delegation to a validator sets the exchange rate to one
		issuedShares = sdkmath.LegacyNewDecFromInt(amount)
	} else {
		shares, err := p.SharesFromTokens(amount)
		if err != nil {
			panic(err)
		}

		issuedShares = shares
	}

	p.Tokens = p.Tokens.Add(amount)
	p.DelegatorShares = p.DelegatorShares.Add(issuedShares)

	return p, issuedShares
}
