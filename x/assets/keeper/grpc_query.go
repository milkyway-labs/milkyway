package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/milkyway-labs/milkyway/v6/x/assets/types"
)

var _ types.QueryServer = queryServer{}

type queryServer struct {
	k *Keeper
}

func NewQueryServer(k *Keeper) types.QueryServer {
	return queryServer{k: k}
}

// Assets queries all the assets store in the module
func (q queryServer) Assets(ctx context.Context, req *types.QueryAssetsRequest) (*types.QueryAssetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Ticker != "" {
		// Validate the ticker
		err := types.ValidateTicker(req.Ticker)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		// Get the assets with the given ticker
		assets, pageRes, err := query.CollectionPaginate(ctx, q.k.TickerIndexes, req.Pagination,
			func(key collections.Pair[string, string], value collections.NoValue) (types.Asset, error) {
				denom := key.K2()
				return q.k.GetAsset(ctx, denom)
			},
			query.WithCollectionPaginationPairPrefix[string, string](req.Ticker),
		)
		if err != nil {
			return nil, err
		}

		return &types.QueryAssetsResponse{Assets: assets, Pagination: pageRes}, nil
	}

	// Get all the assets
	assets, pageRes, err := query.CollectionPaginate(ctx, q.k.Assets, req.Pagination, func(_ string, asset types.Asset) (types.Asset, error) {
		return asset, nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryAssetsResponse{Assets: assets, Pagination: pageRes}, nil
}

// Asset queries a single asset by its denom
func (q queryServer) Asset(ctx context.Context, req *types.QueryAssetRequest) (*types.QueryAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Validate the denom
	err := sdk.ValidateDenom(req.Denom)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid denom: %s", req.Denom)
	}

	// Get the asset
	asset, err := q.k.Assets.Get(ctx, req.Denom)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "asset for denom %s not registered", req.Denom)
		}
		return nil, err
	}

	return &types.QueryAssetResponse{Asset: asset}, nil
}
