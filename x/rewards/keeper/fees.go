package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// PayRegistrationFees pays the registration fees for the user and sends the funds to the community pool.
// If there are multiple fee denoms set inside the params, the user is requested to have enough balance
// of at least one of them.
func (k *Keeper) PayRegistrationFees(ctx context.Context, user string) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	// If there's nothing to pay, do nothing
	if params.RewardsPlanCreationFee.IsZero() {
		return nil
	}

	// Parse the user address
	userAddress, err := k.accountKeeper.AddressCodec().StringToBytes(user)
	if err != nil {
		return err
	}

	// Get the user's balance
	balance := k.bankKeeper.GetAllBalances(ctx, userAddress)

	for _, feeCoin := range params.RewardsPlanCreationFee {
		// Check if the user has enough balance of this fee denom
		if balance.AmountOf(feeCoin.Denom).GTE(feeCoin.Amount) {
			// Pay the fee with this denom
			return k.communityPoolKeeper.FundCommunityPool(ctx, sdk.NewCoins(feeCoin), userAddress)
		}
	}

	return errors.Wrap(sdkerrors.ErrInsufficientFunds, "not enough balance to pay the registration fees")
}
