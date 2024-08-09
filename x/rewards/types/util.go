package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

type DelegationTarget struct {
	restakingtypes.DelegationTarget
}

func NewDelegationTarget(target restakingtypes.DelegationTarget) *DelegationTarget {
	return &DelegationTarget{target}
}

func (target DelegationTarget) Type() restakingtypes.DelegationType {
	switch target.DelegationTarget.(type) {
	case *poolstypes.Pool:
		return restakingtypes.DELEGATION_TYPE_POOL
	case *operatorstypes.Operator:
		return restakingtypes.DELEGATION_TYPE_OPERATOR
	case *servicestypes.Service:
		return restakingtypes.DELEGATION_TYPE_SERVICE
	default:
		panic("unknown delegation target type")
	}
}

func (target DelegationTarget) Tokens() sdk.Coins {
	switch target := target.DelegationTarget.(type) {
	case *poolstypes.Pool:
		return sdk.NewCoins(sdk.NewCoin(target.Denom, target.Tokens))
	case *operatorstypes.Operator:
		return target.Tokens
	case *servicestypes.Service:
		return target.Tokens
	default:
		panic("unknown delegation target type")
	}
}

func (target DelegationTarget) DelegatorShares() sdk.DecCoins {
	switch target := target.DelegationTarget.(type) {
	case *poolstypes.Pool:
		sharesDenom := target.GetSharesDenom(target.Denom)
		return sdk.NewDecCoins(sdk.NewDecCoinFromDec(sharesDenom, target.DelegatorShares))
	case *operatorstypes.Operator:
		return target.DelegatorShares
	case *servicestypes.Service:
		return target.DelegatorShares
	default:
		panic("unknown delegation target type")
	}
}

// TokensFromShares calculates the token worth of provided shares
func (target DelegationTarget) TokensFromShares(shares sdk.DecCoins) sdk.DecCoins {
	tokens := target.Tokens()
	delShares := target.DelegatorShares()
	return utils.ComputeTokensFromShares(shares, tokens, delShares)
}

// TokensFromSharesTruncated calculates the token worth of provided shares, truncated
func (target DelegationTarget) TokensFromSharesTruncated(shares sdk.DecCoins) sdk.DecCoins {
	tokens := target.Tokens()
	delShares := target.DelegatorShares()
	return utils.ComputeTokensFromSharesTruncated(shares, tokens, delShares)
}
