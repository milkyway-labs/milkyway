package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

func (k *Keeper) getDelegationKeyBuilders(delegation types.Delegation) (types.DelegationKeyBuilder, types.DelegationByTargetIDBuilder, error) {
	switch delegation.Type {
	case types.DELEGATION_TYPE_POOL:
		return types.UserPoolDelegationStoreKey, types.DelegationByPoolIDStoreKey, nil

	case types.DELEGATION_TYPE_OPERATOR:
		return types.UserOperatorDelegationStoreKey, types.DelegationByOperatorIDStoreKey, nil

	case types.DELEGATION_TYPE_SERVICE:
		return types.UserServiceDelegationStoreKey, types.DelegationByServiceIDStoreKey, nil

	default:
		return nil, nil, types.ErrInvalidDelegationType
	}
}

// SetDelegation stores the given delegation in the store
func (k *Keeper) SetDelegation(ctx sdk.Context, delegation types.Delegation) error {
	store := ctx.KVStore(k.storeKey)

	// Get the keys builders
	getDelegationKey, getDelegationByTargetID, err := k.getDelegationKeyBuilders(delegation)
	if err != nil {
		return err
	}

	// Marshal and store the delegation
	delegationBz := types.MustMarshalDelegation(k.cdc, delegation)
	store.Set(getDelegationKey(delegation.UserAddress, delegation.TargetID), delegationBz)

	// Store the delegation in the delegations by pool ID store
	store.Set(getDelegationByTargetID(delegation.TargetID, delegation.UserAddress), []byte{})

	return nil
}

// GetDelegationForTarget returns the delegation for the given delegator and target.
func (k *Keeper) GetDelegationForTarget(
	ctx sdk.Context, target types.DelegationTarget, delegator string,
) (types.Delegation, bool) {
	switch target.(type) {
	case *poolstypes.Pool:
		return k.GetPoolDelegation(ctx, target.GetID(), delegator)
	case *operatorstypes.Operator:
		return k.GetOperatorDelegation(ctx, target.GetID(), delegator)
	case *servicestypes.Service:
		return k.GetServiceDelegation(ctx, target.GetID(), delegator)
	default:
		return types.Delegation{}, false
	}
}

// RemoveDelegation removes the given delegation from the store
func (k *Keeper) RemoveDelegation(ctx sdk.Context, delegation types.Delegation) {
	switch delegation.Type {
	case types.DELEGATION_TYPE_POOL:
		k.RemovePoolDelegation(ctx, delegation)
	case types.DELEGATION_TYPE_OPERATOR:
		k.RemoveOperatorDelegation(ctx, delegation)
	case types.DELEGATION_TYPE_SERVICE:
		k.RemoveServiceDelegation(ctx, delegation)
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (k *Keeper) getUnbondingDelegationKeyBuilder(ud types.UnbondingDelegation) (types.UnbondingDelegationKeyBuilder, error) {
	switch ud.Type {
	case types.UNBONDING_DELEGATION_TYPE_POOL:
		return types.UserPoolUnbondingDelegationKey, nil

	case types.UNBONDING_DELEGATION_TYPE_OPERATOR:
		return types.UserOperatorUnbondingDelegationKey, nil

	case types.UNBONDING_DELEGATION_TYPE_SERVICE:
		return types.UserServiceUnbondingDelegationKey, nil

	default:
		return nil, types.ErrInvalidDelegationType
	}
}

func (k *Keeper) SetUnbondingDelegation(ctx sdk.Context, ud types.UnbondingDelegation, entryID uint64) error {
	// Get the key to be used to store the unbonding delegation
	getUnbondingDelegation, err := k.getUnbondingDelegationKeyBuilder(ud)
	if err != nil {
		return err
	}
	unbondingDelegationKey := getUnbondingDelegation(ud.DelegatorAddress, ud.TargetID)

	// Store the unbonding delegation
	store := ctx.KVStore(k.storeKey)
	store.Set(unbondingDelegationKey, types.MustMarshalUnbondingDelegation(k.cdc, ud))

	// Set the index allowing to lookup the UnbondingDelegation by the unbondingID of an
	// UnbondingDelegationEntry that it contains
	store.Set(types.GetUnbondingIndexKey(entryID), unbondingDelegationKey)

	// Set the type of the unbonding delegation so that we know how to deserialize id
	store.Set(types.GetUnbondingTypeKey(entryID), utils.Uint32ToBytes(ud.TargetID))

	return nil
}

func (k *Keeper) GetUnbondingDelegation(ctx sdk.Context, delegatorAddress string, target types.DelegationTarget) (types.UnbondingDelegation, bool) {
	switch target.(type) {
	case *poolstypes.Pool:
		return k.GetPoolUnbondingDelegation(ctx, target.GetID(), delegatorAddress)
	case *operatorstypes.Operator:
		return k.GetOperatorUnbondingDelegation(ctx, delegatorAddress, target.GetID())
	case *servicestypes.Service:
		return k.GetServiceUnbondingDelegation(ctx, delegatorAddress, target.GetID())
	default:
		return types.UnbondingDelegation{}, types.ErrInvalidDelegationTarget
	}
}

// --------------------------------------------------------------------------------------------------------------------

// PerformDelegation performs a delegation of the given amount from the delegator to the receiver.
// It sends the coins to the receiver address and updates the delegation object and returns the new
// shares of the delegation.
// NOTE: This is done so that if we implement other delegation types in the future we can have a single
// function that performs common operations for all of them.
func (k *Keeper) PerformDelegation(ctx sdk.Context, data types.DelegationData) (sdk.DecCoins, error) {
	// Get the data
	receiver := data.Target
	delegator := data.Delegator
	hooks := data.Hooks

	// In some situations, the exchange rate becomes invalid, e.g. if
	// the receives loses all tokens due to slashing. In this case,
	// make all future delegations invalid.
	if receiver.InvalidExRate() {
		return nil, types.ErrDelegatorShareExRateInvalid
	}

	// Get or create the delegation object and call the appropriate hook if present
	delegation, found := k.GetDelegationForTarget(ctx, receiver, delegator)

	if found {
		// Delegation was found
		err := hooks.BeforeDelegationSharesModified(ctx, receiver.GetID(), delegator)
		if err != nil {
			return nil, err
		}
	} else {
		// Delegation was not found
		delegation = data.BuildDelegation(receiver.GetID(), delegator, sdk.NewDecCoins())
		err := hooks.BeforeDelegationCreated(ctx, receiver.GetID(), delegator)
		if err != nil {
			return nil, err
		}
	}

	// Convert the addresses to sdk.AccAddress
	delegatorAddress, err := k.accountKeeper.AddressCodec().StringToBytes(delegator)
	if err != nil {
		return nil, err
	}
	receiverAddress, err := k.accountKeeper.AddressCodec().StringToBytes(receiver.GetAddress())
	if err != nil {
		return nil, err
	}

	// Send the coins to the receiver address
	err = k.bankKeeper.SendCoins(ctx, delegatorAddress, receiverAddress, data.Amount)
	if err != nil {
		return nil, err
	}

	// Update the delegation
	newShares, err := data.UpdateDelegation(ctx, delegation)
	if err != nil {
		return nil, err
	}

	// Call the after-modification hook
	err = hooks.AfterDelegationModified(ctx, receiver.GetID(), delegator)
	if err != nil {
		return nil, err
	}

	return newShares, nil
}
