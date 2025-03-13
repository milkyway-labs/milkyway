package keeper

import (
	"context"

	"github.com/milkyway-labs/milkyway/v10/x/investors/types"
)

// DistributeVestingInvestorsRewards distributes the redirected staking rewards
// back to the vesting investors after deducting the amount based on the
// parameter. The remaining rewards are transferred to the community pool.
func (k *Keeper) DistributeVestingInvestorsRewards(ctx context.Context) error {
	err := k.VestingInvestorsRewards.Walk(ctx, nil, func(delegator string, rewards types.VestingInvestorRewards) (stop bool, err error) {
		delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(delegator)
		if err != nil {
			return true, err
		}

		// Get the withdraw address again, as the delegator might have changed it after
		// withdrawing the rewards
		withdrawAddr, err := k.distrKeeper.GetDelegatorWithdrawAddr(ctx, delAddr)
		if err != nil {
			return true, err
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, rewards.Rewards)
		if err != nil {
			return true, err
		}
		return false, nil
	})
	if err != nil {
		return err
	}

	// Clear the entire investors rewards map every block
	err = k.VestingInvestorsRewards.Clear(ctx, nil)
	if err != nil {
		return err
	}

	// Transfer the remaining rewards to the community pool
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	moduleBalances := k.bankKeeper.GetAllBalances(ctx, moduleAddr)
	if !moduleBalances.IsZero() {
		err = k.distrKeeper.FundCommunityPool(ctx, moduleBalances, moduleAddr)
		if err != nil {
			return err
		}
	}
	return nil
}
