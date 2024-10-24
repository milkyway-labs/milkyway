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
}

type RestakingKeeper interface {
	GetOperatorJoinedServices(ctx sdk.Context, operatorID uint32) (types.OperatorJoinedServices, error)
	AddServiceToOperator(ctx sdk.Context, operatorID uint32, serviceID uint32) error
	AddOperatorToServiceAllowList(ctx sdk.Context, serviceID uint32, operatorID uint32) error
	AddPoolToServiceSecuringPools(ctx sdk.Context, serviceID uint32, poolID uint32) error
}
