package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

func (suite *KeeperTestSuite) TestKeeper_SetNextServiceID() {
	testCases := []struct {
		name  string
		store func(ctx sdk.Context)
		id    uint32
		check func(ctx sdk.Context)
	}{
		{
			name: "next service id is saved correctly",
			id:   1,
			check: func(ctx sdk.Context) {
				store := ctx.KVStore(suite.storeKey)
				serviceIDBz := store.Get(types.NextServiceIDKey)
				suite.Require().Equal(uint32(1), types.GetServiceIDFromBytes(serviceIDBz))
			},
		},
		{
			name: "next service id is overridden properly",
			store: func(ctx sdk.Context) {
				suite.k.SetNextServiceID(ctx, 1)
			},
			id: 2,
			check: func(ctx sdk.Context) {
				store := ctx.KVStore(suite.storeKey)
				serviceIDBz := store.Get(types.NextServiceIDKey)
				suite.Require().Equal(uint32(2), types.GetServiceIDFromBytes(serviceIDBz))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			suite.k.SetNextServiceID(ctx, tc.id)
			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_GetNextServiceID() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		shouldErr bool
		expNext   uint32
	}{
		{
			name:      "non existing next service returns error",
			shouldErr: true,
		},
		{
			name: "exiting next service id is returned properly",
			store: func(ctx sdk.Context) {
				suite.k.SetNextServiceID(ctx, 1)
			},
			expNext: 1,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			next, err := suite.k.GetNextServiceID(ctx)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expNext, next)
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *KeeperTestSuite) TestKeeper_CreateService() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		service   types.Service
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "user without enough funds to pay for registration fees returns error",
			store: func(ctx sdk.Context) {
				// Set the params
				suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
				))

				// Fund the user account
				userBalance := sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(50_000_000)))
				suite.fundAccount(ctx, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", userBalance)
			},
			service: types.NewService(
				1,
				types.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "service is created properly",
			store: func(ctx sdk.Context) {
				// Set the params
				suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
				))

				// Fund the user account
				userBalance := sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(150_000_000)))
				suite.fundAccount(ctx, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", userBalance)
			},
			service: types.NewService(
				1,
				types.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				// Make sure the user balance has been reduced
				userAddress, err := sdk.AccAddressFromBech32("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd")
				suite.Require().NoError(err)

				userBalance := suite.bk.GetBalance(ctx, userAddress, "uatom")
				suite.Require().Equal(sdk.NewCoin("uatom", sdkmath.NewInt(50_000_000)), userBalance)

				// Make sure the community pool has been funded
				poolBalance := suite.bk.GetBalance(ctx, authtypes.NewModuleAddress(authtypes.FeeCollectorName), "uatom")
				suite.Require().Equal(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000)), poolBalance)

				// Make sure the service has been created
				service, found := suite.k.GetService(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				), service)

				// Make sure the hook was called
				suite.Require().True(suite.hooks.CalledMap["AfterServiceCreated"])
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.setup != nil {
				tc.setup()
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			err := suite.k.CreateService(ctx, tc.service)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}
