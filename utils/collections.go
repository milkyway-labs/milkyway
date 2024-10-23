package utils

import (
	"errors"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MapGetOrDefault gets a value from the map with the given key. If the key
// is not found within the map, the default value is returned.
func MapGetOrDefault[K, V any](
	ctx sdk.Context,
	collectionMap collections.Map[K, V],
	key K,
	defaultValueProvider func() V,
) (V, error) {
	value, err := collectionMap.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return *new(V), nil
		}
		return defaultValueProvider(), err
	}

	return value, nil
}
