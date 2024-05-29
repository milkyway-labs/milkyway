package move

var (
	VMStorePrefix = []byte{0x21} // prefix for vm

	ResourceSeparator = byte(2)
)

// GetResourceKey returns the store key of the Move resource
func GetResourceKey(addr AccountAddress, structTag StructTag) ([]byte, error) {
	bz, err := structTag.BcsSerialize()
	if err != nil {
		return nil, err
	}

	return append(append(addr.Bytes(), ResourceSeparator), bz...), nil
}
