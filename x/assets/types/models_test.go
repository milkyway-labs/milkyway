package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v9/x/assets/types"
)

func TestAsset_Validate(t *testing.T) {
	testCases := []struct {
		name   string
		asset  types.Asset
		expErr bool
	}{
		{
			name:   "invalid denom returns error",
			asset:  types.NewAsset("@#$", "bitcoin", 1),
			expErr: true,
		},
		{
			name:   "invalid ticker returns error",
			asset:  types.NewAsset("btc", "@#$%", 1),
			expErr: true,
		},
		{
			name:   "valid asset returns no error",
			asset:  types.NewAsset("btc", "bitcoin", 1),
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.asset.Validate()
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
