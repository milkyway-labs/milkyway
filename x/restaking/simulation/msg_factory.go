package simulation

import (
	"math/rand"
	"slices"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/milkyway-labs/milkyway/v10/testutils/simtesting"
	operatorskeeper "github.com/milkyway-labs/milkyway/v10/x/operators/keeper"
	operatorssimulation "github.com/milkyway-labs/milkyway/v10/x/operators/simulation"
	operatorstypes "github.com/milkyway-labs/milkyway/v10/x/operators/types"
	poolskeeper "github.com/milkyway-labs/milkyway/v10/x/pools/keeper"
	poolssimulation "github.com/milkyway-labs/milkyway/v10/x/pools/simulation"
	"github.com/milkyway-labs/milkyway/v10/x/restaking/keeper"
	"github.com/milkyway-labs/milkyway/v10/x/restaking/types"
	serviceskeeper "github.com/milkyway-labs/milkyway/v10/x/services/keeper"
	servicessimulation "github.com/milkyway-labs/milkyway/v10/x/services/simulation"
	servicestypes "github.com/milkyway-labs/milkyway/v10/x/services/types"
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
	pk *poolskeeper.Keeper,
	opk *operatorskeeper.Keeper,
	sk *serviceskeeper.Keeper,
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
		simulation.NewWeightedOperation(weightMsgJoinService, SimulateMsgJoinService(ak, bk, opk, sk, k)),
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
	opk *operatorskeeper.Keeper,
	sk *serviceskeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgJoinService{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random operator
		operator, found := operatorssimulation.GetRandomExistingOperator(r, ctx, opk, func(o operatorstypes.Operator) bool {
			return o.IsActive()
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get operator"), nil, nil
		}

		// Get a random service
		service, found := servicessimulation.GetRandomExistingService(r, ctx, sk, func(s servicestypes.Service) bool {
			if !s.IsActive() {
				return false
			}
			configured, err := k.IsServiceOperatorsAllowListConfigured(ctx, s.ID)
			if err != nil {
				panic(err)
			}
			if !configured {
				return true
			}
			isAllowed, err := k.IsOperatorInServiceAllowList(ctx, s.ID, operator.ID)
			if err != nil {
				panic(err)
			}
			return isAllowed
		})
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
	opk *operatorskeeper.Keeper,
	sk *serviceskeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgLeaveService{}

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
		operator, found := operatorssimulation.GetRandomExistingOperator(r, ctx, opk, func(o operatorstypes.Operator) bool {
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
	opk *operatorskeeper.Keeper,
	sk *serviceskeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgAddOperatorToAllowList{}

		// Get a random operator
		operator, found := operatorssimulation.GetRandomExistingOperator(r, ctx, opk, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get operators"), nil, nil
		}

		// Get a random service
		service, found := servicessimulation.GetRandomExistingService(r, ctx, sk, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get service"), nil, nil
		}

		// Ensure that the operator is not in the service allow list
		isAllowed, err := k.IsOperatorInServiceAllowList(ctx, service.ID, operator.ID)
		if err != nil {
			panic(err)
		}
		if isAllowed {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "operator is already in the service allow list"), nil, nil
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
	opk *operatorskeeper.Keeper,
	sk *serviceskeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgRemoveOperatorFromAllowlist{}

		// Get a random operator
		operator, found := operatorssimulation.GetRandomExistingOperator(r, ctx, opk, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get operators"), nil, nil
		}

		// Get a random service
		service, found := servicessimulation.GetRandomExistingService(r, ctx, sk, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get service"), nil, nil
		}

		// Ensure that the operator is in the service allow list
		isAllowed, err := k.IsOperatorInServiceAllowList(ctx, service.ID, operator.ID)
		if err != nil {
			panic(err)
		}
		if !isAllowed {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "operator is not in the service allow list"), nil, nil
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
	pk *poolskeeper.Keeper,
	sk *serviceskeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgBorrowPoolSecurity{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random pool
		pool, found := poolssimulation.GetRandomExistingPool(r, ctx, pk, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "error while getting pools"), nil, nil
		}

		// Get a random service
		service, found := servicessimulation.GetRandomExistingService(r, ctx, sk, func(s servicestypes.Service) bool {
			return s.IsActive()
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get service"), nil, nil
		}

		isServiceSecured, err := k.IsServiceSecuredByPool(ctx, service.ID, pool.ID)
		if err != nil {
			panic(err)
		}
		if isServiceSecured {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "service already secured by pool"), nil, nil
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
	pk *poolskeeper.Keeper,
	sk *serviceskeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgCeasePoolSecurityBorrow{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random pool
		pool, found := poolssimulation.GetRandomExistingPool(r, ctx, pk, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "error while getting pools"), nil, nil
		}

		// Get a random service
		service, found := servicessimulation.GetRandomExistingService(r, ctx, sk, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get service"), nil, nil
		}

		isServiceSecured, err := k.IsServiceSecuredByPool(ctx, service.ID, pool.ID)
		if err != nil {
			panic(err)
		}
		if !isServiceSecured {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "service not secured by pool"), nil, nil
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
		delegator, coins, skip := randomDelegatorAndAmount(r, ctx, accs, k, bk, ak)

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
	opk *operatorskeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgDelegateOperator{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random operator
		operator, found := operatorssimulation.GetRandomExistingOperator(r, ctx, opk, func(o operatorstypes.Operator) bool {
			return o.IsActive()
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get operator"), nil, nil
		}

		// Get a random delegator with a random amount
		delegator, coins, skip := randomDelegatorAndAmount(r, ctx, accs, k, bk, ak)

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
	sk *serviceskeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgDelegateService{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random service
		service, found := servicessimulation.GetRandomExistingService(r, ctx, sk, func(s servicestypes.Service) bool {
			return s.IsActive()
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get service"), nil, nil
		}

		// Get a random delegator with a random amount
		delegator, coins, skip := randomDelegatorAndAmount(r, ctx, accs, k, bk, ak)
		// Filter the coins to only those that can be restaked toward the service
		coins = filterServiceRestakableCoins(sk, ctx, service.ID, coins)

		// If coins slice is empty, we can not create valid msg
		if len(coins) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "empty coins slice"), nil, nil
		}

		if skip {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "skip delegate"), nil, nil
		}

		msg = types.NewMsgDelegateService(service.ID, coins, delegator.Address.String())

		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, delegator)
	}
}

func SimulateMsgSetUserPreferences(
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	sk *serviceskeeper.Keeper,
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
	r *rand.Rand, ctx sdk.Context, accs []simtypes.Account, k *keeper.Keeper, bk bankkeeper.Keeper, ak authkeeper.AccountKeeper,
) (simtypes.Account, sdk.Coins, bool) {
	delegator, _ := simtypes.RandomAcc(r, accs)

	acc := ak.GetAccount(ctx, delegator.Address)
	if acc == nil {
		return delegator, nil, true
	}

	spendable := bk.SpendableCoins(ctx, acc.GetAddress())
	// Filter the spendable coins to only include the restakable ones
	restakableCoins := sdk.NewCoins()
	for _, coin := range spendable {
		isRestakable, err := k.IsDenomRestakable(ctx, coin.Denom)
		if err != nil {
			panic(err)
		}
		if isRestakable {
			restakableCoins = restakableCoins.Add(coin)
		}
	}

	coins := simtypes.RandSubsetCoins(r, restakableCoins)
	if coins.Empty() {
		return delegator, nil, true
	}

	return delegator, coins, false
}

func filterServiceRestakableCoins(sk *serviceskeeper.Keeper, ctx sdk.Context, serviceID uint32, coins sdk.Coins) sdk.Coins {
	serviceParams, err := sk.GetServiceParams(ctx, serviceID)
	if err != nil {
		panic(err)
	}
	// If empty allows all denoms
	if len(serviceParams.AllowedDenoms) == 0 {
		return coins
	}

	filteredCoins := sdk.NewCoins()
	for _, coin := range coins {
		if slices.Contains(serviceParams.AllowedDenoms, coin.Denom) {
			filteredCoins = filteredCoins.Add(coin)
		}
	}

	return filteredCoins
}
