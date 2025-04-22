package types

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper defines the expected account keeper.
type AccountKeeper interface {
	AddressCodec() address.Codec
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
}

type BankKeeper interface {
	SendCoinsFromModuleToAccount(
		ctx context.Context,
		senderModule string,
		recipientAddr sdk.AccAddress,
		amt sdk.Coins,
	) error
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
}

type StakingKeeper interface {
	ValidatorAddressCodec() address.Codec
	IterateDelegatorDelegations(
		ctx context.Context,
		delegator sdk.AccAddress,
		cb func(delegation stakingtypes.Delegation) (stop bool),
	) error
}

type DistrKeeper interface {
	GetDelegatorWithdrawAddr(ctx context.Context, delAddr sdk.AccAddress) (sdk.AccAddress, error)
	IterateDelegatorWithdrawAddrs(ctx context.Context, handler func(del, addr sdk.AccAddress) (stop bool))
	WithdrawDelegationRewards(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error)
	FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error
}
