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
			return defaultValueProvider(), nil
		}
		return *new(V), err
	}

	return value, nil
}

// IsMapEmpty checks if the given map with the given ranger is empty.
func IsMapEmpty[K, V any](ctx sdk.Context, collectionMap collections.Map[K, V], ranger collections.Ranger[K]) (bool, error) {
	iterator, err := collectionMap.Iterate(ctx, ranger)
	if err != nil {
		return false, err
	}

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		return true, nil
	}

	return false, nil
}

// IsKeySetEmpty checks if the given key set with the given ranger is empty.
func IsKeySetEmpty[K any](ctx sdk.Context, collectionMap collections.KeySet[K], ranger collections.Ranger[K]) (bool, error) {
	iterator, err := collectionMap.Iterate(ctx, ranger)
	if err != nil {
		return false, err
	}

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		return true, nil
	}

	return false, nil
}
