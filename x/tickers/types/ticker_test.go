package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/tickers/types"
)

func TestValidateTicker(t *testing.T) {
	for _, tc := range []struct {
		name        string
		ticker      string
		expectedErr string
	}{
		{
			"happy case",
			"MILK",
			"",
		},
		{
			"lowercase letters are accepted",
			"milkINIT",
			"",
		},
		{
			"empty",
			"",
			"empty ticker",
		},
		{
			"invalid characters",
			"WOW!",
			"bad ticker format: WOW!",
		},
		{
			"too long",
			"WHATALONGTICKER",
			"ticker too long",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := types.ValidateTicker(tc.ticker)
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}
