package v3

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Migrate2To3(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.Codec,
	restakingKeer RestakingKeeper,
) error {
	// Migrate the params from the old format to the new one.
	err := migrateParams(ctx, storeKey, cdc, restakingKeer)
	if err != nil {
		return err
	}

	return nil
}
