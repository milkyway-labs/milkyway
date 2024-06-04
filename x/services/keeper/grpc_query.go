package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

var _ types.QueryServer = &Keeper{}

// Services implements the Query/Services gRPC method
func (k Keeper) Services(ctx context.Context, request *types.QueryServicesRequest) (*types.QueryServicesResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := sdkCtx.KVStore(k.storeKey)
	servicesStore := prefix.NewStore(store, types.AVSPrefix)

	var services []types.Service
	pageRes, err := query.Paginate(servicesStore, request.Pagination, func(key []byte, value []byte) error {
		var service types.Service
		if err := k.cdc.Unmarshal(value, &service); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		services = append(services, service)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryServicesResponse{
		Services:   services,
		Pagination: pageRes,
	}, nil
}

// Service implements the Query/Service gRPC method
func (k Keeper) Service(ctx context.Context, request *types.QueryServiceRequest) (*types.QueryServiceResponse, error) {
	if request.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid service ID")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	service, found := k.GetAVS(sdkCtx, request.ServiceId)
	if !found {
		return nil, status.Error(codes.NotFound, "service not found")
	}

	return &types.QueryServiceResponse{Service: service}, nil
}

// Params implements the Query/Params gRPC method
func (k Keeper) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(sdkCtx)
	return &types.QueryParamsResponse{Params: params}, nil
}
