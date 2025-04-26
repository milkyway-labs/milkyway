package v10

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/v12/app/keepers"
	ibchookstypes "github.com/milkyway-labs/milkyway/v12/x/ibc-hooks/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configuration module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		vm, err := mm.RunMigrations(ctx, configuration, fromVM)
		if err != nil {
			return nil, err
		}

		// Set the default ibchooks parameters
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		keepers.IBCHooksKeeper.SetParams(sdkCtx, ibchookstypes.DefaultParams())

		return vm, nil
	}
}
