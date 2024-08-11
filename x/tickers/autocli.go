package tickers

import (
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"github.com/cosmos/cosmos-sdk/version"

	tickersv1 "github.com/milkyway-labs/milkyway/api/milkyway/tickers/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: tickersv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the current tickers parameters",
				},
				{
					RpcMethod: "Assets",
					Use:       "assets",
					Short:     "Query all assets",
					Example:   fmt.Sprintf("$ %s query tickers assets --ticker MILK", version.AppName),
				},
				{
					RpcMethod: "Asset",
					Use:       "asset [denom]",
					Short:     "Query a specific asset by its denomination",
					Example:   fmt.Sprintf(`$ %s query tickers asset umilk`, version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "denom"},
					},
				},
			},
			EnhanceCustomCommand: true,
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: tickersv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true,
				},
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
