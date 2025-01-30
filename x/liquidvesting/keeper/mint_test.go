package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v8/x/liquidvesting/types"
)

func (suite *KeeperTestSuite) TestKeeper_MintLockedRepresentation() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		to        string
		amount    sdk.Coins
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name:      "can't mint locked representation of a locked representation",
			to:        "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			amount:    sdk.NewCoins(sdk.NewInt64Coin("locked/stake", 1000)),
			shouldErr: true,
		},
		{
			name:      "mint properly",
			to:        "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			amount:    sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				denom, err := types.GetLockedRepresentationDenom("stake")
				suite.Assert().NoError(err)
				accAddr := sdk.MustAccAddressFromBech32("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd")
				balance := suite.bk.GetBalance(ctx, accAddr, "locked/stake")
				suite.Assert().Equal(sdk.NewInt64Coin(denom, 1000), balance)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			userAddr, err := sdk.AccAddressFromBech32(tc.to)
			suite.Require().NoError(err)

			_, err = suite.k.MintLockedRepresentation(ctx, userAddr, tc.amount)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				tc.check(ctx)
			}
		})
	}
}
