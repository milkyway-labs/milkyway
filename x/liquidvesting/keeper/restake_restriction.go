package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
)

func (k *Keeper) RestakeRestrictionFn(ctx context.Context, restakerAddress string, amount sdk.Coins, _ restakingtypes.DelegationTarget) error {
	restakedVestedRepresentation := sdk.NewCoins()
	for _, coin := range amount {
		if types.IsVestedRepresentationDenom(coin.Denom) {
			restakedVestedRepresentation = restakedVestedRepresentation.Add(coin)
		}
	}

	// The user is not restaking any vested representation
	// allow the restake operation
	if restakedVestedRepresentation.IsZero() {
		return nil
	}

	// We have some vested representation to restake, ensure the user
	// insurance fund can cover the amount

	// Get the current vested representations that are being coverd
	// by the insurance fund
	totalVestedRepresentations, err := k.GetAllUserActiveVestedRepresentations(ctx, restakerAddress)
	if err != nil {
		return err
	}

	// Add the vested representations that are being restaked to the total
	totalVestedRepresentations = totalVestedRepresentations.Add(sdk.NewDecCoinsFromCoins(restakedVestedRepresentation...)...)

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
	canCover, _, err := userInsuranceFund.CanCoverDecCoins(params.InsurancePercentage, totalVestedRepresentations)
	if err != nil {
		return err
	}

	if !canCover {
		return types.ErrInsufficientBalance
	}

	return nil
}
