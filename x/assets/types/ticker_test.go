package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v11/x/assets/types"
)

func TestValidateTicker(t *testing.T) {
	testCases := []struct {
		name   string
		ticker string
		expErr bool
	}{
		{
			name:   "uppercase ticker returns no error",
			ticker: "MILK",
			expErr: false,
		},
		{
			name:   "lowercase letters are accepted",
			ticker: "milkINIT",
			expErr: false,
		},
		{
			name:   "empty ticker returns error",
			ticker: "",
			expErr: true,
		},
		{
			name:   "invalid characters return error",
			ticker: "WOW!",
			expErr: true,
		},
		{
			name:   "too long ticker returns error",
			ticker: "WHATALONGTICKER",
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := types.ValidateTicker(tc.ticker)
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
