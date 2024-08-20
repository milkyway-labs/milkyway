package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

func TestParams_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		params      types.Params
		expectedErr string
	}{
		{
			name:        "default params returns no error",
			params:      types.DefaultParams(),
			expectedErr: "",
		},
		{
			name:        "invalid rewards plan creation fee returns error",
			params:      types.NewParams(sdk.Coins{sdk.NewInt64Coin("umilk", 0)}),
			expectedErr: "invalid rewards plan creation fee: coin 0umilk amount is not positive",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}
