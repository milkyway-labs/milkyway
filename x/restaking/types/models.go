package types

import (
	sdkmath "cosmossdk.io/math"
)

func NewPoolDelegation(poolID uint32, userAddress string, shares sdkmath.LegacyDec) PoolDelegation {
	return PoolDelegation{
		PoolID:      poolID,
		UserAddress: userAddress,
		Shares:      shares,
	}
}
