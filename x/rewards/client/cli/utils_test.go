package cli_test

import (
	"os"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	milkywayapp "github.com/milkyway-labs/milkyway/v7/app"
	"github.com/milkyway-labs/milkyway/v7/x/rewards/client/cli"
	"github.com/milkyway-labs/milkyway/v7/x/rewards/types"
)

func TestCLIUtils_parseRewardsPlan(t *testing.T) {
	cdc, _ := milkywayapp.MakeCodecs()

	testCases := []struct {
		name      string
		jsonFile  *os.File
		shouldErr bool
		expected  cli.ParsedRewardsPlan
	}{
		{
			name: "parse basic distribution json",
			jsonFile: testutil.WriteToNewTempFile(t, `{
	    "service_id": 1,
	    "description": "test plan",
	    "amount_per_day": "1000uinit",
	    "start_time": "2024-01-01T00:00:00Z",
	    "end_time": "2024-12-31T23:59:59Z",
	    "pools_distribution": {
	        "weight": 1,
	        "type": {
	            "@type":"/milkyway.rewards.v2.DistributionTypeBasic"
	        }
	    },
	    "operators_distribution": {
	        "weight": 2,
	        "type": {
	            "@type": "/milkyway.rewards.v2.DistributionTypeBasic"
	        }
	    },
	    "users_distribution": {
	        "weight": 3,
	        "type": {
	            "@type": "/milkyway.rewards.v2.UsersDistributionTypeBasic"
	        }
	    },
		"fee_amount": "100uinit"
	}`),
			shouldErr: false,
			expected: cli.ParsedRewardsPlan{
				RewardsPlan: types.NewRewardsPlan(
					1,
					"test plan",
					1,
					sdk.NewCoin("uinit", math.NewInt(1000)),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
					types.NewBasicPoolsDistribution(1),
					types.NewBasicOperatorsDistribution(2),
					types.NewBasicUsersDistribution(3),
				),
				FeeAmount: sdk.NewCoins(sdk.NewInt64Coin("uinit", 100)),
			},
		},
		{
			name: "parse egalitarian distribution json",
			jsonFile: testutil.WriteToNewTempFile(t, `{
	    "service_id": 1,
	    "description": "test plan",
	    "amount_per_day": "1000uinit",
	    "start_time": "2024-01-01T00:00:00Z",
	    "end_time": "2024-12-31T23:59:59Z",
	    "pools_distribution": {
	        "weight": 1,
	        "type": {
	            "@type":"/milkyway.rewards.v2.DistributionTypeBasic"
	        }
	    },
	    "operators_distribution": {
	        "weight": 2,
	        "type": {
	            "@type": "/milkyway.rewards.v2.DistributionTypeEgalitarian"
	        }
	    },
	    "users_distribution": {
	        "weight": 3,
	        "type": {
	            "@type": "/milkyway.rewards.v2.UsersDistributionTypeBasic"
	        }
	    },
		"fee_amount": "100uinit"
	}`),
			shouldErr: false,
			expected: cli.ParsedRewardsPlan{
				RewardsPlan: types.NewRewardsPlan(
					1,
					"test plan",
					1,
					sdk.NewCoin("uinit", math.NewInt(1000)),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
					types.NewBasicPoolsDistribution(1),
					types.NewEgalitarianOperatorsDistribution(2),
					types.NewBasicUsersDistribution(3),
				),
				FeeAmount: sdk.NewCoins(sdk.NewInt64Coin("uinit", 100)),
			},
		},
		{
			name: "parse weighted distribution json",
			jsonFile: testutil.WriteToNewTempFile(t, `{
	    "service_id": 1,
	    "description": "test plan",
	    "amount_per_day": "1000uinit",
	    "start_time": "2024-01-01T00:00:00Z",
	    "end_time": "2024-12-31T23:59:59Z",
	    "pools_distribution": {
	        "weight": 1,
	        "type": {
	            "@type":"/milkyway.rewards.v2.DistributionTypeWeighted",
                "weights": [{ "delegation_target_id": 1, "weight": 1 }]	
	        }
	    },
	    "operators_distribution": {
	        "weight": 2,
	        "type": {
	            "@type": "/milkyway.rewards.v2.DistributionTypeWeighted",
                 "weights": [{ "delegation_target_id": 2, "weight": 2 }]	
	        }
	    },
	    "users_distribution": {
	        "weight": 3,
	        "type": {
	            "@type": "/milkyway.rewards.v2.UsersDistributionTypeBasic"
	        }
	    },
		"fee_amount": "100uinit"
	}`),
			shouldErr: false,
			expected: cli.ParsedRewardsPlan{
				RewardsPlan: types.NewRewardsPlan(
					1,
					"test plan",
					1,
					sdk.NewCoin("uinit", math.NewInt(1000)),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
					types.NewWeightedPoolsDistribution(1, []types.DistributionWeight{types.NewDistributionWeight(1, 1)}),
					types.NewWeightedOperatorsDistribution(2, []types.DistributionWeight{types.NewDistributionWeight(2, 2)}),
					types.NewBasicUsersDistribution(3),
				),
				FeeAmount: sdk.NewCoins(sdk.NewInt64Coin("uinit", 100)),
			},
		},
		{
			name: "parse edit rewards plan json",
			jsonFile: testutil.WriteToNewTempFile(t, `{
	    "description": "test plan",
	    "amount_per_day": "1000uinit",
	    "start_time": "2024-01-01T00:00:00Z",
	    "end_time": "2024-12-31T23:59:59Z",
	    "pools_distribution": {
	        "weight": 1,
	        "type": {
	            "@type":"/milkyway.rewards.v2.DistributionTypeBasic"
	        }
	    },
	    "operators_distribution": {
	        "weight": 2,
	        "type": {
	            "@type": "/milkyway.rewards.v2.DistributionTypeBasic"
	        }
	    },
	    "users_distribution": {
	        "weight": 3,
	        "type": {
	            "@type": "/milkyway.rewards.v2.UsersDistributionTypeBasic"
	        }
	    }
	}`),
			shouldErr: false,
			expected: cli.ParsedRewardsPlan{
				RewardsPlan: types.NewRewardsPlan(
					1,
					"test plan",
					0,
					sdk.NewCoin("uinit", math.NewInt(1000)),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
					types.NewBasicPoolsDistribution(1),
					types.NewBasicOperatorsDistribution(2),
					types.NewBasicUsersDistribution(3),
				),
				FeeAmount: sdk.NewCoins(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NotNil(t, tc.jsonFile)
			plan, err := cli.ParseRewardsPlan(cdc, tc.jsonFile.Name())
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, plan)
			}
		})
	}
}
