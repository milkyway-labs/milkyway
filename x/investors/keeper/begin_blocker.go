package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) RemoveVestingEndedInvestors(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currTime := sdkCtx.BlockHeader().Time.Unix()
	var keysToRemove []collections.Pair[int64, sdk.AccAddress]
	err := k.InvestorsVestingQueue.Walk(
		ctx,
		collections.NewPrefixUntilPairRange[int64, sdk.AccAddress](currTime),
		func(key collections.Pair[int64, sdk.AccAddress]) (stop bool, err error) {
			keysToRemove = append(keysToRemove, key)
			return false, nil
		},
	)
	if err != nil {
		return err
	}

	for _, key := range keysToRemove {
		err = k.InvestorsVestingQueue.Remove(ctx, key)
		if err != nil {
			return err
		}

		investorAddr := key.K2()
		err = k.VestingInvestors.Remove(ctx, investorAddr)
		if err != nil {
			return err
		}
	}
	return nil
}
