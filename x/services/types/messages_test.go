package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
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
	expected := `{"type":"services/MsgCreateService","value":{"description":"MilkyWay is an AVS of a restaking platform","name":"MilkyWay","picture_url":"https://milkyway.com/logo.png","sender":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd","website":"https://milkyway.com"}}`
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
	expected := `{"type":"services/MsgUpdateService","value":{"description":"MilkyWay is an AVS of a restaking platform","name":"MilkyWay","picture_url":"https://milkyway.com/logo.png","sender":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd","service_id":1,"website":"https://milkyway.com"}}`
	require.Equal(t, expected, string(msgUpdateService.GetSignBytes()))
}

func TestMsgUpdateService_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgUpdateService.Sender)
	require.Equal(t, []sdk.AccAddress{addr}, msgUpdateService.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgActivateService = types.NewMsgActivateService(
	1,
	"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
)

func TestMsgActivateService_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgActivateService
		shouldErr bool
	}{
		{
			name: "invalid service id returns error",
			msg: types.NewMsgActivateService(
				0,
				msgActivateService.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid sender address returns error",
			msg: types.NewMsgActivateService(
				msgActivateService.ServiceID,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name: "valid message returns no error",
			msg:  msgActivateService,
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

func TestMsgActivateService_GetSignBytes(t *testing.T) {
	expected := `{"type":"services/MsgActivateService","value":{"sender":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd","service_id":1}}`
	require.Equal(t, expected, string(msgActivateService.GetSignBytes()))
}

func TestMsgActivateService_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgActivateService.Sender)
	require.Equal(t, []sdk.AccAddress{addr}, msgActivateService.GetSigners())
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
			name: "invalid service id returns error",
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
	expected := `{"type":"services/MsgDeactivateService","value":{"sender":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd","service_id":1}}`
	require.Equal(t, expected, string(msgDeactivateService.GetSignBytes()))
}

func TestMsgDeactivateService_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgDeactivateService.Sender)
	require.Equal(t, []sdk.AccAddress{addr}, msgDeactivateService.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgUpdateParams = types.NewMsgUpdateParams(
	types.DefaultParams(),
	"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
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
				types.NewParams(sdk.Coins{sdk.Coin{Denom: "invalid!", Amount: sdkmath.NewInt(100_000_000)}}),
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
			name: "valid message returns no error",
			msg:  msgUpdateParams,
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
	expected := `{"type":"services/MsgUpdateParams","value":{"authority":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd","params":{"service_registration_fee":[]}}}`
	require.Equal(t, expected, string(msgUpdateParams.GetSignBytes()))
}

func TestMsgUpdateParams_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgUpdateParams.Authority)
	require.Equal(t, []sdk.AccAddress{addr}, msgDeactivateService.GetSigners())
}
