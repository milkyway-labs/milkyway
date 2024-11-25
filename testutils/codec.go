package testutils

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectestutil "github.com/cosmos/cosmos-sdk/codec/testutil"
)

// MakeCodecs constructs the *codec.Codec and *codec.LegacyAmino instances that can be used inside tests
func MakeCodecs() (codec.Codec, *codec.LegacyAmino) {
	interfaceRegistry := codectestutil.CodecOptions{AccAddressPrefix: "cosmos", ValAddressPrefix: "cosmosvaloper"}.NewInterfaceRegistry()
	return codec.NewProtoCodec(interfaceRegistry), codec.NewLegacyAmino()
}
