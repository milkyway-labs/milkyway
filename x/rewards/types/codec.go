package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateRewardsPlan{}, "milkyway/MsgCreateRewardsPlan")
	legacy.RegisterAminoMsg(cdc, &MsgEditRewardsPlan{}, "milkyway/MsgEditRewardsPlan")
	legacy.RegisterAminoMsg(cdc, &MsgSetWithdrawAddress{}, "milkyway/MsgSetWithdrawAddress")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawDelegatorReward{}, "milkyway/MsgWithdrawDelegatorReward")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawOperatorCommission{}, "milkyway/MsgWithdrawOperatorCommission")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "milkyway/rewards/MsgUpdateParams")
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateRewardsPlan{},
		&MsgEditRewardsPlan{},
		&MsgSetWithdrawAddress{},
		&MsgWithdrawDelegatorReward{},
		&MsgWithdrawOperatorCommission{},
		&MsgUpdateParams{},
	)
	registry.RegisterInterface(
		"milkyway.rewards.v1.DistributionType",
		(*DistributionType)(nil),
		&DistributionTypeBasic{},
		&DistributionTypeWeighted{},
		&DistributionTypeEgalitarian{},
	)
	registry.RegisterInterface(
		"milkyway.rewards.v1.UsersDistributionType",
		(*UsersDistributionType)(nil),
		&UsersDistributionTypeBasic{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
