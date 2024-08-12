package types

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

type AccountKeeper interface {
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	HasAccount(ctx context.Context, addr sdk.AccAddress) bool
	SetAccount(ctx context.Context, acc sdk.AccountI)
	AddressCodec() address.Codec
	GetModuleAddress(moduleName string) sdk.AccAddress
}

type BankKeeper interface {
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	BlockedAddr(addr sdk.AccAddress) bool
}

type CommunityPoolKeeper interface {
	FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error
}

type OracleKeeper interface {
	GetPriceWithNonceForCurrencyPair(ctx sdk.Context, cp slinkytypes.CurrencyPair) (oracletypes.QuotePriceWithNonce, error)
	GetDecimalsForCurrencyPair(ctx sdk.Context, cp slinkytypes.CurrencyPair) (decimals uint64, err error)
}

type PoolsKeeper interface {
	GetParams(ctx sdk.Context) poolstypes.Params
	GetPool(ctx sdk.Context, poolID uint32) (poolstypes.Pool, bool)
	GetPools(ctx sdk.Context) []poolstypes.Pool
}

type OperatorsKeeper interface {
	GetOperator(ctx sdk.Context, operatorID uint32) (operatorstypes.Operator, bool)
	IterateOperators(ctx sdk.Context, cb func(operator operatorstypes.Operator) (stop bool))
}

type ServicesKeeper interface {
	GetService(ctx sdk.Context, serviceID uint32) (servicestypes.Service, bool)
}

type RestakingKeeper interface {
	GetOperatorParams(ctx sdk.Context, operatorID uint32) restakingtypes.OperatorParams
	GetServiceParams(ctx sdk.Context, serviceID uint32) restakingtypes.ServiceParams
	GetPoolDelegation(ctx sdk.Context, poolID uint32, userAddress string) (restakingtypes.Delegation, bool)
	GetOperatorDelegation(ctx sdk.Context, operatorID uint32, userAddress string) (restakingtypes.Delegation, bool)
	GetServiceDelegation(ctx sdk.Context, serviceID uint32, userAddress string) (restakingtypes.Delegation, bool)
	IterateUserPoolDelegations(ctx sdk.Context, userAddress string, cb func(del restakingtypes.Delegation) (stop bool, err error)) error
	IterateUserOperatorDelegations(
		ctx sdk.Context, userAddress string, cb func(del restakingtypes.Delegation) (stop bool, err error)) error
	IterateUserServiceDelegations(
		ctx sdk.Context, userAddress string, cb func(del restakingtypes.Delegation) (stop bool, err error)) error
}

type TickersKeeper interface {
	GetTicker(ctx context.Context, denom string) (ticker string, err error)
}
