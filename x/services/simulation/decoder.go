package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/milkyway-labs/milkyway/x/services/keeper"
	"github.com/milkyway-labs/milkyway/x/services/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding services type.
func NewDecodeStore(cdc codec.Codec, keeper *keeper.Keeper) func(kvA, kvB kv.Pair) string {
	collectionsDecoder := simtypes.NewStoreDecoderFuncFromCollectionsSchema(keeper.Schema)

	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.ServicePrefix):
			var serviceA, serviceB types.Service
			if err := cdc.Unmarshal(kvA.Value, &serviceA); err != nil {
				panic(err)
			}
			if err := cdc.Unmarshal(kvB.Value, &serviceB); err != nil {
				panic(err)
			}
			return fmt.Sprintf("%v\n%v", serviceA, serviceB)

		case bytes.Equal(kvA.Key[:1], types.NextServiceIDKey):
			idA := types.GetServiceIDFromBytes(kvA.Value)
			idB := types.GetServiceIDFromBytes(kvB.Value)
			return fmt.Sprintf("%v\n%v", idA, idB)

		case bytes.Equal(kvA.Key[:1], types.ServiceAddressSetPrefix):
			return collectionsDecoder(kvA, kvB)

		case bytes.Equal(kvA.Key[:1], types.ServiceParamsPrefix):
			return collectionsDecoder(kvA, kvB)

		default:
			panic(fmt.Sprintf("invalid services key prefix %X", kvA.Key[:1]))
		}
	}
}
