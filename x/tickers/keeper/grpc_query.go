package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/milkyway-labs/milkyway/x/tickers/types"
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

func (q queryServer) Ticker(ctx context.Context, req *types.QueryTickerRequest) (*types.QueryTickerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if err := sdk.ValidateDenom(req.Denom); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid denom: %s", req.Denom)
	}
	ticker, err := q.k.Tickers.Get(ctx, req.Denom)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "ticker for denom %s not registered", req.Denom)
		}
		return nil, err
	}
	return &types.QueryTickerResponse{Ticker: ticker}, nil
}

func (q queryServer) Denoms(ctx context.Context, req *types.QueryDenomsRequest) (*types.QueryDenomsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if err := types.ValidateTicker(req.Ticker); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid ticker: %s", err)
	}
	denoms, pageRes, err := query.CollectionPaginate(ctx, q.k.TickerIndexes, req.Pagination, func(key collections.Pair[string, string], _ collections.NoValue) (string, error) {
		return key.K2(), nil
	}, query.WithCollectionPaginationPairPrefix[string, string](req.Ticker))
	if err != nil {
		return nil, err
	}
	return &types.QueryDenomsResponse{Denoms: denoms, Pagination: pageRes}, nil
}
