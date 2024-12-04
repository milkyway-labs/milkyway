package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/milkyway-labs/milkyway/v2/x/restaking/keeper"
)

// Simulation operation weights constants
const (
	DefaultWeightMsgJoinService                int = 80
	DefaultWeightMsgLeaveService               int = 30
	DefaultWeightMsgAddOperatorToAllowList     int = 50
	DefaultWeightMsgRemoveOperatorFromAlloList int = 50
	DefaultWeightMsgBorrowPoolSecurity         int = 80
	DefaultWeightMsgCeasePoolSecurityBorrow    int = 30
	DefaultWeightMsgDelegatePool               int = 80
	DefaultWeightMsgUndelegatePool             int = 30
	DefaultWeightMsgDelegateOperator           int = 80
	DefaultWeightMsgUndelegateOperator         int = 30
	DefaultWeightMsgDelegateService            int = 80
	DefaultWeightMsgUndelegateService          int = 30
	DefaultWeightMsgSetUserPreferences         int = 20

	OperationWeightMsgJoinService                 = "op_weight_msg_join_service"
	OperationWeightMsgLeaveService                = "op_weight_msg_leave_service"
	OperationWeightMsgAddOperatorToAlloList       = "op_weight_msg_add_operator_to_allow_list"
	OperationWeightMsgRemoveOperatorFromAllowList = "op_weight_msg_remove_operator_from_allow_list"
	OperationWeightMsgBorrowPoolSecurity          = "op_weight_msg_borrow_pool_security"
	OperationWeightMsgCeasePoolSecurityBorrow     = "op_weight_msg_cease_pool_security_borrow"
	OperationWeightMsgDelegatePool                = "Op_weight_msg_delegate_pool"
	OperationWeightMsgUndelegatePool              = "op_weight_msg_undelegate_pool"
	OperationWeightMsgDelegateOperator            = "op_weight_msg_delegate_operator"
	OperationWeightMsgUndelegateOperator          = "op_weight_msg_undelegate_operator"
	OperationWeightMsgDelegateService             = "op_weight_msg_delegate_service"
	OperationWeightMsgUndelegateService           = "op_weight_msg_undelegate_service"
	OperationWeightMsgSetUserPreferences          = "op_weight_msg_set_user_preferences"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	k *keeper.Keeper,
) simulation.WeightedOperations {
	var (
		weightMsgJoinService                  int
		weightMsgLeaveService                 int
		weightMsgAddOperatorToAllowList       int
		weightMsgRemoveOperatorhFromAllowList int
		weightMsgBorrowPoolSecurity           int
		weightMsgCeasePoolSecurityBorrow      int
		weightMsgDelegatePool                 int
		weightMsgUndelegatePool               int
		weightMsgDelegateOperator             int
		weightMsgUndelegateOperator           int
		weightMsgDelegateService              int
		weightMsgUndelegateService            int
		weightMsgSetUserPreferences           int
	)

	// Generate the weights
	appParams.GetOrGenerate(OperationWeightMsgJoinService, &weightMsgJoinService, nil, func(_ *rand.Rand) {
		weightMsgJoinService = DefaultWeightMsgJoinService
	})

	appParams.GetOrGenerate(OperationWeightMsgLeaveService, &weightMsgLeaveService, nil, func(_ *rand.Rand) {
		weightMsgLeaveService = DefaultWeightMsgLeaveService
	})

	appParams.GetOrGenerate(OperationWeightMsgAddOperatorToAlloList, &weightMsgAddOperatorToAllowList, nil, func(_ *rand.Rand) {
		weightMsgAddOperatorToAllowList = DefaultWeightMsgAddOperatorToAllowList
	})

	appParams.GetOrGenerate(OperationWeightMsgRemoveOperatorFromAllowList, &weightMsgRemoveOperatorhFromAllowList, nil, func(_ *rand.Rand) {
		weightMsgRemoveOperatorhFromAllowList = DefaultWeightMsgRemoveOperatorFromAlloList
	})

	appParams.GetOrGenerate(OperationWeightMsgBorrowPoolSecurity, &weightMsgBorrowPoolSecurity, nil, func(_ *rand.Rand) {
		weightMsgBorrowPoolSecurity = DefaultWeightMsgBorrowPoolSecurity
	})

	appParams.GetOrGenerate(OperationWeightMsgCeasePoolSecurityBorrow, &weightMsgCeasePoolSecurityBorrow, nil, func(_ *rand.Rand) {
		weightMsgCeasePoolSecurityBorrow = DefaultWeightMsgCeasePoolSecurityBorrow
	})

	appParams.GetOrGenerate(OperationWeightMsgDelegatePool, &weightMsgDelegatePool, nil, func(_ *rand.Rand) {
		weightMsgDelegatePool = DefaultWeightMsgDelegatePool
	})

	appParams.GetOrGenerate(OperationWeightMsgUndelegatePool, &weightMsgUndelegatePool, nil, func(_ *rand.Rand) {
		weightMsgUndelegatePool = DefaultWeightMsgUndelegatePool
	})

	appParams.GetOrGenerate(OperationWeightMsgDelegateOperator, &weightMsgDelegateOperator, nil, func(_ *rand.Rand) {
		weightMsgDelegateOperator = DefaultWeightMsgDelegateOperator
	})

	appParams.GetOrGenerate(OperationWeightMsgUndelegateOperator, &weightMsgUndelegateOperator, nil, func(_ *rand.Rand) {
		weightMsgUndelegateOperator = DefaultWeightMsgUndelegateOperator
	})

	appParams.GetOrGenerate(OperationWeightMsgDelegateService, &weightMsgDelegateService, nil, func(_ *rand.Rand) {
		weightMsgDelegateService = DefaultWeightMsgDelegateService
	})

	appParams.GetOrGenerate(OperationWeightMsgUndelegateService, &weightMsgUndelegateService, nil, func(_ *rand.Rand) {
		weightMsgUndelegateService = DefaultWeightMsgDelegateService
	})

	appParams.GetOrGenerate(OperationWeightMsgSetUserPreferences, &weightMsgSetUserPreferences, nil, func(_ *rand.Rand) {
		weightMsgSetUserPreferences = DefaultWeightMsgSetUserPreferences
	})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(weightMsgJoinService, SimulateMsgJoinService()),
		simulation.NewWeightedOperation(weightMsgLeaveService, SimulateMsgLeaveService()),
		simulation.NewWeightedOperation(weightMsgAddOperatorToAllowList, SimulateMsgAddOperatorToAllowList()),
		simulation.NewWeightedOperation(weightMsgRemoveOperatorhFromAllowList, SimulateMsgRemoveOperatorFromAllowlist()),
		simulation.NewWeightedOperation(weightMsgBorrowPoolSecurity, SimulateMsgBorrowPoolSecurity()),
		simulation.NewWeightedOperation(weightMsgCeasePoolSecurityBorrow, SimulateMsgCeasePoolSecurityBorrow()),
		simulation.NewWeightedOperation(weightMsgDelegatePool, SimulateMsgDelegatePool()),
		simulation.NewWeightedOperation(weightMsgUndelegatePool, SimulateMsgUndelegatePool()),
		simulation.NewWeightedOperation(weightMsgDelegateOperator, SimulateMsgDelegateOperator()),
		simulation.NewWeightedOperation(weightMsgUndelegateOperator, SimulateMsgUndelegateOperator()),
		simulation.NewWeightedOperation(weightMsgDelegateService, SimulateMsgDelegateService()),
		simulation.NewWeightedOperation(weightMsgUndelegateService, SimulateMsgUndelegateService()),
		simulation.NewWeightedOperation(weightMsgSetUserPreferences, SimulateMsgSetUserPreferences()),
	}
}

func SimulateMsgJoinService() simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		return simtypes.OperationMsg{}, nil, nil
	}
}

func SimulateMsgLeaveService() simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		return simtypes.OperationMsg{}, nil, nil
	}
}

func SimulateMsgAddOperatorToAllowList() simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		return simtypes.OperationMsg{}, nil, nil
	}
}

func SimulateMsgRemoveOperatorFromAllowlist() simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		return simtypes.OperationMsg{}, nil, nil
	}
}

func SimulateMsgBorrowPoolSecurity() simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		return simtypes.OperationMsg{}, nil, nil
	}
}

func SimulateMsgCeasePoolSecurityBorrow() simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		return simtypes.OperationMsg{}, nil, nil
	}
}

func SimulateMsgDelegatePool() simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		return simtypes.OperationMsg{}, nil, nil
	}
}

func SimulateMsgUndelegatePool() simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		return simtypes.OperationMsg{}, nil, nil
	}
}

func SimulateMsgDelegateOperator() simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		return simtypes.OperationMsg{}, nil, nil
	}
}

func SimulateMsgUndelegateOperator() simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		return simtypes.OperationMsg{}, nil, nil
	}
}

func SimulateMsgDelegateService() simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		return simtypes.OperationMsg{}, nil, nil
	}
}

func SimulateMsgUndelegateService() simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		return simtypes.OperationMsg{}, nil, nil
	}
}

func SimulateMsgSetUserPreferences() simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		return simtypes.OperationMsg{}, nil, nil
	}
}
