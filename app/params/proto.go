package params

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
)

// MakeEncodingConfig creates an EncodingConfig for an amino based test configuration.
func MakeEncodingConfig() EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := testutil.CodecOptions{AccAddressPrefix: "milk", ValAddressPrefix: "milkvaloper"}.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}

// MakeCodecs creates the necessary testing codecs for Amino and Protobuf
func MakeCodecs() (codec.Codec, *codec.LegacyAmino) {
	encodingConfig := MakeEncodingConfig()
	return encodingConfig.Marshaler, encodingConfig.Amino
}
