package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/milkyway-labs/milkyway/v7/x/services/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(k *Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

// CreateService defines the rpc method for Msg/CreateService
func (k msgServer) CreateService(goCtx context.Context, msg *types.MsgCreateService) (*types.MsgCreateServiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the next service id
	serviceID, err := k.GetNextServiceID(ctx)
	if err != nil {
		return nil, err
	}

	// Create the Service and validate it
	service := types.NewService(
		serviceID,
		types.SERVICE_STATUS_CREATED,
		msg.Name,
		msg.Description,
		msg.Website,
		msg.PictureURL,
		msg.Sender,
		false,
	)

	// Validate the service before storing
	err = service.Validate()
	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Charge for the creation
	// We do not place this inside the CreateService method to avoid charging fees during genesis
	// init and other places that use that method
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	if !params.ServiceRegistrationFee.IsZero() {
		// Make sure the specified fees are enough
		if !msg.FeeAmount.IsAnyGTE(params.ServiceRegistrationFee) {
			return nil, errors.Wrapf(sdkerrors.ErrInsufficientFunds, "insufficient funds: %s < %s", msg.FeeAmount, params.ServiceRegistrationFee)
		}

		userAddress, err := sdk.AccAddressFromBech32(service.Admin)
		if err != nil {
			return nil, errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid service admin address: %s", service.Admin)
		}

		err = k.poolKeeper.FundCommunityPool(ctx, msg.FeeAmount, userAddress)
		if err != nil {
			return nil, err
		}
	}

	// Create the service
	err = k.Keeper.CreateService(ctx, service)
	if err != nil {
		return nil, err
	}

	// Update the ID for the next service
	err = k.SetNextServiceID(ctx, service.ID+1)
	if err != nil {
		return nil, err
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateService,
			sdk.NewAttribute(types.AttributeKeyServiceID, fmt.Sprintf("%d", service.ID)),
		),
	})

	return &types.MsgCreateServiceResponse{
		NewServiceID: service.ID,
	}, nil
}

// UpdateService defines the rpc method for Msg/UpdateService
func (k msgServer) UpdateService(ctx context.Context, msg *types.MsgUpdateService) (*types.MsgUpdateServiceResponse, error) {
	// Check if the service exists
	service, err := k.GetService(ctx, msg.ServiceID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, types.ErrServiceNotFound
		}
		return nil, err
	}

	// Make sure the user that is updating the service is the admin
	if service.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only the admin can update the service")
	}

	// Update the service
	updated := service.Update(types.NewServiceUpdate(msg.Name, msg.Description, msg.Website, msg.PictureURL))

	// Validate the updated service
	err = updated.Validate()
	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Save the service
	if err := k.SaveService(ctx, updated); err != nil {
		return nil, err
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateService,
			sdk.NewAttribute(types.AttributeKeyServiceID, fmt.Sprintf("%d", msg.ServiceID)),
		),
	})

	return &types.MsgUpdateServiceResponse{}, nil
}

func (k msgServer) ActivateService(ctx context.Context, msg *types.MsgActivateService) (*types.MsgActivateServiceResponse, error) {
	// Check if the service exists
	service, err := k.GetService(ctx, msg.ServiceID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, types.ErrServiceNotFound
		}
		return nil, err
	}

	// Make sure the user that is activating the service is the admin
	if service.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only the admin can activate the service")
	}

	// Activate the service
	err = k.Keeper.ActivateService(ctx, msg.ServiceID)
	if err != nil {
		return nil, err
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeActivateService,
			sdk.NewAttribute(types.AttributeKeyServiceID, fmt.Sprintf("%d", msg.ServiceID)),
		),
	})

	return &types.MsgActivateServiceResponse{}, nil
}

// DeactivateService defines the rpc method for Msg/DeactivateService
func (k msgServer) DeactivateService(ctx context.Context, msg *types.MsgDeactivateService) (*types.MsgDeactivateServiceResponse, error) {
	// Check if the service exists
	service, err := k.GetService(ctx, msg.ServiceID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, types.ErrServiceNotFound
		}
		return nil, err
	}

	// Make sure the user that is deactivating the service is the admin
	if service.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only the admin can deactivate the service")
	}

	// Deactivate the service
	err = k.Keeper.DeactivateService(ctx, msg.ServiceID)
	if err != nil {
		return nil, err
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDeactivateService,
			sdk.NewAttribute(types.AttributeKeyServiceID, fmt.Sprintf("%d", msg.ServiceID)),
		),
	})

	return &types.MsgDeactivateServiceResponse{}, nil
}

func (k msgServer) DeleteService(ctx context.Context, msg *types.MsgDeleteService) (*types.MsgDeleteServiceResponse, error) {
	// Check if the service exists
	service, err := k.GetService(ctx, msg.ServiceID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, types.ErrServiceNotFound
		}
		return nil, err
	}

	// Make sure the user that is deleting the service is the admin
	if service.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only the admin can delete the service")
	}

	// Delete the service from the store
	err = k.Keeper.DeleteService(ctx, msg.ServiceID)
	if err != nil {
		return nil, err
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDeleteService,
			sdk.NewAttribute(types.AttributeKeyServiceID, fmt.Sprintf("%d", msg.ServiceID)),
		),
	})

	return &types.MsgDeleteServiceResponse{}, nil
}

// TransferServiceOwnership defines the rpc method for Msg/TransferServiceOwnership
func (k msgServer) TransferServiceOwnership(ctx context.Context, msg *types.MsgTransferServiceOwnership) (*types.MsgTransferServiceOwnershipResponse, error) {
	// Check if the service exists
	service, err := k.GetService(ctx, msg.ServiceID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, types.ErrServiceNotFound
		}
		return nil, err
	}

	// Make sure only the admin can transfer the service ownership
	if service.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only the admin can transfer the service ownership")
	}

	// Update the service admin
	service.Admin = msg.NewAdmin
	if err := k.SaveService(ctx, service); err != nil {
		return nil, err
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferServiceOwnership,
			sdk.NewAttribute(types.AttributeKeyServiceID, fmt.Sprintf("%d", msg.ServiceID)),
			sdk.NewAttribute(types.AttributeKeyNewAdmin, msg.NewAdmin),
		),
	})

	return &types.MsgTransferServiceOwnershipResponse{}, nil
}

// SetServiceParams define the rpc method for Msg/SetServiceParams
func (k msgServer) SetServiceParams(ctx context.Context, msg *types.MsgSetServiceParams) (*types.MsgSetServiceParamsResponse, error) {
	// Get the service whose params are being set
	service, err := k.GetService(ctx, msg.ServiceID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, types.ErrServiceNotFound
		}
		return nil, err
	}

	// Ensure the sender is the service admin
	if msg.Sender != service.Admin {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "sender must be the service admin")
	}

	// Set the service params
	err = k.Keeper.SetServiceParams(ctx, service.ID, msg.ServiceParams)
	if err != nil {
		return nil, err
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSetServiceParams,
			sdk.NewAttribute(types.AttributeKeyServiceID, fmt.Sprintf("%d", msg.ServiceID)),
		),
	})

	return &types.MsgSetServiceParamsResponse{}, nil
}

// UpdateParams defines the rpc method for Msg/UpdateParams
func (k msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	// Check the authority
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	// Update the params
	err := k.SetParams(ctx, msg.Params)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

// AccreditService defines the rpc method for Msg/AccreditService
func (k msgServer) AccreditService(ctx context.Context, msg *types.MsgAccreditService) (*types.MsgAccreditServiceResponse, error) {
	// Check the authority
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	err := k.Keeper.SetServiceAccredited(ctx, msg.ServiceID, true)
	if err != nil {
		return nil, err
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAccreditService,
			sdk.NewAttribute(types.AttributeKeyServiceID, fmt.Sprintf("%d", msg.ServiceID)),
		),
	})

	return &types.MsgAccreditServiceResponse{}, nil
}

// RevokeServiceAccreditation defines the rpc method for Msg/RevokeServiceAccreditation
func (k msgServer) RevokeServiceAccreditation(ctx context.Context, msg *types.MsgRevokeServiceAccreditation) (*types.MsgRevokeServiceAccreditationResponse, error) {
	// Check the authority
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	err := k.Keeper.SetServiceAccredited(ctx, msg.ServiceID, false)
	if err != nil {
		return nil, err
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRevokeServiceAccreditation,
			sdk.NewAttribute(types.AttributeKeyServiceID, fmt.Sprintf("%d", msg.ServiceID)),
		),
	})

	return &types.MsgRevokeServiceAccreditationResponse{}, nil
}
