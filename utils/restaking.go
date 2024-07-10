package utils

import (
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetSharesDenomFromTokenDenom returns the shares denom from the token denom.
// The returned shares denom will be in the format "{prefix}/{id}/{tokenDenom}".
func GetSharesDenomFromTokenDenom(prefix string, id uint32, tokenDenom string) string {
	return fmt.Sprintf("%s/%d/%s", prefix, id, tokenDenom)
}

// GetTokenDenomFromSharesDenom returns the token denom from the shares denom.
// It expects the shares denom to be in the format "{xxxxxx}/{xxxxxx}/{tokenDenom}".
func GetTokenDenomFromSharesDenom(sharesDenom string) string {
	parts := strings.Split(sharesDenom, "/")
	if len(parts) != 3 {
		return ""
	}
	return parts[2]
}

// IsInvalidExRate returns true if the delegated tokens are zero and the delegators shares are positive.
func IsInvalidExRate(delegatedTokens sdk.Coins, delegatorsShares sdk.DecCoins) bool {
	for _, token := range delegatedTokens {
		if token.IsZero() && delegatorsShares.AmountOf(token.Denom).IsPositive() {
			return true
		}
	}
	return false
}

// ComputeTokensFromShares calculates the token worth of provided shares.
func ComputeTokensFromShares(shares sdk.DecCoins, delegatedTokens sdk.Coins, delegatorShares sdk.DecCoins) sdk.DecCoins {
	tokens := sdk.NewDecCoins()
	for _, share := range shares {
		tokenDenom := GetTokenDenomFromSharesDenom(share.Denom)

		operatorTokenAmount := delegatedTokens.AmountOf(tokenDenom)
		delegatorSharesAmount := delegatorShares.AmountOf(share.Denom)

		tokenAmount := share.Amount.MulInt(operatorTokenAmount).Quo(delegatorSharesAmount)

		tokens = tokens.Add(sdk.NewDecCoinFromDec(tokenDenom, tokenAmount))
	}

	return tokens
}

// ShareDenomGetter represents a function that returns the shares denom given a token denom.
type ShareDenomGetter func(tokenDenom string) (shareDenom string)

// SharesFromTokens returns the shares of a delegation given a bond amount.
func SharesFromTokens(tokens sdk.Coin, getShareDenom ShareDenomGetter, delegatedTokens sdk.Coins, delegatorsShares sdk.DecCoins) (sdkmath.LegacyDec, error) {
	sharesDenom := getShareDenom(tokens.Denom)
	delegatorTokenShares := delegatorsShares.AmountOf(sharesDenom)
	operatorTokenAmount := delegatedTokens.AmountOf(tokens.Denom)
	return delegatorTokenShares.MulInt(tokens.Amount).QuoInt(operatorTokenAmount), nil
}

// IssueShares calculates the shares to issue for a delegation of the given amount.
func IssueShares(amount sdk.Coins, getShareDenom ShareDenomGetter, delegatedTokens sdk.Coins, delegatorsShares sdk.DecCoins) sdk.DecCoins {
	// calculate the shares to issue
	issuedShares := sdk.NewDecCoins()
	for _, token := range amount {
		var tokenShares sdk.DecCoin
		sharesDenom := getShareDenom(token.Denom)

		delegatorShares := delegatorsShares.AmountOf(sharesDenom)
		if delegatorShares.IsZero() {
			// The first delegation to an operator sets the exchange rate to one
			tokenShares = sdk.NewDecCoin(sharesDenom, token.Amount)
		} else {
			shares, err := SharesFromTokens(token, getShareDenom, delegatedTokens, delegatorsShares)
			if err != nil {
				panic(err)
			}
			tokenShares = sdk.NewDecCoinFromDec(sharesDenom, shares)
		}

		issuedShares = issuedShares.Add(tokenShares)
	}

	return issuedShares
}
