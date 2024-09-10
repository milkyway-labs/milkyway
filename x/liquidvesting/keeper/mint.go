package keeper

import (
	"slices"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

// IsMinter tells if a user have the permissions to mint tokens.
func (k *Keeper) IsMinter(ctx sdk.Context, user sdk.AccAddress) (bool, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return false, err
	}

	return slices.Contains(params.Minters, user.String()), nil
}

// MintStakingRepresentation mints the staking representation of the provided
// amount to the user.
func (k *Keeper) MintStakingRepresentation(
	ctx sdk.Context,
	user sdk.AccAddress,
	amount sdk.Coins,
) error {
	var toMintTokens sdk.Coins
	for _, coin := range amount {
		newTokenDenom, err := types.GetVestedRepresentationDenom(coin.Denom)
		if err != nil {
			return err
		}
		toMintTokens = append(toMintTokens, sdk.NewCoin(newTokenDenom, coin.Amount))
	}

	// Mint the tokens to the module
	err := k.BankKeeper.MintCoins(goCtx, types.ModuleName, toMintTokens)
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
