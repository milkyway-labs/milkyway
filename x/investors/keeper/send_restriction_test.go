package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/milkyway-labs/milkyway/v10/app/testutil"
	"github.com/milkyway-labs/milkyway/v10/utils"
	distrkeeper "github.com/milkyway-labs/milkyway/v10/x/distribution/keeper"
	"github.com/milkyway-labs/milkyway/v10/x/investors/types"
)

func (suite *KeeperTestSuite) TestSendRestrictionFn() {
	// normal send shouldn't be affected
	// if the withdraw address is a vesting investor, it shouldn't be affected
	// it should get the delegator address by the withdraw address
	// if it's a vesting account but not a vesting investor, it shouldn't be affected
	// happy path
	ctx, _ := suite.ctx.CacheContext()
	ctx = suite.ctx.WithBlockTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

	err := suite.k.SetInvestorsRewardRatio(ctx, utils.MustParseDec("0.5")) // 50%
	suite.Require().NoError(err)

	// Initial balances: 2000000stake
	normalAddr := testutil.TestAddress(1)
	suite.fundAccount(ctx, normalAddr.String(), utils.MustParseCoins("2000000stake"))

	// Initial balances: 2000000stake(1000000 spendable, 1000000 vesting)
	normalVestingAddr := testutil.TestAddress(2)
	suite.createVestingAccount(
		ctx,
		testutil.TestAddress(10000).String(),
		normalVestingAddr.String(),
		utils.MustParseCoins("1000000stake"),
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		false,
		true,
	)
	suite.fundAccount(ctx, normalVestingAddr.String(), utils.MustParseCoins("1000000stake"))

	// Initial balances: 2000000stake(1000000 spendable, 1000000 vesting)
	vestingInvestorAddr := testutil.TestAddress(3)
	suite.createVestingAccount(
		ctx,
		testutil.TestAddress(10000).String(),
		vestingInvestorAddr.String(),
		utils.MustParseCoins("1000000stake"),
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		false,
		true,
	)
	suite.fundAccount(ctx, vestingInvestorAddr.String(), utils.MustParseCoins("1000000stake"))
	err = suite.k.SetVestingInvestor(ctx, vestingInvestorAddr.String())
	suite.Require().NoError(err)

	// Fund the distribution module account to simulate rewards distribution
	suite.fundModuleAccount(ctx, distrtypes.ModuleName, utils.MustParseCoins("100000000stake"))

	testCases := []struct {
		name  string
		setup func(ctx sdk.Context)
		send  func(ctx sdk.Context)
		check func(ctx sdk.Context)
	}{
		{
			name: "vesting investor's staking rewards should be redirected",
			send: func(ctx sdk.Context) {
				err := suite.bk.SendCoinsFromModuleToAccount(
					ctx,
					distrtypes.ModuleName,
					vestingInvestorAddr,
					utils.MustParseCoins("1000000stake"),
				)
				suite.Require().NoError(err)
			},
			check: func(ctx sdk.Context) {
				// Rewards are redirected to the investors module account entirely
				balances := suite.bk.GetAllBalances(ctx, vestingInvestorAddr)
				suite.Assert().Equal("2000000stake", balances.String())

				balances = suite.bk.GetAllBalances(ctx, suite.ak.GetModuleAddress(types.ModuleName))
				suite.Assert().Equal("1000000stake", balances.String())
			},
		},
		{
			name: "normal transfers should not be affected",
			send: func(ctx sdk.Context) {
				err := suite.bk.SendCoinsFromModuleToAccount(
					ctx,
					distrtypes.ModuleName,
					normalAddr,
					utils.MustParseCoins("1000000stake"),
				)
				suite.Require().NoError(err)
			},
			check: func(ctx sdk.Context) {
				balances := suite.bk.GetAllBalances(ctx, normalAddr)
				suite.Assert().Equal("3000000stake", balances.String())
			},
		},
		{
			name: "sending from a vesting investor should not be affected",
			send: func(ctx sdk.Context) {
				err := suite.bk.SendCoinsFromAccountToModule(
					ctx,
					vestingInvestorAddr,
					distrtypes.ModuleName,
					utils.MustParseCoins("1000000stake"),
				)
				suite.Require().NoError(err)
			},
			check: func(ctx sdk.Context) {
				balances := suite.bk.GetAllBalances(ctx, suite.ak.GetModuleAddress(distrtypes.ModuleName))
				suite.Assert().Equal("101000000stake", balances.String())
			},
		},
		{
			name: "sending from normal accounts to a vesting investor should not be affected",
			send: func(ctx sdk.Context) {
				err := suite.bk.SendCoins(
					ctx,
					normalAddr,
					vestingInvestorAddr,
					utils.MustParseCoins("1000000stake"),
				)
				suite.Require().NoError(err)
			},
			check: func(ctx sdk.Context) {
				balances := suite.bk.GetAllBalances(ctx, vestingInvestorAddr)
				suite.Assert().Equal("3000000stake", balances.String())
			},
		},
		{
			name: "normal vesting account should not be affected",
			send: func(ctx sdk.Context) {
				err := suite.bk.SendCoinsFromModuleToAccount(
					ctx,
					distrtypes.ModuleName,
					normalVestingAddr,
					utils.MustParseCoins("1000000stake"),
				)
				suite.Require().NoError(err)
			},
			check: func(ctx sdk.Context) {
				balances := suite.bk.GetAllBalances(ctx, normalVestingAddr)
				suite.Assert().Equal("3000000stake", balances.String())
			},
		},
		{
			name: "withdraw address is properly handled",
			setup: func(ctx sdk.Context) {
				msgServer := distrkeeper.NewMsgServerImpl(suite.dk)
				_, err := msgServer.SetWithdrawAddress(
					ctx,
					distrtypes.NewMsgSetWithdrawAddress(vestingInvestorAddr, normalAddr),
				)
				suite.Require().NoError(err)
			},
			send: func(ctx sdk.Context) {
				err := suite.bk.SendCoinsFromModuleToAccount(
					ctx,
					distrtypes.ModuleName,
					normalAddr,
					utils.MustParseCoins("1000000stake"),
				)
				suite.Require().NoError(err)
			},
			check: func(ctx sdk.Context) {
				// Rewards are redirected
				balances := suite.bk.GetAllBalances(ctx, normalAddr)
				suite.Assert().Equal("2000000stake", balances.String())

				balances = suite.bk.GetAllBalances(ctx, suite.ak.GetModuleAddress(types.ModuleName))
				suite.Assert().Equal("1000000stake", balances.String())
			},
		},
		{
			name: "if the withdraw address is a vesting investor's, it should not be affected",
			setup: func(ctx sdk.Context) {
				msgServer := distrkeeper.NewMsgServerImpl(suite.dk)
				_, err := msgServer.SetWithdrawAddress(
					ctx,
					distrtypes.NewMsgSetWithdrawAddress(normalAddr, vestingInvestorAddr),
				)
				suite.Require().NoError(err)
			},
			send: func(ctx sdk.Context) {
				err := suite.bk.SendCoinsFromModuleToAccount(
					ctx,
					distrtypes.ModuleName,
					vestingInvestorAddr,
					utils.MustParseCoins("1000000stake"),
				)
				suite.Require().NoError(err)
			},
			check: func(ctx sdk.Context) {
				balances := suite.bk.GetAllBalances(ctx, vestingInvestorAddr)
				suite.Assert().Equal("3000000stake", balances.String())
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := ctx.CacheContext()

			if tc.setup != nil {
				tc.setup(ctx)
			}

			tc.send(ctx)
			tc.check(ctx)
		})
	}
}
