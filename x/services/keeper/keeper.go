package keeper

import (
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec
	hooks    types.ServicesHooks

	poolKeeper types.CommunityPoolKeeper

	// authority represents the address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority string
}

// NewKeeper creates a new keeper
func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey, poolKeeper types.CommunityPoolKeeper, authority string) *Keeper {
	return &Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		poolKeeper: poolKeeper,
		authority:  authority,
	}
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// SetHooks allows to set the reactions hooks
func (k *Keeper) SetHooks(rs types.ServicesHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set services hooks twice")
	}

	k.hooks = rs
	return k
}
