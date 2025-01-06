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
	legacy.RegisterAminoMsg(cdc, &MsgCreateRewardsPlan{}, "milkyway/MsgCreateRewardsPlan")
	legacy.RegisterAminoMsg(cdc, &MsgEditRewardsPlan{}, "milkyway/MsgEditRewardsPlan")
	legacy.RegisterAminoMsg(cdc, &MsgSetWithdrawAddress{}, "milkyway/MsgSetWithdrawAddress")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawDelegatorReward{}, "milkyway/MsgWithdrawDelegatorReward")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawOperatorCommission{}, "milkyway/MsgWithdrawOperatorCommission")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "milkyway/rewards/MsgUpdateParams")

	cdc.RegisterInterface((*DistributionType)(nil), nil)
	cdc.RegisterConcrete(&DistributionTypeBasic{}, "milkyway/DistributionTypeBasic", nil)
	cdc.RegisterConcrete(&DistributionTypeWeighted{}, "milkyway/DistributionTypeWeighted", nil)
	cdc.RegisterConcrete(&DistributionTypeEgalitarian{}, "milkyway/DistributionTypeEgalitarian", nil)

	cdc.RegisterInterface((*UsersDistributionType)(nil), nil)
	cdc.RegisterConcrete(&UsersDistributionTypeBasic{}, "milkyway/UsersDistributionTypeBasic", nil)
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
		"milkyway.rewards.v2.DistributionType",
		(*DistributionType)(nil),
		&DistributionTypeBasic{},
		&DistributionTypeWeighted{},
		&DistributionTypeEgalitarian{},
	)
	registry.RegisterInterface(
		"milkyway.rewards.v2.UsersDistributionType",
		(*UsersDistributionType)(nil),
		&UsersDistributionTypeBasic{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// AminoCdc references the global x/rewards module codec. Note, the codec should
// ONLY be used in certain instances of tests and for JSON encoding as Amino is
// still used for that purpose.
//
// The actual codec used for serialization should be provided to x/rewards and
// defined at the application level.
var AminoCdc = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(AminoCdc)
	cryptocodec.RegisterCrypto(AminoCdc)
	sdk.RegisterLegacyAminoCodec(AminoCdc)
}
