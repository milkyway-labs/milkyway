package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/milkyway-labs/milkyway/v7/x/vestingreward/types"
)

type Querier struct {
	k *Keeper
}

var _ types.QueryServer = Querier{}

func NewQuerier(keeper *Keeper) Querier {
	return Querier{k: keeper}
}

// VestingAccountsRewardRatio implements the Query/VestingAccountsRewardRatio gRPC method
func (q Querier) VestingAccountsRewardRatio(ctx context.Context, _ *types.QueryVestingAccountsRewardRatioRequest) (*types.QueryVestingAccountsRewardRatioResponse, error) {
	ratio, err := q.k.VestingAccountsRewardRatio.Get(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryVestingAccountsRewardRatioResponse{VestingAccountsRewardRatio: ratio}, nil
}
