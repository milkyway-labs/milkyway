package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.Codec

	authority string
}

func NewKeeper(cdc codec.Codec, storeKey storetypes.StoreKey, authority string) *Keeper {
	return &Keeper{
		storeKey:  storeKey,
		cdc:       cdc,
		authority: authority,
	}
}
