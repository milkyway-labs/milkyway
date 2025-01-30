package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v8/app/testutil"
	"github.com/milkyway-labs/milkyway/v8/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/v8/x/operators/types"
	"github.com/milkyway-labs/milkyway/v8/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v8/x/services/types"
)

func (suite *KeeperTestSuite) TestKeeper_ValidateRestakeRestakingCap() {
	ctx := suite.ctx

	// Register MILK and TIA.
	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, sdkmath.LegacyNewDec(5))
	suite.RegisterCurrency(ctx, "utia", "TIA", 6, sdkmath.LegacyNewDec(8))

	// Delegate $500 worth of MILK to pool.
	addr1 := testutil.TestAddress(1).String()
	suite.fundAccount(ctx, addr1, utils.MustParseCoins("100_000000umilk"))
	_, err := suite.k.DelegateToPool(ctx, utils.MustParseCoin("100_000000umilk"), addr1)
	suite.Require().NoError(err)

	// Create an operator.
	err = suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
		1, operatorstypes.OPERATOR_STATUS_ACTIVE,
		"MilkyWay Operator",
		"https://milkyway.com",
		"https://milkyway.com/picture",
		"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
	))
	suite.Require().NoError(err)

	// Delegate $1600 worth of TIA to the operator.
	addr2 := testutil.TestAddress(2).String()
	suite.fundAccount(ctx, addr2, utils.MustParseCoins("200_000000utia"))
	_, err = suite.k.DelegateToOperator(ctx, 1, utils.MustParseCoins("200_000000utia"), addr2)
	suite.Require().NoError(err)

	// Create a service.
	err = suite.sk.SaveService(ctx, servicestypes.NewService(
		1, servicestypes.SERVICE_STATUS_ACTIVE,
		"MilkyWay",
		"MilkyWay is a restaking platform",
		"https://milkyway.com",
		"https://milkyway.com/logo.png",
		"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
		false,
	))
	suite.Require().NoError(err)

	// Delegate $1500 worth of MILK to the service.
	addr3 := testutil.TestAddress(3).String()
	suite.fundAccount(ctx, addr3, utils.MustParseCoins("300_000000umilk"))
	_, err = suite.k.DelegateToService(ctx, 1, utils.MustParseCoins("300_000000umilk"), addr3)
	suite.Require().NoError(err)

	// Now we have total $3600 worth of assets restaked.

	// Fund the restaker account.
	restaker := testutil.TestAddress(4).String()
	suite.fundAccount(ctx, restaker, utils.MustParseCoins("1000_000000umilk,1000_000000utia"))

	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		amount    sdk.Coins
		shouldErr bool
	}{
		{
			name:      "no restaking cap returns no error",
			amount:    utils.MustParseCoins("100_000000umilk,100_000000utia"),
			shouldErr: false,
		},
		{
			name: "exceeding restaking cap returns an error",
			store: func(ctx sdk.Context) {
				err = suite.k.SetParams(ctx, types.NewParams(
					types.DefaultUnbondingTime,
					nil,
					sdkmath.LegacyNewDec(5000),
				))
				suite.Require().NoError(err)
			},
			amount:    utils.MustParseCoins("300_000000umilk"), // $1500
			shouldErr: true,
		},
		{
			name: "not exceeding restaking cap returns no error",
			store: func(ctx sdk.Context) {
				// Set restaking cap
				err := suite.k.SetParams(ctx, types.NewParams(
					types.DefaultUnbondingTime, nil, sdkmath.LegacyNewDec(5000)),
				)
				suite.Require().NoError(err)
			},
			amount:    utils.MustParseCoins("175_000000utia"), // $1400
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			// We can pass nil DelegationTarget since we're not testing restaking
			// restriction.
			err = suite.k.ValidateRestake(ctx, restaker, tc.amount, nil)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
