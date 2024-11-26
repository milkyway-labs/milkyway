package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/milkyway-labs/milkyway/v2/x/pools/types"
)

var _ types.QueryServer = &Keeper{}

// PoolByID implements the Query/PoolById gRPC method
func (k *Keeper) PoolByID(ctx context.Context, request *types.QueryPoolByIdRequest) (*types.QueryPoolResponse, error) {
	if request.PoolId == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid pool id")
	}

	pool, found, err := k.GetPool(ctx, request.PoolId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !found {
		return nil, status.Error(codes.NotFound, "pool not found")
	}

	return &types.QueryPoolResponse{Pool: pool}, nil
}

// PoolByDenom implements the Query/PoolByDenom gRPC method
func (k *Keeper) PoolByDenom(ctx context.Context, request *types.QueryPoolByDenomRequest) (*types.QueryPoolResponse, error) {
	if err := sdk.ValidateDenom(request.Denom); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid denom")
	}

	pool, found, err := k.GetPoolByDenom(ctx, request.Denom)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !found {
		return nil, status.Error(codes.NotFound, "pool not found")
	}

	return &types.QueryPoolResponse{Pool: pool}, nil
}

// Pools implements the Query/Pools gRPC method
func (k *Keeper) Pools(ctx context.Context, request *types.QueryPoolsRequest) (*types.QueryPoolsResponse, error) {
	pools, pageRes, err := query.CollectionPaginate(ctx, k.pools, request.Pagination, func(key uint32, value types.Pool) (types.Pool, error) {
		return value, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPoolsResponse{
		Pools:      pools,
		Pagination: pageRes,
	}, nil
}
