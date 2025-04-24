package investors

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"

	"github.com/milkyway-labs/milkyway/v11/x/investors/keeper"
	"github.com/milkyway-labs/milkyway/v11/x/investors/types"
)

// BeginBlocker is called every block and is responsible for removing the
// investors that have ended their vesting period from the vesting investors
// queue.
func BeginBlocker(ctx context.Context, k *keeper.Keeper) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	return k.RemoveVestingEndedInvestors(ctx)
}
