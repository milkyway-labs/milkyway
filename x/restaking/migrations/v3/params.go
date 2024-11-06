package v3

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	legacytypes "github.com/milkyway-labs/milkyway/x/restaking/legacy/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

func migrateParams(ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.Codec,
	restakingKeeper RestakingKeeper,
) error {
	// Read the legacy params
	var legacyParams legacytypes.Params
	store := ctx.KVStore(storeKey)
	bz := store.Get(types.LegacyParamsKey)
	if bz == nil {
		// Set the default parameters
		restakingKeeper.SetParams(ctx, types.DefaultParams())
		return nil
	}
	// Decode the readed data
	err := cdc.Unmarshal(bz, &legacyParams)
	if err != nil {
		return err
	}

	// Create a new Params instance with the same unbonding time
	// and an empty allowed denoms list to allow all denoms.
	newParams := types.NewParams(legacyParams.UnbondingTime, nil)
	return restakingKeeper.SetParams(ctx, newParams)
}
