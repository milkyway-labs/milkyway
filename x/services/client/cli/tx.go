package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/milkyway-labs/milkyway/v12/x/services/types"
)

// GetTxCmd returns a new command to perform services transactions
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Services transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdSetServiceParams(),
		// The other commands are generated through the auto CLI
	)

	return txCmd
}

// GetCmdSetServiceParams returns the command allowing to set an existing service's
// parameters
func GetCmdSetServiceParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-params [service-id] [allowed-denoms]",
		Args:    cobra.ExactArgs(2),
		Short:   "Set the parameters of the service with the given id",
		Example: fmt.Sprintf(`%s tx %s set-params 1 utia,umilk --from alice`, version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := types.ParseServiceID(args[0])
			if err != nil {
				return err
			}

			allowedDenoms := strings.Split(args[1], ",")

			// Create the service params
			serviceParams := types.NewServiceParams(allowedDenoms)

			sender := clientCtx.FromAddress.String()

			// Create and validate the message
			msg := types.NewMsgSetServiceParams(id, serviceParams, sender)
			err = msg.ValidateBasic()
			if err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
