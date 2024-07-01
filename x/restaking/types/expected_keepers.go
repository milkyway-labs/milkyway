package types

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

type AccountKeeper interface {
	AddressCodec() address.Codec
}

type BankKeeper interface {
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
}

type PoolsKeeper interface {
	CreateOrGetPoolByDenom(ctx sdk.Context, denom string) (poolstypes.Pool, error)
	SavePool(ctx sdk.Context, pool poolstypes.Pool) error
}

type ServicesKeeper interface {
	GetService(ctx sdk.Context, serviceID uint32) (servicestypes.Service, bool)
}

type OperatorsKeeper interface {
	GetOperator(ctx sdk.Context, operatorID uint32) (operatorstypes.Operator, bool)
	SaveOperator(ctx sdk.Context, operator operatorstypes.Operator)
}
