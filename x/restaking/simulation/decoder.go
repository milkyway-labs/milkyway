package simulation

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/milkyway-labs/milkyway/v10/x/restaking/keeper"
	"github.com/milkyway-labs/milkyway/v10/x/restaking/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding restaking type.
func NewDecodeStore(cdc codec.BinaryCodec, keeper *keeper.Keeper) func(kvA kv.Pair, kvB kv.Pair) string {
	collectionsDecoder := simtypes.NewStoreDecoderFuncFromCollectionsSchema(keeper.Schema)

	return func(kvA, kvB kv.Pair) string {
		switch {

		case bytes.Equal(kvA.Key[:1], types.UnbondingIDKey):
			idA := binary.BigEndian.Uint64(kvA.Value)
			idB := binary.BigEndian.Uint64(kvB.Value)
			return fmt.Sprintf("%d\n%d", idA, idB)

		case bytes.Equal(kvA.Key[:1], types.UnbondingIndexKey):
			return fmt.Sprintf("%v\n%v", kvA.Value, kvB.Value)

		case bytes.Equal(kvA.Key[:1], types.UnbondingTypeKey):
			typeA := binary.BigEndian.Uint32(kvA.Value)
			typeB := binary.BigEndian.Uint32(kvB.Value)
			return fmt.Sprintf("%d\n%d", typeA, typeB)

		case bytes.Equal(kvA.Key[:1], types.PoolDelegationPrefix):
			delegationA := types.MustUnmarshalDelegation(cdc, kvA.Value)
			delegationB := types.MustUnmarshalDelegation(cdc, kvB.Value)
			return fmt.Sprintf("%v\n%v", delegationA, delegationB)

		case bytes.Equal(kvA.Key[:1], types.PoolUnbondingDelegationPrefix):
			undelegationA := types.MustUnmarshalUnbondingDelegation(cdc, kvA.Value)
			undelegationB := types.MustUnmarshalUnbondingDelegation(cdc, kvB.Value)
			return fmt.Sprintf("%v\n%v", undelegationA, undelegationB)

		case bytes.Equal(kvA.Key[:1], types.OperatorDelegationPrefix):
			delegationA := types.MustUnmarshalDelegation(cdc, kvA.Value)
			delegationB := types.MustUnmarshalDelegation(cdc, kvB.Value)
			return fmt.Sprintf("%v\n%v", delegationA, delegationB)

		case bytes.Equal(kvA.Key[:1], types.OperatorUnbondingDelegationPrefix):
			undelegationA := types.MustUnmarshalUnbondingDelegation(cdc, kvA.Value)
			undelegationB := types.MustUnmarshalUnbondingDelegation(cdc, kvB.Value)
			return fmt.Sprintf("%v\n%v", undelegationA, undelegationB)

		case bytes.Equal(kvA.Key[:1], types.ServiceDelegationPrefix):
			delegationA := types.MustUnmarshalDelegation(cdc, kvA.Value)
			delegationB := types.MustUnmarshalDelegation(cdc, kvB.Value)
			return fmt.Sprintf("%v\n%v", delegationA, delegationB)

		case bytes.Equal(kvA.Key[:1], types.ServiceUnbondingDelegationPrefix):
			undelegationA := types.MustUnmarshalUnbondingDelegation(cdc, kvA.Value)
			undelegationB := types.MustUnmarshalUnbondingDelegation(cdc, kvB.Value)
			return fmt.Sprintf("%v\n%v", undelegationA, undelegationB)

		case bytes.Equal(kvA.Key[:1], types.UnbondingQueueKey):
			var listA, listB types.DTDataList
			cdc.MustUnmarshal(kvA.Value, &listA)
			cdc.MustUnmarshal(kvB.Value, &listB)
			return fmt.Sprintf("%v\n%v", listA, listB)

		// Collections
		case bytes.Equal(kvA.Key[:1], types.OperatorJoinedServicesPrefix),
			bytes.Equal(kvA.Key[:1], types.ServiceOperatorsAllowListPrefix),
			bytes.Equal(kvA.Key[:1], types.ServiceSecuringPoolsPrefix),
			bytes.Equal(kvA.Key[:1], types.UserPreferencesPrefix):
			return collectionsDecoder(kvA, kvB)

		default:
			panic(fmt.Sprintf("invalid restaking key prefix %X", kvA.Key[:1]))
		}
	}
}
