package keeper

import (
	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
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
	servicesKeeper  types.ServicesKeeper

	// Keeper data
	schema           collections.Schema
	operatorServices collections.Map[uint32, types.OperatorSecuredServices]

	hooks types.RestakingHooks
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	storeService corestoretypes.KVStoreService,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	poolsKeeper types.PoolsKeeper,
	operatorsKeeper types.OperatorsKeeper,
	servicesKeeper types.ServicesKeeper,
	authority string,
) *Keeper {
	// Ensure that authority is a valid AccAddress
	if _, err := accountKeeper.AddressCodec().StringToBytes(authority); err != nil {
		panic("authority is not a valid account address")
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := &Keeper{
		storeKey: storeKey,
		cdc:      cdc,

		accountKeeper:   accountKeeper,
		bankKeeper:      bankKeeper,
		poolsKeeper:     poolsKeeper,
		operatorsKeeper: operatorsKeeper,
		servicesKeeper:  servicesKeeper,

		operatorServices: collections.NewMap(
			sb, types.OperatorServicesPrefix,
			"operator_services",
			collections.Uint32Key,
			codec.CollValue[types.OperatorSecuredServices](cdc),
		),
		authority: authority,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.schema = schema

	return k
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
