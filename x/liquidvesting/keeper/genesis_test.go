package keeper_test

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
)

func (suite *KeeperTestSuite) TestKeeper_ExportGenesis() {
	user1 := "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"
	user2 := "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"
	vestedStake, err := types.GetVestedRepresentationDenom("stake")
	suite.Assert().NoError(err)
	blockTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	testCases := []struct {
		name       string
		setupCtx   func(sdk.Context) sdk.Context
		store      func(sdk.Context)
		shouldErr  bool
		expGenesis *types.GenesisState
	}{
		{
			name: "params are exported correctly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.DefaultParams())
				suite.Require().NoError(err)
			},
			shouldErr: false,
			expGenesis: &types.GenesisState{
				Params:             types.DefaultParams(),
				BurnCoins:          nil,
				UserInsuranceFunds: nil,
			},
		},
		{
			name: "insurance funds are exported correctly",
			store: func(ctx sdk.Context) {
				// Fund the users' insurance fund
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 2)))
				suite.fundAccountInsuranceFund(ctx, user2, sdk.NewCoins(sdk.NewInt64Coin("stake", 2)))
			},
			expGenesis: &types.GenesisState{
				Params:    types.DefaultParams(),
				BurnCoins: nil,
				UserInsuranceFunds: []types.UserInsuranceFundState{
					types.NewUserInsuranceFundState(user1, types.NewInsuranceFund(sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 2)), nil)),
					types.NewUserInsuranceFundState(user2, types.NewInsuranceFund(sdk.NewCoins(sdk.NewInt64Coin("stake", 2)), nil)),
				},
			},
		},
		{
			name: "burn coins funds are exported correctly",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.
					WithBlockHeight(10).
					WithBlockTime(blockTime)
			},
			store: func(ctx sdk.Context) {
				// Set the unbonding delegation time to 7 days
				suite.rk.SetParams(ctx, restakingtypes.NewParams(7*24*time.Hour, nil))

				// Fund the users' insurance fund
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 2)))
				suite.fundAccountInsuranceFund(ctx, user2, sdk.NewCoins(sdk.NewInt64Coin("stake", 2)))

				// Mint the staked representations
				suite.mintVestedRepresentation(user1, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 100)))
				suite.mintVestedRepresentation(user2, sdk.NewCoins(sdk.NewInt64Coin("stake", 100)))

				// Delegate the tokens so that they will be scheduled for burn after the
				// unbonding period
				suite.createPool(1, vestedIBCDenom)
				_, err = suite.rk.DelegateToPool(ctx, sdk.NewInt64Coin(vestedIBCDenom, 100), user1)
				suite.Assert().NoError(err)

				suite.createPool(2, vestedStake)
				_, err = suite.rk.DelegateToPool(ctx, sdk.NewInt64Coin(vestedStake, 100), user2)
				suite.Assert().NoError(err)

				// Burn the coins
				suite.Assert().NoError(suite.k.BurnVestedRepresentation(ctx, sdk.MustAccAddressFromBech32(user1), sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 100))))
				suite.Assert().NoError(suite.k.BurnVestedRepresentation(ctx, sdk.MustAccAddressFromBech32(user2), sdk.NewCoins(sdk.NewInt64Coin(vestedStake, 100))))
			},
			expGenesis: &types.GenesisState{
				Params: types.DefaultParams(),
				BurnCoins: []types.BurnCoins{
					types.NewBurnCoins(user1, blockTime.Add(7*24*time.Hour), sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 100))),
					types.NewBurnCoins(user2, blockTime.Add(7*24*time.Hour), sdk.NewCoins(sdk.NewInt64Coin(vestedStake, 100))),
				},
				UserInsuranceFunds: []types.UserInsuranceFundState{
					types.NewUserInsuranceFundState(user1, types.NewInsuranceFund(
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 2)),
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 2)))),
					types.NewUserInsuranceFundState(user2, types.NewInsuranceFund(
						sdk.NewCoins(sdk.NewInt64Coin("stake", 2)),
						sdk.NewCoins(sdk.NewInt64Coin("stake", 2)))),
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx := suite.ctx
			if tc.setupCtx != nil {
				ctx = tc.setupCtx(ctx)
			}

			if tc.store != nil {
				tc.store(ctx)
			}

			genState, err := suite.k.ExportGenesis(ctx)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expGenesis, genState)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeepr_InitGenesis() {
	user1 := "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"
	testCases := []struct {
		name      string
		setup     func(ctx sdk.Context)
		genesis   *types.GenesisState
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "genesis is initialized properly",
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				nil,
				nil,
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				params, _ := suite.k.GetParams(ctx)
				suite.Assert().Equal(types.DefaultParams(), params)
			},
		},
		{
			name: "should block negative insurance fund percentage",
			genesis: types.NewGenesisState(
				types.NewParams(math.LegacyNewDec(-1), nil, nil),
				nil, nil,
			),
			shouldErr: true,
		},
		{
			name: "should block 0 insurance fund percentage",
			genesis: types.NewGenesisState(
				types.NewParams(math.LegacyNewDec(0), nil, nil),
				nil, nil,
			),
			shouldErr: true,
		},
		{
			name: "should allow 100 insurance fund percentage",
			genesis: types.NewGenesisState(
				types.NewParams(math.LegacyNewDec(100), nil, nil),
				nil, nil,
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				params, _ := suite.k.GetParams(ctx)
				suite.Assert().Equal(math.LegacyNewDec(100), params.InsurancePercentage)
			},
		},
		{
			name: "should block > 100 insurance fund percentage",
			genesis: types.NewGenesisState(
				types.NewParams(math.LegacyNewDec(101), nil, nil),
				nil, nil,
			),
			shouldErr: true,
		},
		{
			name: "should block invalid minter address",
			genesis: types.NewGenesisState(
				types.NewParams(math.LegacyNewDec(2), nil, []string{"cosmos1fdsfd"}),
				nil, nil,
			),
			shouldErr: true,
		},
		{
			name: "should block invalid burners address",
			genesis: types.NewGenesisState(
				types.NewParams(math.LegacyNewDec(2), []string{"cosmos1fdsfd"}, nil),
				nil, nil,
			),
			shouldErr: true,
		},
		{
			name: "should block insurance fund initialization if module don't have enough tokens",
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				nil,
				[]types.UserInsuranceFundState{
					types.NewUserInsuranceFundState(user1, types.NewInsuranceFund(sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 100)), sdk.NewCoins())),
				},
			),
			shouldErr: true,
		},
		{
			name: "should block insurance fund initialization if user don't have delegations",
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				nil,
				[]types.UserInsuranceFundState{
					types.NewUserInsuranceFundState(user1, types.NewInsuranceFund(
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 100)),
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 2)))),
				},
			),
			shouldErr: true,
		},
		{
			name: "should block insurance fund initialization if insurance fund don't cover delegations and undelegations",
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				nil,
				[]types.UserInsuranceFundState{
					types.NewUserInsuranceFundState(user1, types.NewInsuranceFund(
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 2)),
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 2)))),
				},
			),
			setup: func(ctx sdk.Context) {
				// Send tokens to the liquid vesting module
				suite.Assert().NoError(
					suite.bk.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 2))))

				// Init the pools module
				testPool := poolstypes.NewPool(1, vestedIBCDenom)
				testPool.Tokens = math.NewIntFromUint64(50)
				testPool.DelegatorShares = math.LegacyNewDecFromInt(math.NewIntFromUint64(50))
				poolsKeeperGenesis := poolstypes.GenesisState{
					Params:     poolstypes.DefaultParams(),
					NextPoolID: 2,
					Pools:      []poolstypes.Pool{testPool},
				}
				suite.pk.InitGenesis(ctx, &poolsKeeperGenesis)

				// Init the restaking module
				restakingKeeperGenesis := &restakingtypes.GenesisState{
					Params: restakingtypes.DefaultParams(),
					UnbondingDelegations: []restakingtypes.UnbondingDelegation{
						restakingtypes.NewPoolUnbondingDelegation(
							user1,
							1,
							10,
							time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
							sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 51)),
							1,
						),
					},
					Delegations: []restakingtypes.Delegation{
						restakingtypes.NewPoolDelegation(1, user1,
							sdk.NewDecCoins(sdk.NewInt64DecCoin(testPool.GetSharesDenom(vestedIBCDenom), 50))),
					},
				}
				suite.rk.InitGenesis(ctx, restakingKeeperGenesis)
			},
			shouldErr: true,
		},
		{
			name: "should initialize insurance fund properly",
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				nil,
				[]types.UserInsuranceFundState{
					types.NewUserInsuranceFundState(user1, types.NewInsuranceFund(sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 100)), sdk.NewCoins())),
				},
			),
			setup: func(ctx sdk.Context) {
				err := suite.bk.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 100)))
				suite.Assert().NoError(err)
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				balance, err := suite.k.GetUserInsuranceFundBalance(ctx, sdk.MustAccAddressFromBech32(user1))
				suite.Assert().NoError(err)
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 100)), balance)
			},
		},
		{
			name: "should initialize insurance with delegations and undelegations",
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				nil,
				[]types.UserInsuranceFundState{
					types.NewUserInsuranceFundState(user1, types.NewInsuranceFund(
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10)),
						// Set 3 as used to cover the delegation and the undelegation
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 3)))),
				},
			),
			setup: func(ctx sdk.Context) {
				// Send tokens to the liquid vesting module
				suite.Assert().NoError(
					suite.bk.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10))))

				// Init the pools module
				testPool := poolstypes.NewPool(1, vestedIBCDenom)
				testPool.Tokens = math.NewIntFromUint64(50)
				testPool.DelegatorShares = math.LegacyNewDecFromInt(math.NewIntFromUint64(50))
				poolsKeeperGenesis := poolstypes.GenesisState{
					Params:     poolstypes.DefaultParams(),
					NextPoolID: 2,
					Pools:      []poolstypes.Pool{testPool},
				}
				suite.pk.InitGenesis(ctx, &poolsKeeperGenesis)

				// Init the restaking module
				restakingKeeperGenesis := &restakingtypes.GenesisState{
					Params: restakingtypes.DefaultParams(),
					UnbondingDelegations: []restakingtypes.UnbondingDelegation{
						restakingtypes.NewPoolUnbondingDelegation(
							user1,
							1,
							10,
							time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
							sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 100)),
							1,
						),
					},
					Delegations: []restakingtypes.Delegation{
						restakingtypes.NewPoolDelegation(1, user1,
							sdk.NewDecCoins(sdk.NewInt64DecCoin(testPool.GetSharesDenom(vestedIBCDenom), 50))),
					},
				}
				suite.rk.InitGenesis(ctx, restakingKeeperGenesis)
			},
			shouldErr: false,
		},
		{
			name: "should block burn coins initialization if the tokens are not unbonding",
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				[]types.BurnCoins{
					types.NewBurnCoins(user1, time.Now(), sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 2))),
				},
				nil,
			),
			shouldErr: true,
		},
		{
			name: "should initialize burn coins properly",
			setup: func(ctx sdk.Context) {
				restakingKeeperGenesis := &restakingtypes.GenesisState{
					Params: restakingtypes.DefaultParams(),
					UnbondingDelegations: []restakingtypes.UnbondingDelegation{
						restakingtypes.NewPoolUnbondingDelegation(
							user1,
							1,
							10,
							time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
							sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 2)),
							1,
						),
					},
				}
				suite.rk.InitGenesis(ctx, restakingKeeperGenesis)
			},
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				[]types.BurnCoins{
					types.NewBurnCoins(user1, time.Now(), sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 2))),
				},
				nil,
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				burnCoins := suite.k.GetAllBurnCoins(ctx)
				suite.Assert().Len(burnCoins, 1)
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 2)),
					burnCoins[0].Amount)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			if tc.setup != nil {
				tc.setup(suite.ctx)
			}
			err := suite.k.InitGenesis(suite.ctx, tc.genesis)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}

			if tc.check != nil {
				tc.check(suite.ctx)
			}
		})
	}
}
