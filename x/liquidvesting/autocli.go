package liquidvesting

import (
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"github.com/cosmos/cosmos-sdk/version"

	liquidvestingv1 "github.com/milkyway-labs/milkyway/v7/api/milkyway/liquidvesting/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (a AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: liquidvestingv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the current liquidvesting parameters",
				},
				{
					RpcMethod: "InsuranceFund",
					Use:       "insurance-fund",
					Short:     "Query all the assets in the insurance fund",
					Example:   fmt.Sprintf("$ %s query liquidvesting insurance-fund", version.AppName),
				},
				{
					RpcMethod: "UserInsuranceFunds",
					Use:       "user-insurance-funds",
					Short:     "Query all the users' insurance fund",
					Example:   fmt.Sprintf("$ %s query liquidvesting user-insurance-funds", version.AppName),
				},
				{
					RpcMethod: "UserInsuranceFund",
					Use:       "user-insurance-fund [user-address]",
					Short:     "Query the assets deposited in the insurance fund by an user",
					Example:   fmt.Sprintf(`$ %s query liquidvesting user-insurance-fund init1....`, version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "user_address"},
					},
				},
				{
					RpcMethod: "UserRestakableAssets",
					Use:       "user-restakable-assets [user-address]",
					Short:     "Query the user's assets that are covered by their insurance fund and that can be restaked",
					Example:   fmt.Sprintf(`$ %s query liquidvesting user-restakable-assets init1...`, version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "user_address"},
					},
				},
			},
			EnhanceCustomCommand: true,
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: liquidvestingv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "MintLockedRepresentation",
					Use:       "mint-locked-representation [sender] [receiver] [amount]",
					Short:     "Mint an user's staked locked tokens representation",
					Example:   fmt.Sprintf(`$ %s tx liquidvesting mint-locked-representation init1... init1... 1000umilk`, version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "sender"},
						{ProtoField: "receiver"},
						{ProtoField: "amount"},
					},
				},
				{
					RpcMethod: "BurnLockedRepresentation",
					Use:       "burn-locked-representation [sender] [user] [amount]",
					Short:     "Burns an user's staked locked tokens representation",
					Example:   fmt.Sprintf(`$ %s tx liquidvesting burn-locked-representation init1... init1... 1000vestd/umilk`, version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "sender"},
						{ProtoField: "user"},
						{ProtoField: "amount"},
					},
				},
				{
					RpcMethod: "WithdrawInsuranceFund",
					Use:       "withdraw-insurance-fund [sender] [amount]",
					Short:     "Withdraws coins from the insurance fund",
					Example:   fmt.Sprintf(`$ %s tx liquidvesting withdraw-insurance-fundh init1... 1000umilk`, version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "sender"},
						{ProtoField: "amount"},
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
