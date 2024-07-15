package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

func TestParseOperatorID(t *testing.T) {
	testCases := []struct {
		name      string
		value     string
		shouldErr bool
		expID     uint32
	}{
		{
			name:  "valid ID returns no error",
			value: "1",
			expID: 1,
		},
		{
			name:      "invalid ID returns error",
			value:     "invalid",
			shouldErr: true,
		},
		{
			name:      "empty ID returns error",
			value:     "",
			shouldErr: true,
		},
		{
			name:      "negative ID returns error",
			value:     "-1",
			shouldErr: true,
		},
		{
			name:      "zero ID returns no error",
			value:     "0",
			shouldErr: false,
			expID:     0,
		},
		{
			name:      "max uint32 returns no error",
			value:     "4294967295",
			shouldErr: false,
			expID:     4294967295,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			id, err := types.ParseOperatorID(tc.value)
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expID, id)
			}
		})
	}
}

func TestOperator_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		operator  types.Operator
		shouldErr bool
	}{
		{
			name: "invalid id returns error",
			operator: types.NewOperator(
				0,
				types.OPERATOR_STATUS_ACTIVE,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "invalid status returns error",
			operator: types.NewOperator(
				1,
				types.OPERATOR_STATUS_UNSPECIFIED,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "invalid moniker returns error",
			operator: types.NewOperator(
				1,
				types.OPERATOR_STATUS_ACTIVE,
				"",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "invalid admin address returns error",
			operator: types.NewOperator(
				1,
				types.OPERATOR_STATUS_ACTIVE,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"",
			),
			shouldErr: true,
		},
		{
			name: "invalid address returns error",
			operator: types.Operator{
				ID:      1,
				Status:  types.OPERATOR_STATUS_ACTIVE,
				Moniker: "MilkyWay Operator",
				Admin:   "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Address: "",
			},
			shouldErr: true,
		},
		{
			name: "valid operator returns no error",
			operator: types.NewOperator(
				1,
				types.OPERATOR_STATUS_ACTIVE,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.operator.Validate()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOperator_Update(t *testing.T) {
	testCases := []struct {
		name      string
		operator  types.Operator
		update    types.OperatorUpdate
		expResult types.Operator
	}{
		{
			name: "update moniker",
			operator: types.NewOperator(
				1,
				types.OPERATOR_STATUS_ACTIVE,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			update: types.NewOperatorUpdate(
				"MilkyWay2",
				types.DoNotModify,
				types.DoNotModify,
			),
			expResult: types.NewOperator(
				1,
				types.OPERATOR_STATUS_ACTIVE,
				"MilkyWay2",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
		},
		{
			name: "update description",
			operator: types.NewOperator(
				1,
				types.OPERATOR_STATUS_ACTIVE,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			update: types.NewOperatorUpdate(
				types.DoNotModify,
				"https://example.com",
				types.DoNotModify,
			),
			expResult: types.NewOperator(
				1,
				types.OPERATOR_STATUS_ACTIVE,
				"MilkyWay Operator",
				"https://example.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
		},
		{
			name: "update picture URL",
			operator: types.NewOperator(
				1,
				types.OPERATOR_STATUS_ACTIVE,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			update: types.NewOperatorUpdate(
				types.DoNotModify,
				types.DoNotModify,
				"https://example.com/picture.jpg",
			),
			expResult: types.NewOperator(
				1,
				types.OPERATOR_STATUS_ACTIVE,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://example.com/picture.jpg",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := tc.operator.Update(tc.update)
			require.Equal(t, tc.expResult, result)
		})
	}
}

func TestOperator_SharesFromTokens(t *testing.T) {
	testCases := []struct {
		name      string
		operator  types.Operator
		tokens    sdk.Coin
		shouldErr bool
		expShares sdkmath.LegacyDec
	}{
		{
			name: "operator with no delegation shares returns error",
			operator: types.Operator{
				ID:              1,
				Address:         types.GetOperatorAddress(1).String(),
				DelegatorShares: sdk.NewDecCoins(),
				Tokens:          sdk.NewCoins(),
			},
			tokens:    sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			shouldErr: true,
		},
		{
			name: "shares are computed properly for non empty operator",
			operator: types.Operator{
				ID:      1,
				Address: types.GetOperatorAddress(1).String(),
				Tokens: sdk.NewCoins(
					sdk.NewCoin("umilk", sdkmath.NewInt(50)),
				),
				DelegatorShares: sdk.NewDecCoins(
					sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(100)),
				),
			},
			tokens:    sdk.NewCoin("umilk", sdkmath.NewInt(20)),
			shouldErr: false,
			expShares: sdkmath.LegacyNewDec(40),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			shares, err := tc.operator.SharesFromTokens(tc.tokens)
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expShares, shares)
			}
		})
	}
}

func TestOperator_TokensFromShares(t *testing.T) {
	testCases := []struct {
		name      string
		operator  types.Operator
		shares    sdk.DecCoins
		expTokens sdk.DecCoins
	}{
		{
			name: "operator with shares returns correct amount",
			operator: types.Operator{
				ID:      1,
				Address: types.GetOperatorAddress(1).String(),
				Tokens: sdk.NewCoins(
					sdk.NewCoin("umilk", sdkmath.NewInt(70)),
				),
				DelegatorShares: sdk.NewDecCoins(
					sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(140)),
				),
			},
			shares: sdk.NewDecCoins(
				sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(40)),
			),
			expTokens: sdk.NewDecCoins(
				sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(20)),
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tokens := tc.operator.TokensFromShares(tc.shares)
			require.Equal(t, tc.expTokens, tokens)
		})
	}
}
