package app

import (
	"testing"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	"github.com/initia-labs/initia/app/params"
)

type EncodingConfigCreator func() params.EncodingConfig

// makeCodecs creates the necessary codecs for Amino and Protobuf using the provided EncodingConfigCreator
func makeCodecs(createConfig EncodingConfigCreator) (codec.Codec, *codec.LegacyAmino) {
	encodingConfig := createConfig()
	return encodingConfig.Codec, encodingConfig.Amino
}

// MakeCodecs creates the necessary codecs for Amino and Protobuf
func MakeCodecs() (codec.Codec, *codec.LegacyAmino) {
	return makeCodecs(MakeEncodingConfig)
}

// makeEncodingConfig creates an EncodingConfig instance by creating a temporary app with the given options
func makeEncodingConfig(appOptions types.AppOptions) params.EncodingConfig {
	tempApp := NewMilkyWayApp(log.NewNopLogger(), dbm.NewMemDB(), dbm.NewMemDB(), nil, true, []wasmkeeper.Option{}, appOptions)
	encodingConfig := params.EncodingConfig{
		InterfaceRegistry: tempApp.InterfaceRegistry(),
		Codec:             tempApp.AppCodec(),
		TxConfig:          tempApp.TxConfig(),
		Amino:             tempApp.LegacyAmino(),
	}

	return encodingConfig
}

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	return makeEncodingConfig(EmptyAppOptions{})
}

// MakeTestCodecs creates the necessary testing codecs for Amino and Protobuf
func MakeTestCodecs(t *testing.T) (codec.Codec, *codec.LegacyAmino) {
	t.Helper()
	return makeCodecs(func() params.EncodingConfig {
		return MakeTestEncodingConfig(t)
	})
}

// MakeTestEncodingConfig creates an EncodingConfig for testing
func MakeTestEncodingConfig(t *testing.T) params.EncodingConfig {
	t.Helper()
	return makeEncodingConfig(simtestutil.NewAppOptionsWithFlagHome(t.TempDir()))
}

func AutoCliOpts() autocli.AppOptions {
	tempApp := NewMilkyWayApp(log.NewNopLogger(), dbm.NewMemDB(), dbm.NewMemDB(), nil, true, []wasmkeeper.Option{}, EmptyAppOptions{})
	modules := make(map[string]appmodule.AppModule, 0)
	for _, m := range tempApp.ModuleManager.Modules {
		if moduleWithName, ok := m.(module.HasName); ok {
			moduleName := moduleWithName.Name()
			if appModule, ok := moduleWithName.(appmodule.AppModule); ok {
				modules[moduleName] = appModule
			}
		}
	}

	return autocli.AppOptions{
		Modules:               modules,
		ModuleOptions:         runtimeservices.ExtractAutoCLIOptions(tempApp.ModuleManager.Modules),
		AddressCodec:          authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		ValidatorAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		ConsensusAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	}
}

func BasicManager() module.BasicManager {
	tempApp := NewMilkyWayApp(log.NewNopLogger(), dbm.NewMemDB(), dbm.NewMemDB(), nil, true, []wasmkeeper.Option{}, EmptyAppOptions{})
	return tempApp.BasicModuleManager
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	if o == flags.FlagHome {
		return DefaultNodeHome
	}

	return nil
}
