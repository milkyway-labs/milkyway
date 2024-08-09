package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateRewardsPlan{}, "rewards/MsgCreateRewardsPlan")
	legacy.RegisterAminoMsg(cdc, &MsgSetWithdrawAddress{}, "rewards/MsgSetWithdrawAddress")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawDelegatorReward{}, "rewards/MsgWithdrawDelegatorReward")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawOperatorCommission{}, "rewards/MsgWithdrawOperatorCommission")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "rewards/MsgUpdateParams")
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateRewardsPlan{},
		&MsgSetWithdrawAddress{},
		&MsgWithdrawDelegatorReward{},
		&MsgWithdrawOperatorCommission{},
		&MsgUpdateParams{},
	)
	registry.RegisterInterface(
		"milkyway.rewards.v1.PoolsDistributionType",
		(*PoolsDistributionType)(nil),
		&PoolsDistributionTypeBasic{},
		&PoolsDistributionTypeWeighted{},
		&PoolsDistributionTypeEgalitarian{},
	)
	registry.RegisterInterface(
		"milkyway.rewards.v1.OperatorsDistributionType",
		(*OperatorsDistributionType)(nil),
		&OperatorsDistributionTypeBasic{},
		&OperatorsDistributionTypeWeighted{},
		&OperatorsDistributionTypeEgalitarian{},
	)
	registry.RegisterInterface(
		"milkyway.rewards.v1.UsersDistributionType",
		(*UsersDistributionType)(nil),
		&UsersDistributionTypeBasic{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
