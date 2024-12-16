package milkyway

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"

	"github.com/milkyway-labs/milkyway/v7/app/keepers"
	"github.com/milkyway-labs/milkyway/v7/app/params"
)

// MakeEncodingConfig creates an EncodingConfig.
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	keepers.AppModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

// MakeCodecs constructs the *codec.Codec and *codec.LegacyAmino instances used by MilkyWayApp.
func MakeCodecs() (codec.Codec, *codec.LegacyAmino) {
	encodingCfg := MakeEncodingConfig()
	return encodingCfg.Marshaler, encodingCfg.Amino
}
