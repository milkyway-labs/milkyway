package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

// PerformDelegation performs a delegation of the given amount from the delegator to the receiver.
// It sends the coins to the receiver address and updates the delegation object and returns the new
// shares of the delegation.
// NOTE: This is done so that if we implement other delegation types in the future we can have a single
// function that performs common operations for all of them.
func (k *Keeper) PerformDelegation(ctx sdk.Context, data types.DelegationData) (sdk.DecCoins, error) {
	// Get the data
	receiver := data.Receiver
	delegator := data.Delegator
	hooks := data.Hooks

	// In some situations, the exchange rate becomes invalid, e.g. if
	// the receives loses all tokens due to slashing. In this case,
	// make all future delegations invalid.
	if receiver.InvalidExRate() {
		return nil, types.ErrDelegatorShareExRateInvalid
	}

	// Get or create the delegation object and call the appropriate hook if present
	delegation, found := data.GetDelegation(ctx, receiver.GetID(), delegator)

	if found {
		// Delegation was found
		err := hooks.BeforeDelegationSharesModified(ctx, receiver.GetID(), delegator)
		if err != nil {
			return nil, err
		}
	} else {
		// Delegation was not found
		delegation = data.BuildDelegation(receiver.GetID(), delegator)
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
