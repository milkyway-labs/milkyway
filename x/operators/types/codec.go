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
	legacy.RegisterAminoMsg(cdc, &MsgRegisterOperator{}, "milkyway/MsgRegisterOperator")
	legacy.RegisterAminoMsg(cdc, &MsgSetOperatorParams{}, "milkyway/MsgSetOperatorParams")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateOperator{}, "milkyway/MsgUpdateOperator")
	legacy.RegisterAminoMsg(cdc, &MsgDeactivateOperator{}, "milkyway/MsgDeactivateOperator")
	legacy.RegisterAminoMsg(cdc, &MsgReactivateOperator{}, "milkyway/MsgReactivateOperator")
	legacy.RegisterAminoMsg(cdc, &MsgDeleteOperator{}, "milkyway/MsgDeleteOperator")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "milkyway/operators/MsgUpdateParams")
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterOperator{},
		&MsgSetOperatorParams{},
		&MsgUpdateOperator{},
		&MsgDeactivateOperator{},
		&MsgReactivateOperator{},
		&MsgDeleteOperator{},
		&MsgUpdateParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// AminoCdc references the global x/avs module codec. Note, the codec should
// ONLY be used in certain instances of tests and for JSON encoding as Amino is
// still used for that purpose.
//
// The actual codec used for serialization should be provided to x/avs and
// defined at the application level.
var AminoCdc = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(AminoCdc)
	cryptocodec.RegisterCrypto(AminoCdc)
	sdk.RegisterLegacyAminoCodec(AminoCdc)
}
