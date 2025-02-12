package keeper

import (
	"context"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	bankkeeper.BaseKeeper

	ak    types.AccountKeeper
	hooks BankHooks
}

// NewKeeper creates a new bank Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	ak types.AccountKeeper,
	blockedAddrs map[string]bool,
	authority string,
	logger log.Logger,
) Keeper {
	return Keeper{
		BaseKeeper: bankkeeper.NewBaseKeeper(cdc, storeService, ak, blockedAddrs, authority, logger),
		ak:         ak,
	}
}

// DelegateCoins performs delegation by deducting amt coins from an account with
// address addr. For vesting accounts, delegations amounts are tracked for both
// vesting and vested coins. The coins are then transferred from the delegator
// address to a ModuleAccount address. If any of the delegation amounts are negative,
// an error is returned.
func (k Keeper) DelegateCoins(ctx context.Context, delegatorAddr, moduleAccAddr sdk.AccAddress, amt sdk.Coins) error {
	// call the TrackBeforeSend hooks and the BlockBeforeSend hooks
	err := k.BlockBeforeSend(ctx, delegatorAddr, moduleAccAddr, amt)
	if err != nil {
		return err
	}

	k.TrackBeforeSend(ctx, delegatorAddr, moduleAccAddr, amt)

	return k.BaseKeeper.DelegateCoins(ctx, delegatorAddr, moduleAccAddr, amt)
}

// UndelegateCoins performs undelegation by crediting amt coins to an account with
// address addr. For vesting accounts, undelegation amounts are tracked for both
// vesting and vested coins. The coins are then transferred from a ModuleAccount
// address to the delegator address. If any of the undelegation amounts are
// negative, an error is returned.
func (k Keeper) UndelegateCoins(ctx context.Context, moduleAccAddr, delegatorAddr sdk.AccAddress, amt sdk.Coins) error {
	// call the TrackBeforeSend hooks and the BlockBeforeSend hooks
	err := k.BlockBeforeSend(ctx, moduleAccAddr, delegatorAddr, amt)
	if err != nil {
		return err
	}

	k.TrackBeforeSend(ctx, moduleAccAddr, delegatorAddr, amt)

	return k.BaseKeeper.UndelegateCoins(ctx, moduleAccAddr, moduleAccAddr, amt)
}

type BankHooks interface {
	TrackBeforeSend(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins)       // Must be before any send is executed
	BlockBeforeSend(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins) error // Must be before any send is executed
}

// TrackBeforeSend executes the TrackBeforeSend hook if registered.
func (k Keeper) TrackBeforeSend(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins) {
	if k.hooks != nil {
		k.hooks.TrackBeforeSend(ctx, from, to, amount)
	}
}

// BlockBeforeSend executes the BlockBeforeSend hook if registered.
func (k Keeper) BlockBeforeSend(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins) error {
	if k.hooks != nil {
		return k.hooks.BlockBeforeSend(ctx, from, to, amount)
	}
	return nil
}

// SetHooks sets the bank hooks.
func (k *Keeper) SetHooks(bh BankHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set bank hooks twice")
	}
	k.hooks = bh
	return k
}

// SendCoins transfers amt coins from a sending account to a receiving account.
// An error is returned upon failure.
func (k Keeper) SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error {
	// BlockBeforeSend hook should always be called before the TrackBeforeSend hook.
	err := k.BlockBeforeSend(ctx, fromAddr, toAddr, amt)
	if err != nil {
		return err
	}
	// call the TrackBeforeSend hooks
	k.TrackBeforeSend(ctx, fromAddr, toAddr, amt)
	return k.BaseSendKeeper.SendCoins(ctx, fromAddr, toAddr, amt)
}

// InputOutputCoins performs multi-send functionality. It accepts an
// input that corresponds to a series of outputs. It returns an error if the
// input and outputs don't line up or if any single transfer of tokens fails.
func (k Keeper) InputOutputCoins(ctx context.Context, input types.Input, outputs []types.Output) error {
	inAddress, err := k.ak.AddressCodec().StringToBytes(input.Address)
	if err != nil {
		return err
	}

	for _, out := range outputs {
		outAddress, err := k.ak.AddressCodec().StringToBytes(out.Address)
		if err != nil {
			return err
		}

		if err := k.BlockBeforeSend(ctx, inAddress, outAddress, out.Coins); err != nil {
			return err
		}

		k.TrackBeforeSend(ctx, inAddress, outAddress, out.Coins)
	}

	return k.BaseKeeper.InputOutputCoins(ctx, input, outputs)
}
