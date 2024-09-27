package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgMintVestedRepresentation{}, "milkyway/MsgMintVestedRepresentation")
	legacy.RegisterAminoMsg(cdc, &MsgBurnVestedRepresentation{}, "milkyway/MsgBurnVestedRepresentation")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawInsuranceFund{}, "milkyway/MsgWithdrawInsuranceFund")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "milkyway/liquidvesting/MsgUpdateParams")
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgMintVestedRepresentation{},
		&MsgBurnVestedRepresentation{},
		&MsgWithdrawInsuranceFund{},
		&MsgUpdateParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
