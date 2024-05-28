package move

import (
	"github.com/aptos-labs/serde-reflection/serde-generate/runtime/golang/bcs"
)

var NewDeserializer = bcs.NewDeserializer

// DeserializeUint64 deserialize BCS bytes
func DeserializeUint64(bz []byte) (uint64, error) {
	d := NewDeserializer(bz)
	return d.DeserializeU64()
}
