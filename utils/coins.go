package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DenomsSubsetOf returns true if receiver's denom set
// is subset of coinsB's denoms.
func DenomsSubsetOf(coins sdk.DecCoins, coinsB sdk.DecCoins) bool {
	// more denoms in B than in receiver
	if len(coins) > len(coinsB) {
		return false
	}

	for _, coin := range coins {
		if coinsB.AmountOf(coin.Denom).IsZero() {
			return false
		}
	}

	return true
}

// IsAllGT returns true if for every denom in coinsB,
// the denom is present at a greater amount in coins.
func IsAllGT(coins, coinsB sdk.DecCoins) bool {
	if len(coins) == 0 {
		return false
	}

	if len(coinsB) == 0 {
		return true
	}

	if !DenomsSubsetOf(coinsB, coins) {
		return false
	}

	for _, coinB := range coinsB {
		amountA, amountB := coins.AmountOf(coinB.Denom), coinB.Amount
		if !amountA.GT(amountB) {
			return false
		}
	}

	return true
}

// IsAllGTE returns false if for any denom in coinsB,
// the denom is present at a smaller amount in coins;
// else returns true.
func IsAllGTE(coins, coinsB sdk.DecCoins) bool {
	if len(coinsB) == 0 {
		return true
	}

	if len(coins) == 0 {
		return false
	}

	for _, coinB := range coinsB {
		if coinB.Amount.GT(coins.AmountOf(coinB.Denom)) {
			return false
		}
	}

	return true
}

// IsAnyGT returns true iff for any denom in coins, the denom is present at a
// greater amount in coinsB.
//
// e.g.
// IsAnyGT({2A, 3B}, {A}) = true
// IsAnyGT({2A, 3B}, {5C}) = false
// IsAnyGT({}, {5C}) = false
// IsAnyGT({2A, 3B}, {}) = false
func IsAnyGT(coins, coinsB sdk.DecCoins) bool {
	if len(coinsB) == 0 {
		return false
	}

	for _, coin := range coins {
		amt := coinsB.AmountOf(coin.Denom)
		if coin.Amount.GT(amt) && !amt.IsZero() {
			return true
		}
	}

	return false
}

// IsAnyLT returns true iff for any denom in coins, the denom is present at a
// smaller amount in coinsB.
func IsAnyLT(coins, coinsB sdk.DecCoins) bool {
	return !IsAllGTE(coins, coinsB)
}

// IntersectCoinsByDenom returns the intersection of two coins.
// e.g.
// IntersectCoinsByDenom({2A, 3B}, {A}) = {2A}
// IntersectCoinsByDenom({2A, 3B}, {5C}) = {}
// IntersectCoinsByDenom({2A, 3B}, {A, B}) = {2A, 3B}
func IntersectCoinsByDenom(coins, coinsB sdk.Coins) sdk.Coins {
	res := sdk.NewCoins()
	for _, coin := range coins {
		amount := coinsB.AmountOf(coin.Denom)
		if !amount.IsZero() {
			res = res.Add(coin)
		}
	}
	return res
}

// IntersectDecCoinsByDenom returns the intersection of two coins.
// e.g.
// IntersectDecCoinsByDenom{2A, 3B}, {A}) = {2A}
// IntersectDecCoinsByDenom({2A, 3B}, {5C}) = {}
// IntersectDecCoinsByDenom({2A, 3B}, {A, B}) = {2A, 3B}
func IntersectDecCoinsByDenom(coins, coinsB sdk.DecCoins) sdk.DecCoins {
	res := sdk.NewDecCoins()
	for _, coin := range coins {
		amount := coinsB.AmountOf(coin.Denom)
		if !amount.IsZero() {
			res = res.Add(coin)
		}
	}
	return res
}
