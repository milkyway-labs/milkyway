package types

import (
	context "context"
	"time"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

type AccountKeeper interface {
	AddressCodec() address.Codec
}

type BankKeeper interface {
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(
		ctx context.Context,
		senderModule string,
		recipientAddr sdk.AccAddress,
		amt sdk.Coins,
	) error
	SendCoinsFromAccountToModule(
		ctx context.Context,
		senderAddr sdk.AccAddress,
		recipientModule string,
		amt sdk.Coins,
	) error
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool)
	SetDenomMetaData(ctx context.Context, metadata banktypes.Metadata)
}

type PoolsKeeper interface {
	IsPoolDelegationsAddress(ctx sdk.Context, address string) bool
	GetPoolByDenom(ctx sdk.Context, denom string) (poolstypes.Pool, bool)
}

type OperatorsKeeper interface {
	IsOperatorDelegationsAddress(ctx sdk.Context, address string) bool
	GetOperator(ctx sdk.Context, operatorID uint32) (operatorstypes.Operator, bool)
}

type ServicesKeeper interface {
	IsServiceDelegationsAddress(ctx sdk.Context, address string) bool
	GetService(ctx sdk.Context, serviceID uint32) (servicestypes.Service, bool)
}

type RestakingKeeper interface {
	GetPoolDelegation(ctx sdk.Context, poolID uint32, userAddress string) (restakingtypes.Delegation, bool)
	IterateUserServiceDelegations(ctx sdk.Context, userAddress string, cb func(restakingtypes.Delegation) (bool, error)) error
	IterateUserOperatorDelegations(ctx sdk.Context, userAddress string, cb func(restakingtypes.Delegation) (bool, error)) error
	UnbondRestakedAssets(ctx sdk.Context, user sdk.AccAddress, amount sdk.Coins) (time.Time, error)
	GetAllUnbondingDelegations(ctx sdk.Context) []restakingtypes.UnbondingDelegation
	GetAllUserRestakedCoins(ctx sdk.Context, userAddress string) (sdk.DecCoins, error)
	GetAllUserUnbondingDelegations(ctx sdk.Context, userAddress string) []restakingtypes.UnbondingDelegation
}
