package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

var msgUpdateOperatorParams = types.NewMsgUpdateOperatorParams(
	1, types.NewOperatorParams(sdkmath.LegacyNewDecWithPrec(1, 1), []uint32{1, 2, 3}),
	"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd")

func TestMsgUpdateOperatorParams_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgUpdateOperatorParams
		shouldErr bool
	}{
		{
			name: "invalid operator id returns error",
			msg: types.NewMsgUpdateOperatorParams(
				0, msgUpdateOperatorParams.OperatorParams, msgUpdateOperatorParams.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid operator params returns error",
			msg: types.NewMsgUpdateOperatorParams(
				msgUpdateOperatorParams.OperatorID, types.NewOperatorParams(sdkmath.LegacyNewDec(2), nil),
				msgUpdateOperatorParams.Sender),
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

func TestMsgUpdateOperatorParams_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgUpdateOperatorParams","value":{"operator_id":1,"operator_params":{"commission_rate":"0.100000000000000000","joined_service_ids":[1,2,3]},"sender":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"}}`
	require.Equal(t, expected, string(msgUpdateOperatorParams.GetSignBytes()))
}

func TestMsgUpdateOperatorParams_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgUpdateOperatorParams.Sender)
	require.Equal(t, []sdk.AccAddress{addr}, msgUpdateOperatorParams.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgUpdateServiceParams = types.NewMsgUpdateServiceParams(
	1, types.NewServiceParams(
		sdkmath.LegacyNewDecWithPrec(5, 2), []uint32{1, 2, 3}, []uint32{4, 5, 6}),
	"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd")

func TestMsgUpdateServiceParams_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgUpdateServiceParams
		shouldErr bool
	}{
		{
			name: "invalid service id returns error",
			msg: types.NewMsgUpdateServiceParams(
				0, msgUpdateServiceParams.ServiceParams, msgUpdateServiceParams.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid service params returns error",
			msg: types.NewMsgUpdateServiceParams(
				msgUpdateServiceParams.ServiceID,
				types.NewServiceParams(sdkmath.LegacyNewDec(2), nil, nil),
				msgUpdateServiceParams.Sender),
			shouldErr: true,
		},
		{
			name: "valid message returns no error",
			msg:  msgUpdateServiceParams,
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

func TestMsgUpdateServiceParams_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgUpdateServiceParams","value":{"sender":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd","service_id":1,"service_params":{"slash_fraction":"0.050000000000000000","whitelisted_operator_ids":[4,5,6],"whitelisted_pool_ids":[1,2,3]}}}`
	require.Equal(t, expected, string(msgUpdateServiceParams.GetSignBytes()))
}

func TestMsgUpdateServiceParams_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgUpdateServiceParams.Sender)
	require.Equal(t, []sdk.AccAddress{addr}, msgUpdateServiceParams.GetSigners())
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
				types.NewParams(0),
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
	expected := `{"type":"milkyway/restaking/MsgUpdateParams","value":{"authority":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd","params":{"unbonding_time":"259200000000000"}}}`
	require.Equal(t, expected, string(msgUpdateParams.GetSignBytes()))
}

func TestMsgUpdateParams_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgUpdateParams.Authority)
	require.Equal(t, []sdk.AccAddress{addr}, msgDelegateOperator.GetSigners())
}
