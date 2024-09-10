package v2_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/app"
	"github.com/milkyway-labs/milkyway/testutils/storetesting"
	v2 "github.com/milkyway-labs/milkyway/x/services/legacy/v2"
	"github.com/milkyway-labs/milkyway/x/services/types"
)

func TestMigrateStore(t *testing.T) {
	cdc, _ := app.MakeCodecs()

	// Build all the necessary keys
	keys := storetypes.NewKVStoreKeys(types.StoreKey)
	memKeys := storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name:      "module params are set properly",
			shouldErr: false,
			check: func(ctx sdk.Context) {
				store := ctx.KVStore(keys[types.StoreKey])

				// Check the module params
				paramsBz := store.Get(types.ParamsKey)
				require.NotNil(t, paramsBz)

				var params types.Params
				err := cdc.Unmarshal(paramsBz, &params)
				require.NoError(t, err)
				require.Equal(t, types.DefaultParams(), params)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := storetesting.BuildContext(keys, nil, memKeys)
			if tc.store != nil {
				tc.store(ctx)
			}

			err := v2.MigrateStore(ctx, keys[types.StoreKey], cdc)
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tc.check != nil {
					tc.check(ctx)
				}
			}
		})
	}
}
