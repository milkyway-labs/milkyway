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
	parts := strings.SplitN(sharesDenom, "/", 3)
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

// ComputeTokensFromSharesTruncated calculates the token worth of provided shares, truncated.
func ComputeTokensFromSharesTruncated(shares sdk.DecCoins, delegatedTokens sdk.Coins, delegatorShares sdk.DecCoins) sdk.DecCoins {
	tokens := sdk.NewDecCoins()
	for _, share := range shares {
		tokenDenom := GetTokenDenomFromSharesDenom(share.Denom)

		operatorTokenAmount := delegatedTokens.AmountOf(tokenDenom)
		delegatorSharesAmount := delegatorShares.AmountOf(share.Denom)

		tokenAmount := share.Amount.MulInt(operatorTokenAmount).QuoTruncate(delegatorSharesAmount)

		tokens = tokens.Add(sdk.NewDecCoinFromDec(tokenDenom, tokenAmount))
	}

	return tokens
}

// ShareDenomGetter represents a function that returns the shares denom given a token denom.
type ShareDenomGetter func(tokenDenom string) (shareDenom string)

// SharesFromTokens returns the shares of a delegation given a bond amount.
func SharesFromTokens(tokens sdk.Coins, getShareDenom ShareDenomGetter, delegatedTokens sdk.Coins, delegatorsShares sdk.DecCoins) (sdk.DecCoins, error) {
	shares := sdk.NewDecCoins()
	for _, token := range tokens {
		sharesDenom := getShareDenom(token.Denom)

		operatorTokenAmount := delegatedTokens.AmountOf(token.Denom)

		var sharesAmount sdkmath.LegacyDec
		if operatorTokenAmount.IsZero() {
			sharesAmount = sdkmath.LegacyNewDecFromInt(token.Amount)
		} else {
			delegatorTokenShares := delegatorsShares.AmountOf(sharesDenom)
			sharesAmount = delegatorTokenShares.MulInt(token.Amount).QuoInt(operatorTokenAmount)
		}

		shares = shares.Add(sdk.NewDecCoinFromDec(sharesDenom, sharesAmount))
	}

	return shares, nil
}

// SharesFromTokensTruncated returns the truncated shares of a delegation given a bond amount.
func SharesFromTokensTruncated(tokens sdk.Coins, getShareDenom ShareDenomGetter, delegatedTokens sdk.Coins, delegatorsShares sdk.DecCoins) (sdk.DecCoins, error) {
	shares := sdk.NewDecCoins()
	for _, token := range tokens {
		sharesDenom := getShareDenom(token.Denom)

		delegatorTokenShares := delegatorsShares.AmountOf(sharesDenom)
		operatorTokenAmount := delegatedTokens.AmountOf(token.Denom)

		var sharesAmount sdkmath.LegacyDec
		if operatorTokenAmount.IsZero() {
			sharesAmount = sdkmath.LegacyNewDecFromInt(token.Amount)
		} else {
			sharesAmount = delegatorTokenShares.MulInt(token.Amount).QuoTruncate(sdkmath.LegacyNewDecFromInt(operatorTokenAmount))
		}

		shares = shares.Add(sdk.NewDecCoinFromDec(sharesDenom, sharesAmount))
	}

	return shares, nil
}

// SharesFromDecCoins returns the shares of a delegation given a bond amount.
func SharesFromDecCoins(tokens sdk.DecCoins, getShareDenom ShareDenomGetter, delegatedTokens sdk.Coins, delegatorsShares sdk.DecCoins) (sdk.DecCoins, error) {
	shares := sdk.NewDecCoins()
	for _, token := range tokens {
		sharesDenom := getShareDenom(token.Denom)

		delegatorTokenShares := delegatorsShares.AmountOf(sharesDenom)
		operatorTokenAmount := delegatedTokens.AmountOf(token.Denom)

		var sharesAmount sdkmath.LegacyDec
		if operatorTokenAmount.IsZero() {
			sharesAmount = token.Amount
		} else {
			sharesAmount = delegatorTokenShares.Mul(token.Amount).QuoTruncate(sdkmath.LegacyNewDecFromInt(operatorTokenAmount))
		}

		shares = shares.Add(sdk.NewDecCoinFromDec(sharesDenom, sharesAmount))
	}

	return shares, nil
}

// IssueShares calculates the shares to issue for a delegation of the given amount.
func IssueShares(amount sdk.Coins, getShareDenom ShareDenomGetter, delegatedTokens sdk.Coins, delegatorsShares sdk.DecCoins) sdk.DecCoins {
	issuedShares := sdk.NewDecCoins()
	if delegatorsShares.IsZero() {
		// The first delegation to an operator sets the exchange rate to one
		for _, token := range amount {
			issuedShares = issuedShares.Add(sdk.NewDecCoin(getShareDenom(token.Denom), token.Amount))
		}
	} else {
		shares, err := SharesFromTokens(amount, getShareDenom, delegatedTokens, delegatorsShares)
		if err != nil {
			panic(err)
		}
		issuedShares = shares
	}

	return issuedShares
}
