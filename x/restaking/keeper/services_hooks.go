package keeper

import (
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

type ServicesHooks struct {
	*Keeper
}

func (k *Keeper) ServicesHooks() servicestypes.ServicesHooks {
	return &ServicesHooks{k}
}

// ------------------------------------------------------------------------------

// AfterServiceDeactivated implements types.ServicesHooks.
func (h *ServicesHooks) AfterServiceDeactivated(ctx sdk.Context, serviceID uint32) error {
	// After the service has been deactivated
	// we remove the data that we keep in the x/restaking
	// associated to this service.

	serviceValidatingOperatorsIter, err := h.operatorJoinedServices.Indexes.Service.Iterate(ctx, collections.NewPrefixedPairRange[uint32, uint32](serviceID))
	if err != nil {
		return err
	}
	defer serviceValidatingOperatorsIter.Close()

	// Iterate over the operator that have joined this service
	// and remove the participation
	for ; serviceValidatingOperatorsIter.Valid(); serviceValidatingOperatorsIter.Next() {
		operatorServicePair, err := serviceValidatingOperatorsIter.PrimaryKey()
		if err != nil {
			return err
		}
		err = h.operatorJoinedServices.Remove(ctx, operatorServicePair)
		if err != nil {
			return err
		}
	}

	// We remove the allow list associated to this service
	serviceAllowListIter, err := h.serviceOperatorsAllowList.Iterate(ctx, collections.NewPrefixedPairRange[uint32, uint32](serviceID))
	if err != nil {
		return err
	}
	defer serviceAllowListIter.Close()

	// Clear the service's operators allow list
	for ; serviceAllowListIter.Valid(); serviceAllowListIter.Next() {
		serviceOperatorPair, err := serviceAllowListIter.Key()
		if err != nil {
			return err
		}
		err = h.serviceOperatorsAllowList.Remove(ctx, serviceOperatorPair)
		if err != nil {
			return err
		}
	}

	return nil
}

// AfterServiceActivated implements types.ServicesHooks.
func (h *ServicesHooks) AfterServiceActivated(ctx sdk.Context, serviceID uint32) error {
	return nil
}

// AfterServiceCreated implements types.ServicesHooks.
func (h *ServicesHooks) AfterServiceCreated(ctx sdk.Context, serviceID uint32) error {
	return nil
}
