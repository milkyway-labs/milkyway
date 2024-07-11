package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

const (
	flagName        = "name"
	flagDescription = "description"
	flagWebsite     = "website"
	flagPicture     = "picture"
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
		GetCmdCreateService(),
		GetCmdUpdateService(),
		GetCmdActivateService(),
		GetCmdDeactivateService(),
	)

	return txCmd
}

// GetCmdCreateService returns the command allowing to create a new service
func GetCmdCreateService() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name]",
		Args:  cobra.ExactArgs(1),
		Short: "Create a new service",
		Long: `Create a new service with the given name. 

You can specify a description, website and a picture URL using the optional flags.
The service will be created with the sender as the owner.`,
		Example: fmt.Sprintf(
			`%s tx %s create MilkyWay --description "MilkyWay AVS" --website https://milkyway.zone --from alice`,
			version.AppName, types.ModuleName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			name := args[0]
			creator := clientCtx.FromAddress.String()

			// Get optional data
			description, err := cmd.Flags().GetString(flagDescription)
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
			msg := types.NewMsgCreateService(name, description, website, picture, creator)
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagDescription, "", "The description of the service (optional)")
	cmd.Flags().String(flagWebsite, "", "The website of the service (optional)")
	cmd.Flags().String(flagPicture, "", "The picture URL of the service (optional)")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdUpdateService returns the command allowing to update an existing service
func GetCmdUpdateService() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Update an existing service",
		Long: `Update an existing service having the provided it. 

You can specify a description, website and a picture URL using the optional flags.
Only the fields that you provide will be updated`,
		Example: fmt.Sprintf(
			`%s tx %s update 1 --description "My new description" --from alice`,
			version.AppName, types.ModuleName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := types.ParseServiceID(args[0])
			if err != nil {
				return err
			}

			creator := clientCtx.FromAddress.String()

			// Get new fields values
			name, err := cmd.Flags().GetString(flagName)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(flagDescription)
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
			msg := types.NewMsgUpdateService(id, name, description, website, picture, creator)
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagName, types.DoNotModify, "The new name of the service (optional)")
	cmd.Flags().String(flagDescription, types.DoNotModify, "The new description of the service (optional)")
	cmd.Flags().String(flagWebsite, types.DoNotModify, "The new website of the service (optional)")
	cmd.Flags().String(flagPicture, types.DoNotModify, "The new picture URL of the service (optional)")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdActivateService returns the command allowing to activate an existing service
func GetCmdActivateService() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "activate [id]",
		Args:    cobra.ExactArgs(1),
		Short:   "Activate an existing service",
		Example: fmt.Sprintf(`%s tx %s activate 1 --from alice`, version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := types.ParseServiceID(args[0])
			if err != nil {
				return err
			}

			creator := clientCtx.FromAddress.String()

			// Create and validate the message
			msg := types.NewMsgActivateService(id, creator)
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdDeactivateService returns the command allowing to deactivate an existing service
func GetCmdDeactivateService() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deactivate [id]",
		Args:    cobra.ExactArgs(1),
		Short:   "Deactivate an existing service",
		Example: fmt.Sprintf(`%s tx %s deactivate 1 --from alice`, version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := types.ParseServiceID(args[0])
			if err != nil {
				return err
			}

			creator := clientCtx.FromAddress.String()

			// Create and validate the message
			msg := types.NewMsgDeactivateService(id, creator)
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
