package keeper

import (
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec
	hooks    types.OperatorsHooks

	accountKeeper types.AccountKeeper
	poolKeeper    types.CommunityPoolKeeper

	authority string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	poolKeeper types.CommunityPoolKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		authority:     authority,
		accountKeeper: accountKeeper,
		poolKeeper:    poolKeeper,
	}
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// SetHooks allows to set the operators hooks
func (k *Keeper) SetHooks(rs types.OperatorsHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set avs hooks twice")
	}

	k.hooks = rs
	return k
}
