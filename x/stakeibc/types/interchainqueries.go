package types

// Copied from github.com/initia-labs/initia/x/move/types/connector.go

import (
	"bytes"

	"cosmossdk.io/math"

	"github.com/milkyway-labs/milkyway/x/stakeibc/move"
)

func MoveBankBalanceKey(addr []byte, denom string) ([]byte, error) {
	userAddr, err := move.NewAccountAddressFromBytes(addr[:])
	if err != nil {
		return nil, err
	}
	metadata, err := move.MetadataAddressFromDenom(denom)
	if err != nil {
		return nil, err
	}
	storeAddr := move.UserDerivedObjectAddress(userAddr, metadata)
	keyBz, err := move.GetResourceKey(storeAddr, move.StructTag{
		Address:  move.StdAddress,
		Module:   move.MoveModuleNameFungibleAsset,
		Name:     move.ResourceNameFungibleStore,
		TypeArgs: []move.TypeTag{},
	})
	if err != nil {
		return nil, err
	}
	return bytes.Join([][]byte{move.VMStorePrefix, keyBz}, nil), nil
}

func UnmarshalAmountFromMoveBankBalanceQuery(bz []byte) (math.Int, error) {
	// skipping reading metadata object
	cursor := len(move.AccountAddress{})

	// read balance
	amount, err := move.DeserializeUint64(bz[cursor : cursor+8])
	if err != nil {
		return math.ZeroInt(), err
	}

	return math.NewIntFromUint64(amount), nil
}
