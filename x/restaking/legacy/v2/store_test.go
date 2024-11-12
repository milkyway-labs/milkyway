package v2_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	operatorskeeper "github.com/milkyway-labs/milkyway/x/operators/keeper"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	restakingkeeper "github.com/milkyway-labs/milkyway/x/restaking/keeper"
	v2 "github.com/milkyway-labs/milkyway/x/restaking/legacy/v2"
	"github.com/milkyway-labs/milkyway/x/restaking/testutils"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	serviceskeeper "github.com/milkyway-labs/milkyway/x/services/keeper"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

func TestMigrationsTestSuite(t *testing.T) {
	suite.Run(t, new(MigrationsTestSuite))
}

type MigrationsTestSuite struct {
	suite.Suite

	ctx      sdk.Context
	storeKey storetypes.StoreKey
	cdc      codec.Codec

	restakingKeeper *restakingkeeper.Keeper
	operatorsKeeper *operatorskeeper.Keeper
	servicesKeeper  *serviceskeeper.Keeper
}

func (suite *MigrationsTestSuite) SetupTest() {
	data := testutils.NewKeeperTestData(suite.T())
	suite.storeKey = data.StoreKey
	suite.ctx = data.Context
	suite.cdc = data.Cdc

	suite.operatorsKeeper = data.OperatorsKeeper
	suite.servicesKeeper = data.ServicesKeeper
	suite.restakingKeeper = data.Keeper
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *MigrationsTestSuite) TestMigrateV1To2() {
	testCases := []struct {
		name      string
		setup     func(ctx sdk.Context)
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "non existing operators have their params deleted",
			setup: func(ctx sdk.Context) {
				sdkStore := ctx.KVStore(suite.storeKey)

				// Set the operator params
				paramsBz, err := suite.cdc.Marshal(&v2.LegacyOperatorParams{
					CommissionRate:    sdkmath.LegacyNewDec(100),
					JoinedServicesIDs: []uint32{1, 2, 3},
				})
				suite.Require().NoError(err)
				sdkStore.Set(v2.OperatorParamsStoreKey(1), paramsBz)

				paramsBz, err = suite.cdc.Marshal(&v2.LegacyOperatorParams{
					CommissionRate:    sdkmath.LegacyNewDec(200),
					JoinedServicesIDs: []uint32{4, 5, 6},
				})
				suite.Require().NoError(err)
				sdkStore.Set(v2.OperatorParamsStoreKey(2), paramsBz)
			},
			check: func(ctx sdk.Context) {
				// Make sure the params are deleted
				params, err := suite.operatorsKeeper.GetOperatorParams(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(operatorstypes.DefaultOperatorParams(), params)

				params, err = suite.operatorsKeeper.GetOperatorParams(ctx, 2)
				suite.Require().NoError(err)
				suite.Require().Equal(operatorstypes.DefaultOperatorParams(), params)

				// Make sure the list of joined services has been moved to the restaking keeper
				services, err := suite.restakingKeeper.GetAllOperatorsJoinedServices(ctx)
				suite.Require().NoError(err)
				suite.Require().Empty(services)
			},
		},
		{
			name: "existing operators params are migrated properly",
			setup: func(ctx sdk.Context) {
				sdkStore := ctx.KVStore(suite.storeKey)

				// Store the operators
				err := suite.operatorsKeeper.SaveOperator(ctx, operatorstypes.NewOperator(
					1,
					operatorstypes.OPERATOR_STATUS_INACTIVE,
					"MilkyWay", "", "", "admin",
				))
				suite.Require().NoError(err)

				err = suite.operatorsKeeper.SaveOperator(ctx, operatorstypes.NewOperator(
					2,
					operatorstypes.OPERATOR_STATUS_INACTIVE,
					"Cosmos", "", "", "admin",
				))
				suite.Require().NoError(err)

				// Set the operator params
				paramsBz, err := suite.cdc.Marshal(&v2.LegacyOperatorParams{
					CommissionRate:    sdkmath.LegacyNewDec(100),
					JoinedServicesIDs: []uint32{1, 2, 3},
				})
				suite.Require().NoError(err)
				sdkStore.Set(v2.OperatorParamsStoreKey(1), paramsBz)

				paramsBz, err = suite.cdc.Marshal(&v2.LegacyOperatorParams{
					CommissionRate:    sdkmath.LegacyNewDec(200),
					JoinedServicesIDs: []uint32{4, 5, 6},
				})
				suite.Require().NoError(err)
				sdkStore.Set(v2.OperatorParamsStoreKey(2), paramsBz)
			},
			check: func(ctx sdk.Context) {
				// Make sure the params are upgraded properly
				params, err := suite.operatorsKeeper.GetOperatorParams(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(operatorstypes.OperatorParams{
					CommissionRate: sdkmath.LegacyNewDec(100),
				}, params)

				params, err = suite.operatorsKeeper.GetOperatorParams(ctx, 2)
				suite.Require().NoError(err)
				suite.Require().Equal(operatorstypes.OperatorParams{
					CommissionRate: sdkmath.LegacyNewDec(200),
				}, params)

				// Make sure the list of joined services has been moved to the restaking keeper
				services, err := suite.restakingKeeper.GetAllOperatorsJoinedServices(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal([]restakingtypes.OperatorJoinedServices{
					{OperatorID: 1, ServiceIDs: []uint32{1, 2, 3}},
					{OperatorID: 2, ServiceIDs: []uint32{4, 5, 6}},
				}, services)
			},
		},
		{
			name: "non existing services have their params deleted",
			setup: func(ctx sdk.Context) {
				sdkStore := ctx.KVStore(suite.storeKey)

				// Set the service params
				paramsBz, err := suite.cdc.Marshal(&v2.LegacyServiceParams{
					WhitelistedOperatorsIDs: []uint32{1, 2, 3},
					WhitelistedPoolsIDs:     []uint32{4, 5, 6},
				})
				suite.Require().NoError(err)
				sdkStore.Set(v2.ServiceParamsStoreKey(1), paramsBz)

				paramsBz, err = suite.cdc.Marshal(&v2.LegacyServiceParams{
					WhitelistedOperatorsIDs: []uint32{7, 8, 9},
					WhitelistedPoolsIDs:     []uint32{10, 11, 12},
				})
				suite.Require().NoError(err)
				sdkStore.Set(v2.ServiceParamsStoreKey(2), paramsBz)
			},
			check: func(ctx sdk.Context) {
				// Make sure the list of whitelisted operators and pools has been moved to the restaking keeper
				pools, err := suite.restakingKeeper.GetAllServicesSecuringPools(ctx)
				suite.Require().NoError(err)
				suite.Require().Empty(pools)

				operators, err := suite.restakingKeeper.GetAllServicesAllowedOperators(ctx)
				suite.Require().NoError(err)
				suite.Require().Empty(operators)
			},
		},
		{
			name: "existing services params are migrated properly",
			setup: func(ctx sdk.Context) {
				sdkStore := ctx.KVStore(suite.storeKey)

				// Store the services
				err := suite.servicesKeeper.SaveService(ctx, servicestypes.NewService(
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

				err = suite.servicesKeeper.SaveService(ctx, servicestypes.Service{
					ID:      2,
					Address: servicestypes.GetServiceAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(150)),
					),
				})
				suite.Require().NoError(err)

				// Set the service params
				paramsBz, err := suite.cdc.Marshal(&v2.LegacyServiceParams{
					WhitelistedOperatorsIDs: []uint32{1, 2, 3},
					WhitelistedPoolsIDs:     []uint32{4, 5, 6},
				})
				suite.Require().NoError(err)
				sdkStore.Set(v2.ServiceParamsStoreKey(1), paramsBz)

				paramsBz, err = suite.cdc.Marshal(&v2.LegacyServiceParams{
					WhitelistedOperatorsIDs: []uint32{7, 8, 9},
					WhitelistedPoolsIDs:     []uint32{10, 11, 12},
				})
				suite.Require().NoError(err)
				sdkStore.Set(v2.ServiceParamsStoreKey(2), paramsBz)
			},
			check: func(ctx sdk.Context) {
				// Make sure the params are upgraded properly
				pools, err := suite.restakingKeeper.GetAllServicesSecuringPools(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal([]restakingtypes.ServiceSecuringPools{
					{ServiceID: 1, PoolIDs: []uint32{4, 5, 6}},
					{ServiceID: 2, PoolIDs: []uint32{10, 11, 12}},
				}, pools)

				operators, err := suite.restakingKeeper.GetAllServicesAllowedOperators(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal([]restakingtypes.ServiceAllowedOperators{
					{ServiceID: 1, OperatorIDs: []uint32{1, 2, 3}},
					{ServiceID: 2, OperatorIDs: []uint32{7, 8, 9}},
				}, operators)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.setup != nil {
				tc.setup(ctx)
			}

			err := v2.Migrate1To2(ctx, suite.storeKey, suite.cdc, suite.restakingKeeper, suite.operatorsKeeper, suite.servicesKeeper)
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
