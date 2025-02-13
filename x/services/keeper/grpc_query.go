package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/milkyway-labs/milkyway/v9/x/services/types"
)

var _ types.QueryServer = &Keeper{}

// Service implements the Query/Service gRPC method
func (k *Keeper) Service(ctx context.Context, request *types.QueryServiceRequest) (*types.QueryServiceResponse, error) {
	if request.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid service ID")
	}

	service, err := k.GetService(ctx, request.ServiceId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "service not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryServiceResponse{Service: service}, nil
}

// Services implements the Query/Services gRPC method
func (k *Keeper) Services(ctx context.Context, request *types.QueryServicesRequest) (*types.QueryServicesResponse, error) {
	services, pageRes, err := query.CollectionPaginate(ctx, k.services, request.Pagination,
		func(key uint32, value types.Service) (types.Service, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.QueryServicesResponse{
		Services:   services,
		Pagination: pageRes,
	}, nil
}

// ServiceParams implements the Query/ServiceParams gRPC method
func (k *Keeper) ServiceParams(ctx context.Context, request *types.QueryServiceParamsRequest) (*types.QueryServiceParamsResponse, error) {
	if request.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid service ID")
	}

	// Ensure the service exists
	_, err := k.GetService(ctx, request.ServiceId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "service not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Get the service params
	serviceParams, err := k.GetServiceParams(ctx, request.ServiceId)
	if err != nil {
		return nil, err
	}

	return &types.QueryServiceParamsResponse{ServiceParams: serviceParams}, nil
}

// Params implements the Query/Params gRPC method
func (k *Keeper) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryParamsResponse{Params: params}, nil
}
