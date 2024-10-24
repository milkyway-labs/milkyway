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

	// Wipe the service's operators allow list
	serviceOperatorsAllowListIter, err := h.serviceOperatorsAllowList.Iterate(ctx, collections.NewPrefixedPairRange[uint32, uint32](serviceID))
	if err != nil {
		return err
	}
	defer serviceOperatorsAllowListIter.Close()

	// Iterate over all the items and remove them from the KeySet
	for ; serviceOperatorsAllowListIter.Valid(); serviceOperatorsAllowListIter.Next() {
		serviceOperatorPair, err := serviceOperatorsAllowListIter.Key()
		if err != nil {
			return err
		}
		err = h.serviceOperatorsAllowList.Remove(ctx, serviceOperatorPair)
		if err != nil {
			return err
		}
	}

	// Wipe the list of polls from which the service is allowed to
	// borrow security
	serviceSecuringPoolsIter, err := h.serviceSecuringPools.Iterate(ctx, collections.NewPrefixedPairRange[uint32, uint32](serviceID))
	if err != nil {
		return err
	}
	defer serviceSecuringPoolsIter.Close()

	// Iterate over all the items and remove them from the KeySet
	for ; serviceSecuringPoolsIter.Valid(); serviceSecuringPoolsIter.Next() {
		servicePoolPair, err := serviceSecuringPoolsIter.Key()
		if err != nil {
			return err
		}
		err = h.serviceSecuringPools.Remove(ctx, servicePoolPair)
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
