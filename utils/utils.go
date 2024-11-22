package utils

import (
	"encoding/json"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MustParseCoin(s string) sdk.Coin {
	c, err := sdk.ParseCoinNormalized(strings.ReplaceAll(s, "_", ""))
	if err != nil {
		panic(err)
	}
	return c
}

func MustParseCoins(s string) sdk.Coins {
	c, err := sdk.ParseCoinsNormalized(strings.ReplaceAll(s, "_", ""))
	if err != nil {
		panic(err)
	}
	return c
}

func MustParseDecCoins(s string) sdk.DecCoins {
	d, err := sdk.ParseDecCoins(strings.ReplaceAll(s, "_", ""))
	if err != nil {
		panic(err)
	}
	return d
}

func MustParseDec(s string) sdkmath.LegacyDec {
	return sdkmath.LegacyMustNewDecFromStr(strings.ReplaceAll(s, "_", ""))
}

// JSONStringHasKey parses the provided data as a json object and checks
// if it contains the provided key.
func JSONStringHasKey(data, key string) (found bool, jsonObject map[string]interface{}) {
	jsonObject = make(map[string]interface{})

	// If there is no data, nothing to do here.
	if len(data) == 0 {
		return false, jsonObject
	}

	// the jsonObject must be a valid JSON object
	err := json.Unmarshal([]byte(data), &jsonObject)
	if err != nil {
		return false, jsonObject
	}

	// Check if the provided key exists in the jsonObject.
	_, ok := jsonObject[key]
	if !ok {
		return false, jsonObject
	}

	return true, jsonObject
}
