package types

import (
	"context"
	"time"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	restakingtypes "github.com/milkyway-labs/milkyway/v6/x/restaking/types"
)

type AccountKeeper interface {
	AddressCodec() address.Codec
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
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
	IsPoolAddress(ctx context.Context, address string) (bool, error)
}

type OperatorsKeeper interface {
	IsOperatorAddress(ctx context.Context, address string) (bool, error)
}

type ServicesKeeper interface {
	IsServiceAddress(ctx context.Context, address string) (bool, error)
}

type RestakingKeeper interface {
	GetPoolDelegation(ctx context.Context, poolID uint32, userAddress string) (restakingtypes.Delegation, bool, error)
	IterateUserServiceDelegations(ctx context.Context, userAddress string, cb func(restakingtypes.Delegation) (bool, error)) error
	IterateUserOperatorDelegations(ctx context.Context, userAddress string, cb func(restakingtypes.Delegation) (bool, error)) error
	UnbondRestakedAssets(ctx context.Context, user sdk.AccAddress, amount sdk.Coins) (time.Time, error)
	GetAllUnbondingDelegations(ctx context.Context) ([]restakingtypes.UnbondingDelegation, error)
	GetAllUserRestakedCoins(ctx context.Context, userAddress string) (sdk.DecCoins, error)
	GetAllUserUnbondingDelegations(ctx context.Context, userAddress string) []restakingtypes.UnbondingDelegation
}
