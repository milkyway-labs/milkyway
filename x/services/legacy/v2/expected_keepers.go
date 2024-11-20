package v2

import (
	"context"

	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/services/types"
)

type ServicesKeeper interface {
	IterateServices(ctx context.Context, cb func(service types.Service) (stop bool, err error)) error
	SaveService(ctx context.Context, service types.Service) error
}

type PoolsKeeper interface {
	GetParams(ctx context.Context) (poolstypes.Params, error)
}
