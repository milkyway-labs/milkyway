package keeper_test

import (
	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v11/x/services/types"
)

func (suite *KeeperTestSuite) TestKeeper_SetNextServiceID() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		id        uint32
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name:      "next service id is saved correctly",
			id:        1,
			shouldErr: false,
			check: func(ctx sdk.Context) {
				nextServiceID, err := suite.k.GetNextServiceID(ctx)
				suite.Require().NoError(err)
				suite.Require().EqualValues(1, nextServiceID)
			},
		},
		{
			name: "next service id is overridden properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextServiceID(ctx, 1)
				suite.Require().NoError(err)
			},
			id:        2,
			shouldErr: false,
			check: func(ctx sdk.Context) {
				nextServiceID, err := suite.k.GetNextServiceID(ctx)
				suite.Require().NoError(err)
				suite.Require().EqualValues(2, nextServiceID)
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

			err := suite.k.SetNextServiceID(ctx, tc.id)
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

func (suite *KeeperTestSuite) TestKeeper_GetNextServiceID() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		shouldErr bool
		expNext   uint32
	}{
		{
			name:      "non existing next service returns 1",
			shouldErr: false,
			expNext:   1,
		},
		{
			name: "exiting next service id is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextServiceID(ctx, 1)
				suite.Require().NoError(err)
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
			name: "service is created properly",
			store: func(ctx sdk.Context) {
				// Set the params
				err := suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
				))
				suite.Require().NoError(err)

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
				false,
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				// Make sure the service account has been created
				hasAccount := suite.ak.HasAccount(ctx, types.GetServiceAddress(1))
				suite.Require().True(hasAccount)

				// Make sure the service has been created
				service, err := suite.k.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
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

func (suite *KeeperTestSuite) TestKeeper_ActivateService() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		serviceID uint32
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name:      "service not found returns error",
			serviceID: 1,
			shouldErr: true,
		},
		{
			name: "already active service returns error",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			serviceID: 1,
			shouldErr: true,
		},
		{
			name: "service is activated properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			serviceID: 1,
			shouldErr: false,
			check: func(ctx sdk.Context) {
				service, err := suite.k.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				), service)

				// Make sure the hook was called
				suite.Require().True(suite.hooks.CalledMap["AfterServiceActivated"])
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

			err := suite.k.ActivateService(ctx, tc.serviceID)
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

func (suite *KeeperTestSuite) TestKeeper_DeactivateService() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		serviceID uint32
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name:      "service not found returns error",
			serviceID: 1,
			shouldErr: true,
		},
		{
			name: "inactive service returns error",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			serviceID: 1,
			shouldErr: true,
		},
		{
			name: "service is deactivated properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			serviceID: 1,
			shouldErr: false,
			check: func(ctx sdk.Context) {
				service, err := suite.k.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewService(
					1,
					types.SERVICE_STATUS_INACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				), service)

				// Make sure the hook was called
				suite.Require().True(suite.hooks.CalledMap["AfterServiceDeactivated"])
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

			err := suite.k.DeactivateService(ctx, tc.serviceID)
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

func (suite *KeeperTestSuite) TestKeeper_SetServiceAccreditation() {
	testCases := []struct {
		name       string
		setup      func()
		store      func(ctx sdk.Context)
		serviceID  uint32
		accredited bool
		shouldErr  bool
		check      func(ctx sdk.Context)
	}{
		{
			name:      "service not found returns error",
			serviceID: 1,
			shouldErr: true,
		},
		{
			name: "service's accreditation doesn't change",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			serviceID:  1,
			accredited: false,
			shouldErr:  false,
			check: func(ctx sdk.Context) {
				// Accreditation didn't change
				service, err := suite.k.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().False(service.Accredited)

				// Make sure the hook wasn't called
				suite.Require().False(suite.hooks.CalledMap["AfterServiceAccreditationModified"])
			},
		},
		{
			name: "service's accreditation changes properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			serviceID:  1,
			accredited: true,
			shouldErr:  false,
			check: func(ctx sdk.Context) {
				// Accreditation changed
				service, err := suite.k.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(service.Accredited)

				// Make sure the hook was called
				suite.Require().True(suite.hooks.CalledMap["AfterServiceAccreditationModified"])
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

			err := suite.k.SetServiceAccredited(ctx, tc.serviceID, tc.accredited)
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

func (suite *KeeperTestSuite) TestKeeper_GetService() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		serviceID  uint32
		expFound   bool
		expService types.Service
	}{
		{
			name:      "service not found returns false",
			serviceID: 1,
			expFound:  false,
		},
		{
			name: "service is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			serviceID: 1,
			expFound:  true,
			expService: types.NewService(
				1,
				types.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				false,
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			service, err := suite.k.GetService(ctx, tc.serviceID)
			if !tc.expFound {
				suite.Require().ErrorIs(err, collections.ErrNotFound)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expService, service)
			}
		})
	}
}
