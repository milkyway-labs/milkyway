package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

var msgJoinRestakingPool = types.NewMsgJoinRestakingPool(
	sdk.NewCoin("umilk", sdkmath.NewInt(100_000_000)),
	"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
)

func TestMsgJoinRestakingPool_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgJoinRestakingPool
		shouldErr bool
	}{
		{
			name: "invalid amount returns error",
			msg: types.NewMsgJoinRestakingPool(
				sdk.Coin{Denom: "invalid!", Amount: sdkmath.NewInt(100_000_000)},
				msgJoinRestakingPool.Delegator,
			),
			shouldErr: true,
		},
		{
			name: "invalid delegator address returns error",
			msg: types.NewMsgJoinRestakingPool(
				msgJoinRestakingPool.Amount,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name: "valid message returns no error",
			msg:  msgJoinRestakingPool,
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

func TestMsgJoinRestakingPool_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgJoinRestakingPool","value":{"amount":{"amount":"100000000","denom":"umilk"},"delegator":"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"}}`
	require.Equal(t, expected, string(msgJoinRestakingPool.GetSignBytes()))
}

func TestMsgJoinRestakingPool_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgJoinRestakingPool.Delegator)
	require.Equal(t, []sdk.AccAddress{addr}, msgJoinRestakingPool.GetSigners())
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
