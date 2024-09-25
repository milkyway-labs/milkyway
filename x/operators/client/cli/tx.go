package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

const (
	flagMoniker = "moniker"
	flagWebsite = "website"
	flagPicture = "picture"
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
		GetCmdRegisterOperator(),
		GetCmdEditOperator(),
		GetCmdDeactivateOperator(),
		GetCmdExecuteMessages(),
		GetCmdTransferOperatorOwnership(),
	)

	return txCmd
}

// GetCmdRegisterOperator returns the command allowing to register a new operator
func GetCmdRegisterOperator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [name]",
		Args:  cobra.ExactArgs(1),
		Short: "Register a new operator",
		Long: `Register a new operator having the given name. 

You can specify a website and a picture URL using the optional flags.
The operator will be created with the sender as the admin.`,
		Example: fmt.Sprintf(
			`%s tx %s create MilkyWay --description "MilkyWay Operator" --website https://milkyway.zone --from alice`,
			version.AppName, types.ModuleName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			moniker := args[0]
			creator := clientCtx.FromAddress.String()

			// Get optional data
			website, err := cmd.Flags().GetString(flagWebsite)
			if err != nil {
				return err
			}

			picture, err := cmd.Flags().GetString(flagPicture)
			if err != nil {
				return err
			}

			// Create and validate the message
			msg := types.NewMsgRegisterOperator(moniker, website, picture, creator)
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagWebsite, "", "The website of the operator (optional)")
	cmd.Flags().String(flagPicture, "", "The picture URL of the operator (optional)")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdEditOperator returns the command allowing to edit an existing operator
func GetCmdEditOperator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Edit an existing operator",
		Long: `Edit an existing operator having the provided it. 

You can specify the moniker, website and picture URL using the optional flags.
Only the fields that you provide will be updated`,
		Example: fmt.Sprintf(
			`%s tx %s update 1 --website https://example.com --from alice`,
			version.AppName, types.ModuleName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := types.ParseOperatorID(args[0])
			if err != nil {
				return err
			}

			creator := clientCtx.FromAddress.String()

			// Get new fields values
			moniker, err := cmd.Flags().GetString(flagMoniker)
			if err != nil {
				return err
			}

			website, err := cmd.Flags().GetString(flagWebsite)
			if err != nil {
				return err
			}

			picture, err := cmd.Flags().GetString(flagPicture)
			if err != nil {
				return err
			}

			// Create and validate the message
			msg := types.NewMsgUpdateOperator(id, moniker, website, picture, creator)
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagMoniker, types.DoNotModify, "The new moniker of the operator (optional)")
	cmd.Flags().String(flagWebsite, types.DoNotModify, "The new website of the operator (optional)")
	cmd.Flags().String(flagPicture, types.DoNotModify, "The new picture URL of the operator (optional)")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdDeactivateOperator returns the command allowing to deactivate an existing operator
func GetCmdDeactivateOperator() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deactivate [id]",
		Args:    cobra.ExactArgs(1),
		Short:   "deactivate an existing operator",
		Example: fmt.Sprintf(`%s tx %s deactivate 1 --from alice`, version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := types.ParseOperatorID(args[0])
			if err != nil {
				return err
			}

			creator := clientCtx.FromAddress.String()

			// Create and validate the message
			msg := types.NewMsgDeactivateOperator(id, creator)
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

			msgs, err := parseExecuteMessages(clientCtx.Codec, args[0])
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

// GetCmdTransferOperatorOwnership returns the command allowing to transfer the
// ownership of an operator
func GetCmdTransferOperatorOwnership() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-ownership [id] [new-admin]",
		Args:  cobra.ExactArgs(2),
		Short: "transfer the ownership of an operator",
		Example: fmt.Sprintf(
			`%s tx %s transfer-ownership 1 cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4 --from alice`,
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

			newAdmin := args[1]

			sender := clientCtx.FromAddress.String()

			// Create and validate the message
			msg := types.NewMsgTransferOperatorOwnership(id, newAdmin, sender)
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
