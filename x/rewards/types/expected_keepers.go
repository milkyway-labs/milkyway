package types

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"

	assetstypes "github.com/milkyway-labs/milkyway/v10/x/assets/types"
	operatorstypes "github.com/milkyway-labs/milkyway/v10/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v10/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v10/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v10/x/services/types"
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
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
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
	GetPool(ctx context.Context, poolID uint32) (poolstypes.Pool, error)
	GetPools(ctx context.Context) ([]poolstypes.Pool, error)
	IteratePools(ctx context.Context, cb func(pool poolstypes.Pool) (stop bool, err error)) error
}

type OperatorsKeeper interface {
	GetOperator(ctx context.Context, operatorID uint32) (operatorstypes.Operator, error)
	IterateOperators(ctx context.Context, cb func(operator operatorstypes.Operator) (stop bool, err error)) error
	GetOperatorParams(ctx context.Context, operatorID uint32) (operatorstypes.OperatorParams, error)
}

type ServicesKeeper interface {
	GetService(ctx context.Context, serviceID uint32) (servicestypes.Service, error)
	GetServiceParams(ctx context.Context, serviceID uint32) (servicestypes.ServiceParams, error)
	GetServices(ctx context.Context) ([]servicestypes.Service, error)
	IterateServices(ctx context.Context, cb func(service servicestypes.Service) (stop bool, err error)) error
}

type RestakingKeeper interface {
	IterateServiceValidatingOperators(ctx context.Context, serviceID uint32, action func(operatorID uint32) (stop bool, err error)) error
	CanOperatorValidateService(ctx context.Context, serviceID uint32, operatorID uint32) (bool, error)
	IsServiceSecuredByPool(ctx context.Context, serviceID uint32, poolID uint32) (bool, error)
	GetRestakableDenoms(ctx context.Context) ([]string, error)
	GetPoolDelegation(ctx context.Context, poolID uint32, userAddress string) (restakingtypes.Delegation, bool, error)
	GetDelegationForTarget(ctx context.Context, target restakingtypes.DelegationTarget, delegator string) (restakingtypes.Delegation, bool, error)

	IterateUserPoolDelegations(ctx context.Context, userAddress string, cb func(del restakingtypes.Delegation) (stop bool, err error)) error
	IterateAllPoolDelegations(ctx context.Context, cb func(del restakingtypes.Delegation) (stop bool, err error)) error

	IterateUserOperatorDelegations(ctx context.Context, userAddress string, cb func(del restakingtypes.Delegation) (stop bool, err error)) error
	IterateAllOperatorDelegations(ctx context.Context, cb func(del restakingtypes.Delegation) (stop bool, err error)) error

	IterateUserServiceDelegations(ctx context.Context, userAddress string, cb func(del restakingtypes.Delegation) (stop bool, err error)) error
	IterateServiceDelegations(ctx context.Context, serviceID uint32, cb func(del restakingtypes.Delegation) (stop bool, err error)) error
	IterateAllServiceDelegations(ctx context.Context, cb func(del restakingtypes.Delegation) (stop bool, err error)) error

	GetUserPreferences(ctx context.Context, userAddress string) (restakingtypes.UserPreferences, error)
}

type AssetsKeeper interface {
	GetAsset(ctx context.Context, denom string) (assetstypes.Asset, error)
}
