package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v9/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/v9/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v9/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v9/x/services/types"
)

// CanCoverDecCoins returns whether the insurance fund can cover the provided dec
// coins. It also returns the amount of coins required to cover the dec coins.
func CanCoverDecCoins(insuranceFund sdk.Coins, insurancePercentage sdkmath.LegacyDec, coins sdk.DecCoins) (bool, sdk.Coins) {
	required := sdk.NewCoins()
	for _, coin := range coins {
		if IsLockedRepresentationDenom(coin.Denom) {
			nativeDenom, err := LockedDenomToNative(coin.Denom)
			if err != nil {
				// There must be no error since we already checked the denom
				panic(err)
			}
			required = required.Add(sdk.NewCoin(nativeDenom, insurancePercentage.Mul(coin.Amount).QuoInt64(100).Ceil().TruncateInt()))
		}
	}

	return insuranceFund.IsAllGTE(required), required
}

// GetCoverableDecCoins returns the amount of dec coins that can be covered by the
// insurance fund. It ignores any invalid denoms in the insurance fund.
func GetCoverableDecCoins(insuranceFund sdk.Coins, insurancePercentage sdkmath.LegacyDec) sdk.DecCoins {
	coverable := sdk.NewDecCoins()
	for _, coin := range insuranceFund {
		lockedDenom, err := GetLockedRepresentationDenom(coin.Denom)
		if err != nil {
			continue
		}
		coverable = coverable.Add(sdk.NewDecCoinFromDec(
			lockedDenom,
			coin.Amount.ToLegacyDec().QuoTruncate(insurancePercentage).MulInt64(100),
		))
	}
	return coverable
}

// GetCoveredLockedShares returns the locked shares that are covered by the
// insurance fund.
func GetCoveredLockedShares(
	target restakingtypes.DelegationTarget,
	delegation restakingtypes.Delegation,
	insuranceFund sdk.Coins,
	insurancePercentage sdkmath.LegacyDec,
	activeLockedTokens sdk.DecCoins,
) (sdk.DecCoins, error) {
	// Exit early if the user doesn't have insurance fund balance
	if insuranceFund.IsZero() {
		return nil, nil
	}

	delegationTokens := target.TokensFromSharesTruncated(delegation.Shares)
	if _, hasNeg := activeLockedTokens.SafeSub(delegationTokens); hasNeg { // sanity check
		panic(fmt.Sprintf("delegation tokens %s > active locked tokens %s", delegationTokens, activeLockedTokens))
	}

	coveredTokens := sdk.NewDecCoins()
	coverableTokens := GetCoverableDecCoins(insuranceFund, insurancePercentage)
	for _, coverableToken := range coverableTokens {
		delegationAmount := delegationTokens.AmountOf(coverableToken.Denom)
		usedAmount := activeLockedTokens.AmountOf(coverableToken.Denom)
		coveredTokens = coveredTokens.Add(sdk.NewDecCoinFromDec(
			coverableToken.Denom,
			sdkmath.LegacyMinDec(
				delegationAmount,
				coverableToken.Amount.MulTruncate(delegationAmount.QuoTruncate(usedAmount)),
			),
		))
	}

	// Convert tokens back to shares
	coveredShares, err := target.SharesFromDecCoins(coveredTokens)
	if err != nil {
		return nil, err
	}
	// Truncate the shares to make the numbers to avoid unnecessary rounding errors
	// in calculations
	truncatedShares, _ := coveredShares.TruncateDecimal()
	return sdk.NewDecCoinsFromCoins(truncatedShares...), nil
}

// UncoveredLockedShares returns the locked shares that are not covered by the
// insurance fund.
func UncoveredLockedShares(shares, coveredLockedShares sdk.DecCoins) sdk.DecCoins {
	res := sdk.NewDecCoins()
	for _, share := range shares {
		tokenDenom := utils.GetTokenDenomFromSharesDenom(share.Denom)
		if tokenDenom == "" || !IsLockedRepresentationDenom(tokenDenom) {
			continue
		}
		coveredAmount := coveredLockedShares.AmountOf(share.Denom)
		res = res.Add(sdk.NewDecCoinFromDec(share.Denom, share.Amount.Sub(coveredAmount)))
	}
	return res
}

// HasLockedShares returns whether the provided shares contain any locked shares.
func HasLockedShares(shares sdk.DecCoins) bool {
	for _, share := range shares {
		tokenDenom := utils.GetTokenDenomFromSharesDenom(share.Denom)
		if tokenDenom == "" {
			continue
		}
		if IsLockedRepresentationDenom(tokenDenom) {
			return true
		}
	}
	return false
}

// DelegationTargetWithDeductedShares returns the delegation target with the
// provided shares removed. It doesn't use DelegationTarget's RemoveDelShares
// to avoid leaving excess tokens in the target.
func DelegationTargetWithDeductedShares(target restakingtypes.DelegationTarget, shares sdk.DecCoins) (restakingtypes.DelegationTarget, error) {
	remainingShares := target.GetDelegatorShares().Sub(shares)
	switch target := target.(type) {
	case poolstypes.Pool:
		if remainingShares.IsZero() {
			target.Tokens = sdkmath.ZeroInt()
		} else {
			tokens, _ := target.TokensFromSharesTruncated(remainingShares).TruncateDecimal()
			target.Tokens = tokens.AmountOf(target.Denom)
		}
		target.DelegatorShares = remainingShares.AmountOf(target.GetSharesDenom(target.Denom))
		return target, nil
	case operatorstypes.Operator:
		if remainingShares.IsZero() {
			target.Tokens = sdk.NewCoins()
		} else {
			target.Tokens, _ = target.TokensFromSharesTruncated(remainingShares).TruncateDecimal()
		}
		target.DelegatorShares = remainingShares
		return target, nil
	case servicestypes.Service:
		if remainingShares.IsZero() {
			target.Tokens = sdk.NewCoins()
		} else {
			target.Tokens, _ = target.TokensFromSharesTruncated(remainingShares).TruncateDecimal()
		}
		target.DelegatorShares = remainingShares
		return target, nil
	default:
		return nil, restakingtypes.ErrInvalidDelegationType.Wrapf("invalid target type %T", target)
	}
}
