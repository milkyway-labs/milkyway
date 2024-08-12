package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/milkyway-labs/milkyway/utils"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

func TestGenesisState_Validate(t *testing.T) {
	ir := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(ir)

	testCases := []struct {
		name        string
		malleate    func(genState *types.GenesisState)
		expectedErr string
	}{
		{
			name:        "default params returns no error",
			malleate:    nil,
			expectedErr: "",
		},
		{
			name: "invalid next rewards plan ID returns error",
			malleate: func(genState *types.GenesisState) {
				genState.NextRewardsPlanID = 0
			},
			expectedErr: "invalid next rewards plan ID: 0",
		},
		{
			name: "invalid rewards plan returns error",
			malleate: func(genState *types.GenesisState) {
				genState.RewardsPlans = []types.RewardsPlan{
					types.NewRewardsPlan(
						1, "Rewards plan", 0, utils.MustParseCoins("10_000000umilk"),
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						types.NewBasicPoolsDistribution(0), types.NewBasicOperatorsDistribution(0),
						types.NewBasicUsersDistribution(0)),
				}
			},
			expectedErr: "invalid rewards plan at index 0: invalid service ID: 0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			genState := types.DefaultGenesis()
			if tc.malleate != nil {
				tc.malleate(genState)
			}
			err := genState.Validate(ir)
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}
