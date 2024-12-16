package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/milkyway-labs/milkyway/v7/x/liquidvesting/types"
)

type Querier struct {
	*Keeper
}

var _ types.QueryServer = Querier{}

func NewQuerier(keeper *Keeper) Querier {
	return Querier{Keeper: keeper}
}

// InsuranceFund implements types.QueryServer.
func (q Querier) InsuranceFund(ctx context.Context, _ *types.QueryInsuranceFundRequest) (*types.QueryInsuranceFundResponse, error) {
	balance, err := q.GetInsuranceFundBalance(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryInsuranceFundResponse{Amount: balance}, nil
}

// Params implements types.QueryServer.
func (q Querier) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params, err := q.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryParamsResponse{Params: params}, nil
}

// UserInsuranceFund implements types.QueryServer.
func (q Querier) UserInsuranceFund(ctx context.Context, req *types.QueryUserInsuranceFundRequest) (*types.QueryUserInsuranceFundResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	insuranceFund, err := q.GetUserInsuranceFund(ctx, req.UserAddress)
	if err != nil {
		return nil, err
	}

	usedInsuranceFund, err := q.GetUserUsedInsuranceFund(ctx, req.UserAddress)
	if err != nil {
		return nil, err
	}

	return &types.QueryUserInsuranceFundResponse{
		Balance: insuranceFund.Balance,
		Used:    usedInsuranceFund,
	}, nil
}

// UserInsuranceFunds implements types.QueryServer.
func (q Querier) UserInsuranceFunds(ctx context.Context, req *types.QueryUserInsuranceFundsRequest) (*types.QueryUserInsuranceFundsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	insuranceFunds, pagination, err := query.CollectionPaginate(ctx, q.insuranceFunds, req.Pagination,
		func(userAddress string, insuranceFund types.UserInsuranceFund) (types.UserInsuranceFundData, error) {
			usedInsuranceFund, err := q.GetUserUsedInsuranceFund(ctx, userAddress)
			if err != nil {
				return types.UserInsuranceFundData{}, err
			}

			return types.NewUserInsuranceFundData(userAddress, insuranceFund.Balance, usedInsuranceFund), nil
		})
	if err != nil {
		return nil, err
	}

	return &types.QueryUserInsuranceFundsResponse{
		InsuranceFunds: insuranceFunds,
		Pagination:     pagination,
	}, nil
}

// UserRestakableAssets implements types.QueryServer.
func (q Querier) UserRestakableAssets(ctx context.Context, req *types.QueryUserRestakableAssetsRequest) (*types.QueryUserRestakableAssetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	balance, err := q.GetUserInsuranceFundBalance(ctx, req.UserAddress)
	if err != nil {
		return nil, err
	}

	params, err := q.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	// Compute the amount of tokens that the user can restake
	restakableCoins := sdk.NewCoins()
	for _, coin := range balance {
		restakableAmount := math.LegacyNewDecFromInt(coin.Amount).MulInt64(100).QuoTruncate(params.InsurancePercentage).TruncateInt()

		lockedDenom, err := types.GetLockedRepresentationDenom(coin.Denom)
		if err != nil {
			return nil, err
		}

		restakableCoins = restakableCoins.Add(sdk.NewCoin(lockedDenom, restakableAmount))
	}

	return &types.QueryUserRestakableAssetsResponse{Amount: restakableCoins}, nil
}
