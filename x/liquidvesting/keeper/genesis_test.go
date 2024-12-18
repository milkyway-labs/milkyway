package keeper_test

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v5/x/liquidvesting/types"
	poolstypes "github.com/milkyway-labs/milkyway/v5/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v5/x/restaking/types"
)

func (suite *KeeperTestSuite) TestKeeper_ExportGenesis() {
	lockedStakeDenom, err := types.GetLockedRepresentationDenom("stake")
	suite.Assert().NoError(err)

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
				suite.fundAccountInsuranceFund(ctx, "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre", sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 2)))
				suite.fundAccountInsuranceFund(ctx, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4", sdk.NewCoins(sdk.NewInt64Coin("stake", 2)))
			},
			expGenesis: &types.GenesisState{
				Params:    types.DefaultParams(),
				BurnCoins: nil,
				UserInsuranceFunds: []types.UserInsuranceFundEntry{
					types.NewUserInsuranceFundEntry(
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						sdk.NewCoins(sdk.NewInt64Coin("stake", 2)),
					),
					types.NewUserInsuranceFundEntry(
						"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 2)),
					),
				},
			},
		},
		{
			name: "burn coins funds are exported correctly",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.
					WithBlockHeight(10).
					WithBlockTime(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
			},
			store: func(ctx sdk.Context) {
				// Set the unbonding delegation time to 7 days
				err = suite.rk.SetParams(ctx, restakingtypes.NewParams(
					7*24*time.Hour,
					nil,
					restakingtypes.DefaultRestakingCap,
				))
				suite.Require().NoError(err)

				// Fund the users' insurance fund
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 2)),
				)
				suite.fundAccountInsuranceFund(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewInt64Coin("stake", 2)),
				)

				// Mint the staked representations
				suite.mintLockedRepresentation(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 100)),
				)
				suite.mintLockedRepresentation(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
				)

				// Delegate the tokens so that they will be scheduled for burn after the
				// unbonding period
				suite.createPool(ctx, 1, LockedIBCDenom)
				_, err = suite.rk.DelegateToPool(ctx,
					sdk.NewInt64Coin(LockedIBCDenom, 100),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				)
				suite.Assert().NoError(err)

				suite.createPool(ctx, 2, lockedStakeDenom)
				_, err = suite.rk.DelegateToPool(ctx,
					sdk.NewInt64Coin(lockedStakeDenom, 100),
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				)
				suite.Assert().NoError(err)

				// Burn the coins
				userAddr, err := sdk.AccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Assert().NoError(err)
				err = suite.k.BurnLockedRepresentation(ctx, userAddr, sdk.NewCoins(sdk.NewInt64Coin(lockedStakeDenom, 100)))
				suite.Assert().NoError(err)

				userAddr, err = sdk.AccAddressFromBech32("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre")
				suite.Assert().NoError(err)
				err = suite.k.BurnLockedRepresentation(ctx, userAddr, sdk.NewCoins(sdk.NewInt64Coin(LockedIBCDenom, 100)))
				suite.Assert().NoError(err)
			},
			expGenesis: &types.GenesisState{
				Params: types.DefaultParams(),
				BurnCoins: []types.BurnCoins{
					types.NewBurnCoins(
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC).Add(7*24*time.Hour),
						sdk.NewCoins(sdk.NewInt64Coin(lockedStakeDenom, 100)),
					),
					types.NewBurnCoins(
						"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
						time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC).Add(7*24*time.Hour),
						sdk.NewCoins(sdk.NewInt64Coin(LockedIBCDenom, 100)),
					),
				},
				UserInsuranceFunds: []types.UserInsuranceFundEntry{
					types.NewUserInsuranceFundEntry(
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						sdk.NewCoins(sdk.NewInt64Coin("stake", 2)),
					),
					types.NewUserInsuranceFundEntry(
						"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 2)),
					),
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
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
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
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
				types.NewParams(math.LegacyNewDec(-1), nil, nil, nil),
				nil,
				nil,
			),
			shouldErr: true,
		},
		{
			name: "should block 0 insurance fund percentage",
			genesis: types.NewGenesisState(
				types.NewParams(math.LegacyNewDec(0), nil, nil, nil),
				nil,
				nil,
			),
			shouldErr: true,
		},
		{
			name: "should allow 100 insurance fund percentage",
			genesis: types.NewGenesisState(
				types.NewParams(math.LegacyNewDec(100), nil, nil, nil),
				nil,
				nil,
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
				types.NewParams(math.LegacyNewDec(101), nil, nil, nil),
				nil,
				nil,
			),
			shouldErr: true,
		},
		{
			name: "should block invalid minter address",
			genesis: types.NewGenesisState(
				types.NewParams(math.LegacyNewDec(2), nil, []string{"cosmos1fdsfd"}, nil),
				nil,
				nil,
			),
			shouldErr: true,
		},
		{
			name: "should block invalid burners address",
			genesis: types.NewGenesisState(
				types.NewParams(math.LegacyNewDec(2), []string{"cosmos1fdsfd"}, nil, nil),
				nil,
				nil,
			),
			shouldErr: true,
		},
		{
			name: "should block invalid allowed depositors address",
			genesis: types.NewGenesisState(
				types.NewParams(math.LegacyNewDec(2), nil, nil, []string{"cosmos1fdsfd"}),
				nil,
				nil,
			),
			shouldErr: true,
		},
		{
			name: "should block insurance fund initialization if module don't have enough tokens",
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				nil,
				[]types.UserInsuranceFundEntry{
					types.NewUserInsuranceFundEntry(
						"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 100)),
					),
				},
			),
			shouldErr: true,
		},
		{
			name: "should block insurance fund initialization if user don't have delegations",
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				nil,
				[]types.UserInsuranceFundEntry{
					types.NewUserInsuranceFundEntry(
						"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 100)),
					),
				},
			),
			shouldErr: true,
		},
		{
			name: "should block insurance fund initialization if insurance fund don't cover delegations and undelegations",
			store: func(ctx sdk.Context) {
				// Send tokens to the liquid vesting module
				err := suite.bk.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 2)))
				suite.Assert().NoError(err)

				// Init the pools module
				testPool := poolstypes.NewPool(1, LockedIBCDenom)
				testPool.Tokens = math.NewIntFromUint64(50)
				testPool.DelegatorShares = math.LegacyNewDecFromInt(math.NewIntFromUint64(50))
				poolsKeeperGenesis := poolstypes.GenesisState{
					NextPoolID: 2,
					Pools:      []poolstypes.Pool{testPool},
				}

				err = suite.pk.InitGenesis(ctx, &poolsKeeperGenesis)
				suite.Assert().NoError(err)

				// Init the restaking module
				restakingKeeperGenesis := &restakingtypes.GenesisState{
					Params: restakingtypes.DefaultParams(),
					UnbondingDelegations: []restakingtypes.UnbondingDelegation{
						restakingtypes.NewPoolUnbondingDelegation(
							"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
							1,
							10,
							time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
							sdk.NewCoins(sdk.NewInt64Coin(LockedIBCDenom, 51)),
							1,
						),
					},
					Delegations: []restakingtypes.Delegation{
						restakingtypes.NewPoolDelegation(
							1,
							"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
							sdk.NewDecCoins(sdk.NewInt64DecCoin(testPool.GetSharesDenom(LockedIBCDenom), 50)),
						),
					},
				}
				err = suite.rk.InitGenesis(ctx, restakingKeeperGenesis)
				suite.Assert().NoError(err)
			},
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				nil,
				[]types.UserInsuranceFundEntry{
					types.NewUserInsuranceFundEntry(
						"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 2)),
					),
				},
			),
			shouldErr: true,
		},
		{
			name: "should initialize insurance fund properly",
			store: func(ctx sdk.Context) {
				err := suite.bk.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 100)))
				suite.Assert().NoError(err)
			},
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				nil,
				[]types.UserInsuranceFundEntry{
					types.NewUserInsuranceFundEntry(
						"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 100)),
					),
				},
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				balance, err := suite.k.GetUserInsuranceFundBalance(ctx, "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre")
				suite.Assert().NoError(err)
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 100)), balance)
			},
		},
		{
			name: "should initialize insurance with delegations and undelegations",
			store: func(ctx sdk.Context) {
				// Send tokens to the liquid vesting module
				err := suite.bk.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10)))
				suite.Assert().NoError(err)

				// Init the pools module
				testPool := poolstypes.NewPool(1, LockedIBCDenom)
				testPool.Tokens = math.NewIntFromUint64(50)
				testPool.DelegatorShares = math.LegacyNewDecFromInt(math.NewIntFromUint64(50))
				poolsKeeperGenesis := poolstypes.GenesisState{
					NextPoolID: 2,
					Pools:      []poolstypes.Pool{testPool},
				}

				err = suite.pk.InitGenesis(ctx, &poolsKeeperGenesis)
				suite.Assert().NoError(err)

				// Init the restaking module
				restakingKeeperGenesis := &restakingtypes.GenesisState{
					Params: restakingtypes.DefaultParams(),
					UnbondingDelegations: []restakingtypes.UnbondingDelegation{
						restakingtypes.NewPoolUnbondingDelegation(
							"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
							1,
							10,
							time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
							sdk.NewCoins(sdk.NewInt64Coin(LockedIBCDenom, 100)),
							1,
						),
					},
					Delegations: []restakingtypes.Delegation{
						restakingtypes.NewPoolDelegation(
							1,
							"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
							sdk.NewDecCoins(sdk.NewInt64DecCoin(testPool.GetSharesDenom(LockedIBCDenom), 50)),
						),
					},
				}
				err = suite.rk.InitGenesis(ctx, restakingKeeperGenesis)
				suite.Assert().NoError(err)
			},
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				nil,
				[]types.UserInsuranceFundEntry{
					types.NewUserInsuranceFundEntry(
						"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10)),
					),
				},
			),
			shouldErr: false,
		},
		{
			name: "should block burn coins initialization if the tokens are not unbonding",
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				[]types.BurnCoins{
					types.NewBurnCoins(
						"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
						time.Now(),
						sdk.NewCoins(sdk.NewInt64Coin(LockedIBCDenom, 2)),
					),
				},
				nil,
			),
			shouldErr: true,
		},
		{
			name: "should initialize burn coins properly",
			store: func(ctx sdk.Context) {
				restakingKeeperGenesis := &restakingtypes.GenesisState{
					Params: restakingtypes.DefaultParams(),
					UnbondingDelegations: []restakingtypes.UnbondingDelegation{
						restakingtypes.NewPoolUnbondingDelegation(
							"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
							1,
							10,
							time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
							sdk.NewCoins(sdk.NewInt64Coin(LockedIBCDenom, 2)),
							1,
						),
					},
				}
				err := suite.rk.InitGenesis(ctx, restakingKeeperGenesis)
				suite.Assert().NoError(err)
			},
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				[]types.BurnCoins{
					types.NewBurnCoins(
						"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
						time.Now(),
						sdk.NewCoins(sdk.NewInt64Coin(LockedIBCDenom, 2)),
					),
				},
				nil,
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				burnCoins := suite.k.GetAllBurnCoins(ctx)
				suite.Assert().Len(burnCoins, 1)
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin(LockedIBCDenom, 2)), burnCoins[0].Amount)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			err := suite.k.InitGenesis(ctx, tc.genesis)
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
