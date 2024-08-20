package keeper

import (
	"context"
)

// BeginBlocker allocates restaking rewards for the previous block.
func (k *Keeper) BeginBlocker(ctx context.Context) error {
	err := k.AllocateRewards(ctx)
	if err != nil {
		return err
	}
	return nil
}
