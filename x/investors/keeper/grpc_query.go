package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/milkyway-labs/milkyway/v11/x/investors/types"
)

type Querier struct {
	k *Keeper
}

var _ types.QueryServer = Querier{}

func NewQuerier(keeper *Keeper) Querier {
	return Querier{k: keeper}
}

// InvestorsRewardRatio implements the Query/InvestorsRewardRatio gRPC method
func (q Querier) InvestorsRewardRatio(ctx context.Context, _ *types.QueryInvestorsRewardRatioRequest) (*types.QueryInvestorsRewardRatioResponse, error) {
	ratio, err := q.k.GetInvestorsRewardRatio(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryInvestorsRewardRatioResponse{InvestorsRewardRatio: ratio}, nil
}

// VestingInvestors implements the Query/VestingInvestors gRPC method
func (q Querier) VestingInvestors(ctx context.Context, req *types.QueryVestingInvestorsRequest) (*types.QueryVestingInvestorsResponse, error) {
	investors, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.VestingInvestors,
		req.Pagination,
		func(key string, value collections.NoValue) (string, error) {
			return key, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.QueryVestingInvestorsResponse{VestingInvestorsAddresses: investors, Pagination: pageRes}, nil
}
