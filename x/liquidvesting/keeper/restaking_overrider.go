package keeper

import (
	"context"

	restakingtypes "github.com/milkyway-labs/milkyway/v10/x/restaking/types"
)

type restakingOverriderState = int

const (
	restakingOverriderStateNone restakingOverriderState = iota
	restakingOverriderStateWait
	restakingOverriderStateOverride
)

type restakingOverrider struct {
	state               restakingOverriderState
	GetDelegationTarget func(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32) (restakingtypes.DelegationTarget, error)
	GetDelegation       func(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32, delegator string) (restakingtypes.Delegation, bool, error)
}

func (k *Keeper) withRestakingOverrider(f func() error) error {
	k.restakingOverrider.state = restakingOverriderStateWait
	defer func() {
		k.restakingOverrider.state = restakingOverriderStateNone
	}()
	return f()
}
