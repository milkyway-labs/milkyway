package cli

import (
	"fmt"

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

// GetTxCmd returns a new command to perform services transactions
func GetTxCmd() *cobra.Command {
	subspacesTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Operators transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	subspacesTxCmd.AddCommand(
		GetCmdRegisterOperator(),
		GetCmdEditOperator(),
		GetCmdDeregisterOperator(),
	)

	return subspacesTxCmd
}

// GetCmdRegisterOperator returns the command allowing to register a new operator
func GetCmdRegisterOperator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [name]",
		Args:  cobra.ExactArgs(1),
		Short: "Register a new service",
		Long: `Register a new service having the given name. 

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

	cmd.Flags().String(flagWebsite, "", "The website of the service (optional)")
	cmd.Flags().String(flagPicture, "", "The picture URL of the service (optional)")

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

	cmd.Flags().String(flagMoniker, types.DoNotModify, "The new moniker of the service (optional)")
	cmd.Flags().String(flagWebsite, types.DoNotModify, "The new website of the service (optional)")
	cmd.Flags().String(flagPicture, types.DoNotModify, "The new picture URL of the service (optional)")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdDeregisterOperator returns the command allowing to deactivate an existing service
func GetCmdDeregisterOperator() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deregister [id]",
		Args:    cobra.ExactArgs(1),
		Short:   "Deregister an existing service",
		Example: fmt.Sprintf(`%s tx %s deregister 1 --from alice`, version.AppName, types.ModuleName),
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
