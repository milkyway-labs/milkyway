package keeper

import (
	"context"
)

func (k *Keeper) BeginBlocker(ctx context.Context) error {
	err := k.AllocateRewards(ctx)
	if err != nil {
		return err
	}
	return nil
}
