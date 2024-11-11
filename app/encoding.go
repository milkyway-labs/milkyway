package milkyway

import (
	"os"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/log"
	"cosmossdk.io/x/tx/signing"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codecaddress "github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/std"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/gogoproto/proto"

	"github.com/initia-labs/initia/app/params"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	// Build the interface registry specifying the validator and address codec
	interfaceRegistry, _ := codectypes.NewInterfaceRegistryWithOptions(codectypes.InterfaceRegistryOptions{
		ProtoFiles: proto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec:          codecaddress.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
			ValidatorAddressCodec: codecaddress.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		},
	})
	cdc := codec.NewProtoCodec(interfaceRegistry)
	amino := codec.NewLegacyAmino()

	// Build the encoding config
	encodingConfig := params.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             cdc,
		TxConfig:          tx.NewTxConfig(cdc, tx.DefaultSignModes),
		Amino:             amino,
	}

	// Register the stdlib and module basics
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return encodingConfig
}

// MakeCodecs creates the necessary testing codecs for Amino and Protobuf
func MakeCodecs() (codec.Codec, *codec.LegacyAmino) {
	encodingConfig := MakeEncodingConfig()
	return encodingConfig.Codec, encodingConfig.Amino
}

// AutoCLIOpts returns the options for the auto-generated CLI
func AutoCLIOpts() (autocli.AppOptions, error) {
	// Create a temp app with a temp home dir
	tempDir, err := os.MkdirTemp(os.TempDir(), "milkyway-autocli-opts")
	if err != nil {
		return autocli.AppOptions{}, err
	}

	tempApp := NewMilkyWayApp(log.NewNopLogger(), dbm.NewMemDB(), dbm.NewMemDB(), nil, true, []wasmkeeper.Option{}, simtestutil.NewAppOptionsWithFlagHome(tempDir))
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
	}, nil
}
