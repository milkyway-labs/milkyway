package investors

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	investorsv1 "github.com/milkyway-labs/milkyway/v11/api/milkyway/investors/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: investorsv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "InvestorsRewardRatio",
					Use:       "investors-reward-ratio",
					Short:     "Query the investors reward ratio parameter",
				},
				{
					RpcMethod: "VestingInvestors",
					Use:       "vesting-investors",
					Short:     "Query all the vesting investors",
				},
			},
			EnhanceCustomCommand: true,
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: investorsv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "AddVestingInvestor",
					Skip:      true,
				},
				{
					RpcMethod: "UpdateInvestorsRewardRatio",
					Skip:      true,
				},
			},
			EnhanceCustomCommand: true,
		},
	}
}
