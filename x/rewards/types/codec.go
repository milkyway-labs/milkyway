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
	legacy.RegisterAminoMsg(cdc, &MsgCreateRewardsPlan{}, "rewards/MsgCreateRewardsPlan")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawDelegationReward{}, "rewards/MsgWithdrawDelegationReward")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawOperatorCommission{}, "rewards/MsgWithdrawOperatorCommission")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "rewards/MsgUpdateParams")
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateRewardsPlan{},
		&MsgWithdrawDelegationReward{},
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

var (
	// AminoCdc references the global x/avs module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/avs and
	// defined at the application level.
	AminoCdc = codec.NewLegacyAmino()
)

func init() {
	RegisterLegacyAminoCodec(AminoCdc)
	cryptocodec.RegisterCrypto(AminoCdc)
	sdk.RegisterLegacyAminoCodec(AminoCdc)
}
