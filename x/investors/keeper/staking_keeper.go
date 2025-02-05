package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/milkyway-labs/milkyway/v7/x/investors/types"
)

type AdjustedStakingKeeper struct {
	types.StakingKeeper
	k *Keeper
}

func (k *Keeper) AdjustedStakingKeeper(stakingKeeper types.StakingKeeper) *AdjustedStakingKeeper {
	return &AdjustedStakingKeeper{
		StakingKeeper: stakingKeeper,
		k:             k,
	}
}

func (sk *AdjustedStakingKeeper) Validator(ctx context.Context, address sdk.ValAddress) (stakingtypes.ValidatorI, error) {
	return sk.k.GetAdjustedValidator(ctx, address)
}

func (sk *AdjustedStakingKeeper) Delegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.DelegationI, error) {
	delegation, isInvestor, err := sk.k.GetDelegation(ctx, delAddr, valAddr)
	if err != nil {
		return nil, err
	}
	if isInvestor {
		rewardRatio, err := sk.k.InvestorsRewardRatio.Get(ctx)
		if err != nil {
			return nil, err
		}
		delegation.Shares = delegation.Shares.MulTruncate(rewardRatio)
	}
	return delegation, nil
}
