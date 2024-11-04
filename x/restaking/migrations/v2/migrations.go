package v2

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Migrate1To2(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.Codec,
	restakingKeeper RestakingKeeper,
	operatorsKeeper OperatorsKeeper,
	servicesKeeper ServicesKeeper,
) error {
	err := migateOperatorParams(ctx, storeKey, cdc, restakingKeeper, operatorsKeeper)
	if err != nil {
		return err
	}

	err = migrateServiceParams(ctx, storeKey, cdc, restakingKeeper, servicesKeeper)
	if err != nil {
		return err
	}

	return nil
}
