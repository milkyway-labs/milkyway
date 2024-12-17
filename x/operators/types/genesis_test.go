package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v4/x/operators/types"
)

func TestGenesisState_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		genesis   *types.GenesisState
		shouldErr bool
	}{
		{
			name: "invalid next operator ID returns error",
			genesis: &types.GenesisState{
				NextOperatorID: 0,
				Operators:      nil,
				Params:         types.DefaultParams(),
			},
			shouldErr: true,
		},
		{
			name: "duplicated operator returns error",
			genesis: &types.GenesisState{
				NextOperatorID: 1,
				Operators: []types.Operator{
					types.NewOperator(
						1,
						types.OPERATOR_STATUS_ACTIVE,
						"MilkyWay Operator",
						"https://milkyway.com",
						"https://milkyway.com/picture",
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					),
					types.NewOperator(
						1,
						types.OPERATOR_STATUS_ACTIVE,
						"MilkyWay Operator",
						"https://milkyway.com",
						"https://milkyway.com/picture",
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					),
				},
				Params: types.DefaultParams(),
			},
			shouldErr: true,
		},
		{
			name: "invalid operator returns error",
			genesis: &types.GenesisState{
				NextOperatorID: 1,
				Operators: []types.Operator{
					types.NewOperator(
						1,
						types.OPERATOR_STATUS_UNSPECIFIED,
						"MilkyWay Operator",
						"https://milkyway.com",
						"https://milkyway.com/picture",
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					),
				},
				Params: types.DefaultParams(),
			},
			shouldErr: true,
		},
		{
			name: "duplicated unbonding operator returns error",
			genesis: &types.GenesisState{
				NextOperatorID: 1,
				UnbondingOperators: []types.UnbondingOperator{
					types.NewUnbondingOperator(
						1,
						time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					),
					types.NewUnbondingOperator(
						1,
						time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					),
				},
				Params: types.DefaultParams(),
			},
			shouldErr: true,
		},
		{
			name: "invalid unbonding operator returns error",
			genesis: &types.GenesisState{
				NextOperatorID: 1,
				UnbondingOperators: []types.UnbondingOperator{
					types.NewUnbondingOperator(
						0,
						time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					),
				},
				Params: types.DefaultParams(),
			},
			shouldErr: true,
		},
		{
			name: "not found unbonding operator returns error",
			genesis: &types.GenesisState{
				NextOperatorID: 1,
				UnbondingOperators: []types.UnbondingOperator{
					types.NewUnbondingOperator(
						1,
						time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					),
				},
			},
			shouldErr: true,
		},
		{
			name: "unbonding operator with status non inactivating returns error",
			genesis: &types.GenesisState{
				NextOperatorID: 1,
				Operators: []types.Operator{
					types.NewOperator(
						1,
						types.OPERATOR_STATUS_ACTIVE,
						"MilkyWay Operator",
						"https://milkyway.com",
						"https://milkyway.com/picture",
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					),
				},
				UnbondingOperators: []types.UnbondingOperator{
					types.NewUnbondingOperator(
						1,
						time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					),
				},
			},
			shouldErr: true,
		},
		{
			name: "invalid params returns error",
			genesis: &types.GenesisState{
				NextOperatorID: 1,
				Params: types.Params{
					DeactivationTime: 0,
				},
			},
			shouldErr: true,
		},
		{
			name:      "default genesis state returns no error",
			genesis:   types.DefaultGenesis(),
			shouldErr: false,
		},
		{
			name: "valid genesis state returns no error",
			genesis: &types.GenesisState{
				NextOperatorID: 2,
				Operators: []types.Operator{
					types.NewOperator(
						1,
						types.OPERATOR_STATUS_INACTIVATING,
						"MilkyWay Operator",
						"https://milkyway.com",
						"https://milkyway.com/picture",
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					),
				},
				UnbondingOperators: []types.UnbondingOperator{
					types.NewUnbondingOperator(
						1,
						time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					),
				},
				Params: types.Params{
					DeactivationTime: 3 * 24 * time.Hour,
				},
			},
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.genesis.Validate()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
