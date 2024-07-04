package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgDelegateService{}, "milkyway/MsgDelegateService")
	legacy.RegisterAminoMsg(cdc, &MsgDelegatePool{}, "milkyway/MsgDelegatePool")
	legacy.RegisterAminoMsg(cdc, &MsgDelegateOperator{}, "milkyway/MsgDelegateOperator")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "milkyway/restaking/MsgUpdateParams")
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDelegatePool{},
		&MsgDelegateService{},
		&MsgDelegateOperator{},
		&MsgUpdateParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	// AminoCdc references the global x/services module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/services and
	// defined at the application level.
	AminoCdc = codec.NewLegacyAmino()
)

func init() {
	RegisterLegacyAminoCodec(AminoCdc)
	cryptocodec.RegisterCrypto(AminoCdc)
	sdk.RegisterLegacyAminoCodec(AminoCdc)
}
