package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"

	"github.com/milkyway-labs/milkyway/x/rewards/types"
	tickerstypes "github.com/milkyway-labs/milkyway/x/tickers/types"
)

func (k *Keeper) GetPrice(ctx context.Context, denom string) (math.LegacyDec, error) {
	ticker, err := k.tickersKeeper.GetTicker(ctx, denom)
	if err != nil {
		// If ticker is not found, then we return 0 as price.
		if errors.IsOf(err, tickerstypes.ErrTickerNotFound) {
			return math.LegacyZeroDec(), nil
		}
		return math.LegacyDec{}, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cp := slinkytypes.NewCurrencyPair(ticker, types.USDTicker)
	qpn, err := k.oracleKeeper.GetPriceWithNonceForCurrencyPair(sdkCtx, cp)
	if err != nil {
		// If currency pair is not found return 0 as well.
		if errors.IsOf(err, collections.ErrNotFound) {
			return math.LegacyZeroDec(), nil
		}
		return math.LegacyDec{}, err
	}
	decimals, err := k.oracleKeeper.GetDecimalsForCurrencyPair(sdkCtx, cp)
	if err != nil {
		return math.LegacyDec{}, err
	}

	// Divide returned quote price by 10^{decimals} gives us the real price in
	// decimal number.
	return math.LegacyNewDecFromIntWithPrec(qpn.Price, int64(decimals)), nil
}

func (k *Keeper) GetCoinValue(ctx context.Context, coin sdk.Coin) (math.LegacyDec, error) {
	price, err := k.GetPrice(ctx, coin.Denom)
	if err != nil {
		return math.LegacyDec{}, err
	}
	return price.MulInt(coin.Amount), nil
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
