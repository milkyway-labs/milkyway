package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v3/x/restaking/types"
)

var msgUpdateOperatorParams = types.NewMsgJoinService(1, 1, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd")

func TestMsgUpdateOperatorParams_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgJoinService
		shouldErr bool
	}{
		{
			name: "invalid operator id returns error",
			msg: types.NewMsgJoinService(
				0,
				1,
				msgUpdateOperatorParams.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid service id returns error",
			msg: types.NewMsgJoinService(
				msgUpdateOperatorParams.OperatorID,
				0,
				msgUpdateOperatorParams.Sender,
			),
			shouldErr: true,
		},
		{
			name: "valid message returns no error",
			msg:  msgUpdateOperatorParams,
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

// --------------------------------------------------------------------------------------------------------------------

var msgDelegatePool = types.NewMsgDelegatePool(
	sdk.NewCoin("umilk", sdkmath.NewInt(100_000_000)),
	"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
)

func TestMsgDelegatePool_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgDelegatePool
		shouldErr bool
	}{
		{
			name: "invalid amount returns error",
			msg: types.NewMsgDelegatePool(
				sdk.Coin{Denom: "invalid!", Amount: sdkmath.NewInt(100_000_000)},
				msgDelegatePool.Delegator,
			),
			shouldErr: true,
		},
		{
			name: "invalid delegator address returns error",
			msg: types.NewMsgDelegatePool(
				msgDelegatePool.Amount,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name: "valid message returns no error",
			msg:  msgDelegatePool,
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

func TestMsgDelegatePool_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgDelegatePool","value":{"amount":{"amount":"100000000","denom":"umilk"},"delegator":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"}}`
	require.Equal(t, expected, string(msgDelegatePool.GetSignBytes()))
}

func TestMsgDelegatePool_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgDelegatePool.Delegator)
	require.Equal(t, []sdk.AccAddress{addr}, msgDelegatePool.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgDelegateOperator = types.NewMsgDelegateOperator(
	1,
	sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100_000_000))),
	"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
)

func TestMsgDelegateOperator_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgDelegateOperator
		shouldErr bool
	}{
		{
			name: "invalid operator id returns error",
			msg: types.NewMsgDelegateOperator(
				0,
				msgDelegateOperator.Amount,
				msgDelegateOperator.Delegator,
			),
			shouldErr: true,
		},
		{
			name: "invalid amount returns error",
			msg: types.NewMsgDelegateOperator(
				msgDelegateOperator.OperatorID,
				sdk.Coins{sdk.Coin{Denom: "invalid!", Amount: sdkmath.NewInt(100_000_000)}},
				msgDelegateOperator.Delegator,
			),
			shouldErr: true,
		},
		{
			name: "invalid delegator address returns error",
			msg: types.NewMsgDelegateOperator(
				msgDelegateOperator.OperatorID,
				msgDelegateOperator.Amount,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name: "valid message returns no error",
			msg:  msgDelegateOperator,
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

func TestMsgDelegateOperator_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgDelegateOperator","value":{"amount":[{"amount":"100000000","denom":"umilk"}],"delegator":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd","operator_id":1}}`
	require.Equal(t, expected, string(msgDelegateOperator.GetSignBytes()))
}

func TestMsgDelegateOperator_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgDelegateOperator.Delegator)
	require.Equal(t, []sdk.AccAddress{addr}, msgDelegateOperator.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgDelegateService = types.NewMsgDelegateService(
	1,
	sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100_000_000))),
	"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
)

func TestMsgDelegateService_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgDelegateService
		shouldErr bool
	}{
		{
			name: "invalid service id returns error",
			msg: types.NewMsgDelegateService(
				0,
				msgDelegateService.Amount,
				msgDelegateService.Delegator,
			),
			shouldErr: true,
		},
		{
			name: "invalid amount returns error",
			msg: types.NewMsgDelegateService(
				msgDelegateService.ServiceID,
				sdk.Coins{sdk.Coin{Denom: "invalid!", Amount: sdkmath.NewInt(100_000_000)}},
				msgDelegateService.Delegator,
			),
			shouldErr: true,
		},
		{
			name: "invalid delegator address returns error",
			msg: types.NewMsgDelegateService(
				msgDelegateService.ServiceID,
				msgDelegateService.Amount,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name: "valid message returns no error",
			msg:  msgDelegateService,
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

func TestMsgDelegateService_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgDelegateService","value":{"amount":[{"amount":"100000000","denom":"umilk"}],"delegator":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd","service_id":1}}`
	require.Equal(t, expected, string(msgDelegateService.GetSignBytes()))
}

func TestMsgDelegateService_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgDelegateService.Delegator)
	require.Equal(t, []sdk.AccAddress{addr}, msgDelegateOperator.GetSigners())
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
				types.NewParams(0, nil, types.DefaultRestakingCap),
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
	expected := `{"type":"milkyway/restaking/MsgUpdateParams","value":{"authority":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd","params":{"restaking_cap":"0.000000000000000000","unbonding_time":"259200000000000"}}}`
	require.Equal(t, expected, string(msgUpdateParams.GetSignBytes()))
}

func TestMsgUpdateParams_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgUpdateParams.Authority)
	require.Equal(t, []sdk.AccAddress{addr}, msgDelegateOperator.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgUndelegatePool = types.NewMsgUndelegatePool(
	sdk.NewCoin("umilk", sdkmath.NewInt(100_000_000)),
	"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
)

func TestMsgUndelegatePool_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgUndelegatePool
		shouldErr bool
	}{
		{
			name: "invalid amount return error",
			msg: types.NewMsgUndelegatePool(
				sdk.Coin{Denom: "umilk", Amount: sdkmath.ZeroInt()},
				msgUndelegatePool.Delegator,
			),
			shouldErr: true,
		},
		{
			name: "invalid delegator address returns error",
			msg: types.NewMsgUndelegatePool(
				msgUndelegatePool.Amount,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name: "valid message returns no error",
			msg:  msgUndelegatePool,
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

func TestMsgUndelegatePool_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgUndelegatePool","value":{"amount":{"amount":"100000000","denom":"umilk"},"delegator":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"}}`
	require.Equal(t, expected, string(msgUndelegatePool.GetSignBytes()))
}

func TestMsgUndelegatePool_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgUndelegatePool.Delegator)
	require.Equal(t, []sdk.AccAddress{addr}, msgDelegateOperator.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgUndelegateOperator = types.NewMsgUndelegateOperator(
	1,
	sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100_000_000))),
	"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
)

func TestMsgUndelegateOperator_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgUndelegateOperator
		shouldErr bool
	}{
		{
			name: "invalid operator id returns error",
			msg: types.NewMsgUndelegateOperator(
				0,
				msgUndelegateOperator.Amount,
				msgUndelegateOperator.Delegator,
			),
			shouldErr: true,
		},
		{
			name: "invalid amount returns error",
			msg: types.NewMsgUndelegateOperator(
				msgUndelegateOperator.OperatorID,
				sdk.Coins{sdk.Coin{Denom: "umilk", Amount: sdkmath.ZeroInt()}},
				msgUndelegateOperator.Delegator,
			),
			shouldErr: true,
		},
		{
			name: "invalid delegator address returns error",
			msg: types.NewMsgUndelegateOperator(
				msgUndelegateOperator.OperatorID,
				msgUndelegateOperator.Amount,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name: "valid message returns no error",
			msg:  msgUndelegateOperator,
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

func TestMsgUndelegateOperator_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgUndelegateOperator","value":{"amount":[{"amount":"100000000","denom":"umilk"}],"delegator":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd","operator_id":1}}`
	require.Equal(t, expected, string(msgUndelegateOperator.GetSignBytes()))
}

func TestMsgUndelegateOperator_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgUndelegateOperator.Delegator)
	require.Equal(t, []sdk.AccAddress{addr}, msgDelegateOperator.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgUndelegateService = types.NewMsgUndelegateService(
	1,
	sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100_000_000))),
	"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
)

func TestMsgUndelegateService_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgUndelegateService
		shouldErr bool
	}{
		{
			name: "invalid service id returns error",
			msg: types.NewMsgUndelegateService(
				0,
				msgUndelegateService.Amount,
				msgUndelegateService.Delegator,
			),
			shouldErr: true,
		},
		{
			name: "invalid amount returns error",
			msg: types.NewMsgUndelegateService(
				msgUndelegateService.ServiceID,
				sdk.Coins{sdk.Coin{Denom: "umilk", Amount: sdkmath.ZeroInt()}},
				msgUndelegateService.Delegator,
			),
			shouldErr: true,
		},
		{
			name: "invalid delegator address returns error",
			msg: types.NewMsgUndelegateService(
				msgUndelegateService.ServiceID,
				msgUndelegateService.Amount,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name: "valid message returns no error",
			msg:  msgUndelegateService,
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

func TestMsgUndelegateService_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgUndelegateService","value":{"amount":[{"amount":"100000000","denom":"umilk"}],"delegator":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd","service_id":1}}`
	require.Equal(t, expected, string(msgUndelegateService.GetSignBytes()))
}

func TestMsgUndelegateService_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgUndelegateService.Delegator)
	require.Equal(t, []sdk.AccAddress{addr}, msgDelegateOperator.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgSetUserPreferences = types.NewMsgSetUserPreferences(
	types.NewUserPreferences(
		true,
		false,
		[]uint32{1, 2, 3},
	),
	"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
)

func TestMsgSetUserPreferences_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgSetUserPreferences
		shouldErr bool
	}{
		{
			name: "invalid preferences returns error",
			msg: types.NewMsgSetUserPreferences(
				types.NewUserPreferences(
					false,
					true,
					[]uint32{0},
				),
				msgSetUserPreferences.User,
			),
			shouldErr: true,
		},
		{
			name: "invalid user address returns error",
			msg: types.NewMsgSetUserPreferences(
				msgSetUserPreferences.Preferences,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name: "valid message returns no error",
			msg:  msgSetUserPreferences,
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

func TestMsgSetUserPreferences_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgSetUserPreferences","value":{"preferences":{"trust_non_accredited_services":true,"trusted_services_ids":[1,2,3]},"user":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"}}`
	require.Equal(t, expected, string(msgSetUserPreferences.GetSignBytes()))
}

func TestMsgSetUserPreferences_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgSetUserPreferences.User)
	require.Equal(t, []sdk.AccAddress{addr}, msgSetUserPreferences.GetSigners())
}
