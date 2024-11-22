package cli

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

// GetTxCmd returns a new command to perform operators transactions
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Operators transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdSetOperatorParams(),
		// The other commands are generated through the auto CLI
	)

	return txCmd
}

// GetCmdSetOperatorParams returns the command allowing to set an existing operator's
// parameters
func GetCmdSetOperatorParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-params [operator-id] [commission-rate]",
		Args:    cobra.ExactArgs(2),
		Short:   "Set the parameters of the operator with the given id",
		Example: fmt.Sprintf(`%s tx %s set-params 1 0.2 --from alice`, version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := types.ParseOperatorID(args[0])
			if err != nil {
				return err
			}

			commissionRete, err := math.LegacyNewDecFromStr(args[1])
			if err != nil {
				return fmt.Errorf("invalid commission rate: %s", err)
			}

			creator := clientCtx.FromAddress.String()

			// Create and validate the message
			msg := types.NewMsgSetOperatorParams(id, types.NewOperatorParams(commissionRete), creator)
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
