package keeper

import (
	"context"

	"github.com/milkyway-labs/milkyway/v9/x/liquidvesting/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
	rewardstypes "github.com/milkyway-labs/milkyway/v9/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/v9/x/services/types"
)

// AdjustedServicesKeeper wraps around rewardstypes.ServicesKeeper and adjusts
// the delegator shares of a service by deducting the uncovered locked shares.
type AdjustedServicesKeeper struct {
	rewardstypes.ServicesKeeper
	k *Keeper
}

// AdjustedServicesKeeper returns a new instance of AdjustedServicesKeeper.
func (k *Keeper) AdjustedServicesKeeper(servicesKeeper rewardstypes.ServicesKeeper) *AdjustedServicesKeeper {
	return &AdjustedServicesKeeper{
		ServicesKeeper: servicesKeeper,
		k:              k,
	}
}

// GetService returns the service with the given serviceID and deducts the
// uncovered locked shares from the delegator shares.
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

// AdjustedRestakingKeeper wraps around rewardstypes.RestakingKeeper and adjusts
// the delegation shares by deducting the uncovered locked shares.
type AdjustedRestakingKeeper struct {
	rewardstypes.RestakingKeeper
	k *Keeper
}

// AdjustedRestakingKeeper returns a new instance of AdjustedRestakingKeeper.
func (k *Keeper) AdjustedRestakingKeeper(restakingKeeper rewardstypes.RestakingKeeper) *AdjustedRestakingKeeper {
	return &AdjustedRestakingKeeper{
		RestakingKeeper: restakingKeeper,
		k:               k,
	}
}

// GetDelegationTarget returns the delegation target with the given targetID and
// deducts the uncovered locked shares from the delegator shares.
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

	target, _, err = types.RemoveDelShares(target, uncoveredLockedShares)
	return target, err
}

// GetDelegation returns the delegation with the given targetID and deducts the
// uncovered locked shares from the delegator shares.
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
