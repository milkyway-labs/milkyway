package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

type OperatorsKeeper interface {
	GetOperator(ctx sdk.Context, operatorID uint32) (operatorstypes.Operator, bool)
	SaveOperatorParams(ctx sdk.Context, operatorID uint32, params operatorstypes.OperatorParams) error
}

type ServicesKeeper interface {
	GetService(ctx sdk.Context, serviceID uint32) (servicestypes.Service, bool)
	SaveServiceParams(ctx sdk.Context, serviceID uint32, params servicestypes.ServiceParams) error
}

type RestakingKeeper interface {
	GetOperatorJoinedServices(ctx sdk.Context, operatorID uint32) (types.OperatorJoinedServices, error)
	SaveOperatorJoinedServices(ctx sdk.Context, operatorID uint32, joinedServices types.OperatorJoinedServices) error
	ServiceWhitelistOperator(ctx sdk.Context, serviceID uint32, operatorID uint32) error
	ServiceWhitelistPool(ctx sdk.Context, serviceID uint32, poolID uint32) error
}
