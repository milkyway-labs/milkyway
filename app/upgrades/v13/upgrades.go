package v13

import (
	"context"
	"fmt"
	"slices"
	"time"

	"cosmossdk.io/collections"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/v13/app/keepers"
	"github.com/milkyway-labs/milkyway/v13/utils"
)

const (
	// Our address to transfer all IBC/milkBABY assets to
	ourAddress    = "milk1aeetghumz0crffr990vah06vk2rj4lfe7a94ep"
	milkBABYDenom = "factory/milk1qg5ega6dykkxc307y25pecuufrjkxkaggkkxh7nad0vhyhtuhw3ssgcye4/umilkBBN"
)

var (
	// Addresses to exclude from transfer to our address
	transferExcludedAddresses = []string{
		// IBC escrow addresses; we don't touch these addresses
		"milk1a53udazy8ayufvy0s434pfwjcedzqv34clymru",
		"milk1kq2rzz6fq2q7fsu75a9g7cpzjeanmk68xkrpnx",
		"milk12k2pyuylm9t7ugdvz67h9pg4gmmvhn5v34rvex",
		"milk1krtrt5dvhnw5j46ljgah7q5f9pm3vvakds8g63",
		"milk17sael2kcmm8npe2pmkxj3un90xfg60vv6xznt0",
		"milk16h2ynrzwhxgjnd0hswkvdvq9nav9kklqzl044t",
	}
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configuration module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		vm, err := mm.RunMigrations(ctx, configuration, fromVM)
		if err != nil {
			return nil, err
		}

		// We only need to care about pools since there's no operators or services.
		dels, err := keepers.RestakingKeeper.GetAllPoolDelegations(ctx)
		if err != nil {
			return nil, err
		}
		for _, del := range dels {
			// No slashing ever happened, so we can simply convert shares to tokens.
			tokens, _ := del.Shares.TruncateDecimal()
			if len(tokens) != 1 { // sanity check
				return nil, fmt.Errorf("expected only 1 token, got %s", tokens)
			}
			tokenDenom := utils.GetTokenDenomFromSharesDenom(tokens[0].Denom)
			_, err = keepers.RestakingKeeper.UndelegateFromPool(ctx, sdk.NewCoin(tokenDenom, tokens[0].Amount), del.UserAddress)
			if err != nil {
				return nil, err
			}
		}

		// Prepare a temporary context with a future block time and forcefully complete
		// all unbondings.
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		restakingParams, err := keepers.RestakingKeeper.GetParams(ctx)
		if err != nil {
			return nil, err
		}
		futureCtx := sdkCtx.WithBlockTime(sdkCtx.BlockTime().Add(restakingParams.UnbondingTime).Add(time.Minute))
		err = keepers.RestakingKeeper.CompleteMatureUnbondingDelegations(futureCtx)
		if err != nil {
			return nil, err
		}

		// Move all IBC assets and milkBABY from all addresses to our address
		var denomsToMove []string
		denomsToMove = append(denomsToMove, milkBABYDenom)
		ourAddr := sdk.MustAccAddressFromBech32(ourAddress)

		// Get all IBC denoms available
		denomTraces := keepers.TransferKeeper.GetAllDenomTraces(sdkCtx)
		for _, denomTrace := range denomTraces {
			// Construct IBC denom from denom trace
			ibcDenom := fmt.Sprintf("ibc/%s", denomTrace.Hash().String())
			denomsToMove = append(denomsToMove, ibcDenom)
		}

		for _, denom := range denomsToMove {
			// Get all addresses that own the denom in their balance
			var denomOwners []sdk.AccAddress
			err = keepers.BankKeeper.Balances.Indexes.Denom.Walk(
				ctx,
				collections.NewPrefixedPairRange[string, sdk.AccAddress](denom),
				func(_ string, addr sdk.AccAddress) (stop bool, err error) {
					denomOwners = append(denomOwners, addr)
					return false, nil
				})
			if err != nil {
				return nil, err
			}

			// Send assets to our address
			for _, denomOwner := range denomOwners {
				owner, err := keepers.AccountKeeper.AddressCodec().BytesToString(denomOwner)
				if err != nil {
					return nil, err
				}
				// If denom owner is one of the excluded addresses, skip
				if slices.Contains(transferExcludedAddresses, owner) {
					continue
				}
				balance := keepers.BankKeeper.GetBalance(ctx, denomOwner, denom)
				err = keepers.BankKeeper.SendCoins(ctx, denomOwner, ourAddr, sdk.NewCoins(balance))
				if err != nil {
					return nil, err
				}
			}
		}

		return vm, nil
	}
}
