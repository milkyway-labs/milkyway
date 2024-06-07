package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CommunityPoolKeeper interface {
	FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error
}
