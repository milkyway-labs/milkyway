package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Type that allows to compose multiple BankHooks
type ComposedBankHooks struct {
	h1 BankHooks
	h2 BankHooks
}

var _ BankHooks = &ComposedBankHooks{}

// NewComposedBankHooks creates a new composed BankHooks
func NewComposedBankHooks(h1 BankHooks, h2 BankHooks) *ComposedBankHooks {
	return &ComposedBankHooks{
		h1: h1,
		h2: h2,
	}
}

func (h *ComposedBankHooks) TrackBeforeSend(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins) {
	h.h1.TrackBeforeSend(ctx, from, to, amount)
	h.h2.TrackBeforeSend(ctx, from, to, amount)
}

func (h *ComposedBankHooks) BlockBeforeSend(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins) error {
	if err := h.h1.BlockBeforeSend(ctx, from, to, amount); err != nil {
		return err
	}
	return h.h2.BlockBeforeSend(ctx, from, to, amount)
}
