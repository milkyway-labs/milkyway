package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/avs/types"
)

func TestParams_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		params    types.Params
		shouldErr bool
	}{
		{
			name: "invalid AVS registration fee returns error",
			params: types.Params{
				AvsRegistrationFee: sdk.Coins{sdk.Coin{Denom: "", Amount: sdkmath.NewInt(100)}},
			},
			shouldErr: true,
		},
		{
			name:      "default params returns no error",
			params:    types.DefaultParams(),
			shouldErr: false,
		},
		{
			name: "valid params returns no error",
			params: types.Params{
				AvsRegistrationFee: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100))),
			},
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
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