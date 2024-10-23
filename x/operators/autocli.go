package operators

import (
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"github.com/cosmos/cosmos-sdk/version"

	operatorsv1 "github.com/milkyway-labs/milkyway/api/milkyway/operators/v1"
	"github.com/milkyway-labs/milkyway/x/operators/types"
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
						`%s tx %s register MilkyWay --website https://milkyway.zone --from alice`,
						version.AppName, types.ModuleName,
					),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "moniker"},
					},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"website":     {DefaultValue: types.DoNotModify},
						"picture_url": {DefaultValue: types.DoNotModify},
					},
				},
				{
					RpcMethod: "UpdateOperator",
					Use:       "update [id]",
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
						"moniker":     {DefaultValue: types.DoNotModify},
						"website":     {DefaultValue: types.DoNotModify},
						"picture_url": {DefaultValue: types.DoNotModify},
					},
				},
				{
					RpcMethod: "DeactivateOperator",
					Use:       "deactivate [id]",
					Short:     "Deactivate an existing operator",
					Example:   fmt.Sprintf(`%s tx %s deactivate 1 --from alice`, version.AppName, types.ModuleName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "operator_id"},
					},
				},
				{
					RpcMethod: "TransferOperatorOwnership",
					Use:       "transfer-ownership [id] [new-admin]",
					Short:     "Transfer the ownership of an operator",
					Example:   fmt.Sprintf(`%s tx %s transfer-ownership 1 init1... --from alice`, version.AppName, types.ModuleName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "operator_id"},
						{ProtoField: "new_admin"},
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
