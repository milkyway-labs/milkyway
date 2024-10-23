package cli

import (
	"fmt"
	"strings"

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
		GetCmdExecuteMessages(),
		// The other commands are generated trough the auto cli
	)

	return txCmd
}

// GetCmdSetOperatorParams returns the command allowing to edit an existing operator
func GetCmdSetOperatorParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-operator-params [id] [commission-rate]",
		Args:  cobra.ExactArgs(2),
		Short: "Sets the parameters of the operator with the given id",
		Example: fmt.Sprintf(
			`%s tx %s set-operator-params 1 0.2 --from alice`,
			version.AppName, types.ModuleName),
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
			msg := types.NewMsgSetOperatorParams(creator, id, types.NewOperatorParams(commissionRete))
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdExecuteMessages returns the command allowing to execute messages as
// an operator, by the admin of the operator.
func GetCmdExecuteMessages() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute-messages [operator-id] [path/to/messages.json]",
		Short: "execute messages as an operator",
		Long: strings.TrimSpace(
			fmt.Sprintf(
				`execute messages as an operator.
They should be defined in a JSON file.

Example:
$ %s tx operators execute-messages 1 path/to/proposal.json

Where proposal.json contains:

{
  // array of proto-JSON-encoded sdk.Msgs
  "messages": [
    {
      "@type": "/cosmos.bank.v1beta1.MsgSend",
      "from_address": "init1...",
      "to_address": "init11...",
      "amount":[{"denom": "umilk","amount": "10"}]
    }
  ],
}
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := types.ParseOperatorID(args[0])
			if err != nil {
				return err
			}

			msgs, err := parseExecuteMessages(clientCtx.Codec, args[1])
			if err != nil {
				return err
			}

			sender := clientCtx.FromAddress.String()
			msg, err := types.NewMsgExecuteMessages(id, msgs, sender)
			if err != nil {
				return fmt.Errorf("invalid message: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
