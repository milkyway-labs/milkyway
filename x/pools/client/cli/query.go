package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/milkyway-labs/milkyway/v8/x/pools/types"
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
		getCmdQueryPoolByID(),
		getCmdQueryPoolByDenom(),
		getCmdQueryPools(),
	)

	return queryCmd
}

// getCmdQueryPoolByID returns the command allowing to query a pool
func getCmdQueryPoolByID() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool [pool-id]",
		Short:   "Query the pool with the given id",
		Example: fmt.Sprintf(`%s query %s pool 1`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			poolID, err := types.ParsePoolID(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.PoolByID(cmd.Context(), types.NewQueryPoolByIDRequest(poolID))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// getCmdQueryPoolByDenom returns the command allowing to query pool by denom
func getCmdQueryPoolByDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool-by-denom [denom]",
		Short:   "Query the pool associated with the given denom",
		Example: fmt.Sprintf(`%s query %s pool umilk`, version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.PoolByDenom(cmd.Context(), types.NewQueryPoolByDenomRequest(args[0]))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// getCmdQueryPools returns the command to query the stored pools
func getCmdQueryPools() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pools",
		Short:   "Query the pools",
		Example: fmt.Sprintf(`%s query %s pools --page=2 --limit=100`, version.AppName, types.ModuleName),
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

			res, err := queryClient.Pools(cmd.Context(), types.NewQueryPoolsRequest(pageReq))
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "pools")

	return cmd
}
