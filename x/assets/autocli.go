package assets

import (
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"github.com/cosmos/cosmos-sdk/version"

	assetsv1 "github.com/milkyway-labs/milkyway/v5/api/milkyway/assets/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: assetsv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Assets",
					Use:       "assets",
					Short:     "Query all assets",
					Example:   fmt.Sprintf("$ %s query assets assets --ticker MILK", version.AppName),
				},
				{
					RpcMethod: "Asset",
					Use:       "asset [denom]",
					Short:     "Query a specific asset by its denomination",
					Example:   fmt.Sprintf(`$ %s query assets asset umilk`, version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "denom"},
					},
				},
			},
			EnhanceCustomCommand: true,
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: assetsv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "RegisterAsset",
					Skip:      true,
				},
				{
					RpcMethod: "DeregisterAsset",
					Skip:      true,
				},
			},
			EnhanceCustomCommand: true,
		},
	}
}
