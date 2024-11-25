package types_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

var msgCreateRewardsPlan = types.NewMsgCreateRewardsPlan(
	1,
	"Test rewards plan",
	sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000))),
	time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
	types.NewBasicPoolsDistribution(1),
	types.NewBasicOperatorsDistribution(1),
	types.NewBasicUsersDistribution(1),
	sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100))),
	"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
)

func TestMsgCreateRewardsPlan_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgCreateRewardsPlan
		shouldErr bool
	}{
		{
			name: "invalid service id returns error",
			msg: types.NewMsgCreateRewardsPlan(
				0,
				msgCreateRewardsPlan.Description,
				msgCreateRewardsPlan.Amount,
				msgCreateRewardsPlan.StartTime,
				msgCreateRewardsPlan.EndTime,
				msgCreateRewardsPlan.PoolsDistribution,
				msgCreateRewardsPlan.OperatorsDistribution,
				msgCreateRewardsPlan.UsersDistribution,
				msgCreateRewardsPlan.FeeAmount,
				msgCreateRewardsPlan.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid amount",
			msg: types.NewMsgCreateRewardsPlan(
				msgCreateRewardsPlan.ServiceID,
				msgCreateRewardsPlan.Description,
				sdk.Coins{sdk.Coin{Denom: "invalid", Amount: sdkmath.NewInt(-100)}},
				msgCreateRewardsPlan.StartTime,
				msgCreateRewardsPlan.EndTime,
				msgCreateRewardsPlan.PoolsDistribution,
				msgCreateRewardsPlan.OperatorsDistribution,
				msgCreateRewardsPlan.UsersDistribution,
				msgCreateRewardsPlan.FeeAmount,
				msgCreateRewardsPlan.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid end time returns error",
			msg: types.NewMsgCreateRewardsPlan(
				msgCreateRewardsPlan.ServiceID,
				msgCreateRewardsPlan.Description,
				msgCreateRewardsPlan.Amount,
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				msgCreateRewardsPlan.PoolsDistribution,
				msgCreateRewardsPlan.OperatorsDistribution,
				msgCreateRewardsPlan.UsersDistribution,
				msgCreateRewardsPlan.FeeAmount,
				msgCreateRewardsPlan.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid fee amount returns error",
			msg: types.NewMsgCreateRewardsPlan(
				msgCreateRewardsPlan.ServiceID,
				msgCreateRewardsPlan.Description,
				msgCreateRewardsPlan.Amount,
				msgCreateRewardsPlan.StartTime,
				msgCreateRewardsPlan.EndTime,
				msgCreateRewardsPlan.PoolsDistribution,
				msgCreateRewardsPlan.OperatorsDistribution,
				msgCreateRewardsPlan.UsersDistribution,
				sdk.Coins{sdk.Coin{Denom: "umilk", Amount: sdkmath.NewInt(-100)}},
				msgCreateRewardsPlan.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid sender returns error",
			msg: types.NewMsgCreateRewardsPlan(
				msgCreateRewardsPlan.ServiceID,
				msgCreateRewardsPlan.Description,
				msgCreateRewardsPlan.Amount,
				msgCreateRewardsPlan.StartTime,
				msgCreateRewardsPlan.EndTime,
				msgCreateRewardsPlan.PoolsDistribution,
				msgCreateRewardsPlan.OperatorsDistribution,
				msgCreateRewardsPlan.UsersDistribution,
				msgCreateRewardsPlan.FeeAmount,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name:      "valid message returns no error",
			msg:       msgCreateRewardsPlan,
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

func TestMsgCreateRewardsPlan_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgCreateRewardsPlan","value":{"amount":[{"amount":"1000","denom":"stake"}],"description":"Test rewards plan","end_time":"2024-12-31T23:59:59Z","fee_amount":[{"amount":"100","denom":"stake"}],"operators_distribution":{"delegation_type":2,"type":{"type":"milkyway/DistributionTypeBasic","value":{}},"weight":1},"pools_distribution":{"delegation_type":1,"type":{"type":"milkyway/DistributionTypeBasic","value":{}},"weight":1},"sender":"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn","service_id":1,"start_time":"2024-01-01T00:00:00Z","users_distribution":{"type":{"type":"milkyway/UsersDistributionTypeBasic","value":{}},"weight":1}}}`
	require.Equal(t, expected, string(msgCreateRewardsPlan.GetSignBytes()))
}

func TestMsgCreateRewardsPlan_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgCreateRewardsPlan.Sender)
	require.Equal(t, []sdk.AccAddress{addr}, msgCreateRewardsPlan.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgEditRewardsPlan = types.NewMsgEditRewardsPlan(
	1,
	"Test rewards plan",
	sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000))),
	time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
	types.NewBasicPoolsDistribution(1),
	types.NewBasicOperatorsDistribution(1),
	types.NewBasicUsersDistribution(1),
	"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
)

func TestMsgEditRewardsPlan_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgEditRewardsPlan
		shouldErr bool
	}{
		{
			name: "invalid id returns error",
			msg: types.NewMsgEditRewardsPlan(
				0,
				msgEditRewardsPlan.Description,
				msgEditRewardsPlan.Amount,
				msgEditRewardsPlan.StartTime,
				msgEditRewardsPlan.EndTime,
				msgEditRewardsPlan.PoolsDistribution,
				msgEditRewardsPlan.OperatorsDistribution,
				msgEditRewardsPlan.UsersDistribution,
				msgEditRewardsPlan.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid amount",
			msg: types.NewMsgEditRewardsPlan(
				msgEditRewardsPlan.ID,
				msgEditRewardsPlan.Description,
				sdk.Coins{sdk.Coin{Denom: "invalid", Amount: sdkmath.NewInt(-100)}},
				msgEditRewardsPlan.StartTime,
				msgEditRewardsPlan.EndTime,
				msgEditRewardsPlan.PoolsDistribution,
				msgEditRewardsPlan.OperatorsDistribution,
				msgEditRewardsPlan.UsersDistribution,
				msgEditRewardsPlan.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid end time returns error",
			msg: types.NewMsgEditRewardsPlan(
				msgEditRewardsPlan.ID,
				msgEditRewardsPlan.Description,
				msgEditRewardsPlan.Amount,
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				msgEditRewardsPlan.PoolsDistribution,
				msgEditRewardsPlan.OperatorsDistribution,
				msgEditRewardsPlan.UsersDistribution,
				msgEditRewardsPlan.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid sender returns error",
			msg: types.NewMsgEditRewardsPlan(
				msgEditRewardsPlan.ID,
				msgEditRewardsPlan.Description,
				msgEditRewardsPlan.Amount,
				msgEditRewardsPlan.StartTime,
				msgEditRewardsPlan.EndTime,
				msgEditRewardsPlan.PoolsDistribution,
				msgEditRewardsPlan.OperatorsDistribution,
				msgEditRewardsPlan.UsersDistribution,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name:      "valid message returns no error",
			msg:       msgEditRewardsPlan,
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

func TestMsgEditRewardsPlan_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgEditRewardsPlan","value":{"amount":[{"amount":"1000","denom":"stake"}],"description":"Test rewards plan","end_time":"2024-12-31T23:59:59Z","id":"1","operators_distribution":{"delegation_type":2,"type":{"type":"milkyway/DistributionTypeBasic","value":{}},"weight":1},"pools_distribution":{"delegation_type":1,"type":{"type":"milkyway/DistributionTypeBasic","value":{}},"weight":1},"sender":"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn","start_time":"2024-01-01T00:00:00Z","users_distribution":{"type":{"type":"milkyway/UsersDistributionTypeBasic","value":{}},"weight":1}}}`
	require.Equal(t, expected, string(msgEditRewardsPlan.GetSignBytes()))
}

func TestMsgEditRewardsPlan_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgEditRewardsPlan.Sender)
	require.Equal(t, []sdk.AccAddress{addr}, msgCreateRewardsPlan.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgSetWithdrawAddress = types.NewMsgSetWithdrawAddress(
	"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
	"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
)

func TestMsgSetWithdrawAddress_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgSetWithdrawAddress
		shouldErr bool
	}{
		{
			name: "invalid withdraw address returns error",
			msg: types.NewMsgSetWithdrawAddress(
				"invalid",
				msgSetWithdrawAddress.WithdrawAddress,
			),
			shouldErr: true,
		},
		{
			name: "invalid sender address returns error",
			msg: types.NewMsgSetWithdrawAddress(
				msgSetWithdrawAddress.Sender,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name:      "valid message returns no error",
			msg:       msgSetWithdrawAddress,
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

func TestMsgSetWithdrawAddress_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgSetWithdrawAddress","value":{"sender":"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4","withdraw_address":"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn"}}`
	require.Equal(t, expected, string(msgSetWithdrawAddress.GetSignBytes()))
}

func TestMsgSetWithdrawAddress_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgSetWithdrawAddress.Sender)
	require.Equal(t, []sdk.AccAddress{addr}, msgSetWithdrawAddress.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgWithdrawDelegatorReward = types.NewMsgWithdrawDelegatorReward(
	restakingtypes.DELEGATION_TYPE_SERVICE,
	1,
	"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
)

func TestMsgWithdrawDelegatorReward_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgWithdrawDelegatorReward
		shouldErr bool
	}{
		{
			name: "invalid delegation type returns error",
			msg: types.NewMsgWithdrawDelegatorReward(
				restakingtypes.DELEGATION_TYPE_UNSPECIFIED,
				msgWithdrawDelegatorReward.DelegationTargetID,
				msgWithdrawDelegatorReward.DelegatorAddress,
			),
			shouldErr: true,
		},
		{
			name: "invalid delegation target ID returns error",
			msg: types.NewMsgWithdrawDelegatorReward(
				msgWithdrawDelegatorReward.DelegationType,
				0,
				msgWithdrawDelegatorReward.DelegatorAddress,
			),
			shouldErr: true,
		},
		{
			name: "invalid delegator address returns error",
			msg: types.NewMsgWithdrawDelegatorReward(
				msgWithdrawDelegatorReward.DelegationType,
				msgWithdrawDelegatorReward.DelegationTargetID,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name:      "valid message returns no error",
			msg:       msgWithdrawDelegatorReward,
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

func TestMsgWithdrawDelegatorReward_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgWithdrawDelegatorReward","value":{"delegation_target_id":1,"delegation_type":3,"delegator_address":"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn"}}`
	require.Equal(t, expected, string(msgWithdrawDelegatorReward.GetSignBytes()))
}

func TestMsgWithdrawDelegatorReward_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgWithdrawDelegatorReward.DelegatorAddress)
	require.Equal(t, []sdk.AccAddress{addr}, msgWithdrawDelegatorReward.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgWithdrawOperatorCommission = types.NewMsgWithdrawOperatorCommission(
	1,
	"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
)

func TestMsgWithdrawOperatorCommission_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgWithdrawOperatorCommission
		shouldErr bool
	}{
		{
			name: "invalid operator ID returns error",
			msg: types.NewMsgWithdrawOperatorCommission(
				0,
				msgWithdrawOperatorCommission.Sender,
			),
			shouldErr: true,
		},
		{
			name: "invalid sender address returns error",
			msg: types.NewMsgWithdrawOperatorCommission(
				msgWithdrawOperatorCommission.OperatorID,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name:      "valid message returns no error",
			msg:       msgWithdrawOperatorCommission,
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

func TestMsgWithdrawOperatorCommission_GetSignBytes(t *testing.T) {
	expected := `{"type":"milkyway/MsgWithdrawOperatorCommission","value":{"operator_id":1,"sender":"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn"}}`
	require.Equal(t, expected, string(msgWithdrawOperatorCommission.GetSignBytes()))
}

func TestMsgWithdrawOperatorCommission_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgWithdrawOperatorCommission.Sender)
	require.Equal(t, []sdk.AccAddress{addr}, msgWithdrawOperatorCommission.GetSigners())
}
