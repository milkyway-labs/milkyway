package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PoolsHooks interface {
	AfterPoolCreated(ctx sdk.Context, poolID uint32) error
}
