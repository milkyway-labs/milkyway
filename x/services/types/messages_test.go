package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

var msgCreateService = types.NewMsgCreateService(
	"MilkyWay",
	"MilkyWay is an AVS of a restaking platform",
	"https://milkyway.com",
	"https://milkyway.com/logo.png",
	"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
)

func TestMsgCreateService_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgCreateService
		shouldErr bool
	}{
		{
			name: "do-not-modify name returns error",
			msg: types.NewMsgCreateService(
				types.DoNotModify,
				msgCreateService.Description,
				msgCreateService.Website,
				msgCreateService.PictureURL,
				msgCreateService.Sender,
			),
			shouldErr: true,
		},
		{
			name: "empty name returns error",
			msg: types.NewMsgCreateService(
				"",
				msgCreateService.Description,
				msgCreateService.Website,
				msgCreateService.PictureURL,
				msgCreateService.Sender,
			),
			shouldErr: true,
		},
		{
			name: "do-not-modify description returns error",
			msg: types.NewMsgCreateService(
				msgCreateService.Name,
				types.DoNotModify,
				msgCreateService.Website,
				msgCreateService.PictureURL,
				msgCreateService.Sender,
			),
			shouldErr: true,
		},
		{
			name: "do-not-modify website returns error",
			msg: types.NewMsgCreateService(
				msgCreateService.Name,
				msgCreateService.Description,
				types.DoNotModify,
				msgCreateService.PictureURL,
				msgCreateService.Sender,
			),
			shouldErr: true,
		},
		{
			name: "do-not-modify picture URL returns error",
			msg: types.NewMsgCreateService(
				msgCreateService.Name,
				msgCreateService.Description,
				msgCreateService.Website,
				types.DoNotModify,
				msgCreateService.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid sender address returns error",
			msg: types.NewMsgCreateService(
				msgCreateService.Name,
				msgCreateService.Description,
				msgCreateService.Website,
				msgCreateService.PictureURL,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name: "valid message returns no error",
			msg:  msgCreateService,
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

func TestMsgCreateService_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgCreateService","value":{"description":"MilkyWay is an AVS of a restaking platform","name":"MilkyWay","picture_url":"https://milkyway.com/logo.png","sender":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd","website":"https://milkyway.com"}}`
	require.Equal(t, expected, string(msgCreateService.GetSignBytes()))
}

func TestMsgCreateService_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgCreateService.Sender)
	require.Equal(t, []sdk.AccAddress{addr}, msgCreateService.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgUpdateService = types.NewMsgUpdateService(
	1,
	"MilkyWay",
	"MilkyWay is an AVS of a restaking platform",
	"https://milkyway.com",
	"https://milkyway.com/logo.png",
	"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
)

func TestMsgUpdateService_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgUpdateService
		shouldErr bool
	}{
		{
			name: "invalid ID returns error",
			msg: types.NewMsgUpdateService(
				0,
				msgUpdateService.Name,
				msgUpdateService.Description,
				msgUpdateService.Website,
				msgUpdateService.PictureURL,
				msgUpdateService.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid sender address returns error",
			msg: types.NewMsgUpdateService(
				msgUpdateService.ServiceID,
				msgUpdateService.Name,
				msgUpdateService.Description,
				msgUpdateService.Website,
				msgUpdateService.PictureURL,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name: "valid message returns no error",
			msg:  msgUpdateService,
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

func TestMsgUpdateService_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgUpdateService","value":{"description":"MilkyWay is an AVS of a restaking platform","name":"MilkyWay","picture_url":"https://milkyway.com/logo.png","sender":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd","service_id":1,"website":"https://milkyway.com"}}`
	require.Equal(t, expected, string(msgUpdateService.GetSignBytes()))
}

func TestMsgUpdateService_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgUpdateService.Sender)
	require.Equal(t, []sdk.AccAddress{addr}, msgUpdateService.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgDeactivateService = types.NewMsgDeactivateService(
	1,
	"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
)

func TestMsgDeactivateService_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgDeactivateService
		shouldErr bool
	}{
		{
			name: "invalid ID returns error",
			msg: types.NewMsgDeactivateService(
				0,
				msgDeactivateService.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid sender address returns error",
			msg: types.NewMsgDeactivateService(
				msgDeactivateService.ServiceID,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name: "valid message returns no error",
			msg:  msgDeactivateService,
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

func TestMsgDeactivateService_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgDeactivateService","value":{"sender":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd","service_id":1}}`
	require.Equal(t, expected, string(msgDeactivateService.GetSignBytes()))
}

func TestMsgDeactivateService_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgDeactivateService.Sender)
	require.Equal(t, []sdk.AccAddress{addr}, msgDeactivateService.GetSigners())
}
