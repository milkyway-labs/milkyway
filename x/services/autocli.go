package services

import (
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"github.com/cosmos/cosmos-sdk/version"

	servicesv1 "github.com/milkyway-labs/milkyway/v10/api/milkyway/services/v1"
	"github.com/milkyway-labs/milkyway/v10/x/services/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: servicesv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Service",
					Use:       "service [service-id]",
					Short:     "Query the service with the given id",
					Example:   fmt.Sprintf(`%s query %s service 1`, version.AppName, types.ModuleName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "service_id"},
					},
				},
				{
					RpcMethod: "Services",
					Use:       "services",
					Short:     "Query the services",
					Example:   fmt.Sprintf(`%s query %s services --page-offset-100 --page-limit=100`, version.AppName, types.ModuleName),
				},
				{
					RpcMethod: "ServiceParams",
					Use:       "service-params [service-id]",
					Short:     "Query the parameters of the service with the given id",
					Example:   fmt.Sprintf(`%s query %s service-params 1`, version.AppName, types.ModuleName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "service_id"},
					},
				},
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the current services parameters",
				},
			},
			EnhanceCustomCommand: true,
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: servicesv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "CreateService",
					Use:       "create [name]",
					Short:     "Create a new service",
					Long: `Create a new service with the given name. 

You can specify a description, website and a picture URL using the optional flags.
The service will be created with the sender as the owner.`,
					Example: fmt.Sprintf(
						`%s tx %s create MilkyWay --description "MilkyWay AVS" --website https://milkyway.zone --from alice`,
						version.AppName, types.ModuleName,
					),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "name"},
					},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"description": {Usage: "description of the service"},
						"website":     {Usage: "website URL of the service"},
						"picture_url": {Usage: "picture URL of the service"},
					},
				},
				{
					RpcMethod: "UpdateService",
					Use:       "update [id]",
					Short:     "Update an existing service",
					Long: `Update an existing service having the provided it. 

You can specify a name, description, website and a picture URL using the optional flags.
Only the fields that you provide will be updated`,
					Example: fmt.Sprintf(
						`%s tx %s update 1 --description "My new description" --from alice`,
						version.AppName, types.ModuleName,
					),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "service_id"},
					},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"name": {
							Usage:        "name of the service",
							DefaultValue: types.DoNotModify,
						},
						"description": {
							Usage:        "description of the service",
							DefaultValue: types.DoNotModify,
						},
						"website": {
							Usage:        "website URL of the service",
							DefaultValue: types.DoNotModify,
						},
						"picture_url": {
							Usage:        "picture URL of the service",
							DefaultValue: types.DoNotModify,
						},
					},
				},
				{
					RpcMethod: "ActivateService",
					Use:       "activate [id]",
					Short:     "Activate an existing service",
					Example:   fmt.Sprintf(`%s tx %s activate 1 --from alice`, version.AppName, types.ModuleName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "service_id"},
					},
				},
				{
					RpcMethod: "DeactivateService",
					Use:       "deactivate [id]",
					Short:     "Deactivate an existing service",
					Example:   fmt.Sprintf(`%s tx %s deactivate 1 --from alice`, version.AppName, types.ModuleName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "service_id"},
					},
				},
				{
					RpcMethod: "DeleteService",
					Use:       "delete [id]",
					Short:     "Delete a deactivated service",
					Example:   fmt.Sprintf(`%s tx %s delete 1 --from alice`, version.AppName, types.ModuleName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "service_id"},
					},
				},
				{
					RpcMethod: "TransferServiceOwnership",
					Use:       "transfer-ownership [id] [new-admin]",
					Short:     "Transfer the ownership of a service",
					Example: fmt.Sprintf(
						`%s tx %s transfer-ownership 1 cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4 --from alice`,
						version.AppName, types.ModuleName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "service_id"},
						{ProtoField: "new_admin"},
					},
				},
				{
					RpcMethod: "SetServiceParams",
					Skip:      true,
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
