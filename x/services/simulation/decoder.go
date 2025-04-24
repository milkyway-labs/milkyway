package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/kv"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/milkyway-labs/milkyway/v11/x/services/keeper"
	"github.com/milkyway-labs/milkyway/v11/x/services/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding services type.
func NewDecodeStore(keeper *keeper.Keeper) func(kvA kv.Pair, kvB kv.Pair) string {
	collectionsDecoder := simtypes.NewStoreDecoderFuncFromCollectionsSchema(keeper.Schema)

	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.NextServiceIDKey):
			return collectionsDecoder(kvA, kvB)

		case bytes.Equal(kvA.Key[:1], types.ServicePrefix):
			return collectionsDecoder(kvA, kvB)

		case bytes.Equal(kvA.Key[:1], types.ServiceAddressSetPrefix):
			return collectionsDecoder(kvA, kvB)

		case bytes.Equal(kvA.Key[:1], types.ServiceParamsPrefix):
			return collectionsDecoder(kvA, kvB)

		case bytes.Equal(kvA.Key[:1], types.ParamsKey):
			return collectionsDecoder(kvA, kvB)

		default:
			panic(fmt.Sprintf("invalid services key prefix %X", kvA.Key[:1]))
		}
	}
}
