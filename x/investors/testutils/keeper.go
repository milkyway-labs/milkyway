package testutils

import (
	"testing"

	corestoretypes "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v12/testutils/storetesting"
	"github.com/milkyway-labs/milkyway/v12/x/investors/keeper"
	"github.com/milkyway-labs/milkyway/v12/x/investors/types"
)

type KeeperTestData struct {
	storetesting.BaseKeeperTestData

	StoreService corestoretypes.KVStoreService

	StakingKeeper *stakingkeeper.Keeper
	DistrKeeper   distrkeeper.Keeper
	Keeper        *keeper.Keeper
}

func NewKeeperTestData(t *testing.T) KeeperTestData {
	// Build the base data
	data := KeeperTestData{
		BaseKeeperTestData: storetesting.NewBaseKeeperTestData(t, []string{
			authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey, distrtypes.StoreKey,
			types.StoreKey,
		}),
	}

	// Setup the keys
	data.StoreService = runtime.NewKVStoreService(data.Keys[types.StoreKey])

	// Build the keepers
	data.StakingKeeper = stakingkeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[stakingtypes.StoreKey]),
		data.AccountKeeper,
		data.BankKeeper,
		data.AuthorityAddress,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)
	data.DistrKeeper = distrkeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[distrtypes.StoreKey]),
		data.AccountKeeper,
		data.BankKeeper,
		data.StakingKeeper,
		authtypes.FeeCollectorName,
		data.AuthorityAddress,
	)
	data.Keeper = keeper.NewKeeper(
		data.Cdc,
		data.StoreService,
		data.AccountKeeper,
		data.BankKeeper,
		data.StakingKeeper,
		&data.DistrKeeper,
		data.AuthorityAddress,
	)

	data.StakingKeeper.SetHooks(data.DistrKeeper.Hooks())
	data.DistrKeeper.SetHooks(data.Keeper.Hooks())

	data.StakingKeeper.InitGenesis(data.Context, stakingtypes.DefaultGenesisState())
	data.DistrKeeper.InitGenesis(data.Context, *distrtypes.DefaultGenesisState())
	require.NoError(t, data.Keeper.InitGenesis(data.Context, types.DefaultGenesisState()))

	return data
}
