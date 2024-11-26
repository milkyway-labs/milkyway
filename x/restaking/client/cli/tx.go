package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/milkyway-labs/milkyway/v2/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/v2/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v2/x/pools/types"
	"github.com/milkyway-labs/milkyway/v2/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v2/x/services/types"
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
		GetUnbondTxCmd(),
		GetOperatorTxCmd(),
		GetServiceTxCmd(),
		GetUserTxCmd(),
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

			amount, err := sdk.ParseCoinNormalized(args[0])
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

// --------------------------------------------------------------------------------------------------------------------

// GetUnbondTxCmd returns the command allowing to unbond tokens
func GetUnbondTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "unbond",
		Short: "Unbond transactions subcommands",
	}

	txCmd.AddCommand(
		GetUnbondFromPoolCmd(),
		GetUnbondFromOperatorCmd(),
		GetUnbondFromServiceCmd(),
	)

	return txCmd
}

// GetUnbondFromPoolCmd returns the command allowing to unbond from a pool
func GetUnbondFromPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool [amount]",
		Args:    cobra.ExactArgs(1),
		Short:   "Unbond the given amount from a pool",
		Example: fmt.Sprintf("%s tx %s unbond pool 1000000umilk --from alice", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			delegator := clientCtx.FromAddress.String()

			// Create and validate the message
			msg := types.NewMsgUndelegatePool(amount, delegator)
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetUnbondFromOperatorCmd returns the command allowing to unbond from an operator
func GetUnbondFromOperatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "operator [operator-id] [amount]",
		Args:    cobra.ExactArgs(2),
		Short:   "Unbond the given amount from an operator",
		Example: fmt.Sprintf("%s tx %s unbond operator 1 1000000milk --from alice", version.AppName, types.ModuleName),
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
			msg := types.NewMsgUndelegateOperator(operatorID, amount, delegator)
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetUnbondFromServiceCmd returns the command allowing to unbond from a service
func GetUnbondFromServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service [service-id] [amount]",
		Args:    cobra.ExactArgs(2),
		Short:   "Unbond the given amount from a service",
		Example: fmt.Sprintf("%s tx %s unbond service 1 1000000milk --from alice", version.AppName, types.ModuleName),
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
			msg := types.NewMsgUndelegateService(serviceID, amount, delegator)
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// --------------------------------------------------------------------------------------------------------------------

// GetOperatorTxCmd returns the command to manage the operators.
func GetOperatorTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "operator",
		Short: "Restaking operator subcommands",
	}

	txCmd.AddCommand(
		GetJoinServiceCmd(),
		GetLeaveServiceCmd(),
	)

	return txCmd
}

// GetJoinServiceCmd returns the command allowing to add a service to the
// list of service joined by an operator.
func GetJoinServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "join-service [operator-id] [service-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "Join a service as a validator",
		Example: fmt.Sprintf("%s tx %s operator join-service 1 1 --from alice", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			operatorID, err := operatorstypes.ParseOperatorID(args[0])
			if err != nil {
				return fmt.Errorf("parse operator id: %w", err)
			}

			serviceID, err := servicestypes.ParseServiceID(args[1])
			if err != nil {
				return fmt.Errorf("parse service id: %w", err)
			}

			// Create and validate the message
			msg := types.NewMsgJoinService(operatorID, serviceID, clientCtx.FromAddress.String())
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetLeaveServiceCmd returns the command allowing to add a service to the
// list of service secured by an operator.
func GetLeaveServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "leave-service [operator-id] [service-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "Leave the service as a validator",
		Example: fmt.Sprintf("%s tx %s operator leave-service 1 1 --from alice", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			operatorID, err := operatorstypes.ParseOperatorID(args[0])
			if err != nil {
				return fmt.Errorf("parse operator id: %w", err)
			}

			serviceID, err := servicestypes.ParseServiceID(args[1])
			if err != nil {
				return fmt.Errorf("parse service id: %w", err)
			}

			// Create and validate the message
			msg := types.NewMsgLeaveService(operatorID, serviceID, clientCtx.FromAddress.String())
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// --------------------------------------------------------------------------------------------------------------------

func GetServiceTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "service",
		Short: "Restaking service subcommands",
	}

	txCmd.AddCommand(
		GetAllowOperatorTxCmd(),
		GetRemoveAllowedOperatorTxCmd(),
		GetBorrowPoolSecurityTxCmd(),
		GetCeasePoolSecurityBorrowTxCmd(),
	)

	return txCmd
}

func GetAllowOperatorTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "allow-operator [service-id] [operator-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "Adds a operator to the list of operators allowed to secure the service",
		Example: fmt.Sprintf("%s tx %s service allow-operator 1 1 --from alice", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			serviceID, err := servicestypes.ParseServiceID(args[0])
			if err != nil {
				return fmt.Errorf("parse service id: %w", err)
			}

			operatorID, err := operatorstypes.ParseOperatorID(args[1])
			if err != nil {
				return fmt.Errorf("parse operator id: %w", err)
			}

			// Create and validate the message
			msg := types.NewMsgAddOperatorToAllowList(serviceID, operatorID, clientCtx.FromAddress.String())
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetRemoveAllowedOperatorTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove-allowed-operator [service-id] [operator-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "Removes a operator from the list of operators allowed to secure the service",
		Example: fmt.Sprintf("%s tx %s service remove-allowed-operator 1 1 --from alice", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			serviceID, err := servicestypes.ParseServiceID(args[0])
			if err != nil {
				return fmt.Errorf("parse service id: %w", err)
			}

			operatorID, err := operatorstypes.ParseOperatorID(args[1])
			if err != nil {
				return fmt.Errorf("parse operator id: %w", err)
			}

			// Create and validate the message
			msg := types.NewMsgRemoveOperatorFromAllowList(serviceID, operatorID, clientCtx.FromAddress.String())
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetBorrowPoolSecurityTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "borrow-pool-security [service-id] [pool-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "Adds a pool to the list of pools from which the service has chosen to borrow security",
		Example: fmt.Sprintf("%s tx %s service borrow-pool-security 1 1 --from alice", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			serviceID, err := servicestypes.ParseServiceID(args[0])
			if err != nil {
				return fmt.Errorf("parse service id: %w", err)
			}

			poolID, err := poolstypes.ParsePoolID(args[1])
			if err != nil {
				return fmt.Errorf("parse pool id: %w", err)
			}

			// Create and validate the message
			msg := types.NewMsgBorrowPoolSecurity(serviceID, poolID, clientCtx.FromAddress.String())
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCeasePoolSecurityBorrowTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cease-pool-security-borrow [service-id] [pool-id]",
		Args:    cobra.ExactArgs(2),
		Short:   "Removes a pool from the list of pools from which the service has chosen to borrow security",
		Example: fmt.Sprintf("%s tx %s service cease-pool-security-borrow 1 1 --from alice", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			serviceID, err := servicestypes.ParseServiceID(args[0])
			if err != nil {
				return fmt.Errorf("parse service id: %w", err)
			}

			poolID, err := poolstypes.ParsePoolID(args[1])
			if err != nil {
				return fmt.Errorf("parse pool id: %w", err)
			}

			// Create and validate the message
			msg := types.NewMsgCeasePoolSecurityBorrow(serviceID, poolID, clientCtx.FromAddress.String())
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// --------------------------------------------------------------------------------------------------------------------

//nolint:gosec // This is not a hardcoded credential
const (
	trustNonAccreditedServicesFlag = "trust-non-accredited-services"
	trustAccreditedServicesFlag    = "trust-accredited-services"
	trustedServicesIDsFlag         = "trusted-services-ids"
)

func GetUserTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "user",
		Short: "Restaking user subcommands",
	}

	txCmd.AddCommand(
		GetSetUserPreferencesCmd(),
	)

	return txCmd
}

func GetSetUserPreferencesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-preferences",
		Args:  cobra.NoArgs,
		Short: "Set your user preferences regarding the restaking module",
		Long: `Set your user preferences regarding the restaking module.

If you are updating your preferences, you must provide all the flags that you want to set 
(i.e. the values you provide will completely override the existing ones)`,
		Example: fmt.Sprintf(`%s tx %s user set-preferences \
--trust-accredited-services \
--trust-non-accredited-services \
--trusted-services-ids 1,2,3 \
--from alice`, version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			trustAccreditedServices, err := cmd.Flags().GetBool(trustAccreditedServicesFlag)
			if err != nil {
				return err
			}

			trustNonAccreditedServices, err := cmd.Flags().GetBool(trustNonAccreditedServicesFlag)
			if err != nil {
				return err
			}

			trustedServices, err := cmd.Flags().GetUintSlice(trustedServicesIDsFlag)
			if err != nil {
				return err
			}
			trustedServicesIDs := utils.Map(trustedServices, func(t uint) uint32 { return uint32(t) })

			// Create and validate the message
			preferences := types.NewUserPreferences(trustNonAccreditedServices, trustAccreditedServices, trustedServicesIDs)
			msg := types.NewMsgSetUserPreferences(preferences, clientCtx.FromAddress.String())
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Bool(trustNonAccreditedServicesFlag, false, "Trust non-accredited services")
	cmd.Flags().Bool(trustAccreditedServicesFlag, false, "Trust accredited services")
	cmd.Flags().UintSlice(trustedServicesIDsFlag, nil, "List of IDs of the services you trust")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
