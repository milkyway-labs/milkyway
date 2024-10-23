package types

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"

	assetstypes "github.com/milkyway-labs/milkyway/x/assets/types"
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
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
	SetModuleAccount(ctx context.Context, macc sdk.ModuleAccountI)
}

type BankKeeper interface {
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, moduleName string, addr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, addr sdk.AccAddress, moduleName string, amt sdk.Coins) error
	BlockedAddr(addr sdk.AccAddress) bool
}

type CommunityPoolKeeper interface {
	FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error
}

type OracleKeeper interface {
	GetPriceWithNonceForCurrencyPair(ctx sdk.Context, cp connecttypes.CurrencyPair) (oracletypes.QuotePriceWithNonce, error)
	GetDecimalsForCurrencyPair(ctx sdk.Context, cp connecttypes.CurrencyPair) (decimals uint64, err error)
}

type PoolsKeeper interface {
	GetParams(ctx sdk.Context) poolstypes.Params
	GetPool(ctx sdk.Context, poolID uint32) (poolstypes.Pool, bool)
	GetPools(ctx sdk.Context) []poolstypes.Pool
	IteratePools(ctx sdk.Context, cb func(pool poolstypes.Pool) (stop bool))
}

type OperatorsKeeper interface {
	GetOperator(ctx sdk.Context, operatorID uint32) (operatorstypes.Operator, bool)
	GetOperators(ctx sdk.Context) []operatorstypes.Operator
	IterateOperators(ctx sdk.Context, cb func(operator operatorstypes.Operator) (stop bool))
	GetOperatorParams(ctx sdk.Context, operatorID uint32) (operatorstypes.OperatorParams, error)
}

type ServicesKeeper interface {
	GetService(ctx sdk.Context, serviceID uint32) (servicestypes.Service, bool)
	IterateServices(ctx sdk.Context, cb func(service servicestypes.Service) (stop bool))
}

type RestakingKeeper interface {
	GetOperatorJoinedServices(ctx sdk.Context, operatorID uint32) (restakingtypes.OperatorJoinedServices, error)
	CanOperatorValidateService(ctx sdk.Context, serviceID uint32, operatorID uint32) (bool, error)
	IsServiceSecuredByPool(ctx sdk.Context, serviceID uint32, operatorID uint32) (bool, error)
	GetPoolDelegation(ctx sdk.Context, poolID uint32, userAddress string) (restakingtypes.Delegation, bool)
	GetOperatorDelegation(ctx sdk.Context, operatorID uint32, userAddress string) (restakingtypes.Delegation, bool)
	GetServiceDelegation(ctx sdk.Context, serviceID uint32, userAddress string) (restakingtypes.Delegation, bool)
	GetDelegationForTarget(ctx sdk.Context, target restakingtypes.DelegationTarget, delegator string) (restakingtypes.Delegation, bool)
	IterateUserPoolDelegations(ctx sdk.Context, userAddress string, cb func(del restakingtypes.Delegation) (stop bool, err error)) error
	IterateUserOperatorDelegations(
		ctx sdk.Context, userAddress string, cb func(del restakingtypes.Delegation) (stop bool, err error)) error
	IterateUserServiceDelegations(
		ctx sdk.Context, userAddress string, cb func(del restakingtypes.Delegation) (stop bool, err error)) error
	IterateAllPoolDelegations(ctx sdk.Context, cb func(del restakingtypes.Delegation) (stop bool))
	IterateAllOperatorDelegations(ctx sdk.Context, cb func(del restakingtypes.Delegation) (stop bool))
	IterateAllServiceDelegations(ctx sdk.Context, cb func(del restakingtypes.Delegation) (stop bool))
}

type AssetsKeeper interface {
	GetAsset(ctx context.Context, denom string) (assetstypes.Asset, error)
}
