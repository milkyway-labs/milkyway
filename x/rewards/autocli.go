package rewards

import (
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"github.com/cosmos/cosmos-sdk/version"

	rewardsv1 "github.com/milkyway-labs/milkyway/api/milkyway/rewards/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: rewardsv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the current rewards parameters.",
				},
				{
					RpcMethod: "RewardsPlans",
					Use:       "plans",
					Short:     "Query all rewards plans",
					Example:   fmt.Sprintf(`$ %s query milkyway plans`, version.AppName),
				},
				{
					RpcMethod: "RewardsPlan",
					Use:       "plan [plan-id]",
					Short:     "Query a specific rewards plan",
					Example:   fmt.Sprintf(`$ %s query milkyway plan [plan-id]`, version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "plan_id"},
					},
				},
				{
					RpcMethod: "PoolOutstandingRewards",
					Use:       "pool-outstanding-rewards [pool-id]",
					Short:     "Query outstanding (un-withdrawn) rewards for a pool and all their delegations",
					Example:   fmt.Sprintf("$ %s query milkyway pool-outstanding-rewards [pool-id]", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "pool_id"},
					},
				},
				{
					RpcMethod: "OperatorOutstandingRewards",
					Use:       "operator-outstanding-rewards [operator-id]",
					Short:     "Query outstanding (un-withdrawn) rewards for a operator and all their delegations",
					Example:   fmt.Sprintf("$ %s query milkyway operator-outstanding-rewards [operator-id]", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "operator_id"},
					},
				},
				{
					RpcMethod: "ServiceOutstandingRewards",
					Use:       "service-outstanding-rewards [service-id]",
					Short:     "Query outstanding (un-withdrawn) rewards for a service and all their delegations",
					Example:   fmt.Sprintf("$ %s query milkyway service-outstanding-rewards [service-id]", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "service_id"},
					},
				},
				{
					RpcMethod: "OperatorCommission",
					Use:       "operator-commission [operator-id]",
					Short:     "Query operator commission",
					Example:   fmt.Sprintf("$ %s query milkyway operator-commission [operator-id]", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "operator_id"},
					},
				},
				{
					RpcMethod: "PoolDelegationRewards",
					Use:       "pool-rewards [delegator-address] [pool-id]",
					Short:     "Query all delegation rewards from a particular pool",
					Example:   fmt.Sprintf("$ %s query milkyway pool-rewards [delegator-address] [pool-id]", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "delegator_address"},
						{ProtoField: "pool_id"},
					},
				},
				{
					RpcMethod: "OperatorDelegationRewards",
					Use:       "operator-rewards [delegator-address] [operator-id]",
					Short:     "Query all delegation rewards from a particular operator",
					Example:   fmt.Sprintf("$ %s query milkyway operator-rewards [delegator-address] [operator-id]", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "delegator_address"},
						{ProtoField: "operator_id"},
					},
				},
				{
					RpcMethod: "ServiceDelegationRewards",
					Use:       "service-rewards [delegator-address] [service-id]",
					Short:     "Query all delegation rewards from a particular service",
					Example:   fmt.Sprintf("$ %s query milkyway service-rewards [delegator-address] [service-id]", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "delegator_address"},
						{ProtoField: "service_id"},
					},
				},
				{
					RpcMethod: "DelegationTotalRewards",
					Use:       "rewards [delegator-address]",
					Short:     "Query all delegator rewards",
					Example:   fmt.Sprintf("$ %s query milkyway rewards [delegator-address]", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "delegator_address"},
					},
				},
				{
					RpcMethod: "DelegatorWithdrawAddress",
					Use:       "delegator-withdraw-addr [delegator-address]",
					Short:     "Query delegator withdraw address",
					Example:   fmt.Sprintf("$ %s query milkyway delegator-withdraw-addr [delegator-address]", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "delegator_address"},
					},
				},
			},
			EnhanceCustomCommand: true,
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: rewardsv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "SetWithdrawAddress",
					Use:       "set-withdraw-addr [withdraw-address]",
					Short:     "Change the default withdraw address for rewards associated with an address",
					Example:   fmt.Sprintf("%s tx milkyway set-withdraw-addr [withdraw-address] --from mykey", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "withdraw_address"},
					},
				},
				{
					RpcMethod: "WithdrawDelegatorReward",
					Use:       "withdraw-rewards [delegation-type] [target-id]",
					Short:     "Withdraw rewards from a given delegation target",
					Example:   fmt.Sprintf("%s tx milkyway withdraw-rewards pool [pool-id] --from mykey", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "delegation_type"},
						{ProtoField: "target_id"},
					},
				},
				{
					RpcMethod: "WithdrawOperatorCommission",
					Use:       "withdraw-operator-commission [operator-id]",
					Short:     "Withdraw commissions from a operator address (must be the operator admin)",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "operator_id"},
					},
				},
				{
					RpcMethod: "UpdateParams",
					Skip:      true,
				},
			},
			EnhanceCustomCommand: true,
		},
	}
}