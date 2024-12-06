package types_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v3/x/operators/types"
)

var msgRegisterOperator = types.NewMsgRegisterOperator(
	"MilkyWay Operator",
	"https://milkyway.com",
	"https://milkyway.com/picture",
	sdk.NewCoins(sdk.NewInt64Coin("uatom", 100_000_000)),
	"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
)

func TestMsgRegisterOperator_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgRegisterOperator
		shouldErr bool
	}{
		{
			name: "do-not-modify moniker returns error",
			msg: types.NewMsgRegisterOperator(
				types.DoNotModify,
				msgRegisterOperator.Website,
				msgRegisterOperator.PictureURL,
				msgRegisterOperator.FeeAmount,
				msgRegisterOperator.Sender,
			),
			shouldErr: true,
		},
		{
			name: "empty moniker returns error",
			msg: types.NewMsgRegisterOperator(
				"",
				msgRegisterOperator.Website,
				msgRegisterOperator.PictureURL,
				msgRegisterOperator.FeeAmount,
				msgRegisterOperator.Sender,
			),
			shouldErr: true,
		},
		{
			name: "do-not-modify website returns error",
			msg: types.NewMsgRegisterOperator(
				msgRegisterOperator.Moniker,
				types.DoNotModify,
				msgRegisterOperator.PictureURL,
				msgRegisterOperator.FeeAmount,
				msgRegisterOperator.Sender,
			),
			shouldErr: true,
		},
		{
			name: "do-not-modify picture URL returns error",
			msg: types.NewMsgRegisterOperator(
				msgRegisterOperator.Moniker,
				msgRegisterOperator.Website,
				types.DoNotModify,
				msgRegisterOperator.FeeAmount,
				msgRegisterOperator.Sender,
			),
			shouldErr: true,
		},
		{
			name: "do-not-modify picture URL returns error",
			msg: types.NewMsgRegisterOperator(
				msgRegisterOperator.Moniker,
				msgRegisterOperator.Website,
				msgRegisterOperator.PictureURL,
				sdk.Coins{sdk.Coin{Denom: "invalid", Amount: sdkmath.NewInt(-100)}},
				msgRegisterOperator.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid sender address returns error",
			msg: types.NewMsgRegisterOperator(
				msgRegisterOperator.Moniker,
				msgRegisterOperator.Website,
				msgRegisterOperator.PictureURL,
				msgRegisterOperator.FeeAmount,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name:      "valid message returns no error",
			msg:       msgRegisterOperator,
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgRegisterOperator_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgRegisterOperator","value":{"fee_amount":[{"amount":"100000000","denom":"uatom"}],"moniker":"MilkyWay Operator","picture_url":"https://milkyway.com/picture","sender":"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4","website":"https://milkyway.com"}}`
	require.Equal(t, expected, string(msgRegisterOperator.GetSignBytes()))
}

func TestMsgRegisterOperator_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgRegisterOperator.Sender)
	require.Equal(t, []sdk.AccAddress{addr}, msgRegisterOperator.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgUpdateOperator = types.NewMsgUpdateOperator(
	1,
	"MilkyWay Operator",
	"https://milkyway.com",
	"https://milkyway.com/picture",
	"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
)

func TestMsgUpdateOperator_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgUpdateOperator
		shouldErr bool
	}{
		{
			name: "invalid operator id returns error",
			msg: types.NewMsgUpdateOperator(
				0,
				types.DoNotModify,
				msgUpdateOperator.Website,
				msgUpdateOperator.PictureURL,
				msgUpdateOperator.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid sender address returns error",
			msg: types.NewMsgUpdateOperator(
				msgUpdateOperator.OperatorID,
				msgUpdateOperator.Moniker,
				msgUpdateOperator.Website,
				msgUpdateOperator.PictureURL,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name:      "valid message returns no error",
			msg:       msgUpdateOperator,
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgUpdateOperator_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgUpdateOperator","value":{"moniker":"MilkyWay Operator","operator_id":1,"picture_url":"https://milkyway.com/picture","sender":"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4","website":"https://milkyway.com"}}`
	require.Equal(t, expected, string(msgUpdateOperator.GetSignBytes()))
}

func TestMsgUpdateOperator_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgUpdateOperator.Sender)
	require.Equal(t, []sdk.AccAddress{addr}, msgUpdateOperator.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgDeactivateOperator = types.NewMsgDeactivateOperator(
	1,
	"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
)

func TestMsgDeactivateOperator_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgDeactivateOperator
		shouldErr bool
	}{
		{
			name: "invalid operator id returns error",
			msg: types.NewMsgDeactivateOperator(
				0,
				msgDeactivateOperator.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid sender address returns error",
			msg: types.NewMsgDeactivateOperator(
				msgDeactivateOperator.OperatorID,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name:      "valid message returns no error",
			msg:       msgDeactivateOperator,
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgDeactivateOperator_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgDeactivateOperator","value":{"operator_id":1,"sender":"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"}}`
	require.Equal(t, expected, string(msgDeactivateOperator.GetSignBytes()))
}

func TestMsgDeactivateOperator_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgDeactivateOperator.Sender)
	require.Equal(t, []sdk.AccAddress{addr}, msgDeactivateOperator.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgTransferOperatorOwnership = types.NewMsgTransferOperatorOwnership(
	1,
	"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
	"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
)

func TestMsgTransferOperatorOwnership_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgTransferOperatorOwnership
		shouldErr bool
	}{
		{
			name: "invalid operator id returns error",
			msg: types.NewMsgTransferOperatorOwnership(
				0,
				msgTransferOperatorOwnership.NewAdmin,
				msgTransferOperatorOwnership.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid new admin address returns error",
			msg: types.NewMsgTransferOperatorOwnership(
				msgTransferOperatorOwnership.OperatorID,
				"invalid",
				msgTransferOperatorOwnership.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid sender address returns error",
			msg: types.NewMsgTransferOperatorOwnership(
				msgTransferOperatorOwnership.OperatorID,
				msgTransferOperatorOwnership.NewAdmin,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name: "valid message returns no error",
			msg:  msgTransferOperatorOwnership,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgTransferOperatorOwnership_GetSignBytes(t *testing.T) {
	expected := `{"new_admin":"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn","operator_id":1,"sender":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"}`
	require.Equal(t, expected, string(msgTransferOperatorOwnership.GetSignBytes()))
}

func TestMsgTransferOperatorOwnership_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgTransferOperatorOwnership.Sender)
	require.Equal(t, []sdk.AccAddress{addr}, msgTransferOperatorOwnership.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgUpdateParams = types.NewMsgUpdateParams(
	types.NewParams(
		sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
		24*time.Hour,
	),
	"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
)

func TestMsgUpdateParams_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgUpdateParams
		shouldErr bool
	}{
		{
			name: "invalid params return error",
			msg: types.NewMsgUpdateParams(
				types.NewParams(nil, 0),
				msgUpdateParams.Authority,
			),
			shouldErr: true,
		},
		{
			name: "invalid authority address returns error",
			msg: types.NewMsgUpdateParams(
				msgUpdateParams.Params,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name:      "valid message returns no error",
			msg:       msgUpdateParams,
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgUpdateParams_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/operators/MsgUpdateParams","value":{"authority":"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4","params":{"deactivation_time":"86400000000000","operator_registration_fee":[{"amount":"100000000","denom":"uatom"}]}}}`
	require.Equal(t, expected, string(msgUpdateParams.GetSignBytes()))
}

func TestMsgUpdateParams_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgUpdateParams.Authority)
	require.Equal(t, []sdk.AccAddress{addr}, msgUpdateParams.GetSigners())
}
