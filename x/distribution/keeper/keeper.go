package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

type Keeper struct {
	distrkeeper.Keeper

	authKeeper    types.AccountKeeper
	stakingKeeper types.StakingKeeper
	hooks         DistrHooks
}

// NewKeeper creates a new distribution Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, storeService store.KVStoreService,
	ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper,
	feeCollectorName, authority string,
) Keeper {
	return Keeper{
		Keeper:        distrkeeper.NewKeeper(cdc, storeService, ak, bk, sk, feeCollectorName, authority),
		authKeeper:    ak,
		stakingKeeper: sk,
	}
}

func (k Keeper) WithdrawDelegationRewards(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error) {
	fmt.Println("Custom WithdrawDelegationRewards")
	err := k.BeforeWithdrawDelegationRewards(ctx, delAddr, valAddr)
	if err != nil {
		return nil, err
	}
	rewards, err := k.Keeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
	if err != nil {
		return nil, err
	}
	err = k.AfterWithdrawDelegationRewards(ctx, delAddr, valAddr, rewards)
	if err != nil {
		return nil, err
	}
	return rewards, nil
}

type DistrHooks interface {
	BeforeWithdrawDelegationRewards(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error
	AfterWithdrawDelegationRewards(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, rewards sdk.Coins) error
}

// SetHooks sets the distr hooks.
func (k *Keeper) SetHooks(dh DistrHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set distr hooks twice")
	}
	k.hooks = dh
	return k
}

// BeforeWithdrawDelegationRewards calls the BeforeWithdrawDelegationRewards hook
// if it is registered.
func (k Keeper) BeforeWithdrawDelegationRewards(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		return k.hooks.BeforeWithdrawDelegationRewards(ctx, delAddr, valAddr)
	}
	return nil
}

// AfterWithdrawDelegationRewards calls the AfterWithdrawDelegationRewards hook
// if it is registered.
func (k Keeper) AfterWithdrawDelegationRewards(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, rewards sdk.Coins) error {
	if k.hooks != nil {
		return k.hooks.AfterWithdrawDelegationRewards(ctx, delAddr, valAddr, rewards)
	}
	return nil
}
