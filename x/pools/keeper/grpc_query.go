package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/milkyway-labs/milkyway/x/pools/types"
)

var _ types.QueryServer = &Keeper{}

// PoolById implements the Query/PoolById gRPC method
func (k *Keeper) PoolById(ctx context.Context, request *types.QueryPoolByIdRequest) (*types.QueryPoolResponse, error) {
	if request.PoolId == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid pool id")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	pool, found := k.GetPool(sdkCtx, request.PoolId)
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

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	pool, found := k.GetPoolByDenom(sdkCtx, request.Denom)
	if !found {
		return nil, status.Error(codes.NotFound, "pool not found")
	}

	return &types.QueryPoolResponse{Pool: pool}, nil
}

// Pools implements the Query/Pools gRPC method
func (k *Keeper) Pools(ctx context.Context, request *types.QueryPoolsRequest) (*types.QueryPoolsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := sdkCtx.KVStore(k.storeKey)
	poolsStore := prefix.NewStore(store, types.PoolPrefix)

	var pools []types.Pool
	pageRes, err := query.Paginate(poolsStore, request.Pagination, func(key []byte, value []byte) error {
		var pool types.Pool
		if err := k.cdc.Unmarshal(value, &pool); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		pools = append(pools, pool)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPoolsResponse{
		Pools:      pools,
		Pagination: pageRes,
	}, nil
}
