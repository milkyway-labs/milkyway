package keeper

import (
	"context"
	"fmt"

	"github.com/milkyway-labs/milkyway/v9/x/liquidvesting/types"
	operatorstypes "github.com/milkyway-labs/milkyway/v9/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v9/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
	rewardstypes "github.com/milkyway-labs/milkyway/v9/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/v9/x/services/types"
)

type AdjustedServicesKeeper struct {
	rewardstypes.ServicesKeeper
	k *Keeper
}

func (k *Keeper) AdjustedServicesKeeper(keeper rewardstypes.ServicesKeeper) *AdjustedServicesKeeper {
	return &AdjustedServicesKeeper{
		ServicesKeeper: keeper,
		k:              k,
	}
}

func (sk *AdjustedServicesKeeper) GetService(ctx context.Context, serviceID uint32) (servicestypes.Service, error) {
	service, err := sk.ServicesKeeper.GetService(ctx, serviceID)
	if err != nil {
		return servicestypes.Service{}, err
	}
	coveredLockedShares, err := sk.k.GetTargetCoveredLockedShares(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID)
	if err != nil {
		return servicestypes.Service{}, err
	}
	uncoveredLockedShares := types.UncoveredLockedShares(service.DelegatorShares, coveredLockedShares)
	service, _ = service.RemoveDelShares(uncoveredLockedShares)
	return service, nil
}

// --------------------------------------------------------------------------------------------------------------------

type AdjustedRestakingKeeper struct {
	rewardstypes.RestakingKeeper
	k *Keeper
}

func (k *Keeper) AdjustedRestakingKeeper(keeper rewardstypes.RestakingKeeper) *AdjustedRestakingKeeper {
	return &AdjustedRestakingKeeper{
		RestakingKeeper: keeper,
		k:               k,
	}
}

func (rk *AdjustedRestakingKeeper) GetDelegationTarget(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32) (restakingtypes.DelegationTarget, error) {
	if rk.k.restakingOverrider.state == restakingOverriderStateOverride {
		return rk.k.restakingOverrider.GetDelegationTarget(ctx, delType, targetID)
	}

	target, err := rk.RestakingKeeper.GetDelegationTarget(ctx, delType, targetID)
	if err != nil {
		return nil, err
	}
	coveredLockedShares, err := rk.k.GetTargetCoveredLockedShares(ctx, delType, targetID)
	if err != nil {
		return nil, err
	}
	uncoveredLockedShares := types.UncoveredLockedShares(target.GetDelegatorShares(), coveredLockedShares)

	switch target := target.(type) {
	case poolstypes.Pool:
		target, _, err = target.RemoveDelShares(uncoveredLockedShares)
		if err != nil {
			return nil, err
		}
		return target, nil
	case operatorstypes.Operator:
		target, _ = target.RemoveDelShares(uncoveredLockedShares)
		return target, nil
	case servicestypes.Service:
		target, _ = target.RemoveDelShares(uncoveredLockedShares)
		return target, nil
	default:
		return nil, fmt.Errorf("invalid target type %T", target)
	}
}

func (rk *AdjustedRestakingKeeper) GetDelegation(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32, delegator string) (restakingtypes.Delegation, bool, error) {
	if rk.k.restakingOverrider.state == restakingOverriderStateOverride {
		return rk.k.restakingOverrider.GetDelegation(ctx, delType, targetID, delegator)
	}

	delegation, found, err := rk.RestakingKeeper.GetDelegation(ctx, delType, targetID, delegator)
	if err != nil || !found {
		return restakingtypes.Delegation{}, found, err
	}

	coveredLockedShares, err := rk.k.GetCoveredLockedShares(ctx, delegation)
	if err != nil {
		return restakingtypes.Delegation{}, false, err
	}
	delegation.Shares = types.DeductUncoveredLockedShares(delegation.Shares, coveredLockedShares)

	// After the first call to GetDelegation, transition the overrider's
	// state to Override
	if rk.k.restakingOverrider.state == restakingOverriderStateWait {
		rk.k.restakingOverrider.state = restakingOverriderStateOverride
	}

	return delegation, true, nil
}
