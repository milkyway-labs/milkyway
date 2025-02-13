package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/hashicorp/go-metrics"

	operatorstypes "github.com/milkyway-labs/milkyway/v9/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v9/x/pools/types"
	"github.com/milkyway-labs/milkyway/v9/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v9/x/services/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// JoinService defines the rpc method for Msg/JoinService
func (k msgServer) JoinService(ctx context.Context, msg *types.MsgJoinService) (*types.MsgJoinServiceResponse, error) {
	operator, err := k.operatorsKeeper.GetOperator(ctx, msg.OperatorID)
	if err != nil {
		return nil, err
	}

	if operator.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only the admin can join the service")
	}

	service, err := k.servicesKeeper.GetService(ctx, msg.ServiceID)
	if err != nil {
		return nil, err
	}

	if !service.IsActive() {
		return nil, errors.Wrapf(servicestypes.ErrServiceNotActive, "service %d is not active", msg.ServiceID)
	}

	// Make sure the operator is allowed to join
	operatorAllowed, err := k.CanOperatorValidateService(ctx, msg.ServiceID, msg.OperatorID)
	if err != nil {
		return nil, err
	}

	if !operatorAllowed {
		return nil, types.ErrOperatorNotAllowed
	}

	err = k.AddServiceToOperatorJoinedServices(ctx, msg.OperatorID, msg.ServiceID)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeJoinService,
			sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, fmt.Sprint(msg.OperatorID)),
			sdk.NewAttribute(servicestypes.AttributeKeyServiceID, fmt.Sprintf("%d", msg.ServiceID)),
		),
	})

	return &types.MsgJoinServiceResponse{}, nil
}

// LeaveService defines the rpc method for Msg/LeaveService
func (k msgServer) LeaveService(ctx context.Context, msg *types.MsgLeaveService) (*types.MsgLeaveServiceResponse, error) {
	operator, err := k.operatorsKeeper.GetOperator(ctx, msg.OperatorID)
	if err != nil {
		return nil, err
	}

	if operator.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only the admin can leave the service")
	}

	_, err = k.servicesKeeper.GetService(ctx, msg.ServiceID)
	if err != nil {
		return nil, err
	}

	err = k.RemoveServiceFromOperatorJoinedServices(ctx, msg.OperatorID, msg.ServiceID)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLeaveService,
			sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, fmt.Sprint(msg.OperatorID)),
			sdk.NewAttribute(servicestypes.AttributeKeyServiceID, fmt.Sprint(msg.ServiceID)),
		),
	})

	return &types.MsgLeaveServiceResponse{}, nil
}

// AddOperatorToAllowList defines the rpc method for Msg/AddOperatorToAllowList
func (k msgServer) AddOperatorToAllowList(ctx context.Context, msg *types.MsgAddOperatorToAllowList) (*types.MsgAddOperatorToAllowListResponse, error) {
	// Ensure that the service exists
	service, err := k.servicesKeeper.GetService(ctx, msg.ServiceID)
	if err != nil {
		return nil, err
	}

	// Ensure the service admin is performing this action
	if service.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only the service admin can allow an operator")
	}

	// Ensure that the operator exists
	_, err = k.operatorsKeeper.GetOperator(ctx, msg.OperatorID)
	if err != nil {
		return nil, err
	}

	// Ensure operator is not allowed before adding it
	operatorAllowed, err := k.IsOperatorInServiceAllowList(ctx, msg.ServiceID, msg.OperatorID)
	if err != nil {
		return nil, err
	}
	if operatorAllowed {
		return nil, types.ErrOperatorAlreadyAllowed
	}

	// If the service didn't have any operators in the allow list, it means the
	// service allowed any operators to join the service. Since we're adding an
	// operator to the allow list, we need make all the operators that have joined
	// the service to leave the service
	configured, err := k.IsServiceOperatorsAllowListConfigured(ctx, msg.ServiceID)
	if err != nil {
		return nil, err
	}
	if !configured {
		var leavingOperatorsIDs []uint32
		err = k.IterateServiceValidatingOperators(ctx, msg.ServiceID, func(operatorID uint32) (stop bool, err error) {
			if operatorID != msg.OperatorID {
				leavingOperatorsIDs = append(leavingOperatorsIDs, operatorID)
			}
			return false, nil
		})
		if err != nil {
			return nil, err
		}

		for _, operatorID := range leavingOperatorsIDs {
			err = k.RemoveServiceFromOperatorJoinedServices(ctx, operatorID, msg.ServiceID)
			if err != nil {
				return nil, err
			}
		}
	}

	// Add the operator to the service's operators allow list
	err = k.AddOperatorToServiceAllowList(ctx, msg.ServiceID, msg.OperatorID)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAllowOperator,
			sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, fmt.Sprint(msg.OperatorID)),
			sdk.NewAttribute(servicestypes.AttributeKeyServiceID, fmt.Sprint(msg.ServiceID)),
		),
	})

	return &types.MsgAddOperatorToAllowListResponse{}, nil
}

// RemoveOperatorFromAllowlist defines the rpc method for Msg/RemoveOperatorFromAllowlist
func (k msgServer) RemoveOperatorFromAllowlist(ctx context.Context, msg *types.MsgRemoveOperatorFromAllowlist) (*types.MsgRemoveOperatorFromAllowlistResponse, error) {
	// Ensure that the service exists
	service, err := k.servicesKeeper.GetService(ctx, msg.ServiceID)
	if err != nil {
		return nil, err
	}

	// Ensure the service admin is performing this action
	if service.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only the service admin can allow an operator")
	}

	// Ensure the operator is allowed before removing
	operatorAllowed, err := k.IsOperatorInServiceAllowList(ctx, msg.ServiceID, msg.OperatorID)
	if err != nil {
		return nil, err
	}
	if !operatorAllowed {
		return nil, types.ErrOperatorNotAllowed
	}

	// Remove the operator from the service's operators allow list
	err = k.RemoveOperatorFromServiceAllowList(ctx, msg.ServiceID, msg.OperatorID)
	if err != nil {
		return nil, err
	}

	// Check if the operator was the last operator in the service's allow list
	// If so, removing the operator from the allow list means the service
	// allows all operators so we don't forcefully remove the service from the
	// operator's joined services list
	configured, err := k.IsServiceOperatorsAllowListConfigured(ctx, msg.ServiceID)
	if err != nil {
		return nil, err
	}
	if configured {
		// Since the operator is removed from the service's allow list, remove the
		// service from the operator's joined services list
		err = k.RemoveServiceFromOperatorJoinedServices(ctx, msg.OperatorID, msg.ServiceID)
		if err != nil {
			return nil, err
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRemoveAllowedOperator,
			sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, fmt.Sprint(msg.OperatorID)),
			sdk.NewAttribute(servicestypes.AttributeKeyServiceID, fmt.Sprint(msg.ServiceID)),
		),
	})

	return &types.MsgRemoveOperatorFromAllowlistResponse{}, nil
}

// BorrowPoolSecurity defines the rpc method for Msg/BorrowPoolSecurity
func (k msgServer) BorrowPoolSecurity(ctx context.Context, msg *types.MsgBorrowPoolSecurity) (*types.MsgBorrowPoolSecurityResponse, error) {
	// Ensure that the service exists
	service, err := k.servicesKeeper.GetService(ctx, msg.ServiceID)
	if err != nil {
		return nil, err
	}

	// Ensure the service admin is performing this action
	if service.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized,
			"only the service admin can borrow the security from a new pool")
	}

	// Ensure the service is active
	if !service.IsActive() {
		return nil, errors.Wrapf(servicestypes.ErrServiceNotActive, "service %d is not active", msg.ServiceID)
	}

	// Ensure that the pool exists
	_, err = k.poolsKeeper.GetPool(ctx, msg.PoolID)
	if err != nil {
		return nil, err
	}

	// Ensure the pool is not in the service securing pools list
	isSecuring, err := k.IsServiceSecuredByPool(ctx, msg.ServiceID, msg.PoolID)
	if err != nil {
		return nil, err
	}
	if isSecuring {
		return nil, types.ErrPoolAlreadySecuringService
	}

	// Add the pool to the list of pools from which the service can borrow
	// security from
	err = k.AddPoolToServiceSecuringPools(ctx, msg.ServiceID, msg.PoolID)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBorrowPoolSecurity,
			sdk.NewAttribute(servicestypes.AttributeKeyServiceID, fmt.Sprint(msg.ServiceID)),
			sdk.NewAttribute(poolstypes.AttributeKeyPoolID, fmt.Sprint(msg.PoolID)),
		),
	})

	return &types.MsgBorrowPoolSecurityResponse{}, nil
}

// CeasePoolSecurityBorrow defines the rpc method for Msg/CeasePoolSecurityBorrow
func (k msgServer) CeasePoolSecurityBorrow(ctx context.Context, msg *types.MsgCeasePoolSecurityBorrow) (*types.MsgCeasePoolSecurityBorrowResponse, error) {
	// Ensure that the service exists
	service, err := k.servicesKeeper.GetService(ctx, msg.ServiceID)
	if err != nil {
		return nil, err
	}

	// Ensure the service admin is performing this action
	if service.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized,
			"only the service admin can cease the pool security borrow")
	}

	// Ensure the pool is in the service securing pools list
	isSecuring, err := k.IsServiceSecuredByPool(ctx, msg.ServiceID, msg.PoolID)
	if err != nil {
		return nil, err
	}
	if !isSecuring {
		return nil, types.ErrPoolNotSecuringService
	}

	// Remove the pool from the service's pools whitelist
	err = k.RemovePoolFromServiceSecuringPools(ctx, msg.ServiceID, msg.PoolID)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCeasePoolSecurityBorrow,
			sdk.NewAttribute(servicestypes.AttributeKeyServiceID, fmt.Sprint(msg.ServiceID)),
			sdk.NewAttribute(poolstypes.AttributeKeyPoolID, fmt.Sprint(msg.PoolID)),
		),
	})

	return &types.MsgCeasePoolSecurityBorrowResponse{}, nil
}

// DelegatePool defines the rpc method for Msg/DelegatePool
func (k msgServer) DelegatePool(ctx context.Context, msg *types.MsgDelegatePool) (*types.MsgDelegatePoolResponse, error) {
	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return nil, errors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid delegation amount",
		)
	}

	newShares, err := k.DelegateToPool(ctx, msg.Amount, msg.Delegator)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, types.ModuleName, "pool_restake")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", sdk.MsgTypeURL(msg)},
				float32(msg.Amount.Amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", msg.Amount.Denom)},
			)
		}()
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegatePool,
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.Delegator),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyNewShares, newShares.String()),
		),
	})

	return &types.MsgDelegatePoolResponse{}, nil
}

// UndelegatePool defines the rpc method for Msg/UndelegatePool
func (k msgServer) UndelegatePool(ctx context.Context, msg *types.MsgUndelegatePool) (*types.MsgUndelegateResponse, error) {
	// Perform the undelegation
	completionTime, err := k.UndelegateFromPool(ctx, msg.Amount, msg.Delegator)
	if err != nil {
		return nil, err
	}

	// Log the undelegation
	if msg.Amount.Amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, types.ModuleName, "undelegate_pool")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", sdk.MsgTypeURL(msg)},
				float32(msg.Amount.Amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", msg.Amount.Denom)},
			)
		}()
	}

	// Emit the undelegation event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnbondPool,
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.Delegator),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
	})

	return &types.MsgUndelegateResponse{
		CompletionTime: completionTime,
	}, nil
}

// DelegateOperator defines the rpc method for Msg/DelegateOperator
func (k msgServer) DelegateOperator(ctx context.Context, msg *types.MsgDelegateOperator) (*types.MsgDelegateOperatorResponse, error) {
	if !msg.Amount.IsValid() || !msg.Amount.IsAllPositive() {
		return nil, errors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid delegation amount",
		)
	}

	newShares, err := k.Keeper.DelegateToOperator(ctx, msg.OperatorID, msg.Amount, msg.Delegator)
	if err != nil {
		return nil, err
	}

	for _, token := range msg.Amount {
		if token.Amount.IsInt64() {
			defer func() {
				telemetry.IncrCounter(1, types.ModuleName, "operator_restake")
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", sdk.MsgTypeURL(msg)},
					float32(token.Amount.Int64()),
					[]metrics.Label{
						telemetry.NewLabel("operator_id", fmt.Sprintf("%d", msg.OperatorID)),
						telemetry.NewLabel("denom", token.Denom),
					},
				)
			}()
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegateOperator,
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.Delegator),
			sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, fmt.Sprintf("%d", msg.OperatorID)),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyNewShares, newShares.String()),
		),
	})

	return &types.MsgDelegateOperatorResponse{}, nil
}

// UndelegateOperator defines the rpc method for Msg/UndelegateOperator
func (k msgServer) UndelegateOperator(ctx context.Context, msg *types.MsgUndelegateOperator) (*types.MsgUndelegateResponse, error) {
	// Perform the undelegation
	completionTime, err := k.Keeper.UndelegateFromOperator(ctx, msg.OperatorID, msg.Amount, msg.Delegator)
	if err != nil {
		return nil, err
	}

	// Log the undelegation
	for _, token := range msg.Amount {
		if token.Amount.IsInt64() {
			defer func() {
				telemetry.IncrCounter(1, types.ModuleName, "undelegate_operator")
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", sdk.MsgTypeURL(msg)},
					float32(token.Amount.Int64()),
					[]metrics.Label{
						telemetry.NewLabel("operator_id", fmt.Sprintf("%d", msg.OperatorID)),
						telemetry.NewLabel("denom", token.Denom),
					},
				)
			}()
		}
	}

	// Emit the undelegation event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnbondOperator,
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.Delegator),
			sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, fmt.Sprintf("%d", msg.OperatorID)),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
	})

	return &types.MsgUndelegateResponse{
		CompletionTime: completionTime,
	}, nil
}

// DelegateService defines the rpc method for Msg/DelegateService
func (k msgServer) DelegateService(ctx context.Context, msg *types.MsgDelegateService) (*types.MsgDelegateServiceResponse, error) {
	if !msg.Amount.IsValid() || !msg.Amount.IsAllPositive() {
		return nil, errors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid delegation amount",
		)
	}

	newShares, err := k.Keeper.DelegateToService(ctx, msg.ServiceID, msg.Amount, msg.Delegator)
	if err != nil {
		return nil, err
	}

	for _, token := range msg.Amount {
		if token.Amount.IsInt64() {
			defer func() {
				telemetry.IncrCounter(1, types.ModuleName, "service restake")
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", sdk.MsgTypeURL(msg)},
					float32(token.Amount.Int64()),
					[]metrics.Label{
						telemetry.NewLabel("service_id", fmt.Sprintf("%d", msg.ServiceID)),
						telemetry.NewLabel("denom", token.Denom),
					},
				)
			}()
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegateService,
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.Delegator),
			sdk.NewAttribute(servicestypes.AttributeKeyServiceID, fmt.Sprintf("%d", msg.ServiceID)),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyNewShares, newShares.String()),
		),
	})

	return &types.MsgDelegateServiceResponse{}, nil
}

// UndelegateService defines the rpc method for Msg/UndelegateService
func (k msgServer) UndelegateService(ctx context.Context, msg *types.MsgUndelegateService) (*types.MsgUndelegateResponse, error) {
	// Perform the undelegation
	completionTime, err := k.Keeper.UndelegateFromService(ctx, msg.ServiceID, msg.Amount, msg.Delegator)
	if err != nil {
		return nil, err
	}

	// Log the undelegation
	for _, token := range msg.Amount {
		if token.Amount.IsInt64() {
			defer func() {
				telemetry.IncrCounter(1, types.ModuleName, "undelegate_service")
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", sdk.MsgTypeURL(msg)},
					float32(token.Amount.Int64()),
					[]metrics.Label{
						telemetry.NewLabel("service_id", fmt.Sprintf("%d", msg.ServiceID)),
						telemetry.NewLabel("denom", token.Denom),
					},
				)
			}()
		}
	}

	// Emit the undelegation event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnbondService,
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.Delegator),
			sdk.NewAttribute(servicestypes.AttributeKeyServiceID, fmt.Sprintf("%d", msg.ServiceID)),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
	})

	return &types.MsgUndelegateResponse{
		CompletionTime: completionTime,
	}, nil
}

// SetUserPreferences defines the rpc method for Msg/SetUserPreferences
func (k msgServer) SetUserPreferences(ctx context.Context, msg *types.MsgSetUserPreferences) (*types.MsgSetUserPreferencesResponse, error) {
	// Make sure that each service exists
	for _, entry := range msg.Preferences.TrustedServices {
		hasService, err := k.servicesKeeper.HasService(ctx, entry.ServiceID)
		if err != nil {
			return nil, err
		}

		if !hasService {
			return nil, errors.Wrapf(servicestypes.ErrServiceNotFound, "service %d does not exist", entry.ServiceID)
		}

		// Make sure the pools exist
		for _, poolID := range entry.PoolsIDs {
			hasPool, err := k.poolsKeeper.HasPool(ctx, poolID)
			if err != nil {
				return nil, err
			}

			if !hasPool {
				return nil, errors.Wrapf(poolstypes.ErrPoolNotFound, "pool %d does not exist", poolID)
			}
		}
	}

	err := k.Keeper.SetUserPreferences(ctx, msg.User, msg.Preferences)
	if err != nil {
		return nil, err
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSetUserPreferences,
			sdk.NewAttribute(types.AttributeKeyUser, msg.User),
		),
	})

	return &types.MsgSetUserPreferencesResponse{}, nil
}

// UpdateParams defines the rpc method for Msg/UpdateParams
func (k msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	// Check the authority
	authority := k.authority
	if authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", authority, msg.Authority)
	}

	// Update the params
	err := k.SetParams(ctx, msg.Params)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
