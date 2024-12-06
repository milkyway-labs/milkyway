package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v3/x/liquidvesting/types"
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
						Amount:    math.NewIntFromUint64(1000),
					},
				},
			),
			shouldErr: true,
		},
		{
			name: "negative amount returns error",
			msg: types.NewMsgDepositInsurance(
				[]types.InsuranceDeposit{
					{
						Depositor: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						Amount:    math.NewInt(-1000),
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
						Amount:    math.NewIntFromUint64(1000),
					},
					{
						Depositor: "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
						Amount:    math.NewIntFromUint64(2000),
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
