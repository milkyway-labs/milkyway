package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/milkyway-labs/milkyway/v7/x/investors/types"
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
	ratio, err := q.k.InvestorsRewardRatio.Get(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryInvestorsRewardRatioResponse{InvestorsRewardRatio: ratio}, nil
}
