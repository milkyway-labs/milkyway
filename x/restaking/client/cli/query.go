package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

// GetQueryCmd returns the command allowing to perform queries
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetPoolsQueryCmd(),
		GetOperatorsQueryCmd(),
		GetServicesQueryCmd(),
		GetDelegatorQueryCmd(),
		GetParamsQueryCmd(),
	)

	return queryCmd
}

// --------------------------------------------------------------------------------------------------------------------

// GetPoolsQueryCmd returns the command allowing to query pools
func GetPoolsQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "pool",
		Short: "Querying commands for a restaking pool",
	}

	queryCmd.AddCommand(
		getPoolDelegationsQueryCmd(),
		getPoolDelegationQueryCmd(),
		getPoolUnbondingDelegationsQueryCmd(),
		getPoolUnbondingDelegationQueryCmd(),
	)

	return queryCmd
}

// getPoolDelegationsQueryCmd returns the command allowing to query delegations of a pool
func getPoolDelegationsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delegations [pool-id]",
		Short:   "Query delegations of a pool",
		Example: fmt.Sprintf(`%s query %s pool delegations 1 --page=2 --limit=100`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			poolID, err := poolstypes.ParsePoolID(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.PoolDelegations(cmd.Context(), types.NewQueryPoolDelegationsRequest(poolID, pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "pool delegations")

	return cmd
}

// getPoolDelegationQueryCmd returns the command allowing to query delegation of a pool
func getPoolDelegationQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delegation [pool-id] [delegator-address]",
		Short:   "Query delegation of a pool",
		Example: fmt.Sprintf(`%s query %s pool delegation 1 init1yu5vratzjspgtd0rnrc0d5a79kkqy0n57rhfyh`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			poolID, err := poolstypes.ParsePoolID(args[0])
			if err != nil {
				return err
			}

			delegatorAddress := args[1]

			res, err := queryClient.PoolDelegation(cmd.Context(), types.NewQueryPoolDelegationRequest(poolID, delegatorAddress))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// getPoolUnbondingDelegationsQueryCmd returns the command allowing to query the unbonding delegations of a pool
func getPoolUnbondingDelegationsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unbonding-delegations [pool-id]",
		Short:   "Query the unbonding delegations of a pool",
		Example: fmt.Sprintf(`%s query %s pool unbonding-delegations 1 --page=2 --limit=100`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			poolID, err := poolstypes.ParsePoolID(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.PoolUnbondingDelegations(cmd.Context(), types.NewQueryPoolUnbondingDelegationsRequest(poolID, pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "pool unbonding delegations")

	return cmd
}

// getPoolUnbondingDelegationQueryCmd returns the command allowing to query an unbonding delegation of a pool
func getPoolUnbondingDelegationQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unbonding-delegation [pool-id] [delegator-address]",
		Short:   "Query the unbonding delegation of a pool",
		Example: fmt.Sprintf(`%s query %s pool unbonding-delegation 1 init1yu5vratzjspgtd0rnrc0d5a79kkqy0n57rhfyh`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			poolID, err := poolstypes.ParsePoolID(args[0])
			if err != nil {
				return err
			}

			delegatorAddress := args[1]

			res, err := queryClient.PoolUnbondingDelegation(cmd.Context(), types.NewQueryPoolUnbondingDelegationRequest(poolID, delegatorAddress))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// --------------------------------------------------------------------------------------------------------------------

// GetOperatorsQueryCmd returns the command allowing to query operators
func GetOperatorsQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "operator",
		Short: "Querying commands for a restaking operator",
	}

	queryCmd.AddCommand(
		getOperatorJoinedServices(),
		getOperatorDelegationsQueryCmd(),
		getOperatorDelegationQueryCmd(),
		getOperatorUnbondingDelegationsQueryCmd(),
		getOperatorUnbondingDelegationQueryCmd(),
	)

	return queryCmd
}

// getOperatorJoinedServices returns the command allowing to query an
// operator's joined services.
func getOperatorJoinedServices() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "joined-services [operator-id]",
		Short:   "Query an operator's joined services",
		Example: fmt.Sprintf(`%s query %s operator joined-services 1`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			operatorID, err := operatorstypes.ParseOperatorID(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.OperatorJoinedServices(cmd.Context(), types.NewQueryOperatorJoinedServicesRequest(operatorID))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// getOperatorDelegationsQueryCmd returns the command allowing to query an operator delegations
func getOperatorDelegationsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delegations [operator-id]",
		Short:   "Query delegations of an operator",
		Example: fmt.Sprintf(`%s query %s operator delegations 1 --page=2 --limit=100`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			operatorID, err := operatorstypes.ParseOperatorID(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.OperatorDelegations(cmd.Context(), types.NewQueryOperatorDelegationsRequest(operatorID, pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "operator delegations")

	return cmd
}

// getOperatorDelegationQueryCmd returns the command allowing to query delegation of an operator
func getOperatorDelegationQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delegation [operator-id] [delegator-address]",
		Short:   "Query delegation of an operator",
		Example: fmt.Sprintf(`%s query %s operator delegation 1 init1yu5vratzjspgtd0rnrc0d5a79kkqy0n57rhfyh`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			operatorID, err := operatorstypes.ParseOperatorID(args[0])
			if err != nil {
				return err
			}

			delegatorAddress := args[1]

			res, err := queryClient.OperatorDelegation(cmd.Context(), types.NewQueryOperatorDelegationRequest(operatorID, delegatorAddress))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// getOperatorUnbondingDelegationsQueryCmd returns the command allowing to query an operator unbonding delegations
func getOperatorUnbondingDelegationsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unbonding-delegations [operator-id]",
		Short:   "Query the unbonding delegations of an operator",
		Example: fmt.Sprintf(`%s query %s operator unbonding-delegations 1 --page=2 --limit=100`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			operatorID, err := operatorstypes.ParseOperatorID(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.OperatorUnbondingDelegations(cmd.Context(), types.NewQueryOperatorUnbondingDelegationsRequest(operatorID, pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "operator unbonding delegations")

	return cmd
}

// getOperatorUnbondingDelegationQueryCmd returns the command allowing to query an unbonding delegation of an operator
func getOperatorUnbondingDelegationQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unbonding-delegation [operator-id] [delegator-address]",
		Short:   "Query an unbonding delegation of an operator",
		Example: fmt.Sprintf(`%s query %s operator unbonding-delegation 1 init1yu5vratzjspgtd0rnrc0d5a79kkqy0n57rhfyh`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			operatorID, err := operatorstypes.ParseOperatorID(args[0])
			if err != nil {
				return err
			}

			delegatorAddress := args[1]

			res, err := queryClient.OperatorUnbondingDelegation(cmd.Context(), types.NewQueryOperatorUnbondingDelegationRequest(operatorID, delegatorAddress))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// --------------------------------------------------------------------------------------------------------------------

// GetServicesQueryCmd returns the command allowing to perform queries
func GetServicesQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "service",
		Short: "Querying commands for a restaking service",
	}

	queryCmd.AddCommand(
		getServiceAllowedOperatorsQueryCmd(),
		getServiceAllowedPoolsQueryCmd(),
		getServiceDelegationsQueryCmd(),
		getServiceDelegationQueryCmd(),
		getServiceUnbondingDelegationsQueryCmd(),
		getServiceUnbondingDelegationQueryCmd(),
	)

	return queryCmd
}

// getServiceAllowedOperatorsQueryCmd returns the command allowing to query
// the list of operators allowed to join a service.
func getServiceAllowedOperatorsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "allowed-operators [service-id]",
		Short:   "Query the list of operators allowed to join a service",
		Example: fmt.Sprintf(`%s query %s service whitelisted-operators 1`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			serviceID, err := servicestypes.ParseServiceID(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.ServiceAllowedOperators(
				cmd.Context(),
				types.NewQueryServiceAllowedOperatorsRequest(serviceID, pageReq),
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "service whitelisted operators")

	return cmd
}

// getServiceAllowedPoolsQueryCmd returns the command allowing to query the list
// of pools from which a service is borrowing its security.
func getServiceAllowedPoolsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "allowed-pools [service-id]",
		Short:   "Query the list of pools from which a service is borrowing its security",
		Example: fmt.Sprintf(`%s query %s service whitelisted-pools 1`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			serviceID, err := servicestypes.ParseServiceID(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.ServiceAllowedPools(
				cmd.Context(),
				types.NewQueryServiceAllowedPoolsRequest(serviceID, pageReq),
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "service whitelisted pools")

	return cmd
}

// getServiceDelegationsQueryCmd returns the command allowing to query delegations of a service
func getServiceDelegationsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delegations [service-id]",
		Short:   "Query delegations of a service",
		Example: fmt.Sprintf(`%s query %s service delegations 1 --page=2 --limit=100`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			serviceID, err := servicestypes.ParseServiceID(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.ServiceDelegations(cmd.Context(), types.NewQueryServiceDelegationsRequest(serviceID, pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "service delegations")

	return cmd
}

// getServiceDelegationQueryCmd returns the command allowing to query delegation of a service
func getServiceDelegationQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delegation [service-id] [delegator-address]",
		Short:   "Query delegation of a service",
		Example: fmt.Sprintf(`%s query %s service delegation 1 init1yu5vratzjspgtd0rnrc0d5a79kkqy0n57rhfyh`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			serviceID, err := servicestypes.ParseServiceID(args[0])
			if err != nil {
				return err
			}

			delegatorAddress := args[1]

			res, err := queryClient.ServiceDelegation(cmd.Context(), types.NewQueryServiceDelegationRequest(serviceID, delegatorAddress))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// getServiceUnbondingDelegationsQueryCmd returns the command allowing to query the unbonding delegations of a service
func getServiceUnbondingDelegationsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unbonding-delegations [service-id]",
		Short:   "Query the unbonding delegations of a service",
		Example: fmt.Sprintf(`%s query %s service unbonding-delegations 1 --page=2 --limit=100`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			serviceID, err := servicestypes.ParseServiceID(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.ServiceUnbondingDelegations(cmd.Context(), types.NewQueryServiceUnbondingDelegationsRequest(serviceID, pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "service delegations")

	return cmd
}

// getServiceUnbondingDelegationQueryCmd returns the command allowing to query an unbonding delegation of a service
func getServiceUnbondingDelegationQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unbonding-delegation [service-id] [delegator-address]",
		Short:   "Query an unbonding delegation of a service",
		Example: fmt.Sprintf(`%s query %s service unbonding-delegation 1 init1yu5vratzjspgtd0rnrc0d5a79kkqy0n57rhfyh`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			serviceID, err := servicestypes.ParseServiceID(args[0])
			if err != nil {
				return err
			}

			delegatorAddress := args[1]

			res, err := queryClient.ServiceUnbondingDelegation(cmd.Context(), types.NewQueryServiceUnbondingDelegationRequest(serviceID, delegatorAddress))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// --------------------------------------------------------------------------------------------------------------------

// GetDelegatorQueryCmd returns the command allowing to perform queries for a delegator
func GetDelegatorQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "delegator",
		Short: "Querying commands for restaking delegator",
	}

	queryCmd.AddCommand(
		getDelegatorPoolDelegationsQueryCmd(),
		getDelegatorPoolUnbondingDelegationsQueryCmd(),
		getDelegatorPoolsQueryCmd(),
		getDelegatorPoolQueryCmd(),
		getDelegatorOperatorDelegationsQueryCmd(),
		getDelegatorOperatorUnbondingDelegationsQueryCmd(),
		getDelegatorOperatorsQueryCmd(),
		getDelegatorOperatorQueryCmd(),
		getDelegatorServiceDelegationsQueryCmd(),
		getDelegatorServiceUnbondingDelegationsQueryCmd(),
		getDelegatorServicesQueryCmd(),
		getDelegatorServiceQueryCmd(),
	)

	return queryCmd
}

// getDelegatorPoolDelegationsQueryCmd returns the command allowing to query all pools delegations of a delegator
func getDelegatorPoolDelegationsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pools-delegations [delegator-address]",
		Short:   "Query all pools delegations of a delegator",
		Example: fmt.Sprintf(`%s query %s delegator pools-delegations init1yu5vratzjspgtd0rnrc0d5a79kkqy0n57rhfyh --page=2 --limit=100`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delegatorAddress := args[0]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.DelegatorPoolDelegations(cmd.Context(), types.NewQueryDelegatorPoolDelegationsRequest(delegatorAddress, pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "delegator pool delegations")

	return cmd
}

// getDelegatorPoolUnbondingDelegationsQueryCmd returns the command allowing to query all pools unbonding delegations of a delegator
func getDelegatorPoolUnbondingDelegationsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pools-unbonding-delegations [delegator-address]",
		Short:   "Query all pools unbonding delegations of a delegator",
		Example: fmt.Sprintf(`%s query %s delegator pools-unbonding-delegations init1yu5vratzjspgtd0rnrc0d5a79kkqy0n57rhfyh --page=2 --limit=100`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delegatorAddress := args[0]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.DelegatorPoolUnbondingDelegations(cmd.Context(), types.NewQueryDelegatorPoolUnbondingDelegationsRequest(delegatorAddress, pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "delegator pool unbonding delegations")

	return cmd
}

// getDelegatorPoolsQueryCmd returns the command allowing to query all pools a delegator is participating in
func getDelegatorPoolsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pools [delegator-address]",
		Short:   "Query all pools a delegator is participating in",
		Example: fmt.Sprintf(`%s query %s delegator pools init1yu5vratzjspgtd0rnrc0d5a79kkqy0n57rhfyh --page=2 --limit=100`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delegatorAddress := args[0]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.DelegatorPools(cmd.Context(), types.NewQueryDelegatorPoolsRequest(delegatorAddress, pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "delegator pools")

	return cmd
}

// getDelegatorPoolQueryCmd returns the command allowing to query a pool a delegator is participating in
func getDelegatorPoolQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool [delegator-address] [pool-id]",
		Short:   "Query a pool a delegator is participating in",
		Example: fmt.Sprintf(`%s query %s delegator pool init1yu5vratzjspgtd0rnrc0d5a79kkqy0n57rhfyh 1`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delegatorAddress := args[0]
			poolID, err := poolstypes.ParsePoolID(args[1])
			if err != nil {
				return err
			}

			res, err := queryClient.DelegatorPool(cmd.Context(), types.NewQueryDelegatorPoolRequest(delegatorAddress, poolID))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// getDelegatorOperatorDelegationsQueryCmd returns the command allowing to query all operators delegations of a delegator
func getDelegatorOperatorDelegationsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "operators-delegations [delegator-address]",
		Short:   "Query all operators delegations of a delegator",
		Example: fmt.Sprintf(`%s query %s delegator operators-delegations init1yu5vratzjspgtd0rnrc0d5a79kkqy0n57rhfyh --page=2 --limit=100`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delegatorAddress := args[0]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.DelegatorOperatorDelegations(cmd.Context(), types.NewQueryDelegatorOperatorDelegationsRequest(delegatorAddress, pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "delegator operator delegations")

	return cmd
}

// getDelegatorOperatorDelegationsQueryCmd returns the command allowing to query all operators unbonding delegations of a delegator
func getDelegatorOperatorUnbondingDelegationsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "operators-unbonding-delegations [delegator-address]",
		Short:   "Query all operators unbonding delegations of a delegator",
		Example: fmt.Sprintf(`%s query %s delegator operators-unbonding-delegations init1yu5vratzjspgtd0rnrc0d5a79kkqy0n57rhfyh --page=2 --limit=100`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delegatorAddress := args[0]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.DelegatorOperatorUnbondingDelegations(cmd.Context(), types.NewQueryDelegatorOperatorUnbondingDelegationsRequest(delegatorAddress, pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "delegator operator unbonding delegations")

	return cmd
}

// getDelegatorOperatorsQueryCmd returns the command allowing to query all operators a delegator has delegated to
func getDelegatorOperatorsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "operators [delegator-address]",
		Short:   "Query all operators a delegator has delegated to",
		Example: fmt.Sprintf(`%s query %s delegator operators init1yu5vratzjspgtd0rnrc0d5a79kkqy0n57rhfyh --page=2 --limit=100`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delegatorAddress := args[0]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.DelegatorOperators(cmd.Context(), types.NewQueryDelegatorOperatorsRequest(delegatorAddress, pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "delegator operators")

	return cmd
}

// getDelegatorOperatorQueryCmd returns the command allowing to query an operator a delegator has delegated to
func getDelegatorOperatorQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "operator [delegator-address] [operator-id]",
		Short:   "Query an operator a delegator has delegated to",
		Example: fmt.Sprintf(`%s query %s delegator operator init1yu5vratzjspgtd0rnrc0d5a79kkqy0n57rhfyh 1`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delegatorAddress := args[0]
			operatorID, err := operatorstypes.ParseOperatorID(args[1])
			if err != nil {
				return err
			}

			res, err := queryClient.DelegatorOperator(cmd.Context(), types.NewQueryDelegatorOperatorRequest(delegatorAddress, operatorID))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// getDelegatorServiceDelegationsQueryCmd returns the command allowing to query all services delegations of a delegator
func getDelegatorServiceDelegationsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "services-delegations [delegator-address]",
		Short:   "Query all services delegations of a delegator",
		Example: fmt.Sprintf(`%s query %s delegator services-delegations init1yu5vratzjspgtd0rnrc0d5a79kkqy0n57rhfyh --page=2 --limit=100`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delegatorAddress := args[0]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.DelegatorServiceDelegations(cmd.Context(), types.NewQueryDelegatorServiceDelegationsRequest(delegatorAddress, pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "delegator service delegations")

	return cmd
}

// getDelegatorServiceDelegationsQueryCmd returns the command allowing to query all services unbonding delegations of a delegator
func getDelegatorServiceUnbondingDelegationsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "services-unbonding-delegations [delegator-address]",
		Short:   "Query all services unbonding delegations of a delegator",
		Example: fmt.Sprintf(`%s query %s delegator services-unbonding-delegations init1yu5vratzjspgtd0rnrc0d5a79kkqy0n57rhfyh --page=2 --limit=100`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delegatorAddress := args[0]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.DelegatorServiceUnbondingDelegations(cmd.Context(), types.NewQueryDelegatorServiceUnbondingDelegationsRequest(delegatorAddress, pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "delegator service unbonding delegations")

	return cmd
}

// getDelegatorServicesQueryCmd returns the command allowing to query all services a delegator has delegated to
func getDelegatorServicesQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "services [delegator-address]",
		Short:   "Query all services a delegator has delegated to",
		Example: fmt.Sprintf(`%s query %s delegator services init1yu5vratzjspgtd0rnrc0d5a79kkqy0n57rhfyh --page=2 --limit=100`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delegatorAddress := args[0]

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.DelegatorServices(cmd.Context(), types.NewQueryDelegatorServicesRequest(delegatorAddress, pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "delegator services")

	return cmd
}

// getDelegatorServiceQueryCmd returns the command allowing to query a service a delegator has delegated to
func getDelegatorServiceQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service [delegator-address] [service-id]",
		Short:   "Query a service a delegator has delegated to",
		Example: fmt.Sprintf(`%s query %s delegator service init1yu5vratzjspgtd0rnrc0d5a79kkqy0n57rhfyh 1`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delegatorAddress := args[0]
			serviceID, err := servicestypes.ParseServiceID(args[1])
			if err != nil {
				return err
			}

			res, err := queryClient.DelegatorService(cmd.Context(), types.NewQueryDelegatorServiceRequest(delegatorAddress, serviceID))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// --------------------------------------------------------------------------------------------------------------------

// GetParamsQueryCmd returns the command to query the module params
func GetParamsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "params",
		Short:   "Query the module parameters",
		Example: fmt.Sprintf(`%s query %s params`, version.AppName, types.ModuleName),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), types.NewQueryParamsRequest())
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
