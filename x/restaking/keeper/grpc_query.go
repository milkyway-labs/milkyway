package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

type Querier struct {
	*Keeper
}

var _ types.QueryServer = Querier{}

func NewQuerier(keeper *Keeper) Querier {
	return Querier{Keeper: keeper}
}

// OperatorJoinedServices queries the services joined by the operator with the
// given ID
func (k Querier) OperatorJoinedServices(goCtx context.Context, req *types.QueryOperatorJoinedServicesRequest) (*types.QueryOperatorJoinedServicesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.OperatorId == 0 {
		return nil, status.Error(codes.InvalidArgument, "operator id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	_, found := k.operatorsKeeper.GetOperator(ctx, req.OperatorId)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "operator not found")
	}

	// Get the operator joined services
	joinedServices, err := k.GetOperatorJoinedServices(ctx, req.OperatorId)
	if err != nil {
		return nil, err
	}

	return &types.QueryOperatorJoinedServicesResponse{ServiceIds: joinedServices.ServiceIDs}, nil
}

// ServiceParams queries the service params for the given service id
func (k Querier) ServiceParams(goCtx context.Context, req *types.QueryServiceParamsRequest) (*types.QueryServiceParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "service id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the service params store
	params, err := k.GetServiceParams(ctx, req.ServiceId)
	if err != nil {
		return nil, err
	}

	return &types.QueryServiceParamsResponse{ServiceParams: params}, nil
}

// PoolDelegations queries the pool delegations for the given pool id
func (k Querier) PoolDelegations(goCtx context.Context, req *types.QueryPoolDelegationsRequest) (*types.QueryPoolDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.PoolId == 0 {
		return nil, status.Error(codes.InvalidArgument, "pool id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the pool delegations store
	store := ctx.KVStore(k.storeKey)
	delegationsStore := prefix.NewStore(store, types.PoolDelegationPrefix)

	// Query the pool delegations for the given pool id
	delegations, pageRes, err := query.GenericFilteredPaginate(k.cdc, delegationsStore, req.Pagination, func(key []byte, delegation *types.Delegation) (*types.Delegation, error) {
		if delegation.TargetID != req.PoolId {
			return nil, nil
		}
		return delegation, nil
	}, func() *types.Delegation {
		return &types.Delegation{}
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	poolDelegations := make([]types.DelegationResponse, len(delegations))
	for i, delegation := range delegations {
		response, err := PoolDelegationToPoolDelegationResponse(ctx, k.Keeper, *delegation)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		poolDelegations[i] = response
	}

	return &types.QueryPoolDelegationsResponse{
		Delegations: poolDelegations,
		Pagination:  pageRes,
	}, nil
}

// PoolDelegation queries the pool delegation for the given pool id and user address
func (k Querier) PoolDelegation(goCtx context.Context, req *types.QueryPoolDelegationRequest) (*types.QueryPoolDelegationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.PoolId == 0 {
		return nil, status.Error(codes.InvalidArgument, "pool id cannot be zero")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	poolDelegation, found := k.GetPoolDelegation(ctx, req.PoolId, req.DelegatorAddress)
	if !found {
		return nil, status.Error(codes.NotFound, "pool delegation not found")
	}

	response, err := PoolDelegationToPoolDelegationResponse(ctx, k.Keeper, poolDelegation)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPoolDelegationResponse{
		Delegation: response,
	}, nil
}

// PoolUnbondingDelegations queries the pool unbonding delegations for the given pool id
func (k Querier) PoolUnbondingDelegations(goCtx context.Context, req *types.QueryPoolUnbondingDelegationsRequest) (*types.QueryPoolUnbondingDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.PoolId == 0 {
		return nil, status.Error(codes.InvalidArgument, "pool id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the pool unbonding delegations store
	store := ctx.KVStore(k.storeKey)
	delegationsStore := prefix.NewStore(store, types.PoolUnbondingDelegationPrefix)

	// Query the pool unbonding delegations for the given pool id
	var unbondingDelegations []types.UnbondingDelegation
	pageRes, err := query.Paginate(delegationsStore, req.Pagination, func(key []byte, value []byte) error {
		unbond, err := types.UnmarshalUnbondingDelegation(k.cdc, value)
		if err != nil {
			return err
		}
		unbondingDelegations = append(unbondingDelegations, unbond)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPoolUnbondingDelegationsResponse{
		UnbondingDelegations: unbondingDelegations,
		Pagination:           pageRes,
	}, nil
}

// PoolUnbondingDelegation queries the pool unbonding delegation for the given pool id and user address
func (k Querier) PoolUnbondingDelegation(goCtx context.Context, req *types.QueryPoolUnbondingDelegationRequest) (*types.QueryPoolUnbondingDelegationResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Errorf(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.PoolId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "pool id cannot be zero")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	unbond, found := k.GetPoolUnbondingDelegation(ctx, req.PoolId, req.DelegatorAddress)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"unbonding delegation with delegator %s not found for pool %d",
			req.DelegatorAddress, req.PoolId,
		)
	}

	return &types.QueryPoolUnbondingDelegationResponse{
		UnbondingDelegation: unbond,
	}, nil
}

// OperatorDelegations queries the operator delegations for the given operator id
func (k Querier) OperatorDelegations(goCtx context.Context, req *types.QueryOperatorDelegationsRequest) (*types.QueryOperatorDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.OperatorId == 0 {
		return nil, status.Error(codes.InvalidArgument, "operator id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the operator delegations store
	store := ctx.KVStore(k.storeKey)
	delegationsStore := prefix.NewStore(store, types.OperatorDelegationPrefix)

	// Query the operator delegations for the given pool id
	delegations, pageRes, err := query.GenericFilteredPaginate(k.cdc, delegationsStore, req.Pagination, func(key []byte, delegation *types.Delegation) (*types.Delegation, error) {
		if delegation.TargetID != req.OperatorId {
			return nil, nil
		}
		return delegation, nil
	}, func() *types.Delegation {
		return &types.Delegation{}
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	operatorDelegations := make([]types.DelegationResponse, len(delegations))
	for i, delegation := range delegations {
		response, err := OperatorDelegationToOperatorDelegationResponse(ctx, k.Keeper, *delegation)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		operatorDelegations[i] = response
	}

	return &types.QueryOperatorDelegationsResponse{
		Delegations: operatorDelegations,
		Pagination:  pageRes,
	}, nil
}

// OperatorDelegation queries the operator delegation for the given operator id and user address
func (k Querier) OperatorDelegation(goCtx context.Context, req *types.QueryOperatorDelegationRequest) (*types.QueryOperatorDelegationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.OperatorId == 0 {
		return nil, status.Error(codes.InvalidArgument, "operator id cannot be zero")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	operatorDelegation, found := k.GetOperatorDelegation(ctx, req.OperatorId, req.DelegatorAddress)
	if !found {
		return nil, status.Error(codes.NotFound, "operator delegation not found")
	}

	response, err := OperatorDelegationToOperatorDelegationResponse(ctx, k.Keeper, operatorDelegation)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOperatorDelegationResponse{
		Delegation: response,
	}, nil
}

// OperatorUnbondingDelegations queries the operator unbonding delegations for the given operator id
func (k Querier) OperatorUnbondingDelegations(goCtx context.Context, req *types.QueryOperatorUnbondingDelegationsRequest) (*types.QueryOperatorUnbondingDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.OperatorId == 0 {
		return nil, status.Error(codes.InvalidArgument, "operator id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the operator unbonding delegations store
	store := ctx.KVStore(k.storeKey)
	delegationsStore := prefix.NewStore(store, types.OperatorUnbondingDelegationPrefix)

	// Query the operator unbonding delegations for the given pool id
	var unbondingDelegations []types.UnbondingDelegation
	pageRes, err := query.Paginate(delegationsStore, req.Pagination, func(key []byte, value []byte) error {
		unbond, err := types.UnmarshalUnbondingDelegation(k.cdc, value)
		if err != nil {
			return err
		}
		unbondingDelegations = append(unbondingDelegations, unbond)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOperatorUnbondingDelegationsResponse{
		UnbondingDelegations: unbondingDelegations,
		Pagination:           pageRes,
	}, nil
}

// OperatorUnbondingDelegation queries the operator unbonding delegation for the given operator id and user address
func (k Querier) OperatorUnbondingDelegation(goCtx context.Context, req *types.QueryOperatorUnbondingDelegationRequest) (*types.QueryOperatorUnbondingDelegationResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Errorf(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.OperatorId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "operator id cannot be zero")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	unbond, found := k.GetOperatorUnbondingDelegation(ctx, req.OperatorId, req.DelegatorAddress)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"unbonding delegation with delegator %s not found for operator %d",
			req.DelegatorAddress, req.OperatorId,
		)
	}

	return &types.QueryOperatorUnbondingDelegationResponse{
		UnbondingDelegation: unbond,
	}, nil
}

// ServiceDelegations queries the service delegations for the given service id
func (k Querier) ServiceDelegations(goCtx context.Context, req *types.QueryServiceDelegationsRequest) (*types.QueryServiceDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "service id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the service delegations store
	store := ctx.KVStore(k.storeKey)
	delegationsStore := prefix.NewStore(store, types.ServiceDelegationPrefix)

	// Query the service delegations for the given pool id
	delegations, pageRes, err := query.GenericFilteredPaginate(k.cdc, delegationsStore, req.Pagination, func(key []byte, delegation *types.Delegation) (*types.Delegation, error) {
		if delegation.TargetID != req.ServiceId {
			return nil, nil
		}
		return delegation, nil
	}, func() *types.Delegation {
		return &types.Delegation{}
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	serviceDelegationResponses := make([]types.DelegationResponse, len(delegations))
	for i, delegation := range delegations {
		response, err := ServiceDelegationToServiceDelegationResponse(ctx, k.Keeper, *delegation)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		serviceDelegationResponses[i] = response
	}

	return &types.QueryServiceDelegationsResponse{
		Delegations: serviceDelegationResponses,
		Pagination:  pageRes,
	}, nil
}

// ServiceDelegation queries the service delegation for the given service id and user address
func (k Querier) ServiceDelegation(goCtx context.Context, req *types.QueryServiceDelegationRequest) (*types.QueryServiceDelegationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "service id cannot be zero")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	serviceDelegation, found := k.GetServiceDelegation(ctx, req.ServiceId, req.DelegatorAddress)
	if !found {
		return nil, status.Error(codes.NotFound, "pool delegation not found")
	}

	response, err := ServiceDelegationToServiceDelegationResponse(ctx, k.Keeper, serviceDelegation)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryServiceDelegationResponse{
		Delegation: response,
	}, nil
}

// ServiceUnbondingDelegations queries the service unbonding delegations for the given service id
func (k Querier) ServiceUnbondingDelegations(goCtx context.Context, req *types.QueryServiceUnbondingDelegationsRequest) (*types.QueryServiceUnbondingDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "service id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the service unbonding delegations store
	store := ctx.KVStore(k.storeKey)
	delegationsStore := prefix.NewStore(store, types.ServiceUnbondingDelegationPrefix)

	// Query the service unbonding delegations for the given pool id
	var unbondingDelegations []types.UnbondingDelegation
	pageRes, err := query.Paginate(delegationsStore, req.Pagination, func(key []byte, value []byte) error {
		unbond, err := types.UnmarshalUnbondingDelegation(k.cdc, value)
		if err != nil {
			return err
		}
		unbondingDelegations = append(unbondingDelegations, unbond)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryServiceUnbondingDelegationsResponse{
		UnbondingDelegations: unbondingDelegations,
		Pagination:           pageRes,
	}, nil
}

// ServiceUnbondingDelegation queries the service unbonding delegation for the given service id and user address
func (k Querier) ServiceUnbondingDelegation(goCtx context.Context, req *types.QueryServiceUnbondingDelegationRequest) (*types.QueryServiceUnbondingDelegationResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Errorf(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.ServiceId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "service id cannot be zero")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	unbond, found := k.GetServiceUnbondingDelegation(ctx, req.ServiceId, req.DelegatorAddress)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"unbonding delegation with delegator %s not found for service %d",
			req.DelegatorAddress, req.ServiceId,
		)
	}

	return &types.QueryServiceUnbondingDelegationResponse{
		UnbondingDelegation: unbond,
	}, nil
}

// DelegatorPoolDelegations queries the pool delegations for the given delegator address
func (k Querier) DelegatorPoolDelegations(goCtx context.Context, req *types.QueryDelegatorPoolDelegationsRequest) (*types.QueryDelegatorPoolDelegationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the user pool delegations store
	store := ctx.KVStore(k.storeKey)
	delStore := prefix.NewStore(store, types.UserPoolDelegationsStorePrefix(req.DelegatorAddress))

	// Get the delegations
	var delegations []types.Delegation
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return err
		}
		delegations = append(delegations, delegation)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	delegationsResponses := make([]types.DelegationResponse, len(delegations))
	for i, delegation := range delegations {
		response, err := PoolDelegationToPoolDelegationResponse(ctx, k.Keeper, delegation)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		delegationsResponses[i] = response
	}

	return &types.QueryDelegatorPoolDelegationsResponse{
		Delegations: delegationsResponses,
		Pagination:  pageRes,
	}, nil
}

// DelegatorPoolUnbondingDelegations queries the pool unbonding delegations for the given delegator address
func (k Querier) DelegatorPoolUnbondingDelegations(goCtx context.Context, req *types.QueryDelegatorPoolUnbondingDelegationsRequest) (*types.QueryDelegatorPoolUnbondingDelegationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the user pool unbonding delegations store
	store := ctx.KVStore(k.storeKey)
	delStore := prefix.NewStore(store, types.PoolUnbondingDelegationsStorePrefix(req.DelegatorAddress))

	// Get the unbonding delegations
	var unbondingDelegations []types.UnbondingDelegation
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		unbonding, err := types.UnmarshalUnbondingDelegation(k.cdc, value)
		if err != nil {
			return err
		}
		unbondingDelegations = append(unbondingDelegations, unbonding)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorPoolUnbondingDelegationsResponse{
		UnbondingDelegations: unbondingDelegations,
		Pagination:           pageRes,
	}, nil
}

// DelegatorOperatorDelegations queries the operator delegations for the given delegator address
func (k Querier) DelegatorOperatorDelegations(goCtx context.Context, req *types.QueryDelegatorOperatorDelegationsRequest) (*types.QueryDelegatorOperatorDelegationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the user operator delegations store
	store := ctx.KVStore(k.storeKey)
	delStore := prefix.NewStore(store, types.UserOperatorDelegationsStorePrefix(req.DelegatorAddress))

	// Get the delegations
	var delegations []types.Delegation
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return err
		}
		delegations = append(delegations, delegation)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	delegationsResponses := make([]types.DelegationResponse, len(delegations))
	for i, delegation := range delegations {
		response, err := OperatorDelegationToOperatorDelegationResponse(ctx, k.Keeper, delegation)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		delegationsResponses[i] = response
	}

	return &types.QueryDelegatorOperatorDelegationsResponse{
		Delegations: delegationsResponses,
		Pagination:  pageRes,
	}, nil
}

// DelegatorOperatorUnbondingDelegations queries the operator unbonding delegations for the given delegator address
func (k Querier) DelegatorOperatorUnbondingDelegations(goCtx context.Context, req *types.QueryDelegatorOperatorUnbondingDelegationsRequest) (*types.QueryDelegatorOperatorUnbondingDelegationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the user operator unbonding delegations store
	store := ctx.KVStore(k.storeKey)
	delStore := prefix.NewStore(store, types.OperatorUnbondingDelegationsStorePrefix(req.DelegatorAddress))

	// Get the unbonding delegations
	var unbondingDelegations []types.UnbondingDelegation
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		unbonding, err := types.UnmarshalUnbondingDelegation(k.cdc, value)
		if err != nil {
			return err
		}
		unbondingDelegations = append(unbondingDelegations, unbonding)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorOperatorUnbondingDelegationsResponse{
		UnbondingDelegations: unbondingDelegations,
		Pagination:           pageRes,
	}, nil
}

// DelegatorServiceDelegations queries the service delegations for the given delegator address
func (k Querier) DelegatorServiceDelegations(goCtx context.Context, req *types.QueryDelegatorServiceDelegationsRequest) (*types.QueryDelegatorServiceDelegationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the user services delegations store
	store := ctx.KVStore(k.storeKey)
	delStore := prefix.NewStore(store, types.UserServiceDelegationsStorePrefix(req.DelegatorAddress))

	// Get the delegations
	var delegations []types.Delegation
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return err
		}
		delegations = append(delegations, delegation)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	delegationsResponses := make([]types.DelegationResponse, len(delegations))
	for i, delegation := range delegations {
		response, err := ServiceDelegationToServiceDelegationResponse(ctx, k.Keeper, delegation)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		delegationsResponses[i] = response
	}

	return &types.QueryDelegatorServiceDelegationsResponse{
		Delegations: delegationsResponses,
		Pagination:  pageRes,
	}, nil
}

// DelegatorServiceUnbondingDelegations queries the service unbonding delegations for the given delegator address
func (k Querier) DelegatorServiceUnbondingDelegations(goCtx context.Context, req *types.QueryDelegatorServiceUnbondingDelegationsRequest) (*types.QueryDelegatorServiceUnbondingDelegationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the user services unbonding delegations store
	store := ctx.KVStore(k.storeKey)
	delStore := prefix.NewStore(store, types.ServiceUnbondingDelegationsStorePrefix(req.DelegatorAddress))

	// Get the unbonding delegations
	var unbondingDelegations []types.UnbondingDelegation
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		unbonding, err := types.UnmarshalUnbondingDelegation(k.cdc, value)
		if err != nil {
			return err
		}
		unbondingDelegations = append(unbondingDelegations, unbonding)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorServiceUnbondingDelegationsResponse{
		UnbondingDelegations: unbondingDelegations,
		Pagination:           pageRes,
	}, nil
}

// DelegatorPools queries the pools for the given delegator address
func (k Querier) DelegatorPools(goCtx context.Context, req *types.QueryDelegatorPoolsRequest) (*types.QueryDelegatorPoolsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the user pools delegations store
	store := ctx.KVStore(k.storeKey)
	delStore := prefix.NewStore(store, types.UserPoolDelegationsStorePrefix(req.DelegatorAddress))

	// Get the pools
	var pools []poolstypes.Pool
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return err
		}

		pool, found := k.poolsKeeper.GetPool(ctx, delegation.TargetID)
		if !found {
			return poolstypes.ErrPoolNotFound
		}

		pools = append(pools, pool)

		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorPoolsResponse{
		Pools:      pools,
		Pagination: pageRes,
	}, nil
}

// DelegatorPool queries the pool for the given delegator address and pool id
func (k Querier) DelegatorPool(goCtx context.Context, req *types.QueryDelegatorPoolRequest) (*types.QueryDelegatorPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.PoolId == 0 {
		return nil, status.Error(codes.InvalidArgument, "pool id cannot be zero")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	delegation, found := k.GetPoolDelegation(ctx, req.PoolId, req.DelegatorAddress)
	if !found {
		return nil, status.Error(codes.NotFound, "pool delegation not found")
	}

	pool, found := k.poolsKeeper.GetPool(ctx, delegation.TargetID)
	if !found {
		return nil, status.Error(codes.NotFound, "pool not found")
	}

	return &types.QueryDelegatorPoolResponse{
		Pool: pool,
	}, nil
}

// DelegatorOperators queries the operators for the given delegator address
func (k Querier) DelegatorOperators(goCtx context.Context, req *types.QueryDelegatorOperatorsRequest) (*types.QueryDelegatorOperatorsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the user operators delegations store
	store := ctx.KVStore(k.storeKey)
	delStore := prefix.NewStore(store, types.UserOperatorDelegationsStorePrefix(req.DelegatorAddress))

	// Get the operators
	var operators []operatorstypes.Operator
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return err
		}

		operator, found := k.operatorsKeeper.GetOperator(ctx, delegation.TargetID)
		if !found {
			return operatorstypes.ErrOperatorNotFound
		}

		operators = append(operators, operator)

		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorOperatorsResponse{
		Operators:  operators,
		Pagination: pageRes,
	}, nil
}

// DelegatorOperator queries the operator for the given delegator address and operator id
func (k Querier) DelegatorOperator(goCtx context.Context, req *types.QueryDelegatorOperatorRequest) (*types.QueryDelegatorOperatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.OperatorId == 0 {
		return nil, status.Error(codes.InvalidArgument, "operator id cannot be zero")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	delegation, found := k.GetOperatorDelegation(ctx, req.OperatorId, req.DelegatorAddress)
	if !found {
		return nil, status.Error(codes.NotFound, "operator delegation not found")
	}

	operator, found := k.operatorsKeeper.GetOperator(ctx, delegation.TargetID)
	if !found {
		return nil, status.Error(codes.NotFound, "operator not found")
	}

	return &types.QueryDelegatorOperatorResponse{
		Operator: operator,
	}, nil
}

// DelegatorServices queries the services for the given delegator address
func (k Querier) DelegatorServices(goCtx context.Context, req *types.QueryDelegatorServicesRequest) (*types.QueryDelegatorServicesResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the user services delegations store
	store := ctx.KVStore(k.storeKey)
	delStore := prefix.NewStore(store, types.UserServiceDelegationsStorePrefix(req.DelegatorAddress))

	// Get the services
	var services []servicestypes.Service
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return err
		}

		pool, found := k.servicesKeeper.GetService(ctx, delegation.TargetID)
		if !found {
			return servicestypes.ErrServiceNotFound
		}

		services = append(services, pool)

		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorServicesResponse{
		Services:   services,
		Pagination: pageRes,
	}, nil
}

// DelegatorService queries the service for the given delegator address and service id
func (k Querier) DelegatorService(goCtx context.Context, req *types.QueryDelegatorServiceRequest) (*types.QueryDelegatorServiceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "service id cannot be zero")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	delegation, found := k.GetServiceDelegation(ctx, req.ServiceId, req.DelegatorAddress)
	if !found {
		return nil, status.Error(codes.NotFound, "service delegation not found")
	}

	service, found := k.servicesKeeper.GetService(ctx, delegation.TargetID)
	if !found {
		return nil, status.Error(codes.NotFound, "service not found")
	}

	return &types.QueryDelegatorServiceResponse{
		Service: service,
	}, nil
}

// Params queries the restaking module parameters
func (k Querier) Params(goCtx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

// --------------------------------------------------------------------------------------------------------------------

// PoolDelegationToPoolDelegationResponse converts a PoolDelegation to a PoolDelegationResponse
func PoolDelegationToPoolDelegationResponse(ctx sdk.Context, k *Keeper, delegation types.Delegation) (types.DelegationResponse, error) {
	pool, found := k.poolsKeeper.GetPool(ctx, delegation.TargetID)
	if !found {
		return types.DelegationResponse{}, poolstypes.ErrPoolNotFound
	}

	truncatedBalance, _ := pool.TokensFromShares(delegation.Shares).TruncateDecimal()
	return types.NewDelegationResponse(delegation, truncatedBalance), nil
}

// OperatorDelegationToOperatorDelegationResponse converts a OperatorDelegation to a OperatorDelegationResponse
func OperatorDelegationToOperatorDelegationResponse(ctx sdk.Context, k *Keeper, delegation types.Delegation) (types.DelegationResponse, error) {
	operator, found := k.operatorsKeeper.GetOperator(ctx, delegation.TargetID)
	if !found {
		return types.DelegationResponse{}, operatorstypes.ErrOperatorNotFound
	}

	truncatedBalance, _ := operator.TokensFromShares(delegation.Shares).TruncateDecimal()
	return types.NewDelegationResponse(delegation, truncatedBalance), nil
}

// ServiceDelegationToServiceDelegationResponse converts a ServiceDelegation to a ServiceDelegationResponse
func ServiceDelegationToServiceDelegationResponse(ctx sdk.Context, k *Keeper, delegation types.Delegation) (types.DelegationResponse, error) {
	service, found := k.servicesKeeper.GetService(ctx, delegation.TargetID)
	if !found {
		return types.DelegationResponse{}, servicestypes.ErrServiceNotFound
	}

	truncatedBalance, _ := service.TokensFromShares(delegation.Shares).TruncateDecimal()
	return types.NewDelegationResponse(delegation, truncatedBalance), nil
}
