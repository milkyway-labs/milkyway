package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v9/x/rewards/types"
)

func TestParams_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		params    types.Params
		shouldErr bool
	}{
		{
			name:      "default params returns no error",
			params:    types.DefaultParams(),
			shouldErr: false,
		},
		{
			name:      "invalid rewards plan creation fee returns error",
			params:    types.NewParams(sdk.Coins{sdk.NewInt64Coin("umilk", 0)}),
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
