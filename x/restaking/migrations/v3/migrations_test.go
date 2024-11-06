package v3_test

import (
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	db "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"

	"github.com/milkyway-labs/milkyway/app"
	appkeepers "github.com/milkyway-labs/milkyway/app/keepers"
	bankkeeper "github.com/milkyway-labs/milkyway/x/bank/keeper"
	operatorskeeper "github.com/milkyway-labs/milkyway/x/operators/keeper"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolskeeper "github.com/milkyway-labs/milkyway/x/pools/keeper"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingkeeper "github.com/milkyway-labs/milkyway/x/restaking/keeper"
	legacytypes "github.com/milkyway-labs/milkyway/x/restaking/legacy/types"
	v3 "github.com/milkyway-labs/milkyway/x/restaking/migrations/v3"
	"github.com/milkyway-labs/milkyway/x/restaking/testutils"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
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
	// Define store keys
	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey,
		poolstypes.StoreKey, operatorstypes.StoreKey, servicestypes.StoreKey,
		restakingtypes.StoreKey,
	)
	suite.storeKey = keys[restakingtypes.StoreKey]

	// Create logger
	logger := log.NewNopLogger()

	// Create an in-memory db
	memDB := db.NewMemDB()
	ms := store.NewCommitMultiStore(memDB, logger, metrics.NewNoOpMetrics())
	for _, key := range keys {
		ms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, memDB)
	}

	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}

	suite.ctx = sdk.NewContext(ms, tmproto.Header{ChainID: "test-chain"}, false, log.NewNopLogger())
	suite.cdc, _ = app.MakeCodecs()

	// Authority address
	authorityAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// Build keepers

	authKeeper := authkeeper.NewAccountKeeper(
		suite.cdc,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		app.GetMaccPerms(),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authorityAddr,
	)
	bankKeeper := bankkeeper.NewKeeper(
		suite.cdc,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		authKeeper,
		nil,
		authorityAddr,
		logger,
	)
	communityPoolKeeper := appkeepers.NewCommunityPoolKeeper(
		bankKeeper,
		authtypes.FeeCollectorName,
	)
	poolsKeeper := poolskeeper.NewKeeper(
		suite.cdc,
		keys[poolstypes.StoreKey],
		runtime.NewKVStoreService(keys[poolstypes.StoreKey]),
		authKeeper,
	)
	suite.operatorsKeeper = operatorskeeper.NewKeeper(
		suite.cdc,
		keys[operatorstypes.StoreKey],
		runtime.NewKVStoreService(keys[operatorstypes.StoreKey]),
		authKeeper,
		communityPoolKeeper,
		authorityAddr,
	)
	suite.servicesKeeper = serviceskeeper.NewKeeper(
		suite.cdc,
		keys[servicestypes.StoreKey],
		runtime.NewKVStoreService(keys[servicestypes.StoreKey]),
		authKeeper,
		communityPoolKeeper,
		authorityAddr,
	)
	suite.restakingKeeper = restakingkeeper.NewKeeper(
		suite.cdc,
		suite.storeKey,
		runtime.NewKVStoreService(keys[restakingtypes.StoreKey]),
		authKeeper,
		bankKeeper,
		poolsKeeper,
		suite.operatorsKeeper,
		suite.servicesKeeper,
		authorityAddr,
	).SetHooks(testutils.NewMockHooks())
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
			name: "legacy params are migrated correctly",
			setup: func(ctx sdk.Context) {
				sdkStore := ctx.KVStore(suite.storeKey)

				// Set the legacy params
				paramsBz, err := suite.cdc.Marshal(&legacytypes.Params{
					UnbondingTime: 32 * time.Hour,
				})
				suite.Require().NoError(err)
				sdkStore.Set(types.LegacyParamsKey, paramsBz)
			},
			check: func(ctx sdk.Context) {
				params, err := suite.restakingKeeper.GetParams(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewParams(32*time.Hour, nil), params)
			},
		},
		{
			name: "not stored legacy params are migrated to default params",
			check: func(ctx sdk.Context) {
				params, err := suite.restakingKeeper.GetParams(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(types.DefaultParams(), params)
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

			err := v3.Migrate2To3(ctx, suite.storeKey, suite.cdc, suite.restakingKeeper)
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
