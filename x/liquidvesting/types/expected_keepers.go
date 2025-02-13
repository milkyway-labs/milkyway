package types

import (
	"context"
	"time"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
	rewardstypes "github.com/milkyway-labs/milkyway/v9/x/rewards/types"
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
	UnbondRestakedAssets(ctx context.Context, user sdk.AccAddress, amount sdk.Coins) (time.Time, error)
	GetAllUnbondingDelegations(ctx context.Context) ([]restakingtypes.UnbondingDelegation, error)
	GetAllUserRestakedCoins(ctx context.Context, userAddress string) (sdk.DecCoins, error)
	GetAllUserUnbondingDelegations(ctx context.Context, userAddress string) []restakingtypes.UnbondingDelegation
	IterateUserDelegations(ctx context.Context, userAddress string, cb func(del restakingtypes.Delegation) (stop bool, err error)) error
	GetDelegationTarget(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32) (restakingtypes.DelegationTarget, error)
	GetDelegation(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32, delegator string) (restakingtypes.Delegation, bool, error)
	GetDelegationForTarget(ctx context.Context, target restakingtypes.DelegationTarget, delegator string) (restakingtypes.Delegation, bool, error)
	GetAllDelegations(ctx context.Context) ([]restakingtypes.Delegation, error)
}

type RewardsKeeper interface {
	WithdrawDelegationRewards(ctx context.Context, delAddr sdk.AccAddress, delType restakingtypes.DelegationType, targetID uint32) (rewardstypes.Pools, error)
}
