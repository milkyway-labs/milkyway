package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/milkyway-labs/milkyway/v9/x/investors/types"
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
	if sk.k.stakingKeeperOverrider.state == stakingOverriderStateOverride {
		return sk.k.stakingKeeperOverrider.Validator(ctx, address)
	}

	return sk.k.GetAdjustedValidator(ctx, address)
}

func (sk *AdjustedStakingKeeper) Delegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.DelegationI, error) {
	if sk.k.stakingKeeperOverrider.state == stakingOverriderStateOverride {
		return sk.k.stakingKeeperOverrider.Delegation(ctx, delAddr, valAddr)
	}

	delegation, err := sk.GetDelegation(ctx, delAddr, valAddr)
	if err != nil {
		return nil, err
	}

	isVestingInvestor, err := sk.k.VestingInvestors.Has(ctx, delAddr)
	if err != nil {
		return nil, err
	}
	if isVestingInvestor {
		investorsRewardRatio, err := sk.k.InvestorsRewardRatio.Get(ctx)
		if err != nil {
			return nil, err
		}
		delegation.Shares = delegation.Shares.MulTruncate(investorsRewardRatio)
	}

	// After the first call to Delegation, transition the injector's state to Inject
	if sk.k.stakingKeeperOverrider.state == stakingOverriderStateWait {
		sk.k.stakingKeeperOverrider.state = stakingOverriderStateOverride
	}

	return delegation, nil
}
