package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v10/x/liquidvesting/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v10/x/restaking/types"
)

func (k *Keeper) RestakeRestrictionFn(ctx context.Context, restakerAddress string, amount sdk.Coins, _ restakingtypes.DelegationTarget) error {
	restakedLockedRepresentation := sdk.NewCoins()
	for _, coin := range amount {
		if types.IsLockedRepresentationDenom(coin.Denom) {
			restakedLockedRepresentation = restakedLockedRepresentation.Add(coin)
		}
	}

	// The user is not restaking any locked representation
	// allow the restake operation
	if restakedLockedRepresentation.IsZero() {
		return nil
	}

	// We have some locked representation to restake, ensure the user
	// insurance fund can cover the amount

	// Get the current locked representations that are being covered
	// by the insurance fund
	totalLockedRepresentations, err := k.GetAllUserActiveLockedRepresentations(ctx, restakerAddress)
	if err != nil {
		return err
	}

	// Add the locked representations that are being restaked to the total
	totalLockedRepresentations = totalLockedRepresentations.Add(sdk.NewDecCoinsFromCoins(restakedLockedRepresentation...)...)

	// Get the user's insurance fund
	userInsuranceFund, err := k.GetUserInsuranceFund(ctx, restakerAddress)
	if err != nil {
		return nil
	}

	// Get the module params
	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	// Check if the insurance fund can cover the newly restaked amount
	canCover, _, err := userInsuranceFund.CanCoverDecCoins(params.InsurancePercentage, totalLockedRepresentations)
	if err != nil {
		return err
	}

	if !canCover {
		return types.ErrInsufficientBalance
	}

	return nil
}
