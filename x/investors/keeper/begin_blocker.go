package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RemoveVestingEndedInvestors removes all investors from the vesting queue that
// have ended their vesting period.
func (k *Keeper) RemoveVestingEndedInvestors(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currTime := sdkCtx.BlockTime().Unix()
	iter, err := k.InvestorsVestingQueue.Iterate(ctx, collections.NewPrefixUntilPairRange[int64, sdk.AccAddress](currTime))
	if err != nil {
		return err
	}
	keysToRemove, err := iter.Keys()
	if err != nil {
		return err
	}

	for _, key := range keysToRemove {
		err = k.InvestorsVestingQueue.Remove(ctx, key)
		if err != nil {
			return err
		}

		investorAddr := key.K2()
		err = k.RemoveVestingInvestor(ctx, investorAddr)
		if err != nil {
			return err
		}
	}
	return nil
}
