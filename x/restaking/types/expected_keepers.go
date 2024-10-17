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
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
}

type PoolsKeeper interface {
	GetPoolByDenom(ctx sdk.Context, denom string) (poolstypes.Pool, bool)
	CreateOrGetPoolByDenom(ctx sdk.Context, denom string) (poolstypes.Pool, error)
	GetPool(ctx sdk.Context, poolID uint32) (poolstypes.Pool, bool)
	SavePool(ctx sdk.Context, pool poolstypes.Pool) error
	IteratePools(ctx sdk.Context, cb func(poolstypes.Pool) bool)
	GetPools(ctx sdk.Context) []poolstypes.Pool
}

type OperatorsKeeper interface {
	GetOperator(ctx sdk.Context, operatorID uint32) (operatorstypes.Operator, bool)
	SaveOperator(ctx sdk.Context, operator operatorstypes.Operator) error
	IterateOperators(ctx sdk.Context, cb func(operatorstypes.Operator) bool)
	GetOperators(ctx sdk.Context) []operatorstypes.Operator
	SaveOperatorParams(ctx sdk.Context, operatorID uint32, params operatorstypes.OperatorParams) error
	GetOperatorParams(ctx sdk.Context, operatorID uint32) (operatorstypes.OperatorParams, error)
}

type ServicesKeeper interface {
	GetService(ctx sdk.Context, serviceID uint32) (servicestypes.Service, bool)
	SaveService(ctx sdk.Context, service servicestypes.Service) error
	IterateServices(ctx sdk.Context, cb func(servicestypes.Service) bool)
	GetServices(ctx sdk.Context) []servicestypes.Service
}
