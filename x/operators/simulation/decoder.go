package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/kv"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/milkyway-labs/milkyway/x/operators/keeper"
	"github.com/milkyway-labs/milkyway/x/operators/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding services type.
func NewDecodeStore(keeper *keeper.Keeper) func(kvA kv.Pair, kvB kv.Pair) string {
	collectionsDecoder := simtypes.NewStoreDecoderFuncFromCollectionsSchema(keeper.Schema)

	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.ParamsKey),
			bytes.Equal(kvA.Key[:1], types.NextOperatorIDKey),
			bytes.Equal(kvA.Key[:1], types.OperatorPrefix),
			bytes.Equal(kvA.Key[:1], types.OperatorAddressSetPrefix),
			bytes.Equal(kvA.Key[:1], types.OperatorParamsMapPrefix):
			return collectionsDecoder(kvA, kvB)

		case bytes.Equal(kvA.Key[:1], types.InactivatingOperatorQueuePrefix):
			valueA := types.GetOperatorIDFromBytes(kvA.Value)
			valueB := types.GetOperatorIDFromBytes(kvB.Value)
			return fmt.Sprintf("operatorIDA: %d\noperatorIDB: %d", valueA, valueB)

		default:
			panic(fmt.Sprintf("invalid operators key prefix %X", kvA.Key[:1]))
		}
	}
}
