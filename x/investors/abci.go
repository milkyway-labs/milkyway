package investors

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"

	"github.com/milkyway-labs/milkyway/v9/x/investors/keeper"
	"github.com/milkyway-labs/milkyway/v9/x/investors/types"
)

// BeginBlocker is called every block and is responsible for removing
func BeginBlocker(ctx context.Context, k *keeper.Keeper) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	return k.RemoveVestingEndedInvestors(ctx)
}
