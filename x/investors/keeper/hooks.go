package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var _ distrtypes.DistrHooks = Hooks{}

type Hooks struct {
	*Keeper
}

func (k *Keeper) Hooks() Hooks {
	return Hooks{k}
}

func (h Hooks) BeforeDelegationRewardsWithdrawn(ctx context.Context, _ stakingtypes.ValidatorI, del stakingtypes.DelegationI, _ sdk.AccAddress, rewards sdk.DecCoins) (finalRewards sdk.DecCoins, err error) {
	// If the delegator is not a vesting investor, skip.
	isVestingInvestor, err := h.IsVestingInvestor(ctx, del.GetDelegatorAddr())
	if err != nil {
		return nil, err
	}
	if !isVestingInvestor {
		return rewards, nil
	}

	investorsRewardRatio, err := h.GetInvestorsRewardRatio(ctx)
	if err != nil {
		return nil, err
	}
	return rewards.MulDecTruncate(investorsRewardRatio), nil
}
