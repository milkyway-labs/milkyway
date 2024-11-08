package testutils

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibchooks "github.com/initia-labs/initia/x/ibc-hooks"
	ibchookskeeper "github.com/initia-labs/initia/x/ibc-hooks/keeper"
	ibchookstypes "github.com/initia-labs/initia/x/ibc-hooks/types"
	"github.com/stretchr/testify/require"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	appkeepers "github.com/milkyway-labs/milkyway/app/keepers"
	"github.com/milkyway-labs/milkyway/testutils/storetesting"
	"github.com/milkyway-labs/milkyway/x/liquidvesting/hooks"
	"github.com/milkyway-labs/milkyway/x/liquidvesting/keeper"
	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
	operatorskeeper "github.com/milkyway-labs/milkyway/x/operators/keeper"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolskeeper "github.com/milkyway-labs/milkyway/x/pools/keeper"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingkeeper "github.com/milkyway-labs/milkyway/x/restaking/keeper"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	serviceskeeper "github.com/milkyway-labs/milkyway/x/services/keeper"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

type KeeperTestData struct {
	storetesting.BaseKeeperTestData

	StoreKey storetypes.StoreKey

	LiquidVestingModuleAddress sdk.AccAddress

	Keeper          *keeper.Keeper
	IBCHooksKeeper  *ibchookskeeper.Keeper
	IBCMiddleware   ibchooks.IBCMiddleware
	OperatorsKeeper *operatorskeeper.Keeper
	PoolsKeeper     *poolskeeper.Keeper
	ServicesKeeper  *serviceskeeper.Keeper
	RestakingKeeper *restakingkeeper.Keeper
}

// NewKeeperTestData returns a new KeeperTestData
func NewKeeperTestData(t *testing.T) KeeperTestData {
	t.Helper()

	var data = KeeperTestData{
		BaseKeeperTestData: storetesting.NewBaseKeeperTestData(t, []string{
			types.StoreKey, authtypes.StoreKey, banktypes.StoreKey,
			operatorstypes.StoreKey, poolstypes.StoreKey, servicestypes.StoreKey,
			restakingtypes.StoreKey, stakingtypes.StoreKey,
			distributiontypes.StoreKey, ibchookstypes.StoreKey,
		}),
	}

	ac := authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	// Module addresses
	data.LiquidVestingModuleAddress = authtypes.NewModuleAddress(types.ModuleName)

	// Build keepers

	data.IBCHooksKeeper = ibchookskeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[ibchookstypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		ac,
	)
	if err := data.IBCHooksKeeper.Params.Set(data.Context, ibchookstypes.DefaultParams()); err != nil {
		panic(err)
	}

	data.PoolsKeeper = poolskeeper.NewKeeper(
		data.Cdc,
		data.Keys[poolstypes.StoreKey],
		runtime.NewKVStoreService(data.Keys[poolstypes.StoreKey]),
		data.AccountKeeper,
	)
	communityPoolKeeper := appkeepers.NewCommunityPoolKeeper(
		data.BankKeeper,
		authtypes.FeeCollectorName,
	)
	data.OperatorsKeeper = operatorskeeper.NewKeeper(
		data.Cdc,
		data.Keys[operatorstypes.StoreKey],
		runtime.NewKVStoreService(data.Keys[operatorstypes.StoreKey]),
		data.AccountKeeper,
		communityPoolKeeper,
		data.AuthorityAddress,
	)
	data.ServicesKeeper = serviceskeeper.NewKeeper(
		data.Cdc,
		data.Keys[servicestypes.StoreKey],
		runtime.NewKVStoreService(data.Keys[servicestypes.StoreKey]),
		data.AccountKeeper,
		communityPoolKeeper,
		data.AuthorityAddress,
	)
	data.RestakingKeeper = restakingkeeper.NewKeeper(
		data.Cdc,
		data.Keys[restakingtypes.StoreKey],
		runtime.NewKVStoreService(data.Keys[restakingtypes.StoreKey]),
		data.AccountKeeper,
		data.BankKeeper,
		data.PoolsKeeper,
		data.OperatorsKeeper,
		data.ServicesKeeper,
		data.AuthorityAddress,
	)
	data.Keeper = keeper.NewKeeper(
		data.Cdc,
		data.Keys[types.StoreKey],
		runtime.NewKVStoreService(data.Keys[types.StoreKey]),
		data.AccountKeeper,
		data.BankKeeper,
		data.OperatorsKeeper,
		data.PoolsKeeper,
		data.ServicesKeeper,
		data.RestakingKeeper,
		data.LiquidVestingModuleAddress.String(),
		data.AuthorityAddress,
	)
	// Set bank hooks
	data.BankKeeper.AppendSendRestriction(data.Keeper.SendRestrictionFn)

	// Set ibc hooks
	mockIBCMiddleware := mockIBCMiddleware{}
	middleware := ibchooks.NewICS4Middleware(mockIBCMiddleware, hooks.NewIBCHooks(data.Keeper))
	data.IBCMiddleware = ibchooks.NewIBCMiddleware(mockIBCMiddleware, middleware, data.IBCHooksKeeper)

	account := data.AccountKeeper.GetModuleAccount(data.Context, types.ModuleName)
	require.NotNil(t, account)
	data.AccountKeeper.SetModuleAccount(data.Context, account)

	return data
}
