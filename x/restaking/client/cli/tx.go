package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
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

	// txCmd.AddCommand(
	// 	GetDelegateTxCmd(),
	// 	GetUnbondTxCmd(),
	// 	GetUpdateTxCmd(),
	// )

	return txCmd
}

// --------------------------------------------------------------------------------------------------------------------

// // GetDelegateTxCmd returns the command allowing to delegate tokens
// func GetDelegateTxCmd() *cobra.Command {
// 	txCmd := &cobra.Command{
// 		Use:   "delegate",
// 		Short: "Delegate transactions subcommands",
// 	}
//
// 	txCmd.AddCommand(
// 		GetDelegateToPoolCmd(),
// 		GetDelegateToOperatorCmd(),
// 		GetDelegateToServiceCmd(),
// 	)
//
// 	return txCmd
// }
//
// // GetDelegateToPoolCmd returns the command allowing to delegate to a pool
// func GetDelegateToPoolCmd() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:     "pool [amount]",
// 		Args:    cobra.ExactArgs(1),
// 		Short:   "Delegate the given amount to a pool",
// 		Example: fmt.Sprintf("%s tx %s delegate pool 1000000milk --from alice", version.AppName, types.ModuleName),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			clientCtx, err := client.GetClientTxContext(cmd)
// 			if err != nil {
// 				return err
// 			}
//
// 			amount, err := sdk.ParseCoinNormalized(args[1])
// 			if err != nil {
// 				return err
// 			}
//
// 			delegator := clientCtx.FromAddress.String()
//
// 			// Create and validate the message
// 			msg := types.NewMsgDelegatePool(amount, delegator)
// 			if err = msg.ValidateBasic(); err != nil {
// 				return fmt.Errorf("message validation failed: %w", err)
// 			}
//
// 			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
// 		},
// 	}
//
// 	flags.AddTxFlagsToCmd(cmd)
//
// 	return cmd
// }
//
// // GetDelegateToOperatorCmd returns the command allowing to delegate to an operator
// func GetDelegateToOperatorCmd() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:     "operator [operator-id] [amount]",
// 		Args:    cobra.ExactArgs(2),
// 		Short:   "Delegate the given amount to an operator",
// 		Example: fmt.Sprintf("%s tx %s delegate operator 1 1000000milk --from alice", version.AppName, types.ModuleName),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			clientCtx, err := client.GetClientTxContext(cmd)
// 			if err != nil {
// 				return err
// 			}
//
// 			operatorID, err := operatorstypes.ParseOperatorID(args[0])
// 			if err != nil {
// 				return err
// 			}
//
// 			amount, err := sdk.ParseCoinsNormalized(args[1])
// 			if err != nil {
// 				return err
// 			}
//
// 			delegator := clientCtx.FromAddress.String()
//
// 			// Create and validate the message
// 			msg := types.NewMsgDelegateOperator(operatorID, amount, delegator)
// 			if err = msg.ValidateBasic(); err != nil {
// 				return fmt.Errorf("message validation failed: %w", err)
// 			}
//
// 			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
// 		},
// 	}
//
// 	flags.AddTxFlagsToCmd(cmd)
//
// 	return cmd
// }
//
// // GetDelegateToServiceCmd returns the command allowing to delegate to a service
// func GetDelegateToServiceCmd() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:     "service [service-id] [amount]",
// 		Args:    cobra.ExactArgs(2),
// 		Short:   "Delegate the given amount to a service",
// 		Example: fmt.Sprintf("%s tx %s delegate service 1 1000000milk --from alice", version.AppName, types.ModuleName),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			clientCtx, err := client.GetClientTxContext(cmd)
// 			if err != nil {
// 				return err
// 			}
//
// 			serviceID, err := servicestypes.ParseServiceID(args[0])
// 			if err != nil {
// 				return err
// 			}
//
// 			amount, err := sdk.ParseCoinsNormalized(args[1])
// 			if err != nil {
// 				return err
// 			}
//
// 			delegator := clientCtx.FromAddress.String()
//
// 			// Create and validate the message
// 			msg := types.NewMsgDelegateService(serviceID, amount, delegator)
// 			if err = msg.ValidateBasic(); err != nil {
// 				return fmt.Errorf("message validation failed: %w", err)
// 			}
//
// 			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
// 		},
// 	}
//
// 	flags.AddTxFlagsToCmd(cmd)
//
// 	return cmd
// }
//
// // --------------------------------------------------------------------------------------------------------------------
//
// // GetUnbondTxCmd returns the command allowing to unbond tokens
// func GetUnbondTxCmd() *cobra.Command {
// 	txCmd := &cobra.Command{
// 		Use:   "unbond",
// 		Short: "Unbond transactions subcommands",
// 	}
//
// 	txCmd.AddCommand(
// 		GetUnbondFromPoolCmd(),
// 		GetUnbondFromOperatorCmd(),
// 		GetUnbondFromServiceCmd(),
// 	)
//
// 	return txCmd
// }
//
// // GetUnbondFromPoolCmd returns the command allowing to unbond from a pool
// func GetUnbondFromPoolCmd() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:     "pool [amount]",
// 		Args:    cobra.ExactArgs(1),
// 		Short:   "Unbond the given amount from a pool",
// 		Example: fmt.Sprintf("%s tx %s unbond pool 1000000umilk --from alice", version.AppName, types.ModuleName),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			clientCtx, err := client.GetClientTxContext(cmd)
// 			if err != nil {
// 				return err
// 			}
//
// 			amount, err := sdk.ParseCoinNormalized(args[1])
// 			if err != nil {
// 				return err
// 			}
//
// 			delegator := clientCtx.FromAddress.String()
//
// 			// Create and validate the message
// 			msg := types.NewMsgUndelegatePool(amount, delegator)
// 			if err = msg.ValidateBasic(); err != nil {
// 				return fmt.Errorf("message validation failed: %w", err)
// 			}
//
// 			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
// 		},
// 	}
//
// 	flags.AddTxFlagsToCmd(cmd)
//
// 	return cmd
// }
//
// // GetUnbondFromOperatorCmd returns the command allowing to unbong from an operator
// func GetUnbondFromOperatorCmd() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:     "operator [operator-id] [amount]",
// 		Args:    cobra.ExactArgs(2),
// 		Short:   "Unbond the given amount from an operator",
// 		Example: fmt.Sprintf("%s tx %s unbond operator 1 1000000milk --from alice", version.AppName, types.ModuleName),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			clientCtx, err := client.GetClientTxContext(cmd)
// 			if err != nil {
// 				return err
// 			}
//
// 			operatorID, err := operatorstypes.ParseOperatorID(args[0])
// 			if err != nil {
// 				return err
// 			}
//
// 			amount, err := sdk.ParseCoinsNormalized(args[1])
// 			if err != nil {
// 				return err
// 			}
//
// 			delegator := clientCtx.FromAddress.String()
//
// 			// Create and validate the message
// 			msg := types.NewMsgUndelegateOperator(operatorID, amount, delegator)
// 			if err = msg.ValidateBasic(); err != nil {
// 				return fmt.Errorf("message validation failed: %w", err)
// 			}
//
// 			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
// 		},
// 	}
//
// 	flags.AddTxFlagsToCmd(cmd)
//
// 	return cmd
// }
//
// // GetUnbondFromServiceCmd returns the command allowing to unbond from a service
// func GetUnbondFromServiceCmd() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:     "service [service-id] [amount]",
// 		Args:    cobra.ExactArgs(2),
// 		Short:   "Unbond the given amount from a service",
// 		Example: fmt.Sprintf("%s tx %s unbond service 1 1000000milk --from alice", version.AppName, types.ModuleName),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			clientCtx, err := client.GetClientTxContext(cmd)
// 			if err != nil {
// 				return err
// 			}
//
// 			serviceID, err := servicestypes.ParseServiceID(args[0])
// 			if err != nil {
// 				return err
// 			}
//
// 			amount, err := sdk.ParseCoinsNormalized(args[1])
// 			if err != nil {
// 				return err
// 			}
//
// 			delegator := clientCtx.FromAddress.String()
//
// 			// Create and validate the message
// 			msg := types.NewMsgUndelegateService(serviceID, amount, delegator)
// 			if err = msg.ValidateBasic(); err != nil {
// 				return fmt.Errorf("message validation failed: %w", err)
// 			}
//
// 			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
// 		},
// 	}
//
// 	flags.AddTxFlagsToCmd(cmd)
//
// 	return cmd
// }
//
// // --------------------------------------------------------------------------------------------------------------------
//
// // GetUpdateTxCmd returns the command allowing to update operator or service
// // params
// func GetUpdateTxCmd() *cobra.Command {
// 	txCmd := &cobra.Command{
// 		Use:   "update",
// 		Short: "Update transactions subcommands",
// 	}
//
// 	txCmd.AddCommand(
// 		GetUpdateOperatorParamsCmd(),
// 		GetUpdateServiceParamsCmd(),
// 	)
//
// 	return txCmd
// }
//
// // GetUpdateOperatorParamsCmd returns the command allowing to update an
// // operator's params.
// func GetUpdateOperatorParamsCmd() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:     "operator-params [operator-id] [commission-rate] [joined-service-ids]",
// 		Args:    cobra.ExactArgs(3),
// 		Short:   "Update an operator's params",
// 		Example: fmt.Sprintf("%s tx %s update operator-params 1 0.05 1,3,4 --from alice", version.AppName, types.ModuleName),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			clientCtx, err := client.GetClientTxContext(cmd)
// 			if err != nil {
// 				return err
// 			}
//
// 			operatorID, err := operatorstypes.ParseOperatorID(args[0])
// 			if err != nil {
// 				return err
// 			}
//
// 			commissionRate, err := math.LegacyNewDecFromStr(args[1])
// 			if err != nil {
// 				return fmt.Errorf("invalid commission rate: %w", err)
// 			}
//
// 			joinedServiceIDs, err := utils.ParseUint32Slice(args[2])
// 			if err != nil {
// 				return fmt.Errorf("parse joined service ids: %w", err)
// 			}
//
// 			params := types.NewOperatorParams(commissionRate, joinedServiceIDs)
//
// 			// Create and validate the message
// 			msg := types.NewMsgJoinService(operatorID, params, clientCtx.FromAddress.String())
// 			if err = msg.ValidateBasic(); err != nil {
// 				return fmt.Errorf("message validation failed: %w", err)
// 			}
//
// 			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
// 		},
// 	}
//
// 	flags.AddTxFlagsToCmd(cmd)
//
// 	return cmd
// }
//
// // GetUpdateServiceParamsCmd returns the command allowing to update a service's
// // params.
// func GetUpdateServiceParamsCmd() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:     "service-params [service-id] [slash-fraction] [whitelisted-pool-ids] [whitelisted-operator-ids]",
// 		Args:    cobra.ExactArgs(4),
// 		Short:   "Update a service's params",
// 		Example: fmt.Sprintf("%s tx %s update service-params 1 0.02 1,3,4 1,2,3,4,5 --from alice", version.AppName, types.ModuleName),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			clientCtx, err := client.GetClientTxContext(cmd)
// 			if err != nil {
// 				return err
// 			}
//
// 			serviceID, err := servicestypes.ParseServiceID(args[0])
// 			if err != nil {
// 				return err
// 			}
//
// 			slashFraction, err := math.LegacyNewDecFromStr(args[1])
// 			if err != nil {
// 				return fmt.Errorf("invalid slash fraction: %w", err)
// 			}
//
// 			whitelistedPoolIDs, err := utils.ParseUint32Slice(args[2])
// 			if err != nil {
// 				return fmt.Errorf("parse whitelisted pool ids: %w", err)
// 			}
//
// 			whitelistedOperatorIDs, err := utils.ParseUint32Slice(args[3])
// 			if err != nil {
// 				return fmt.Errorf("parse whitelisted operator ids: %w", err)
// 			}
//
// 			params := types.NewServiceParams(slashFraction, whitelistedPoolIDs, whitelistedOperatorIDs)
//
// 			// Create and validate the message
// 			msg := types.NewMsgUpdateServiceParams(serviceID, params, clientCtx.FromAddress.String())
// 			if err = msg.ValidateBasic(); err != nil {
// 				return fmt.Errorf("message validation failed: %w", err)
// 			}
//
// 			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
// 		},
// 	}
//
// 	flags.AddTxFlagsToCmd(cmd)
//
// 	return cmd
// }
