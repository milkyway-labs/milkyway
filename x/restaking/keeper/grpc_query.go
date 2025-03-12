package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	operatorstypes "github.com/milkyway-labs/milkyway/v10/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v10/x/pools/types"
	"github.com/milkyway-labs/milkyway/v10/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v10/x/services/types"
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
func (k Querier) OperatorJoinedServices(ctx context.Context, req *types.QueryOperatorJoinedServicesRequest) (*types.QueryOperatorJoinedServicesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.OperatorId == 0 {
		return nil, status.Error(codes.InvalidArgument, "operator id cannot be 0")
	}

	_, err := k.operatorsKeeper.GetOperator(ctx, req.OperatorId)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "operator with id %d not found", req.OperatorId)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Get the operator joined services
	serviceIDs, pageResponse, err := query.CollectionPaginate(ctx, k.operatorJoinedServices, req.Pagination,
		func(key collections.Pair[uint32, uint32], _ collections.NoValue) (uint32, error) {
			return key.K2(), nil
		}, query.WithCollectionPaginationPairPrefix[uint32, uint32](req.OperatorId))
	if err != nil {
		return nil, err
	}

	return &types.QueryOperatorJoinedServicesResponse{ServiceIds: serviceIDs, Pagination: pageResponse}, nil
}

// ServiceAllowedOperators queries the allowed operators for a given service.
func (k Querier) ServiceAllowedOperators(
	ctx context.Context,
	req *types.QueryServiceAllowedOperatorsRequest,
) (*types.QueryServiceAllowedOperatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "service id cannot be 0")
	}

	// Get the service's allowed operators
	operatorIDs, pageResponse, err := query.CollectionPaginate(ctx, k.serviceOperatorsAllowList, req.Pagination,
		func(key collections.Pair[uint32, uint32], _ collections.NoValue) (uint32, error) {
			return key.K2(), nil
		}, query.WithCollectionPaginationPairPrefix[uint32, uint32](req.ServiceId))
	if err != nil {
		return nil, err
	}

	return &types.QueryServiceAllowedOperatorsResponse{
		OperatorIds: operatorIDs,
		Pagination:  pageResponse,
	}, nil
}

// ServiceSecuringPools queries the pools that are securing a given service.
func (k Querier) ServiceSecuringPools(
	ctx context.Context,
	req *types.QueryServiceSecuringPoolsRequest,
) (*types.QueryServiceSecuringPoolsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "service id cannot be 0")
	}

	// Get the service securing pools
	poolIDs, pageResponse, err := query.CollectionPaginate(ctx, k.serviceSecuringPools, req.Pagination,
		func(key collections.Pair[uint32, uint32], _ collections.NoValue) (uint32, error) {
			return key.K2(), nil
		}, query.WithCollectionPaginationPairPrefix[uint32, uint32](req.ServiceId))
	if err != nil {
		return nil, err
	}

	return &types.QueryServiceSecuringPoolsResponse{
		PoolIds:    poolIDs,
		Pagination: pageResponse,
	}, nil
}

func (k Querier) ServiceOperators(ctx context.Context, req *types.QueryServiceOperatorsRequest) (*types.QueryServiceOperatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "service id cannot be 0")
	}

	_, err := k.servicesKeeper.GetService(ctx, req.ServiceId)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "service with id %d not found", req.ServiceId)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	eligibleOperators, pageResponse, err := query.CollectionFilteredPaginate(ctx, k.operatorJoinedServices.Indexes.Service, req.Pagination,
		// Filter to return only the operators that have joined the service and
		// that are allowed to validate it
		func(key collections.Pair[uint32, uint32], _ collections.NoValue) (bool, error) {
			// Here is k2 the operator id since the Service index provides association
			// between a service and the operator securing it
			operatorID := key.K2()
			return k.CanOperatorValidateService(ctx, operatorID, req.ServiceId)
		},
		func(key collections.Pair[uint32, uint32], _ collections.NoValue) (operatorstypes.Operator, error) {
			// Here is k2 the operator id since the Service index provides association
			// between a service and the operator securing it
			operatorID := key.K2()
			operator, err := k.operatorsKeeper.GetOperator(ctx, operatorID)
			if err != nil {
				if errors.IsOf(err, collections.ErrNotFound) {
					return operatorstypes.Operator{}, status.Errorf(codes.NotFound, "operator %d not found", operatorID)
				}
				return operatorstypes.Operator{}, err
			}
			return operator, nil
		}, query.WithCollectionPaginationPairPrefix[uint32, uint32](req.ServiceId))
	if err != nil {
		return nil, err
	}

	return &types.QueryServiceOperatorsResponse{Operators: eligibleOperators, Pagination: pageResponse}, nil
}

// PoolDelegations queries the pool delegations for the given pool id
func (k Querier) PoolDelegations(ctx context.Context, req *types.QueryPoolDelegationsRequest) (*types.QueryPoolDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.PoolId == 0 {
		return nil, status.Error(codes.InvalidArgument, "pool id cannot be 0")
	}

	// Get the pool delegations store
	store := k.storeService.OpenKVStore(ctx)
	keyPrefix := types.DelegationsByPoolIDStorePrefix(req.PoolId)
	indexStore := prefix.NewStore(runtime.KVStoreAdapter(store), keyPrefix)

	// Query the pool delegations for the given pool id
	var delegations []types.Delegation
	pageRes, err := query.Paginate(indexStore, req.Pagination, func(key, _ []byte) error {
		_, userAddress := types.ParseDelegationByPoolIDStoreKey(append(keyPrefix, key...))
		delegation, found, err := k.GetPoolDelegation(ctx, req.PoolId, userAddress)
		if err != nil {
			return err
		}
		if !found {
			return types.ErrDelegationNotFound
		}
		delegations = append(delegations, delegation)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	poolDelegations := make([]types.DelegationResponse, len(delegations))
	for i, delegation := range delegations {
		response, err := PoolDelegationToPoolDelegationResponse(ctx, k.Keeper, delegation)
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
func (k Querier) PoolDelegation(ctx context.Context, req *types.QueryPoolDelegationRequest) (*types.QueryPoolDelegationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.PoolId == 0 {
		return nil, status.Error(codes.InvalidArgument, "pool id cannot be zero")
	}

	poolDelegation, found, err := k.GetPoolDelegation(ctx, req.PoolId, req.DelegatorAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

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
func (k Querier) PoolUnbondingDelegations(ctx context.Context, req *types.QueryPoolUnbondingDelegationsRequest) (*types.QueryPoolUnbondingDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.PoolId == 0 {
		return nil, status.Error(codes.InvalidArgument, "pool id cannot be 0")
	}

	// Get the pool unbonding delegations store
	store := k.storeService.OpenKVStore(ctx)
	delegationsStore := prefix.NewStore(runtime.KVStoreAdapter(store), types.PoolUnbondingDelegationPrefix)

	// Query the pool unbonding delegations for the given pool id
	unbondingDelegations, pageRes, err := query.GenericFilteredPaginate(k.cdc, delegationsStore, req.Pagination, func(_ []byte, unbond *types.UnbondingDelegation) (*types.UnbondingDelegation, error) {
		if unbond.TargetID != req.PoolId {
			return nil, nil
		}
		return unbond, nil
	}, func() *types.UnbondingDelegation {
		return &types.UnbondingDelegation{}
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	unbondingDels := make([]types.UnbondingDelegation, len(unbondingDelegations))
	for i, unbond := range unbondingDelegations {
		unbondingDels[i] = *unbond
	}

	return &types.QueryPoolUnbondingDelegationsResponse{
		UnbondingDelegations: unbondingDels,
		Pagination:           pageRes,
	}, nil
}

// PoolUnbondingDelegation queries the pool unbonding delegation for the given pool id and user address
func (k Querier) PoolUnbondingDelegation(ctx context.Context, req *types.QueryPoolUnbondingDelegationRequest) (*types.QueryPoolUnbondingDelegationResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Errorf(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.PoolId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "pool id cannot be zero")
	}

	unbond, found, err := k.GetPoolUnbondingDelegation(ctx, req.PoolId, req.DelegatorAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

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
func (k Querier) OperatorDelegations(ctx context.Context, req *types.QueryOperatorDelegationsRequest) (*types.QueryOperatorDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.OperatorId == 0 {
		return nil, status.Error(codes.InvalidArgument, "operator id cannot be 0")
	}

	// Get the operator delegations store
	store := k.storeService.OpenKVStore(ctx)
	keyPrefix := types.DelegationsByOperatorIDStorePrefix(req.OperatorId)
	indexStore := prefix.NewStore(runtime.KVStoreAdapter(store), keyPrefix)

	// Query the operator delegations for the given operator id
	var delegations []types.Delegation
	pageRes, err := query.Paginate(indexStore, req.Pagination, func(key, _ []byte) error {
		_, userAddress := types.ParseDelegationByOperatorIDStoreKey(append(keyPrefix, key...))
		delegation, found, err := k.GetOperatorDelegation(ctx, req.OperatorId, userAddress)
		if err != nil {
			return err
		}
		if !found {
			return types.ErrDelegationNotFound
		}
		delegations = append(delegations, delegation)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	operatorDelegations := make([]types.DelegationResponse, len(delegations))
	for i, delegation := range delegations {
		response, err := OperatorDelegationToOperatorDelegationResponse(ctx, k.Keeper, delegation)
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
func (k Querier) OperatorDelegation(ctx context.Context, req *types.QueryOperatorDelegationRequest) (*types.QueryOperatorDelegationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.OperatorId == 0 {
		return nil, status.Error(codes.InvalidArgument, "operator id cannot be zero")
	}

	operatorDelegation, found, err := k.GetOperatorDelegation(ctx, req.OperatorId, req.DelegatorAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

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
func (k Querier) OperatorUnbondingDelegations(ctx context.Context, req *types.QueryOperatorUnbondingDelegationsRequest) (*types.QueryOperatorUnbondingDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.OperatorId == 0 {
		return nil, status.Error(codes.InvalidArgument, "operator id cannot be 0")
	}

	// Get the operator unbonding delegations store
	store := k.storeService.OpenKVStore(ctx)
	delegationsStore := prefix.NewStore(runtime.KVStoreAdapter(store), types.OperatorUnbondingDelegationPrefix)

	// Query the operator unbonding delegations for the given operator id
	unbondingDelegations, pageRes, err := query.GenericFilteredPaginate(k.cdc, delegationsStore, req.Pagination, func(_ []byte, unbond *types.UnbondingDelegation) (*types.UnbondingDelegation, error) {
		if unbond.TargetID != req.OperatorId {
			return nil, nil
		}
		return unbond, nil
	}, func() *types.UnbondingDelegation {
		return &types.UnbondingDelegation{}
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	unbondingDels := make([]types.UnbondingDelegation, len(unbondingDelegations))
	for i, unbond := range unbondingDelegations {
		unbondingDels[i] = *unbond
	}

	return &types.QueryOperatorUnbondingDelegationsResponse{
		UnbondingDelegations: unbondingDels,
		Pagination:           pageRes,
	}, nil
}

// OperatorUnbondingDelegation queries the operator unbonding delegation for the given operator id and user address
func (k Querier) OperatorUnbondingDelegation(ctx context.Context, req *types.QueryOperatorUnbondingDelegationRequest) (*types.QueryOperatorUnbondingDelegationResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Errorf(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.OperatorId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "operator id cannot be zero")
	}

	unbond, found, err := k.GetOperatorUnbondingDelegation(ctx, req.OperatorId, req.DelegatorAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

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
func (k Querier) ServiceDelegations(ctx context.Context, req *types.QueryServiceDelegationsRequest) (*types.QueryServiceDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "service id cannot be 0")
	}

	// Get the service delegations store
	store := k.storeService.OpenKVStore(ctx)
	keyPrefix := types.DelegationsByServiceIDStorePrefix(req.ServiceId)
	indexStore := prefix.NewStore(runtime.KVStoreAdapter(store), keyPrefix)

	// Query the service delegations for the given service id
	var delegations []types.Delegation
	pageRes, err := query.Paginate(indexStore, req.Pagination, func(key, _ []byte) error {
		_, userAddress := types.ParseDelegationByServiceIDStoreKey(append(keyPrefix, key...))
		delegation, found, err := k.GetServiceDelegation(ctx, req.ServiceId, userAddress)
		if err != nil {
			return err
		}
		if !found {
			return types.ErrDelegationNotFound
		}
		delegations = append(delegations, delegation)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	serviceDelegationResponses := make([]types.DelegationResponse, len(delegations))
	for i, delegation := range delegations {
		response, err := ServiceDelegationToServiceDelegationResponse(ctx, k.Keeper, delegation)
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
func (k Querier) ServiceDelegation(ctx context.Context, req *types.QueryServiceDelegationRequest) (*types.QueryServiceDelegationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "service id cannot be zero")
	}

	serviceDelegation, found, err := k.GetServiceDelegation(ctx, req.ServiceId, req.DelegatorAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

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
func (k Querier) ServiceUnbondingDelegations(ctx context.Context, req *types.QueryServiceUnbondingDelegationsRequest) (*types.QueryServiceUnbondingDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "service id cannot be 0")
	}

	// Get the service unbonding delegations store
	store := k.storeService.OpenKVStore(ctx)
	delegationsStore := prefix.NewStore(runtime.KVStoreAdapter(store), types.ServiceUnbondingDelegationPrefix)

	// Query the service unbonding delegations for the given service id
	unbondingDelegations, pageRes, err := query.GenericFilteredPaginate(k.cdc, delegationsStore, req.Pagination, func(_ []byte, unbond *types.UnbondingDelegation) (*types.UnbondingDelegation, error) {
		if unbond.TargetID != req.ServiceId {
			return nil, nil
		}
		return unbond, nil
	}, func() *types.UnbondingDelegation {
		return &types.UnbondingDelegation{}
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	unbondingDels := make([]types.UnbondingDelegation, len(unbondingDelegations))
	for i, unbond := range unbondingDelegations {
		unbondingDels[i] = *unbond
	}

	return &types.QueryServiceUnbondingDelegationsResponse{
		UnbondingDelegations: unbondingDels,
		Pagination:           pageRes,
	}, nil
}

// ServiceUnbondingDelegation queries the service unbonding delegation for the given service id and user address
func (k Querier) ServiceUnbondingDelegation(ctx context.Context, req *types.QueryServiceUnbondingDelegationRequest) (*types.QueryServiceUnbondingDelegationResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Errorf(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.ServiceId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "service id cannot be zero")
	}

	unbond, found, err := k.GetServiceUnbondingDelegation(ctx, req.ServiceId, req.DelegatorAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

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
func (k Querier) DelegatorPoolDelegations(ctx context.Context, req *types.QueryDelegatorPoolDelegationsRequest) (*types.QueryDelegatorPoolDelegationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	// Get the user pool delegations store
	store := k.storeService.OpenKVStore(ctx)
	delStore := prefix.NewStore(runtime.KVStoreAdapter(store), types.UserPoolDelegationsStorePrefix(req.DelegatorAddress))

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
func (k Querier) DelegatorPoolUnbondingDelegations(ctx context.Context, req *types.QueryDelegatorPoolUnbondingDelegationsRequest) (*types.QueryDelegatorPoolUnbondingDelegationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	// Get the user pool unbonding delegations store
	store := k.storeService.OpenKVStore(ctx)
	delStore := prefix.NewStore(runtime.KVStoreAdapter(store), types.PoolUnbondingDelegationsStorePrefix(req.DelegatorAddress))

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
func (k Querier) DelegatorOperatorDelegations(ctx context.Context, req *types.QueryDelegatorOperatorDelegationsRequest) (*types.QueryDelegatorOperatorDelegationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	// Get the user operator delegations store
	store := k.storeService.OpenKVStore(ctx)
	delStore := prefix.NewStore(runtime.KVStoreAdapter(store), types.UserOperatorDelegationsStorePrefix(req.DelegatorAddress))

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
func (k Querier) DelegatorOperatorUnbondingDelegations(ctx context.Context, req *types.QueryDelegatorOperatorUnbondingDelegationsRequest) (*types.QueryDelegatorOperatorUnbondingDelegationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	// Get the user operator unbonding delegations store
	store := k.storeService.OpenKVStore(ctx)
	delStore := prefix.NewStore(runtime.KVStoreAdapter(store), types.OperatorUnbondingDelegationsStorePrefix(req.DelegatorAddress))

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
func (k Querier) DelegatorServiceDelegations(ctx context.Context, req *types.QueryDelegatorServiceDelegationsRequest) (*types.QueryDelegatorServiceDelegationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	// Get the user services delegations store
	store := k.storeService.OpenKVStore(ctx)
	delStore := prefix.NewStore(runtime.KVStoreAdapter(store), types.UserServiceDelegationsStorePrefix(req.DelegatorAddress))

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
func (k Querier) DelegatorServiceUnbondingDelegations(ctx context.Context, req *types.QueryDelegatorServiceUnbondingDelegationsRequest) (*types.QueryDelegatorServiceUnbondingDelegationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	// Get the user services unbonding delegations store
	store := k.storeService.OpenKVStore(ctx)
	delStore := prefix.NewStore(runtime.KVStoreAdapter(store), types.ServiceUnbondingDelegationsStorePrefix(req.DelegatorAddress))

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
func (k Querier) DelegatorPools(ctx context.Context, req *types.QueryDelegatorPoolsRequest) (*types.QueryDelegatorPoolsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	// Get the user pools delegations store
	store := k.storeService.OpenKVStore(ctx)
	delStore := prefix.NewStore(runtime.KVStoreAdapter(store), types.UserPoolDelegationsStorePrefix(req.DelegatorAddress))

	// Get the pools
	var pools []poolstypes.Pool
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return err
		}

		pool, err := k.poolsKeeper.GetPool(ctx, delegation.TargetID)
		if err != nil {
			if errors.IsOf(err, collections.ErrNotFound) {
				return status.Errorf(codes.NotFound, "pool with id %d not found", delegation.TargetID)
			}
			return err
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
func (k Querier) DelegatorPool(ctx context.Context, req *types.QueryDelegatorPoolRequest) (*types.QueryDelegatorPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.PoolId == 0 {
		return nil, status.Error(codes.InvalidArgument, "pool id cannot be zero")
	}

	delegation, found, err := k.GetPoolDelegation(ctx, req.PoolId, req.DelegatorAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !found {
		return nil, status.Error(codes.NotFound, "pool delegation not found")
	}

	pool, err := k.poolsKeeper.GetPool(ctx, delegation.TargetID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "pool with id %d not found", delegation.TargetID)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorPoolResponse{
		Pool: pool,
	}, nil
}

// DelegatorOperators queries the operators for the given delegator address
func (k Querier) DelegatorOperators(ctx context.Context, req *types.QueryDelegatorOperatorsRequest) (*types.QueryDelegatorOperatorsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	// Get the user operators delegations store
	store := k.storeService.OpenKVStore(ctx)
	delStore := prefix.NewStore(runtime.KVStoreAdapter(store), types.UserOperatorDelegationsStorePrefix(req.DelegatorAddress))

	// Get the operators
	var operators []operatorstypes.Operator
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return err
		}

		operator, err := k.operatorsKeeper.GetOperator(ctx, delegation.TargetID)
		if err != nil {
			if errors.IsOf(err, collections.ErrNotFound) {
				return status.Errorf(codes.NotFound, "operator with id %d not found", delegation.TargetID)
			}
			return err
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
func (k Querier) DelegatorOperator(ctx context.Context, req *types.QueryDelegatorOperatorRequest) (*types.QueryDelegatorOperatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.OperatorId == 0 {
		return nil, status.Error(codes.InvalidArgument, "operator id cannot be zero")
	}

	delegation, found, err := k.GetOperatorDelegation(ctx, req.OperatorId, req.DelegatorAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !found {
		return nil, status.Error(codes.NotFound, "operator delegation not found")
	}

	operator, err := k.operatorsKeeper.GetOperator(ctx, delegation.TargetID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "operator with id %d not found", delegation.TargetID)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorOperatorResponse{
		Operator: operator,
	}, nil
}

// DelegatorServices queries the services for the given delegator address
func (k Querier) DelegatorServices(ctx context.Context, req *types.QueryDelegatorServicesRequest) (*types.QueryDelegatorServicesResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}

	// Get the user services delegations store
	store := k.storeService.OpenKVStore(ctx)
	delStore := prefix.NewStore(runtime.KVStoreAdapter(store), types.UserServiceDelegationsStorePrefix(req.DelegatorAddress))

	// Get the services
	var services []servicestypes.Service
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return err
		}

		pool, err := k.servicesKeeper.GetService(ctx, delegation.TargetID)
		if err != nil {
			return err
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
func (k Querier) DelegatorService(ctx context.Context, req *types.QueryDelegatorServiceRequest) (*types.QueryDelegatorServiceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "service id cannot be zero")
	}

	delegation, found, err := k.GetServiceDelegation(ctx, req.ServiceId, req.DelegatorAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !found {
		return nil, status.Error(codes.NotFound, "service delegation not found")
	}

	service, err := k.servicesKeeper.GetService(ctx, delegation.TargetID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "service with id %d not found", delegation.TargetID)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorServiceResponse{
		Service: service,
	}, nil
}

// UserPreferences queries the user preferences for the given user address
func (k Querier) UserPreferences(ctx context.Context, req *types.QueryUserPreferencesRequest) (*types.QueryUserPreferencesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.UserAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "user address cannot be empty")
	}

	preferences, err := k.GetUserPreferences(ctx, req.UserAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryUserPreferencesResponse{
		Preferences: preferences,
	}, nil
}

// Params queries the restaking module parameters
func (k Querier) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryParamsResponse{Params: params}, nil
}

// --------------------------------------------------------------------------------------------------------------------

// PoolDelegationToPoolDelegationResponse converts a PoolDelegation to a PoolDelegationResponse
func PoolDelegationToPoolDelegationResponse(ctx context.Context, k *Keeper, delegation types.Delegation) (types.DelegationResponse, error) {
	pool, err := k.poolsKeeper.GetPool(ctx, delegation.TargetID)
	if err != nil {
		return types.DelegationResponse{}, err
	}

	truncatedBalance, _ := pool.TokensFromShares(delegation.Shares).TruncateDecimal()
	return types.NewDelegationResponse(delegation, truncatedBalance), nil
}

// OperatorDelegationToOperatorDelegationResponse converts a OperatorDelegation to a OperatorDelegationResponse
func OperatorDelegationToOperatorDelegationResponse(ctx context.Context, k *Keeper, delegation types.Delegation) (types.DelegationResponse, error) {
	operator, err := k.operatorsKeeper.GetOperator(ctx, delegation.TargetID)
	if err != nil {
		return types.DelegationResponse{}, err
	}

	truncatedBalance, _ := operator.TokensFromShares(delegation.Shares).TruncateDecimal()
	return types.NewDelegationResponse(delegation, truncatedBalance), nil
}

// ServiceDelegationToServiceDelegationResponse converts a ServiceDelegation to a ServiceDelegationResponse
func ServiceDelegationToServiceDelegationResponse(ctx context.Context, k *Keeper, delegation types.Delegation) (types.DelegationResponse, error) {
	service, err := k.servicesKeeper.GetService(ctx, delegation.TargetID)
	if err != nil {
		return types.DelegationResponse{}, err
	}

	truncatedBalance, _ := service.TokensFromShares(delegation.Shares).TruncateDecimal()
	return types.NewDelegationResponse(delegation, truncatedBalance), nil
}
