package cli

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/milkyway-labs/milkyway/x/services/types"
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
	)

	return txCmd
}

// GetCmdSetServiceParams returns the command allowing to transfer the
// ownership of a service
func GetCmdSetServiceParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-service-params [id] [slash-fraction]",
		Args:  cobra.ExactArgs(2),
		Short: "sets the parameters of the service with the given id",
		Example: fmt.Sprintf(
			`%s tx %s set-service-params 1 0.02 --from alice`,
			version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			serviceID, err := types.ParseServiceID(args[0])
			if err != nil {
				return err
			}

			slashFraction, err := math.LegacyNewDecFromStr(args[1])
			if err != nil {
				return fmt.Errorf("parse slash fraction: %w", err)
			}

			params := types.NewServiceParams(slashFraction)

			// Create and validate the message
			msg := types.NewMsgSetServiceParams(serviceID, params, clientCtx.FromAddress.String())
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
