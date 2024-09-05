package keeper

import (
	"context"
)

// BeginBlocker is called every block and is used to terminate ended rewards
// plans and allocate restaking rewards for the previous block.
func (k *Keeper) BeginBlocker(ctx context.Context) error {
	err := k.TerminateEndedRewardsPlans(ctx)
	if err != nil {
		return err
	}

	err = k.AllocateRewards(ctx)
	if err != nil {
		return err
	}
	return nil
}
