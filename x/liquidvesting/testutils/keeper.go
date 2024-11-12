package testutils

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibchookstypes "github.com/initia-labs/initia/x/ibc-hooks/types"
	"github.com/stretchr/testify/require"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	appkeepers "github.com/milkyway-labs/milkyway/app/keepers"
	"github.com/milkyway-labs/milkyway/testutils/storetesting"
	"github.com/milkyway-labs/milkyway/x/liquidvesting"
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
	IBCMiddleware   porttypes.IBCModule
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

	// Module addresses
	data.LiquidVestingModuleAddress = authtypes.NewModuleAddress(types.ModuleName)

	// Build keepers
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
	var ibcStack porttypes.IBCModule = mockIBCMiddleware{}
	data.IBCMiddleware = liquidvesting.NewIBCMiddleware(ibcStack, data.Keeper)

	account := data.AccountKeeper.GetModuleAccount(data.Context, types.ModuleName)
	require.NotNil(t, account)
	data.AccountKeeper.SetModuleAccount(data.Context, account)

	return data
}
