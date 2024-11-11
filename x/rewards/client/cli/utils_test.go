package cli_test

import (
	"os"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/rewards/types"

	"github.com/cosmos/cosmos-sdk/testutil"

	"github.com/stretchr/testify/require"

	milkyway "github.com/milkyway-labs/milkyway/app"
	"github.com/milkyway-labs/milkyway/x/rewards/client/cli"
)

func TestCliUtils_parseRewardsPlan(t *testing.T) {
	encodingConfig := milkyway.MakeEncodingConfig()
	codec := encodingConfig.Marshaler

	testCases := []struct {
		name      string
		jsonFile  *os.File
		shouldErr bool
		expected  types.RewardsPlan
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
	            "@type":"/milkyway.rewards.v1.DistributionTypeBasic"
	        }
	    },
	    "operators_distribution": {
	        "weight": 2,
	        "type": {
	            "@type": "/milkyway.rewards.v1.DistributionTypeBasic"
	        }
	    },
	    "users_distribution": {
	        "weight": 3,
	        "type": {
	            "@type": "/milkyway.rewards.v1.UsersDistributionTypeBasic"
	        }
	    }
	}`),
			shouldErr: false,
			expected: types.NewRewardsPlan(
				1,
				"test plan",
				1,
				sdk.NewCoins(sdk.NewCoin("uinit", math.NewInt(1000))),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
				types.NewBasicPoolsDistribution(1),
				types.NewBasicOperatorsDistribution(2),
				types.NewBasicUsersDistribution(3),
			),
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
	            "@type":"/milkyway.rewards.v1.DistributionTypeBasic"
	        }
	    },
	    "operators_distribution": {
	        "weight": 2,
	        "type": {
	            "@type": "/milkyway.rewards.v1.DistributionTypeEgalitarian"
	        }
	    },
	    "users_distribution": {
	        "weight": 3,
	        "type": {
	            "@type": "/milkyway.rewards.v1.UsersDistributionTypeBasic"
	        }
	    }
	}`),
			shouldErr: false,
			expected: types.NewRewardsPlan(
				1,
				"test plan",
				1,
				sdk.NewCoins(sdk.NewCoin("uinit", math.NewInt(1000))),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
				types.NewBasicPoolsDistribution(1),
				types.NewEgalitarianOperatorsDistribution(2),
				types.NewBasicUsersDistribution(3),
			),
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
	            "@type":"/milkyway.rewards.v1.DistributionTypeWeighted",
                "weights": [{ "delegation_target_id": 1, "weight": 1 }]	
	        }
	    },
	    "operators_distribution": {
	        "weight": 2,
	        "type": {
	            "@type": "/milkyway.rewards.v1.DistributionTypeWeighted",
                 "weights": [{ "delegation_target_id": 2, "weight": 2 }]	
	        }
	    },
	    "users_distribution": {
	        "weight": 3,
	        "type": {
	            "@type": "/milkyway.rewards.v1.UsersDistributionTypeBasic"
	        }
	    }
	}`),
			shouldErr: false,
			expected: types.NewRewardsPlan(
				1,
				"test plan",
				1,
				sdk.NewCoins(sdk.NewCoin("uinit", math.NewInt(1000))),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
				types.NewWeightedPoolsDistribution(1, []types.DistributionWeight{types.NewDistributionWeight(1, 1)}),
				types.NewWeightedOperatorsDistribution(2, []types.DistributionWeight{types.NewDistributionWeight(2, 2)}),
				types.NewBasicUsersDistribution(3),
			),
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
	            "@type":"/milkyway.rewards.v1.DistributionTypeBasic"
	        }
	    },
	    "operators_distribution": {
	        "weight": 2,
	        "type": {
	            "@type": "/milkyway.rewards.v1.DistributionTypeBasic"
	        }
	    },
	    "users_distribution": {
	        "weight": 3,
	        "type": {
	            "@type": "/milkyway.rewards.v1.UsersDistributionTypeBasic"
	        }
	    }
	}`),
			shouldErr: false,
			expected: types.NewRewardsPlan(
				1,
				"test plan",
				0,
				sdk.NewCoins(sdk.NewCoin("uinit", math.NewInt(1000))),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
				types.NewBasicPoolsDistribution(1),
				types.NewBasicOperatorsDistribution(2),
				types.NewBasicUsersDistribution(3),
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NotNil(t, tc.jsonFile)
			plan, err := cli.ParseRewardsPlan(codec, tc.jsonFile.Name())
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, plan)
			}
		})
	}
}
