package v2_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	poolskeeper "github.com/milkyway-labs/milkyway/x/pools/keeper"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/services/keeper"
	v2 "github.com/milkyway-labs/milkyway/x/services/legacy/v2"
	"github.com/milkyway-labs/milkyway/x/services/testutils"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

func TestMigrateStoreTestSuite(t *testing.T) {
	suite.Run(t, new(MigrateStoreTestSuite))
}

type MigrateStoreTestSuite struct {
	suite.Suite

	ctx sdk.Context

	keeper      *keeper.Keeper
	poolsKeeper *poolskeeper.Keeper
}

func (suite *MigrateStoreTestSuite) SetupTest() {
	data := testutils.NewKeeperTestData(suite.T())

	suite.ctx = data.Context
	suite.keeper = data.Keeper
	suite.poolsKeeper = data.PoolsKeeper
}

func (suite *MigrateStoreTestSuite) TestMigrateStore() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "no services to set as accredited",
			store: func(ctx sdk.Context) {
				// Create a service
				err := suite.keeper.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)

				// Set the default params (empty list of allowed services)
				suite.poolsKeeper.SetParams(ctx, poolstypes.DefaultParams())
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				// Get the service
				service, found, err := suite.keeper.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(found)

				// Check that the service has not been set as accredited
				suite.Require().False(service.Accredited)
			},
		},
		{
			name: "non existing service is skipped",
			store: func(ctx sdk.Context) {
				// Set the default params (empty list of allowed services)
				suite.poolsKeeper.SetParams(ctx, poolstypes.Params{
					AllowedServicesIDs: []uint32{1},
				})
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				// Get the service
				_, found, err := suite.keeper.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().False(found)
			},
		},
		{
			name: "service is set as accredited",
			store: func(ctx sdk.Context) {
				// Create a service
				err := suite.keeper.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)

				// Set the default params (empty list of allowed services)
				suite.poolsKeeper.SetParams(ctx, poolstypes.Params{
					AllowedServicesIDs: []uint32{1},
				})
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				// Get the service
				service, found, err := suite.keeper.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(found)

				// Check that the service has been set as accredited
				suite.Require().True(service.Accredited)
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

			err := v2.MigrateStore(ctx, suite.keeper, suite.poolsKeeper)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
