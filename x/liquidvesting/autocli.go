package liquidvesting

import (
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"github.com/cosmos/cosmos-sdk/version"

	liquidvestingv1 "github.com/milkyway-labs/milkyway/api/milkyway/liquidvesting/v1"
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
					Use:       "assets",
					Short:     "Query all the assets in the insurance fund",
					Example:   fmt.Sprintf("$ %s query insurancefund assets", version.AppName),
				},
				{
					RpcMethod: "UserInsuranceFund",
					Use:       "user-assets [user-address]",
					Short:     "Query the assets deposited in the insurance fund by an user",
					Example:   fmt.Sprintf(`$ %s query liquidvesting user-assets init1....`, version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "user_address"},
					},
				},
				{
					RpcMethod: "UserRestakableAssets",
					Use:       "user-restakable-assets [user-address]",
					Short:     "Query the user's assets that are coverd by their insurance fund and that can be restaked",
					Example:   fmt.Sprintf(`$ %s query assets asset umilk`, version.AppName),
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
					RpcMethod: "MintVestedRepresentation",
					Use:       "mint-vested-representation [sender] [receiver] [amount]",
					Short:     "Mint an user's staked vested tokens representation",
					Example:   fmt.Sprintf(`$ %s tx liquidvesting mint-vested-representation init1... init1... 1000umilk`, version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "sender"},
						{ProtoField: "receiver"},
						{ProtoField: "amount"},
					},
				},
				{
					RpcMethod: "BurnVestedRepresentation",
					Use:       "burn-vested-representation [sender] [user] [amount]",
					Short:     "Burns an user's staked vested tokens representation",
					Example:   fmt.Sprintf(`$ %s tx liquidvesting burn-vested-representation init1... init1... 1000vestd/umilk`, version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "sender"},
						{ProtoField: "user"},
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
