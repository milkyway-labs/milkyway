package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type stakingOverriderState = int

const (
	stakingOverriderStateNone stakingOverriderState = iota
	stakingOverriderStateWait
	stakingOverriderStateOverride
)

// stakingKeeperOverrider is used to override Validator and Delegation methods
// that are used inside the distrKeeper.WithdrawDelegationRewards method. Inside
// distrKeeper.WithdrawDelegationRewards, the distribution keeper first gets the
// validator and the delegation from the adjusted staking keeper to calculate the
// rewards accumulated so far. We need to use the previous states(vesting
// investors reward ratio, validator investor shares, delegation shares) up to
// this point. But after withdrawing rewards, inside initializeDelegation, the
// new states must be taken into account in order to update the delegator
// starting info properly and that's why we use stakingKeeperOverrider.
// stakingKeeperOverrider doesn't override the methods immediately when state is
// set to stakingOverriderStateWait. Instead, the methods will be overridden after the
// next call of stakingKeeper.Delegation ensuring that the overrides occur after
// distrKeeper.WithdrawDelegationRewards has completed and just before
// distrKeeper.initializeDelegation is executed.
//
// # state: none -> wait
// val = stakingKeeper.Validator() # uses prev states
// del = stakingKeeper.Delegation() # uses prev states
// # state: wait -> override
// distrKeeper.withdrawDelegationRewards(val, del)
// distrKeeper.initializeDelegation()
// - val = stakingKeeper.Validator() # uses new states
// - del = stakingKeeper.Delegation() # uses new states
// # state: override -> none
// - set delegator starting info
type stakingKeeperOverrider struct {
	state      stakingOverriderState
	Validator  func(ctx context.Context, address sdk.ValAddress) (stakingtypes.ValidatorI, error)
	Delegation func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.DelegationI, error)
}

func (k *Keeper) withOverrider(f func() error) error {
	k.stakingKeeperOverrider.state = stakingOverriderStateWait
	defer func() {
		k.stakingKeeperOverrider.state = stakingOverriderStateNone
	}()
	return f()
}
