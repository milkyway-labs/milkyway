package types

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"

	assetstypes "github.com/milkyway-labs/milkyway/v5/x/assets/types"
	operatorstypes "github.com/milkyway-labs/milkyway/v5/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v5/x/pools/types"
	servicestypes "github.com/milkyway-labs/milkyway/v5/x/services/types"
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
	GetPoolByDenom(ctx context.Context, denom string) (poolstypes.Pool, bool, error)
	CreateOrGetPoolByDenom(ctx context.Context, denom string) (poolstypes.Pool, error)
	GetPool(ctx context.Context, poolID uint32) (poolstypes.Pool, error)
	SavePool(ctx context.Context, pool poolstypes.Pool) error
	IteratePools(ctx context.Context, cb func(poolstypes.Pool) (bool, error)) error
	GetPools(ctx context.Context) ([]poolstypes.Pool, error)
}

type OperatorsKeeper interface {
	GetOperator(ctx context.Context, operatorID uint32) (operatorstypes.Operator, error)
	SaveOperator(ctx context.Context, operator operatorstypes.Operator) error
	IterateOperators(ctx context.Context, cb func(operatorstypes.Operator) (bool, error)) error
	GetOperators(ctx context.Context) ([]operatorstypes.Operator, error)
	SaveOperatorParams(ctx context.Context, operatorID uint32, params operatorstypes.OperatorParams) error
	GetOperatorParams(ctx context.Context, operatorID uint32) (operatorstypes.OperatorParams, error)
}

type ServicesKeeper interface {
	HasService(ctx context.Context, serviceID uint32) (bool, error)
	GetService(ctx context.Context, serviceID uint32) (servicestypes.Service, error)
	SaveService(ctx context.Context, service servicestypes.Service) error
	IterateServices(ctx context.Context, cb func(servicestypes.Service) (bool, error)) error
	GetServices(ctx context.Context) ([]servicestypes.Service, error)
	DeactivateService(ctx context.Context, serviceID uint32) error
	GetServiceParams(ctx context.Context, serviceID uint32) (servicestypes.ServiceParams, error)
}

type OracleKeeper interface {
	GetPriceWithNonceForCurrencyPair(ctx sdk.Context, cp connecttypes.CurrencyPair) (oracletypes.QuotePriceWithNonce, error)
	GetDecimalsForCurrencyPair(ctx sdk.Context, cp connecttypes.CurrencyPair) (decimals uint64, err error)
}

type AssetsKeeper interface {
	GetAsset(ctx context.Context, denom string) (assetstypes.Asset, error)
}
