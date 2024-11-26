package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	servicestypes "github.com/milkyway-labs/milkyway/v2/x/services/types"
)

func (suite *KeeperTestSuite) TestOperatorHooks_BeforeOperatorDeleted() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		check      func(ctx sdk.Context)
		operatorID uint32
		shouldErr  bool
	}{
		{
			name: "operator services associations is removed correctly",
			store: func(ctx sdk.Context) {
				err := suite.k.AddServiceToOperatorJoinedServices(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 1, 2)
				suite.Require().NoError(err)

				// Add some other to ensure we don't eliminate the wrong data
				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 2, 1)
				suite.Require().NoError(err)
				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 2, 2)
				suite.Require().NoError(err)
			},
			check: func(ctx sdk.Context) {
				joined, err := suite.k.HasOperatorJoinedService(ctx, 1, 1)
				suite.Assert().NoError(err)
				suite.Assert().False(joined)
				joined, err = suite.k.HasOperatorJoinedService(ctx, 1, 2)
				suite.Assert().NoError(err)
				suite.Assert().False(joined)

				// Check that we didn't remove other data
				joined, err = suite.k.HasOperatorJoinedService(ctx, 2, 1)
				suite.Assert().NoError(err)
				suite.Assert().True(joined)
				joined, err = suite.k.HasOperatorJoinedService(ctx, 2, 2)
				suite.Assert().NoError(err)
				suite.Assert().True(joined)
			},
			operatorID: 1,
			shouldErr:  false,
		},
		{
			name: "service status doesn't change if it didn't have an operators allowed list configured",
			store: func(ctx sdk.Context) {
				err := suite.sk.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					2,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)

				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 1, 1)
				suite.Require().NoError(err)
			},
			check: func(ctx sdk.Context) {
				joined, err := suite.k.HasOperatorJoinedService(ctx, 1, 1)
				suite.Assert().NoError(err)
				suite.Assert().False(joined)

				// Ensure that the service is status has not changed
				service, found, err := suite.sk.GetService(ctx, 2)
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(servicestypes.SERVICE_STATUS_ACTIVE, service.Status)
			},
			operatorID: 1,
			shouldErr:  false,
		},
		{
			name: "service status doesn't change if the operator was not part of the allow list",
			store: func(ctx sdk.Context) {
				err := suite.sk.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					2,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)

				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 2, 2)
				suite.Require().NoError(err)
			},
			check: func(ctx sdk.Context) {
				joined, err := suite.k.HasOperatorJoinedService(ctx, 1, 1)
				suite.Assert().NoError(err)
				suite.Assert().False(joined)

				joined, err = suite.k.HasOperatorJoinedService(ctx, 2, 2)
				suite.Assert().NoError(err)
				suite.Assert().True(joined)

				// Ensure that the service is status has not changed
				service, found, err := suite.sk.GetService(ctx, 2)
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(servicestypes.SERVICE_STATUS_ACTIVE, service.Status)
			},
			operatorID: 1,
			shouldErr:  false,
		},
		{
			name: "service is deactivated once last operator is removed from allow list",
			store: func(ctx sdk.Context) {
				err := suite.sk.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)

				// Add the operator to the service allow list
				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 1)
				suite.Require().NoError(err)
			},
			check: func(ctx sdk.Context) {
				joined, err := suite.k.HasOperatorJoinedService(ctx, 1, 1)
				suite.Assert().NoError(err)
				suite.Assert().False(joined)

				// Ensure that the service is now inactive
				service, found, err := suite.sk.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(servicestypes.SERVICE_STATUS_INACTIVE, service.Status)
			},
			operatorID: 1,
			shouldErr:  false,
		},
		{
			name: "service is not deactivated if an operator is still present in its allow list",
			store: func(ctx sdk.Context) {
				err := suite.sk.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)

				// Add the operator to the service allow list
				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 1)
				suite.Require().NoError(err)
				// Add a second operator to service the allow list
				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			check: func(ctx sdk.Context) {
				joined, err := suite.k.HasOperatorJoinedService(ctx, 1, 1)
				suite.Assert().NoError(err)
				suite.Assert().False(joined)

				// Ensure the service status has not changed
				service, found, err := suite.sk.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(servicestypes.SERVICE_STATUS_ACTIVE, service.Status)
			},
			operatorID: 1,
			shouldErr:  false,
		},
		{
			name: "service status don't change when removing last operator if service status is created",
			store: func(ctx sdk.Context) {
				err := suite.sk.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)

				// Add the operator to the service allow list
				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 1)
				suite.Require().NoError(err)
			},
			check: func(ctx sdk.Context) {
				joined, err := suite.k.HasOperatorJoinedService(ctx, 1, 1)
				suite.Assert().NoError(err)
				suite.Assert().False(joined)

				// Ensure the service status has not changed
				service, found, err := suite.sk.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(servicestypes.SERVICE_STATUS_CREATED, service.Status)
			},
			operatorID: 1,
			shouldErr:  false,
		},
		{
			name: "service status don't change when removing last operator if service status is inactive",
			store: func(ctx sdk.Context) {
				err := suite.sk.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_INACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)

				// Add the operator to the service allow list
				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 1)
				suite.Require().NoError(err)
			},
			check: func(ctx sdk.Context) {
				joined, err := suite.k.HasOperatorJoinedService(ctx, 1, 1)
				suite.Assert().NoError(err)
				suite.Assert().False(joined)

				// Ensure the service status has not changed
				service, found, err := suite.sk.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(servicestypes.SERVICE_STATUS_INACTIVE, service.Status)
			},
			operatorID: 1,
			shouldErr:  false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx := suite.ctx
			if tc.store != nil {
				tc.store(ctx)
			}

			hooks := suite.k.OperatorsHooks()
			err := hooks.BeforeOperatorDeleted(ctx, tc.operatorID)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				if tc.check != nil {
					tc.check(ctx)
				}
			}
		})
	}
}
