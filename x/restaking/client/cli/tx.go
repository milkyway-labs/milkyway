package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

// GetTxCmd returns a new command to perform restaking transactions
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Restaking transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetDelegateTxCmd(),
	)

	return txCmd
}

// --------------------------------------------------------------------------------------------------------------------

// GetDelegateTxCmd returns the command allowing to delegate tokens
func GetDelegateTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "delegate",
		Short: "Delegate transactions subcommands",
	}

	txCmd.AddCommand(
		GetDelegateToPoolCmd(),
		GetDelegateToOperatorCmd(),
		GetDelegateToServiceCmd(),
	)

	return txCmd
}

// GetDelegateToPoolCmd returns the command allowing to delegate to a pool
func GetDelegateToPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool [amount]",
		Args:    cobra.ExactArgs(1),
		Short:   "Delegate the given amount to a pool",
		Example: fmt.Sprintf("%s tx %s delegate pool 1000000milk --from alice", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			delegator := clientCtx.FromAddress.String()

			// Create and validate the message
			msg := types.NewMsgDelegatePool(amount, delegator)
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetDelegateToOperatorCmd returns the command allowing to delegate to an operator
func GetDelegateToOperatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "operator [operator-id] [amount]",
		Args:    cobra.ExactArgs(2),
		Short:   "Delegate the given amount to an operator",
		Example: fmt.Sprintf("%s tx %s delegate operator 1 1000000milk --from alice", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			operatorID, err := operatorstypes.ParseOperatorID(args[0])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			delegator := clientCtx.FromAddress.String()

			// Create and validate the message
			msg := types.NewMsgDelegateOperator(operatorID, amount, delegator)
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetDelegateToServiceCmd returns the command allowing to delegate to a service
func GetDelegateToServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service [service-id] [amount]",
		Args:    cobra.ExactArgs(2),
		Short:   "Delegate the given amount to a service",
		Example: fmt.Sprintf("%s tx %s delegate service 1 1000000milk --from alice", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			serviceID, err := servicestypes.ParseServiceID(args[0])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			delegator := clientCtx.FromAddress.String()

			// Create and validate the message
			msg := types.NewMsgDelegateService(serviceID, amount, delegator)
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
