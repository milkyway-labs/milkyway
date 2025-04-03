package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v10/utils"
	"github.com/milkyway-labs/milkyway/v10/x/investors/types"
)

func TestValidateInvestorsRewardRatio(t *testing.T) {
	testCases := []struct {
		name      string
		ratio     sdkmath.LegacyDec
		shouldErr bool
	}{
		{
			name:      "valid ratio",
			ratio:     utils.MustParseDec("0.5"), // 50%
			shouldErr: false,
		},
		{
			name:      "zero ratio",
			ratio:     utils.MustParseDec("0"), // 0%
			shouldErr: false,
		},
		{
			name:      "one ratio",
			ratio:     utils.MustParseDec("1"), // 100%
			shouldErr: false,
		},
		{
			name:      "negative ratio",
			ratio:     utils.MustParseDec("-0.1"), // -10%
			shouldErr: true,
		},
		{
			name:      "greater than one",
			ratio:     utils.MustParseDec("1.1"), // 110%
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := types.ValidateInvestorsRewardRatio(tc.ratio)
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
