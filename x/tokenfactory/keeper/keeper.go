package keeper

import (
	"fmt"

	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v12/x/tokenfactory/types"
)

type (
	Keeper struct {
		cdc          codec.Codec
		storeService corestoretypes.KVStoreService

		accountKeeper       types.AccountKeeper
		bankKeeper          types.BankKeeper
		contractKeeper      types.ContractKeeper
		communityPoolKeeper types.CommunityPoolKeeper

		authority string
	}
)

// NewKeeper returns a new instance of the x/tokenfactory keeper
func NewKeeper(
	cdc codec.Codec,
	storeService corestoretypes.KVStoreService,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	communityPoolKeeper types.CommunityPoolKeeper,
	authority string,
) Keeper {
	return Keeper{
		cdc:                 cdc,
		storeService:        storeService,
		accountKeeper:       accountKeeper,
		bankKeeper:          bankKeeper,
		communityPoolKeeper: communityPoolKeeper,
		authority:           authority,
	}
}

// Logger returns a logger for the x/tokenfactory module
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetAuthority returns the x/tokenfactory module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// GetDenomPrefixStore returns the substore for a specific denom
func (k Keeper) GetDenomPrefixStore(ctx sdk.Context, denom string) storetypes.KVStore {
	store := k.storeService.OpenKVStore(ctx)
	return prefix.NewStore(runtime.KVStoreAdapter(store), types.GetDenomPrefixStore(denom))
}

// GetCreatorPrefixStore returns the substore for a specific creator address
func (k Keeper) GetCreatorPrefixStore(ctx sdk.Context, creator string) storetypes.KVStore {
	store := k.storeService.OpenKVStore(ctx)
	return prefix.NewStore(runtime.KVStoreAdapter(store), types.GetCreatorPrefix(creator))
}

// GetCreatorsPrefixStore returns the substore that contains a list of creators
func (k Keeper) GetCreatorsPrefixStore(ctx sdk.Context) storetypes.KVStore {
	store := k.storeService.OpenKVStore(ctx)
	return prefix.NewStore(runtime.KVStoreAdapter(store), types.GetCreatorsPrefix())
}

// Set the wasm keeper.
func (k *Keeper) SetContractKeeper(contractKeeper types.ContractKeeper) {
	k.contractKeeper = contractKeeper
}

// CreateModuleAccount creates a module account with minting and burning capabilities
// This account isn't intended to store any coins,
// it purely mints and burns them on behalf of the admin of respective denoms,
// and sends to the relevant address.
func (k Keeper) CreateModuleAccount(ctx sdk.Context) {
	k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}
