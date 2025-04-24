package keeper

import (
	"context"
	"slices"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/milkyway-labs/milkyway/v11/x/liquidvesting/types"
)

// IsMinter tells if a user have the permissions to mint tokens.
func (k *Keeper) IsMinter(ctx context.Context, user sdk.AccAddress) (bool, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return false, err
	}

	stringAddr, err := k.accountKeeper.AddressCodec().BytesToString(user)
	if err != nil {
		return false, err
	}

	return slices.Contains(params.Minters, stringAddr), nil
}

// MintLockedRepresentation mints the locked staked representation of the provided
// amount to the user.
func (k *Keeper) MintLockedRepresentation(ctx context.Context, user sdk.AccAddress, amount sdk.Coins) (sdk.Coins, error) {
	var toMintTokens sdk.Coins
	for _, coin := range amount {
		// Create the locked representation for the received denom
		lockedRepresentationDenom, err := types.GetLockedRepresentationDenom(coin.Denom)
		if err != nil {
			return sdk.Coins{}, err
		}

		// Check if we have the metadata for the locked representation
		_, lockedDenomMetadataFound := k.bankKeeper.GetDenomMetaData(ctx, lockedRepresentationDenom)
		if !lockedDenomMetadataFound {
			// We don't have the metadata for the locked representation
			// we should create it
			denomMetadata := banktypes.Metadata{
				DenomUnits: []*banktypes.DenomUnit{{
					Denom:    lockedRepresentationDenom,
					Exponent: 0,
				}},
				Base:        lockedRepresentationDenom,
				Name:        lockedRepresentationDenom,
				Symbol:      lockedRepresentationDenom,
				Display:     lockedRepresentationDenom,
				Description: "Locked representation of " + coin.Denom,
			}
			k.bankKeeper.SetDenomMetaData(ctx, denomMetadata)
		}

		toMintTokens = append(toMintTokens, sdk.NewCoin(lockedRepresentationDenom, coin.Amount))
	}

	// Mint the tokens to the module
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, toMintTokens)
	if err != nil {
		return sdk.Coins{}, err
	}

	// Transfer the minted tokens to the user
	err = k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		user,
		toMintTokens,
	)
	return toMintTokens, err
}
