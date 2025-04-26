package liquidvesting

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v12/x/liquidvesting/keeper"
	"github.com/milkyway-labs/milkyway/v12/x/liquidvesting/types"
)

// EndBlocker is called every block and is responsible for maturing unbonding delegations
func EndBlocker(ctx sdk.Context, k *keeper.Keeper) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	return k.CompleteBurnCoins(ctx)
}
