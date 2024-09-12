package keeper

import (
	"slices"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

// IsMinter tells if a user have the permissions to mint tokens.
func (k *Keeper) IsMinter(ctx sdk.Context, user sdk.AccAddress) (bool, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return false, err
	}

	return slices.Contains(params.Minters, user.String()), nil
}

// MintVestedRepresentation mints the vested staked representation of the provided
// amount to the user.
func (k *Keeper) MintVestedRepresentation(
	ctx sdk.Context,
	user sdk.AccAddress,
	amount sdk.Coins,
) error {
	var toMintTokens sdk.Coins
	for _, coin := range amount {
		// Create the vested representation for the received denom
		vestedRepresentationDenom, err := types.GetVestedRepresentationDenom(coin.Denom)
		if err != nil {
			return err
		}

		// Check if we have the metadata for the vested representation
		_, vestedDenomMetadataFound := k.BankKeeper.GetDenomMetaData(ctx, vestedRepresentationDenom)
		if !vestedDenomMetadataFound {
			// We don't have the metadata for the vested representation
			// we should create it
			denomMetadata := banktypes.Metadata{
				DenomUnits: []*banktypes.DenomUnit{{
					Denom:    vestedRepresentationDenom,
					Exponent: 0,
				}},
				Base:        vestedRepresentationDenom,
				Name:        vestedRepresentationDenom,
				Symbol:      vestedRepresentationDenom,
				Display:     vestedRepresentationDenom,
				Description: "Vested representation of " + vestedRepresentationDenom,
			}
			k.BankKeeper.SetDenomMetaData(ctx, denomMetadata)
		}

		toMintTokens = append(toMintTokens, sdk.NewCoin(vestedRepresentationDenom, coin.Amount))
	}

	// Mint the tokens to the module
	err := k.BankKeeper.MintCoins(ctx, types.ModuleName, toMintTokens)
	if err != nil {
		return err
	}

	// Transfer the minted tokens to the user
	return k.BankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		user,
		toMintTokens,
	)
}
