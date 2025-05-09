package operators

import (
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"github.com/cosmos/cosmos-sdk/version"

	operatorsv1 "github.com/milkyway-labs/milkyway/v12/api/milkyway/operators/v1"
	"github.com/milkyway-labs/milkyway/v12/x/operators/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: operatorsv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Operator",
					Use:       "operator [operator-id]",
					Short:     "Query the operator with the given id",
					Example:   fmt.Sprintf(`%s query %s operator 1`, version.AppName, types.ModuleName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "operator_id"},
					},
				},
				{
					RpcMethod: "Operators",
					Use:       "operators",
					Short:     "Query the operators",
					Example:   fmt.Sprintf(`%s query %s operators --page=2 --limit=100`, version.AppName, types.ModuleName),
				},
				{
					RpcMethod: "OperatorParams",
					Use:       "operator-params [operator-id]",
					Short:     "Query the parameters of the operator with the given id",
					Example:   fmt.Sprintf(`%s query %s operators-params 1`, version.AppName, types.ModuleName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "operator_id"},
					},
				},
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the module parameters",
					Example:   fmt.Sprintf(`%s query %s params`, version.AppName, types.ModuleName),
				},
			},
			EnhanceCustomCommand: true,
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: operatorsv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "RegisterOperator",
					Use:       "register [moniker]",
					Short:     "Register a new operator",
					Long: `Register a new operator having the given name.

You can specify a website and a picture URL using the optional flags.
The operator will be created with the sender as the admin.`,
					Example: fmt.Sprintf(
						`%s tx %s create MilkyWay --website https://milkyway.zone --from alice`,
						version.AppName, types.ModuleName,
					),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "moniker"},
					},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"website":     {Usage: "website URL of the operator"},
						"picture_url": {Usage: "picture URL of the operator"},
					},
				},
				{
					RpcMethod: "UpdateOperator",
					Use:       "update [operator-id]",
					Short:     "Update an existing operator",
					Long: `Update an existing operator having the provided it.

You can specify the moniker, website and picture URL using the optional flags.
Only the fields that you provide will be updated`,
					Example: fmt.Sprintf(
						`%s tx %s update 1 --website https://example.com --from alice`,
						version.AppName, types.ModuleName,
					),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "operator_id"},
					},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"moniker": {
							Usage:        "moniker of the operator",
							DefaultValue: types.DoNotModify,
						},
						"website": {
							Usage:        "website URL of the operator",
							DefaultValue: types.DoNotModify,
						},
						"picture_url": {
							Usage:        "picture URL of the operator",
							DefaultValue: types.DoNotModify,
						},
					},
				},
				{
					RpcMethod: "SetOperatorParams",
					Skip:      true,
				},
				{
					RpcMethod: "TransferOperatorOwnership",
					Use:       "transfer-ownership [operator-id] [new-owner]",
					Short:     "Transfer the ownership of an operator to another user",
					Example: fmt.Sprintf(`%s tx %s transfer-ownership 1 cosmos1example --from alice`,
						version.AppName, types.ModuleName,
					),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "operator_id"},
					},
				},
				{
					RpcMethod: "DeactivateOperator",
					Use:       "deactivate [operator-id]",
					Short:     "Deactivate an existing operator",
					Example:   fmt.Sprintf(`%s tx %s deactivate 1 --from alice`, version.AppName, types.ModuleName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "operator_id"},
					},
				},
				{
					RpcMethod: "DeleteOperator",
					Use:       "delete [operator-id]",
					Short:     "Delete an inactive operator",
					Example:   fmt.Sprintf(`%s tx %s delete 1 --from alice`, version.AppName, types.ModuleName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "operator_id"},
					},
				},
				{
					RpcMethod: "ReactivateOperator",
					Use:       "reactivate [operator-id]",
					Short:     "Reactivate an operator that has been deactivated",
					Example:   fmt.Sprintf(`%s tx %s reactivate 1 --from alice`, version.AppName, types.ModuleName),
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
