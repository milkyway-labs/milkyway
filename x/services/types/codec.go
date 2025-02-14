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
	legacy.RegisterAminoMsg(cdc, &MsgCreateService{}, "milkyway/MsgCreateService")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateService{}, "milkyway/MsgUpdateService")
	legacy.RegisterAminoMsg(cdc, &MsgActivateService{}, "milkyway/MsgActivateService")
	legacy.RegisterAminoMsg(cdc, &MsgDeactivateService{}, "milkyway/MsgDeactivateService")
	legacy.RegisterAminoMsg(cdc, &MsgDeleteService{}, "milkyway/MsgDeleteService")
	legacy.RegisterAminoMsg(cdc, &MsgSetServiceParams{}, "milkyway/MsgSetServiceParams")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "milkyway/services/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgAccreditService{}, "milkyway/MsgAccreditService")
	legacy.RegisterAminoMsg(cdc, &MsgRevokeServiceAccreditation{}, "milkyway/MsgRevokeServiceAccreditation")
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateService{},
		&MsgUpdateService{},
		&MsgActivateService{},
		&MsgDeactivateService{},
		&MsgDeleteService{},
		&MsgSetServiceParams{},
		&MsgUpdateParams{},
		&MsgAccreditService{},
		&MsgRevokeServiceAccreditation{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// AminoCdc references the global x/services module codec. Note, the codec should
// ONLY be used in certain instances of tests and for JSON encoding as Amino is
// still used for that purpose.
//
// The actual codec used for serialization should be provided to x/services and
// defined at the application level.
var AminoCdc = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(AminoCdc)
	cryptocodec.RegisterCrypto(AminoCdc)
	sdk.RegisterLegacyAminoCodec(AminoCdc)
}
