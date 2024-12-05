package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/milkyway-labs/milkyway/v2/testutils/simtesting"
	operatorskeeper "github.com/milkyway-labs/milkyway/v2/x/operators/keeper"
	operatorssimulation "github.com/milkyway-labs/milkyway/v2/x/operators/simulation"
	operatorstypes "github.com/milkyway-labs/milkyway/v2/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v2/x/pools/types"
	"github.com/milkyway-labs/milkyway/v2/x/restaking/keeper"
	"github.com/milkyway-labs/milkyway/v2/x/restaking/types"
	serviceskeeper "github.com/milkyway-labs/milkyway/v2/x/services/keeper"
	servicessimulation "github.com/milkyway-labs/milkyway/v2/x/services/simulation"
	servicestypes "github.com/milkyway-labs/milkyway/v2/x/services/types"
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
	DefaultWeightMsgDelegateOperator           int = 80
	DefaultWeightMsgDelegateService            int = 80
	DefaultWeightMsgSetUserPreferences         int = 20

	OperationWeightMsgJoinService                 = "op_weight_msg_join_service"
	OperationWeightMsgLeaveService                = "op_weight_msg_leave_service"
	OperationWeightMsgAddOperatorToAlloList       = "op_weight_msg_add_operator_to_allow_list"
	OperationWeightMsgRemoveOperatorFromAllowList = "op_weight_msg_remove_operator_from_allow_list"
	OperationWeightMsgBorrowPoolSecurity          = "op_weight_msg_borrow_pool_security"
	OperationWeightMsgCeasePoolSecurityBorrow     = "op_weight_msg_cease_pool_security_borrow"
	OperationWeightMsgDelegatePool                = "Op_weight_msg_delegate_pool"
	OperationWeightMsgDelegateOperator            = "op_weight_msg_delegate_operator"
	OperationWeightMsgDelegateService             = "op_weight_msg_delegate_service"
	OperationWeightMsgSetUserPreferences          = "op_weight_msg_set_user_preferences"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	pk types.PoolsKeeper,
	opk types.OperatorsKeeper,
	sk types.ServicesKeeper,
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
		weightMsgDelegateOperator             int
		weightMsgDelegateService              int
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

	appParams.GetOrGenerate(OperationWeightMsgDelegateOperator, &weightMsgDelegateOperator, nil, func(_ *rand.Rand) {
		weightMsgDelegateOperator = DefaultWeightMsgDelegateOperator
	})

	appParams.GetOrGenerate(OperationWeightMsgDelegateService, &weightMsgDelegateService, nil, func(_ *rand.Rand) {
		weightMsgDelegateService = DefaultWeightMsgDelegateService
	})

	appParams.GetOrGenerate(OperationWeightMsgSetUserPreferences, &weightMsgSetUserPreferences, nil, func(_ *rand.Rand) {
		weightMsgSetUserPreferences = DefaultWeightMsgSetUserPreferences
	})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(weightMsgJoinService, SimulateMsgJoinService(ak, bk, opk, sk)),
		simulation.NewWeightedOperation(weightMsgLeaveService, SimulateMsgLeaveService(ak, bk, opk, sk, k)),
		simulation.NewWeightedOperation(weightMsgAddOperatorToAllowList, SimulateMsgAddOperatorToAllowList(ak, bk, opk, sk, k)),
		simulation.NewWeightedOperation(weightMsgRemoveOperatorhFromAllowList, SimulateMsgRemoveOperatorFromAllowlist(ak, bk, opk, sk, k)),
		simulation.NewWeightedOperation(weightMsgBorrowPoolSecurity, SimulateMsgBorrowPoolSecurity(ak, bk, pk, sk, k)),
		simulation.NewWeightedOperation(weightMsgCeasePoolSecurityBorrow, SimulateMsgCeasePoolSecurityBorrow(ak, bk, pk, sk, k)),
		simulation.NewWeightedOperation(weightMsgDelegatePool, SimulateMsgDelegatePool(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgDelegateOperator, SimulateMsgDelegateOperator(ak, bk, opk, k)),
		simulation.NewWeightedOperation(weightMsgDelegateService, SimulateMsgDelegateService(ak, bk, sk, k)),
		simulation.NewWeightedOperation(weightMsgSetUserPreferences, SimulateMsgSetUserPreferences(ak, bk, sk, k)),
	}
}

func SimulateMsgJoinService(
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	opk types.OperatorsKeeper,
	sk types.ServicesKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgJoinService{}

		operatorsKeeper := opk.(*operatorskeeper.Keeper)
		servicesKeeper := sk.(*serviceskeeper.Keeper)

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random operator
		operator, found := operatorssimulation.GetRandomExistingOperator(r, ctx, operatorsKeeper, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get operator"), nil, nil
		}

		// Get a random service
		service, found := servicessimulation.GetRandomExistingService(r, ctx, servicesKeeper, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get service"), nil, nil
		}

		msg = types.NewMsgJoinService(operator.ID, service.ID, operator.Admin)

		// Get the admin account that should sign the transaction
		adminAddress, err := sdk.AccAddressFromBech32(operator.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "error while parsing admin address"), nil, nil
		}
		signer, found := simtesting.GetSimAccount(adminAddress, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "admin account not found"), nil, nil
		}

		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateMsgLeaveService(
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	opk types.OperatorsKeeper,
	sk types.ServicesKeeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgLeaveService{}

		operatorsKeeper := opk.(*operatorskeeper.Keeper)

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get all the services
		services, err := sk.GetServices(ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get services"), nil, nil
		}

		// Get a operator that has joined a service
		var service servicestypes.Service
		operator, found := operatorssimulation.GetRandomExistingOperator(r, ctx, operatorsKeeper, func(o operatorstypes.Operator) bool {
			// Search a service that the operator has joined
			for _, s := range services {
				hasJoined, _ := k.HasOperatorJoinedService(ctx, o.ID, s.ID)
				if hasJoined {
					service = s
					return true
				}
			}
			return false
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get operator"), nil, nil
		}

		msg = types.NewMsgLeaveService(operator.ID, service.ID, operator.Admin)

		// Get the admin account that should sign the transaction
		adminAddress, err := sdk.AccAddressFromBech32(operator.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "error while parsing admin address"), nil, nil
		}
		signer, found := simtesting.GetSimAccount(adminAddress, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "admin account not found"), nil, nil
		}

		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateMsgAddOperatorToAllowList(
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	opk types.OperatorsKeeper,
	sk types.ServicesKeeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgAddOperatorToAllowList{}

		servicesKeeper := sk.(*serviceskeeper.Keeper)

		// Get all the operators
		operators, err := opk.GetOperators(ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get operators"), nil, nil
		}

		// Get a service and an operator that is not allowed
		var operator operatorstypes.Operator
		service, found := servicessimulation.GetRandomExistingService(r, ctx, servicesKeeper, func(s servicestypes.Service) bool {
			for _, o := range operators {
				isAllowed, _ := k.IsOperatorInServiceAllowList(ctx, s.ID, o.ID)
				if !isAllowed {
					operator = o
					return true
				}
			}
			return false
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get service"), nil, nil
		}

		msg = types.NewMsgAddOperatorToAllowList(service.ID, operator.ID, service.Admin)

		// Get the admin account that should sign the transaction
		adminAddress, err := sdk.AccAddressFromBech32(service.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "error while parsing admin address"), nil, nil
		}
		signer, found := simtesting.GetSimAccount(adminAddress, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "admin account not found"), nil, nil
		}

		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateMsgRemoveOperatorFromAllowlist(
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	opk types.OperatorsKeeper,
	sk types.ServicesKeeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgRemoveOperatorFromAllowlist{}

		servicesKeeper := sk.(*serviceskeeper.Keeper)

		// Get all the operators
		operators, err := opk.GetOperators(ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get operators"), nil, nil
		}

		// Get a service and an operator that is not allowed
		var operator operatorstypes.Operator
		service, found := servicessimulation.GetRandomExistingService(r, ctx, servicesKeeper, func(s servicestypes.Service) bool {
			for _, o := range operators {
				isAllowed, _ := k.IsOperatorInServiceAllowList(ctx, s.ID, o.ID)
				if isAllowed {
					operator = o
					return true
				}
			}
			return false
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get service"), nil, nil
		}

		msg = types.NewMsgRemoveOperatorFromAllowList(service.ID, operator.ID, service.Admin)

		// Get the admin account that should sign the transaction
		adminAddress, err := sdk.AccAddressFromBech32(service.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "error while parsing admin address"), nil, nil
		}
		signer, found := simtesting.GetSimAccount(adminAddress, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "admin account not found"), nil, nil
		}

		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateMsgBorrowPoolSecurity(
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	pk types.PoolsKeeper,
	sk types.ServicesKeeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgBorrowPoolSecurity{}

		servicesKeeper := sk.(*serviceskeeper.Keeper)

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get all the pools
		pools, err := pk.GetPools(ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "error while getting pools"), nil, nil
		}

		var pool poolstypes.Pool
		service, found := servicessimulation.GetRandomExistingService(r, ctx, servicesKeeper, func(s servicestypes.Service) bool {
			for _, p := range pools {
				isBorrowing, _ := k.IsServiceSecuredByPool(ctx, s.ID, p.ID)
				if !isBorrowing {
					pool = p
					return true
				}
			}
			return false
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get service"), nil, nil
		}

		msg = types.NewMsgBorrowPoolSecurity(service.ID, pool.ID, service.Admin)

		// Get the admin account that should sign the transaction
		adminAddress, err := sdk.AccAddressFromBech32(service.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "error while parsing admin address"), nil, nil
		}
		signer, found := simtesting.GetSimAccount(adminAddress, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "admin account not found"), nil, nil
		}

		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateMsgCeasePoolSecurityBorrow(
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	pk types.PoolsKeeper,
	sk types.ServicesKeeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgCeasePoolSecurityBorrow{}

		servicesKeeper := sk.(*serviceskeeper.Keeper)

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get all the pools
		pools, err := pk.GetPools(ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "error while getting pools"), nil, nil
		}

		var pool poolstypes.Pool
		service, found := servicessimulation.GetRandomExistingService(r, ctx, servicesKeeper, func(s servicestypes.Service) bool {
			for _, p := range pools {
				isBorrowing, _ := k.IsServiceSecuredByPool(ctx, s.ID, p.ID)
				if isBorrowing {
					pool = p
					return true
				}
			}
			return false
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get service"), nil, nil
		}

		msg = types.NewMsgCeasePoolSecurityBorrow(service.ID, pool.ID, service.Admin)

		// Get the admin account that should sign the transaction
		adminAddress, err := sdk.AccAddressFromBech32(service.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "error while parsing admin address"), nil, nil
		}
		signer, found := simtesting.GetSimAccount(adminAddress, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "admin account not found"), nil, nil
		}

		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateMsgDelegatePool(
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgDelegatePool{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random delegator with a random amount
		delegator, coins, skip := randomDelegatorAndAmount(r, ctx, accs, bk, ak)

		// If coins slice is empty, we can not create valid msg
		if len(coins) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "empty coins slice"), nil, nil
		}

		if skip {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "skip delegate"), nil, nil
		}

		msg = types.NewMsgDelegatePool(coins[0], delegator.Address.String())

		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, delegator)
	}
}

func SimulateMsgDelegateOperator(
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	opk types.OperatorsKeeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgDelegateOperator{}

		operatorsKeeper := opk.(*operatorskeeper.Keeper)

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random operator
		operator, found := operatorssimulation.GetRandomExistingOperator(r, ctx, operatorsKeeper, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get operator"), nil, nil
		}

		// Get a random delegator with a random amount
		delegator, coins, skip := randomDelegatorAndAmount(r, ctx, accs, bk, ak)

		// If coins slice is empty, we can not create valid msg
		if len(coins) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "empty coins slice"), nil, nil
		}

		if skip {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "skip delegate"), nil, nil
		}

		msg = types.NewMsgDelegateOperator(operator.ID, coins, delegator.Address.String())

		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, delegator)
	}
}

func SimulateMsgDelegateService(
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	sk types.ServicesKeeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgDelegateService{}

		servicesKeeper := sk.(*serviceskeeper.Keeper)

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random service
		operator, found := servicessimulation.GetRandomExistingService(r, ctx, servicesKeeper, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get service"), nil, nil
		}

		// Get a random delegator with a random amount
		delegator, coins, skip := randomDelegatorAndAmount(r, ctx, accs, bk, ak)

		// If coins slice is empty, we can not create valid msg
		if len(coins) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "empty coins slice"), nil, nil
		}

		if skip {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "skip delegate"), nil, nil
		}

		msg = types.NewMsgDelegateService(operator.ID, coins, delegator.Address.String())

		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, delegator)
	}
}

func SimulateMsgSetUserPreferences(
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	sk types.ServicesKeeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgSetUserPreferences{}

		delegator, _ := simtypes.RandomAcc(r, accs)
		acc := ak.GetAccount(ctx, delegator.Address)
		if acc == nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		services, err := sk.GetServices(ctx)
		if err != nil {
			panic(err)
		}
		msg = types.NewMsgSetUserPreferences(RandomUserPreferences(r, services), delegator.Address.String())

		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, delegator)
	}
}

// Get a random account with an amount that can be delegated from the provided
// account.
func randomDelegatorAndAmount(
	r *rand.Rand, ctx sdk.Context, accs []simtypes.Account, bk bankkeeper.Keeper, ak authkeeper.AccountKeeper,
) (simtypes.Account, sdk.Coins, bool) {
	delegator, _ := simtypes.RandomAcc(r, accs)

	acc := ak.GetAccount(ctx, delegator.Address)
	if acc == nil {
		return delegator, nil, true
	}

	spendable := bk.SpendableCoins(ctx, acc.GetAddress())

	sendCoins := simtypes.RandSubsetCoins(r, spendable)
	if sendCoins.Empty() {
		return delegator, nil, true
	}

	return delegator, sendCoins, false
}
