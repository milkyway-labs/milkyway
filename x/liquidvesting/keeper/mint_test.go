package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

func (suite *KeeperTestSuite) TestKeepr_MintVestedRepresentation() {
	testCases := []struct {
		name      string
		setup     func()
		to        string
		amount    sdk.Coins
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name:      "can't mint vested representation of a vested representation",
			to:        "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			amount:    sdk.NewCoins(sdk.NewInt64Coin("vested/stake", 1000)),
			shouldErr: true,
		},
		{
			name:      "mint properly",
			to:        "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			amount:    sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				denom, err := types.GetVestedRepresentationDenom("stake")
				suite.Assert().NoError(err)
				accAddr := sdk.MustAccAddressFromBech32("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd")
				balance := suite.bk.GetBalance(ctx, accAddr, "vested/stake")
				suite.Assert().Equal(sdk.NewInt64Coin(denom, 1000), balance)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.setup != nil {
				tc.setup()
			}
			accAddr := sdk.MustAccAddressFromBech32(tc.to)
			_, err := suite.k.MintVestedRepresentation(ctx, accAddr, tc.amount)

			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				tc.check(ctx)
			}
		})
	}
}
