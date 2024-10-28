package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestOperatorHooks_AfterOperatorInactivatingCompleted() {
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
			err := hooks.AfterOperatorDeleted(ctx, tc.operatorID)
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
