package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/services/types"
)

type ServicesKeeper interface {
	IterateServices(ctx sdk.Context, cb func(service types.Service) (stop bool))
	SaveService(ctx sdk.Context, service types.Service) error
}

type PoolsKeeper interface {
	GetParams(ctx sdk.Context) (params poolstypes.Params)
}
