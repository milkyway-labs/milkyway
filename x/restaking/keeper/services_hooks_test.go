package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestServicesHooks_AfterServiceDeleted() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		check     func(ctx sdk.Context)
		serviceID uint32
		shouldErr bool
	}{
		{
			name: "operator participation is removed",
			store: func(ctx sdk.Context) {
				err := suite.k.AddServiceToOperatorJoinedServices(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 1, 2)
				suite.Require().NoError(err)

				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 2, 1)
				suite.Require().NoError(err)
				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 2, 2)
				suite.Require().NoError(err)
			},
			check: func(ctx sdk.Context) {
				joined, err := suite.k.HasOperatorJoinedService(ctx, 1, 1)
				suite.Assert().NoError(err)
				suite.Assert().False(joined)
				joined, err = suite.k.HasOperatorJoinedService(ctx, 2, 1)
				suite.Assert().NoError(err)
				suite.Assert().False(joined)

				// Check that we didn't remove other data
				joined, err = suite.k.HasOperatorJoinedService(ctx, 1, 2)
				suite.Assert().NoError(err)
				suite.Assert().True(joined)
				joined, err = suite.k.HasOperatorJoinedService(ctx, 2, 2)
				suite.Assert().NoError(err)
				suite.Assert().True(joined)
			},
			serviceID: 1,
			shouldErr: false,
		},
		{
			name: "service's operators allow list is wiped",
			store: func(ctx sdk.Context) {
				err := suite.k.AddOperatorToServiceAllowList(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			check: func(ctx sdk.Context) {
				configured, err := suite.k.IsServiceOpertorsAllowListConfigured(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().False(configured)
			},
			serviceID: 1,
			shouldErr: false,
		},
		{
			name: "service's securing pools list is wiped",
			store: func(ctx sdk.Context) {
				err := suite.k.AddPoolToServiceSecuringPools(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.AddPoolToServiceSecuringPools(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			check: func(ctx sdk.Context) {
				configured, err := suite.k.IsServiceSecuringPoolsConfigured(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().False(configured)
			},
			serviceID: 1,
			shouldErr: false,
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

			hooks := suite.k.ServicesHooks()
			err := hooks.AfterServiceDeleted(ctx, tc.serviceID)
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
