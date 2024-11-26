package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v2/testutils"
	"github.com/milkyway-labs/milkyway/v2/utils"
	"github.com/milkyway-labs/milkyway/v2/x/rewards/types"
)

func TestGenesisState_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		genesis   *types.GenesisState
		shouldErr bool
	}{
		{
			name: "invalid next rewards plan ID returns error",
			genesis: &types.GenesisState{
				NextRewardsPlanID: 0,
			},
			shouldErr: true,
		},
		{
			name: "invalid rewards plan returns error",
			genesis: &types.GenesisState{
				NextRewardsPlanID: 1,
				RewardsPlans: []types.RewardsPlan{
					types.NewRewardsPlan(
						1, "Rewards plan", 0, utils.MustParseCoins("10_000000umilk"),
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						types.NewBasicPoolsDistribution(0), types.NewBasicOperatorsDistribution(0),
						types.NewBasicUsersDistribution(0),
					),
				},
			},
			shouldErr: true,
		},
		{
			name: "invalid pool service total delegator shares pool ID returns error",
			genesis: &types.GenesisState{
				NextRewardsPlanID: 1,
				PoolServiceTotalDelegatorShares: []types.PoolServiceTotalDelegatorShares{
					types.NewPoolServiceTotalDelegatorShares(0, 1, utils.MustParseDecCoins("10_000000umilk")),
				},
			},
			shouldErr: true,
		},
		{
			name: "invalid pool service total delegator shares service ID returns error",
			genesis: &types.GenesisState{
				NextRewardsPlanID: 1,
				PoolServiceTotalDelegatorShares: []types.PoolServiceTotalDelegatorShares{
					types.NewPoolServiceTotalDelegatorShares(1, 0, utils.MustParseDecCoins("10_000000umilk")),
				},
			},
			shouldErr: true,
		},
		{
			name:      "default genesis returns no error",
			genesis:   types.DefaultGenesis(),
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cdc, _ := testutils.MakeCodecs()

			err := tc.genesis.Validate(cdc)
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
