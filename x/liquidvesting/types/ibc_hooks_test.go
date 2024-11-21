package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

func TestMsgDepositInsurance_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgDepositInsurance
		shouldErr bool
	}{
		{
			name: "invalid depositor address returns error",
			msg: types.NewMsgDepositInsurance(
				[]types.InsuranceDeposit{
					{
						Depositor: "invalid",
						Amount:    sdk.NewInt64Coin("stake", 1000),
					},
				},
			),
			shouldErr: true,
		},
		{
			name: "multiple denoms are not allowed",
			msg: types.NewMsgDepositInsurance(
				[]types.InsuranceDeposit{
					{
						Depositor: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						Amount:    sdk.NewInt64Coin("stake", 1000),
					},
					{
						Depositor: "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
						Amount:    sdk.NewInt64Coin("umilk", 1000),
					},
				},
			),
			shouldErr: true,
		},
		{
			name: "validates correctly",
			msg: types.NewMsgDepositInsurance(
				[]types.InsuranceDeposit{
					{
						Depositor: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						Amount:    sdk.NewInt64Coin("stake", 1000),
					},
					{
						Depositor: "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
						Amount:    sdk.NewInt64Coin("stake", 3000),
					},
				},
			),
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
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
