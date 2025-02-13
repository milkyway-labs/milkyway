package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v9/utils"
	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
)

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

func GetCoveredLockedShares(target restakingtypes.DelegationTarget, delegation restakingtypes.Delegation, insuranceFund sdk.Coins, insurancePercentage sdkmath.LegacyDec) (sdk.DecCoins, error) {
	// Exit early if the user doesn't have insurance fund balance
	if insuranceFund.IsZero() {
		return nil, nil
	}

	coverable := GetCoverableDecCoins(insuranceFund, insurancePercentage)

	// Calculate covered locked shares
	tokens := target.TokensFromSharesTruncated(delegation.Shares)
	coveredTokens := tokens.Intersect(coverable)
	return target.SharesFromDecCoins(coveredTokens)
}

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

func DeductUncoveredLockedShares(shares, coveredLockedShares sdk.DecCoins) sdk.DecCoins {
	uncovered := UncoveredLockedShares(shares, coveredLockedShares)
	return shares.Sub(uncovered)
}
