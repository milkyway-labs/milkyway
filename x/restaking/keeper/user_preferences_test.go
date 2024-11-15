package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

func (suite *KeeperTestSuite) TestKeeper_SetUserPreferences() {
	testCases := []struct {
		name        string
		store       func(ctx sdk.Context)
		userAddress string
		preferences types.UserPreferences
		shouldErr   bool
		check       func(ctx sdk.Context)
	}{
		{
			name:        "User preferences are saved correctly",
			userAddress: "cosmos1jseuux3pktht0kkhlcsv4kqff3mql65udqs4jw",
			preferences: types.NewUserPreferences(
				true,
				false,
				[]uint32{1, 2, 3},
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				stored, err := suite.k.GetUserPreferences(ctx, "cosmos1jseuux3pktht0kkhlcsv4kqff3mql65udqs4jw")
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewUserPreferences(
					true,
					false,
					[]uint32{1, 2, 3},
				), stored)
			},
		},
		{
			name: "existing preferences are overridden properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetUserPreferences(ctx, "cosmos1jseuux3pktht0kkhlcsv4kqff3mql65udqs4jw", types.NewUserPreferences(
					true,
					false,
					[]uint32{1, 2, 3},
				))
				suite.Require().NoError(err)
			},
			userAddress: "cosmos1jseuux3pktht0kkhlcsv4kqff3mql65udqs4jw",
			preferences: types.NewUserPreferences(
				false,
				true,
				[]uint32{7},
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				stored, err := suite.k.GetUserPreferences(ctx, "cosmos1jseuux3pktht0kkhlcsv4kqff3mql65udqs4jw")
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewUserPreferences(
					false,
					true,
					[]uint32{7},
				), stored)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			err := suite.k.SetUserPreferences(ctx, tc.userAddress, tc.preferences)
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

func (suite *KeeperTestSuite) TestKeeper_GetUserPreferences() {
	testCases := []struct {
		name           string
		store          func(ctx sdk.Context)
		userAddress    string
		shouldErr      bool
		expPreferences types.UserPreferences
	}{
		{
			name:           "user without custom preferences returns default ones",
			userAddress:    "cosmos1jseuux3pktht0kkhlcsv4kqff3mql65udqs4jw",
			shouldErr:      false,
			expPreferences: types.DefaultUserPreferences(),
		},
		{
			name: "custom preferences are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetUserPreferences(ctx, "cosmos1jseuux3pktht0kkhlcsv4kqff3mql65udqs4jw", types.NewUserPreferences(
					true,
					false,
					[]uint32{1, 2, 3},
				))
				suite.Require().NoError(err)
			},
			userAddress: "cosmos1jseuux3pktht0kkhlcsv4kqff3mql65udqs4jw",
			shouldErr:   false,
			expPreferences: types.NewUserPreferences(
				true,
				false,
				[]uint32{1, 2, 3},
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			preferences, err := suite.k.GetUserPreferences(ctx, tc.userAddress)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}

			suite.Require().Equal(tc.expPreferences, preferences)
		})
	}
}
