package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgMintStakingRepresentation{}, "milkyway/MsgMintStakingRepresentation")
	legacy.RegisterAminoMsg(cdc, &MsgBurnStakingRepresentation{}, "milkyway/MsgBurnStakingRepresentation")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "milkyway/liquidvesting/MsgUpdateParams")
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgMintStakingRepresentation{},
		&MsgBurnStakingRepresentation{},
		&MsgUpdateParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
