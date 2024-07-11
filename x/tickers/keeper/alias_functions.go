package keeper

import (
	"context"

	"cosmossdk.io/collections"

	"github.com/milkyway-labs/milkyway/x/tickers/types"
)

func (k *Keeper) SetTicker(ctx context.Context, denom, ticker string) error {
	if err := k.Tickers.Set(ctx, denom, ticker); err != nil {
		return err
	}
	if err := k.TickerIndexes.Set(ctx, collections.Join(ticker, denom)); err != nil {
		return err
	}
	return nil
}

func (k *Keeper) GetTicker(ctx context.Context, denom string) (string, error) {
	ticker, err := k.Tickers.Get(ctx, denom)
	if err != nil {
		return "", types.ErrTickerNotFound
	}
	return ticker, nil
}

func (k *Keeper) RemoveTicker(ctx context.Context, denom string) error {
	ticker, err := k.GetTicker(ctx, denom)
	if err != nil {
		return err
	}
	if err := k.Tickers.Remove(ctx, denom); err != nil {
		return err
	}
	return k.TickerIndexes.Remove(ctx, collections.Join(ticker, denom))
}
