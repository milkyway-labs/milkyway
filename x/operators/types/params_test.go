package types_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v2/x/operators/types"
)

func TestParams_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		params    types.Params
		shouldErr bool
	}{
		{
			name: "invalid registration fee returns error",
			params: types.NewParams(
				sdk.Coins{sdk.Coin{Denom: "", Amount: sdkmath.NewInt(100_000_000)}},
				24*time.Hour,
			),
			shouldErr: true,
		},
		{
			name: "invalid deactivation time returns error",
			params: types.NewParams(
				nil,
				0,
			),
			shouldErr: true,
		},
		{
			name:      "default params do not return errors",
			params:    types.DefaultParams(),
			shouldErr: false,
		},
		{
			name: "valid params do not return errors",
			params: types.NewParams(
				sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
				60*time.Second,
			),
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
