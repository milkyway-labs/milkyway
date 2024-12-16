package types

import (
	"fmt"
	"strconv"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/milkyway-labs/milkyway/v7/utils"
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

func (p Pool) GetID() uint32 {
	return p.ID
}

func (p Pool) GetDenom() string {
	return p.Denom
}

func (p Pool) GetAddress() string {
	return p.Address
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

// GetTokens returns the pool's tokens as sdk.Coins.
func (p Pool) GetTokens() sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin(p.Denom, p.Tokens))
}

// GetDelegatorShares returns delegator shares of the pool as sdk.DecCoins.
func (p Pool) GetDelegatorShares() sdk.DecCoins {
	return sdk.NewDecCoins(sdk.NewDecCoinFromDec(p.GetSharesDenom(p.Denom), p.DelegatorShares))
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
func (p Pool) TokensFromShares(shares sdk.DecCoins) sdk.DecCoins {
	return utils.ComputeTokensFromShares(shares, p.GetTokens(), p.GetDelegatorShares())
}

// TokensFromSharesTruncated calculates the token worth of provided shares, truncated
func (p Pool) TokensFromSharesTruncated(shares sdk.DecCoins) sdk.DecCoins {
	return utils.ComputeTokensFromSharesTruncated(shares, p.GetTokens(), p.GetDelegatorShares())
}

// SharesFromTokens returns the shares of a delegation given a bond amount. It
// returns an error if the pool has no tokens.
func (p Pool) SharesFromTokens(amt sdk.Coins) (sdk.DecCoins, error) {
	if p.Tokens.IsZero() {
		return sdk.NewDecCoins(), ErrInsufficientShares
	}

	shares := sdk.NewDecCoins()
	for _, coin := range amt {
		if coin.Denom != p.Denom {
			return sdk.NewDecCoins(), ErrInvalidDenom
		}

		shareDenom := p.GetSharesDenom(coin.Denom)
		shareAmount := p.DelegatorShares.MulInt(coin.Amount).QuoInt(p.Tokens)

		shares = shares.Add(sdk.NewDecCoinFromDec(shareDenom, shareAmount))
	}

	return shares, nil
}

// SharesFromTokensTruncated returns the truncated shares of a delegation given
// a bond amount. It returns an error if the pool has no tokens.
func (p Pool) SharesFromTokensTruncated(amt sdk.Coins) (sdk.DecCoins, error) {
	if p.Tokens.IsZero() {
		return sdk.NewDecCoins(), ErrInsufficientShares
	}

	shares := sdk.NewDecCoins()
	for _, coin := range amt {
		if coin.Denom != p.Denom {
			return sdk.NewDecCoins(), ErrInvalidDenom
		}

		shareDenom := p.GetSharesDenom(coin.Denom)
		shareAmount := p.DelegatorShares.MulInt(coin.Amount).QuoTruncate(sdkmath.LegacyNewDecFromInt(p.Tokens))

		shares = shares.Add(sdk.NewDecCoinFromDec(shareDenom, shareAmount))
	}

	return shares, nil
}

// SharesFromDecCoins returns the shares of a delegation given a bond amount. It
// returns an error if the service has no tokens.
func (p Pool) SharesFromDecCoins(coins sdk.DecCoins) (sdk.DecCoins, error) {
	if p.Tokens.IsZero() {
		return sdk.NewDecCoins(), ErrInsufficientShares
	}

	shares := sdk.NewDecCoins()
	for _, coin := range coins {
		if coin.Denom != p.Denom {
			return sdk.NewDecCoins(), ErrInvalidDenom
		}

		shareDenom := p.GetSharesDenom(coin.Denom)
		shareAmount := p.DelegatorShares.Mul(coin.Amount).QuoInt(p.Tokens)

		shares = shares.Add(sdk.NewDecCoinFromDec(shareDenom, shareAmount))
	}

	return shares, nil
}

// AddTokensFromDelegation adds the given amount of tokens to the pool's total tokens,
// also updating the pool's delegator shares.
// It returns the updated pool and the shares issued.
func (p Pool) AddTokensFromDelegation(amount sdk.Coin) (Pool, sdk.DecCoin, error) {
	// calculate the shares to issue
	var issuedShares sdk.DecCoins
	if p.DelegatorShares.IsZero() {
		// the first delegation to a validator sets the exchange rate to one
		issuedShares = sdk.NewDecCoinsFromCoins(sdk.NewCoin(p.GetSharesDenom(amount.Denom), amount.Amount))
	} else {
		shares, err := p.SharesFromTokens(sdk.NewCoins(amount))
		if err != nil {
			return p, sdk.DecCoin{}, err
		}
		issuedShares = shares
	}

	p.Tokens = p.Tokens.Add(amount.Amount)
	p.DelegatorShares = p.DelegatorShares.Add(issuedShares[0].Amount)

	return p, issuedShares[0], nil
}

// RemoveDelShares removes delegator shares from an operator and returns
// the amount of tokens that should be issued for those shares.
// NOTE: Because token fractions are left in the operator,
// the exchange rate of future shares of this validator can increase.
func (p Pool) RemoveDelShares(shares sdk.DecCoins) (Pool, sdk.Coins, error) {
	if len(shares) > 1 {
		return p, sdk.Coins{}, ErrInvalidShares
	}

	for _, share := range shares {
		if share.Denom != p.GetSharesDenom(p.Denom) {
			return p, sdk.Coins{}, ErrInvalidShares
		}
	}

	delShares := shares.AmountOf(p.GetSharesDenom(p.Denom))
	remainingShares := p.DelegatorShares.Sub(delShares)

	var issuedTokens sdk.Coins
	if remainingShares.IsZero() {
		// Last delegation share gets any trimmings
		issuedTokens = sdk.NewCoins(sdk.NewCoin(p.Denom, p.Tokens))
		p.Tokens = sdkmath.ZeroInt()
	} else {
		// Leave excess tokens in the operator
		// However, fully use all the delegator shares
		issuedTokens, _ = p.TokensFromShares(shares).TruncateDecimal()
		p.Tokens = p.Tokens.Sub(issuedTokens.AmountOf(p.Denom))

		if p.Tokens.IsNegative() {
			panic("attempting to remove more tokens than available in operator")
		}
	}

	p.DelegatorShares = remainingShares

	return p, issuedTokens, nil
}
