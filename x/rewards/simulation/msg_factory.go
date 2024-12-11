package simulation

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/milkyway-labs/milkyway/v3/testutils/simtesting"
	"github.com/milkyway-labs/milkyway/v3/utils"
	operatorskeeper "github.com/milkyway-labs/milkyway/v3/x/operators/keeper"
	operatorssimulation "github.com/milkyway-labs/milkyway/v3/x/operators/simulation"
	operatorstypes "github.com/milkyway-labs/milkyway/v3/x/operators/types"
	poolskeeper "github.com/milkyway-labs/milkyway/v3/x/pools/keeper"
	poolssimulation "github.com/milkyway-labs/milkyway/v3/x/pools/simulation"
	restakingtypes "github.com/milkyway-labs/milkyway/v3/x/restaking/types"
	"github.com/milkyway-labs/milkyway/v3/x/rewards/keeper"
	"github.com/milkyway-labs/milkyway/v3/x/rewards/types"
	serviceskeeper "github.com/milkyway-labs/milkyway/v3/x/services/keeper"
	servicessimulation "github.com/milkyway-labs/milkyway/v3/x/services/simulation"
)

// Simulation operation weights constants
const (
	DefaultWeightMsgCreateRewardsPlan          int = 80
	DefaultWeightMsgEditRewardsPlan            int = 40
	DefaultWeightMsgSetWithdrawAddress         int = 20
	DefaultWeightMsgWithdrawDelegatorReward    int = 30
	DefaultWeightMsgWithdrawOperatorCommission int = 10

	OperationWeightMsgCreateRewardsPlan          = "op_weight_msg_create_rewards_plan"
	OperationWeightMsgEditRewardsPlan            = "op_weight_msg_edit_rewards_plan"
	OperationWeightMsgSetWithdrawAddress         = "op_weight_msg_set_withdraw_address"
	OperationWeightMsgWithdrawDelegatorReward    = "op_weight_msg_withdraw_delegator_reward"
	OperationWeightMsgWithdrawOperatorCommission = "op_weight_msg_withdraw_operator_commission"
)

func WeightedOperations(
	appParams simtypes.AppParams,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	pk *poolskeeper.Keeper,
	ok *operatorskeeper.Keeper,
	sk *serviceskeeper.Keeper,
	k *keeper.Keeper,
) simulation.WeightedOperations {
	var (
		weightMsgCreateRewardsPlan          int
		weightMsgEditRewardsPlan            int
		weightMsgSetWithdrawAddress         int
		weightMsgWithdrawDelegatorReward    int
		weightMsgWithdrawOperatorCommission int
	)

	// Generate the weights for the messages
	appParams.GetOrGenerate(OperationWeightMsgCreateRewardsPlan, &weightMsgCreateRewardsPlan, nil, func(_ *rand.Rand) {
		weightMsgCreateRewardsPlan = DefaultWeightMsgCreateRewardsPlan
	})

	appParams.GetOrGenerate(OperationWeightMsgEditRewardsPlan, &weightMsgEditRewardsPlan, nil, func(_ *rand.Rand) {
		weightMsgEditRewardsPlan = DefaultWeightMsgEditRewardsPlan
	})

	appParams.GetOrGenerate(OperationWeightMsgSetWithdrawAddress, &weightMsgSetWithdrawAddress, nil, func(_ *rand.Rand) {
		weightMsgSetWithdrawAddress = DefaultWeightMsgEditRewardsPlan
	})

	appParams.GetOrGenerate(OperationWeightMsgWithdrawDelegatorReward, &weightMsgWithdrawDelegatorReward, nil, func(_ *rand.Rand) {
		weightMsgWithdrawDelegatorReward = DefaultWeightMsgWithdrawDelegatorReward
	})

	appParams.GetOrGenerate(OperationWeightMsgWithdrawOperatorCommission, &weightMsgWithdrawOperatorCommission, nil, func(_ *rand.Rand) {
		weightMsgWithdrawOperatorCommission = DefaultWeightMsgWithdrawOperatorCommission
	})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(weightMsgCreateRewardsPlan, SimulateMsgCreateRewardsPlan(ak, bk, pk, ok, sk, k)),
		simulation.NewWeightedOperation(weightMsgEditRewardsPlan, SimulateMsgEditRewardsPlan(ak, bk, pk, ok, sk, k)),
		simulation.NewWeightedOperation(weightMsgSetWithdrawAddress, SimulateMsgSetWithdrawAddress(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgWithdrawDelegatorReward, SimulateMsgWithdrawDelegatorReward(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgWithdrawOperatorCommission, SimulateMsgWithdrawOperatorCommission(ak, bk, ok, k)),
	}
}

// --------------------------------------------------------------------------------------------------------------------

func SimulateMsgCreateRewardsPlan(
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	pk *poolskeeper.Keeper,
	ok *operatorskeeper.Keeper,
	sk *serviceskeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgCreateRewardsPlan{}

		// Get a random service
		service, found := servicessimulation.GetRandomExistingService(r, ctx, sk, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no service found"), nil, nil
		}

		// Get the module parameters to get the fees required to create a rewards plan
		rewardsParams, err := k.GetParams(ctx)
		if err != nil {
			panic(err)
		}

		// Get a random account
		adminAddress, err := sdk.AccAddressFromBech32(service.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "invalid admin address"), nil, nil
		}

		// Ensure the sender has enough balance to create a rewards plan
		senderSpendableBalance := bk.SpendableCoins(ctx, adminAddress)
		if !senderSpendableBalance.IsAllGTE(rewardsParams.RewardsPlanCreationFee) {
			return simtypes.NoOpMsg(
				types.ModuleName,
				sdk.MsgTypeURL(msg),
				fmt.Sprintf("sender: %s don't have enough balance to create rewards plan, required: %s, available: %s",
					service.Admin,
					rewardsParams.RewardsPlanCreationFee.String(),
					senderSpendableBalance.String(),
				),
			), nil, nil
		}

		// Get a random pool that we will use to configure the pool distribution
		pool, found := poolssimulation.GetRandomExistingPool(r, ctx, pk, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no pool found"), nil, nil
		}

		// Get a random operator that we will use to configure the operator distribution
		operator, found := operatorssimulation.GetRandomExistingOperator(r, ctx, ok, func(o operatorstypes.Operator) bool {
			return o.IsActive()
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no operator found"), nil, nil
		}

		// Compute some random start/end time
		rewardsStart := ctx.BlockTime().Add(time.Hour * time.Duration(r.Intn(10)+1))
		rewardsEnd := rewardsStart.Add(time.Hour * time.Duration(r.Intn(96)+1))

		// Get a random rewards plan amount
		amount, err := simtypes.RandomFees(r, ctx, senderSpendableBalance)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), err.Error()), nil, nil
		}

		signer, found := simtesting.GetSimAccount(adminAddress, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "admin account not found"), nil, nil
		}

		msg = types.NewMsgCreateRewardsPlan(
			service.ID,
			simtypes.RandStringOfLength(r, 32),
			amount,
			rewardsStart,
			rewardsEnd,
			RandomDistribution(r, restakingtypes.DELEGATION_TYPE_POOL, pool),
			RandomDistribution(r, restakingtypes.DELEGATION_TYPE_OPERATOR, operator),
			RandomUsersDistribution(r),
			rewardsParams.RewardsPlanCreationFee,
			service.Admin,
		)
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateMsgEditRewardsPlan(
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	pk *poolskeeper.Keeper,
	ok *operatorskeeper.Keeper,
	sk *serviceskeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgEditRewardsPlan{}

		plan, found := GetRandomExistingRewardsPlan(r, ctx, k)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no rewards plan found"), nil, nil
		}

		// Get a random pool that we will use to configure the pool distribution
		pool, found := poolssimulation.GetRandomExistingPool(r, ctx, pk, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no pool found"), nil, nil
		}

		// Get a random operator that we will use to configure the operator distribution
		operator, found := operatorssimulation.GetRandomExistingOperator(r, ctx, ok, func(o operatorstypes.Operator) bool {
			return o.IsActive()
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no operator found"), nil, nil
		}

		// Get the service admin
		service, found, err := sk.GetService(ctx, plan.ServiceID)
		if err != nil {
			panic(err)
		}
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "service not found"), nil, nil
		}

		adminAddress, err := sdk.AccAddressFromBech32(service.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "invalid admin address"), nil, nil
		}

		// Get a random rewards plan amount
		senderSpendableBalance := bk.SpendableCoins(ctx, adminAddress)
		amount, err := simtypes.RandomFees(r, ctx, senderSpendableBalance)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), err.Error()), nil, nil
		}

		// Compute some random start/end time
		rewardsStart := ctx.BlockTime().Add(time.Hour * time.Duration(r.Intn(10)+1))
		rewardsEnd := rewardsStart.Add(time.Hour * time.Duration(r.Intn(96)+1))

		signer, found := simtesting.GetSimAccount(adminAddress, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "admin account not found"), nil, nil
		}

		msg = types.NewMsgEditRewardsPlan(
			plan.ID,
			simtypes.RandStringOfLength(r, 32),
			amount,
			rewardsStart,
			rewardsEnd,
			RandomDistribution(r, restakingtypes.DELEGATION_TYPE_POOL, pool),
			RandomDistribution(r, restakingtypes.DELEGATION_TYPE_OPERATOR, operator),
			RandomUsersDistribution(r),
			service.Admin,
		)
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateMsgSetWithdrawAddress(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgSetWithdrawAddress{}

		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		sender, _ := simtypes.RandomAcc(r, accs)
		msg = types.NewMsgSetWithdrawAddress(sender.Address.String(), sender.Address.String())
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, sender)
	}
}

func SimulateMsgWithdrawDelegatorReward(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgWithdrawDelegatorReward{}

		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}
		sender, _ := simtypes.RandomAcc(r, accs)

		// Get the user's pending rewards
		queryServer := keeper.NewQueryServer(k)
		res, err := queryServer.DelegatorTotalRewards(ctx, &types.QueryDelegatorTotalRewardsRequest{
			DelegatorAddress: sender.Address.String(),
		})
		if err != nil {
			panic(err)
		}

		// Get delegation reward so that we can withdraw the rewards
		delegatorRewards, found := utils.Find(res.Rewards, func(d types.DelegationDelegatorReward) bool {
			return !d.Reward.IsEmpty()
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no rewards"), nil, nil
		}

		msg = types.NewMsgWithdrawDelegatorReward(
			delegatorRewards.DelegationType,
			delegatorRewards.DelegationTargetID,
			sender.Address.String(),
		)
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, sender)
	}
}

func SimulateMsgWithdrawOperatorCommission(
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	ok *operatorskeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgWithdrawOperatorCommission{}

		operator, found := operatorssimulation.GetRandomExistingOperator(r, ctx, ok, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "operator not found"), nil, nil
		}

		adminAddress, err := sdk.AccAddressFromBech32(operator.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "invalid admin address"), nil, nil
		}

		signer, found := simtesting.GetSimAccount(adminAddress, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "admin account not found"), nil, nil
		}

		commission, err := k.GetOperatorAccumulatedCommission(ctx, operator.ID)
		if err != nil {
			panic(err)
		}

		if commission.Commissions.IsEmpty() {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no commissions to withdraw"), nil, nil
		}

		msg = types.NewMsgWithdrawOperatorCommission(operator.ID, operator.Admin)
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}
