package keeper

import (
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.Codec

	authority string

	accountKeeper   types.AccountKeeper
	bankKeeper      types.BankKeeper
	poolsKeeper     types.PoolsKeeper
	operatorsKeeper types.OperatorsKeeper

	hooks types.RestakingHooks
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	poolsKeeper types.PoolsKeeper,
	operatorsKeeper types.OperatorsKeeper,
	authority string,
) *Keeper {

	// Ensure that authority is a valid AccAddress
	if _, err := accountKeeper.AddressCodec().StringToBytes(authority); err != nil {
		panic("authority is not a valid account address")
	}

	return &Keeper{
		storeKey: storeKey,
		cdc:      cdc,

		accountKeeper:   accountKeeper,
		bankKeeper:      bankKeeper,
		poolsKeeper:     poolsKeeper,
		operatorsKeeper: operatorsKeeper,

		authority: authority,
	}
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// SetHooks allows to set the reactions hooks
func (k *Keeper) SetHooks(rs types.RestakingHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set services hooks twice")
	}

	k.hooks = rs
	return k
}
