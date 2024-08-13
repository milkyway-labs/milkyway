package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/milkyway-labs/milkyway/x/assets/types"
)

var _ types.QueryServer = queryServer{}

type queryServer struct {
	k *Keeper
}

func NewQueryServer(k *Keeper) types.QueryServer {
	return queryServer{k: k}
}

func (q queryServer) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params, err := q.k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryParamsResponse{Params: params}, nil
}

func (q queryServer) Assets(ctx context.Context, req *types.QueryAssetsRequest) (*types.QueryAssetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Ticker != "" {
		err := types.ValidateTicker(req.Ticker)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		assets, pageRes, err := query.CollectionPaginate(ctx, q.k.TickerIndexes, req.Pagination,
			func(key collections.Pair[string, string], value collections.NoValue) (types.Asset, error) {
				denom := key.K2()
				return q.k.GetAsset(ctx, denom)
			},
			query.WithCollectionPaginationPairPrefix[string, string](req.Ticker))
		if err != nil {
			return nil, err
		}
		return &types.QueryAssetsResponse{Assets: assets, Pagination: pageRes}, nil
	}

	assets, pageRes, err := query.CollectionPaginate(
		ctx, q.k.Assets, req.Pagination,
		func(_ string, asset types.Asset) (types.Asset, error) {
			return asset, nil
		})
	if err != nil {
		return nil, err
	}
	return &types.QueryAssetsResponse{Assets: assets, Pagination: pageRes}, nil
}

func (q queryServer) Asset(ctx context.Context, req *types.QueryAssetRequest) (*types.QueryAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if err := sdk.ValidateDenom(req.Denom); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid denom: %s", req.Denom)
	}
	asset, err := q.k.Assets.Get(ctx, req.Denom)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "asset for denom %s not registered", req.Denom)
		}
		return nil, err
	}
	return &types.QueryAssetResponse{Asset: asset}, nil
}
