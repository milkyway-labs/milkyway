package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/milkyway-labs/milkyway/x/operators/types"
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
		getCmdQueryOperator(),
		getCmdQueryOperators(),
		getCmdQueryParams(),
	)

	return queryCmd
}

// getCmdQueryOperator returns the command allowing to query an operator
func getCmdQueryOperator() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "operator [operator-id]",
		Short:   "Query the operator with the given id",
		Example: fmt.Sprintf(`%s query %s operator 1`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			operatorID, err := types.ParseOperatorID(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.Operator(context.Background(), types.NewQueryOperatorRequest(operatorID))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// getCmdQueryOperators returns the command allowing to query operators
func getCmdQueryOperators() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "operators",
		Short:   "Query the operators",
		Example: fmt.Sprintf(`%s query %s operators --page=2 --limit=100`, version.AppName, types.ModuleName),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Operators(context.Background(), types.NewQueryOperatorsRequest(pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "operators")

	return cmd
}

// GetCmdQueryParams returns the command to query the module params
func getCmdQueryParams() *cobra.Command {
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

			res, err := queryClient.Params(context.Background(), types.NewQueryParamsRequest())
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
