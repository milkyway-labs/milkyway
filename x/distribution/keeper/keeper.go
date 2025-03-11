package keeper

import (
	"context"

	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

type Keeper struct {
	distrkeeper.Keeper

	authKeeper types.AccountKeeper
	hooks      DistrHooks
}

// NewKeeper creates a new distribution Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, storeService store.KVStoreService,
	ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper,
	feeCollectorName, authority string,
) Keeper {
	return Keeper{
		Keeper:     distrkeeper.NewKeeper(cdc, storeService, ak, bk, sk, feeCollectorName, authority),
		authKeeper: ak,
	}
}

type DistrHooks interface {
	AfterSetWithdrawAddress(ctx context.Context, delAddr, withdrawAddr sdk.AccAddress) error
}

// SetHooks sets the distr hooks.
func (k *Keeper) SetHooks(dh DistrHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set distr hooks twice")
	}
	k.hooks = dh
	return k
}

// AfterSetWithdrawAddress calls the AfterSetWithdrawAddress hook if it is
// registered.
func (k Keeper) AfterSetWithdrawAddress(ctx context.Context, delAddr, withdrawAddr sdk.AccAddress) error {
	if k.hooks != nil {
		return k.hooks.AfterSetWithdrawAddress(ctx, delAddr, withdrawAddr)
	}
	return nil
}
