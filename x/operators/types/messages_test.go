package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

var msgRegisterOperator = types.NewMsgRegisterOperator(
	"MilkyWay Operator",
	"https://milkyway.com",
	"https://milkyway.com/picture",
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
	expected := `{"type":"milkyway/MsgRegisterOperator","value":{"moniker":"MilkyWay Operator","picture_url":"https://milkyway.com/picture","sender":"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4","website":"https://milkyway.com"}}`
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
				msgUpdateOperator.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid sender address returns error",
			msg: types.NewMsgDeactivateOperator(
				msgUpdateOperator.OperatorID,
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