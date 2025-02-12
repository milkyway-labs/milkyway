package keeper

import (
	"context"
	stdmath "math"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"

	assetstypes "github.com/milkyway-labs/milkyway/v9/x/assets/types"
	"github.com/milkyway-labs/milkyway/v9/x/rewards/types"
)

// These code snippets are copied from x/rewards/keeper/oracle.go
// TODO: remove redundant code

func (k *Keeper) GetAssetAndPrice(ctx context.Context, denom string) (assetstypes.Asset, math.LegacyDec, error) {
	asset, err := k.assetsKeeper.GetAsset(ctx, denom)
	if err != nil {
		// If asset is not found, then we return 0 as price.
		if errors.IsOf(err, collections.ErrNotFound) {
			return assetstypes.Asset{}, math.LegacyZeroDec(), nil
		}
		return assetstypes.Asset{}, math.LegacyDec{}, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cp := connecttypes.NewCurrencyPair(asset.Ticker, types.USDTicker)
	qpn, err := k.oracleKeeper.GetPriceWithNonceForCurrencyPair(sdkCtx, cp)
	if err != nil {
		// If currency pair is not found return 0 as well.
		if errors.IsOf(err, collections.ErrNotFound) {
			return asset, math.LegacyZeroDec(), nil
		}
		return asset, math.LegacyDec{}, err
	}

	decimals, err := k.oracleKeeper.GetDecimalsForCurrencyPair(sdkCtx, cp)
	if err != nil {
		return asset, math.LegacyDec{}, err
	}

	// Divide returned quote price by 10^{decimals} gives us the real price in
	// decimal number.
	return asset, math.LegacyNewDecFromIntWithPrec(qpn.Price, int64(decimals)), nil
}

func (k *Keeper) GetCoinValue(ctx context.Context, coin sdk.Coin) (math.LegacyDec, error) {
	asset, price, err := k.GetAssetAndPrice(ctx, coin.Denom)
	if err != nil {
		return math.LegacyDec{}, err
	}

	if price.IsZero() {
		return math.LegacyZeroDec(), nil
	}

	return price.MulInt(coin.Amount).QuoInt64(int64(stdmath.Pow10(int(asset.Exponent)))), nil
}

func (k *Keeper) GetCoinsValue(ctx context.Context, coins sdk.Coins) (math.LegacyDec, error) {
	totalValue := math.LegacyZeroDec()
	for _, coin := range coins {
		value, err := k.GetCoinValue(ctx, coin)
		if err != nil {
			return math.LegacyDec{}, err
		}
		totalValue = totalValue.Add(value)
	}
	return totalValue, nil
}
