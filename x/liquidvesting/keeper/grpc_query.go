package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

type Querier struct {
	*Keeper
}

var _ types.QueryServer = Querier{}

func NewQuerier(keeper *Keeper) Querier {
	return Querier{Keeper: keeper}
}

// InsuranceFund implements types.QueryServer.
func (q Querier) InsuranceFund(goCtx context.Context, _ *types.QueryInsuranceFundRequest) (*types.QueryInsuranceFundResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(goCtx)

	balance, err := q.GetInsuranceFundBalance(sdkCtx)
	if err != nil {
		return nil, err
	}
	return &types.QueryInsuranceFundResponse{Amount: balance}, nil
}

// Params implements types.QueryServer.
func (q Querier) Params(goCtx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(goCtx)

	params, err := q.GetParams(sdkCtx)
	if err != nil {
		return nil, err
	}

	return &types.QueryParamsResponse{Params: params}, nil
}

// UserInsuranceFund implements types.QueryServer.
func (q Querier) UserInsuranceFund(goCtx context.Context, req *types.QueryUserInsuranceFundRequest) (*types.QueryUserInsuranceFundResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(goCtx)

	accAddr, err := sdk.AccAddressFromBech32(req.UserAddress)
	if err != nil {
		return nil, err
	}

	insuranceFund, err := q.GetUserInsuranceFund(sdkCtx, accAddr)
	if err != nil {
		return nil, err
	}
	usedInsuranceFund, err := q.GetUserUsedInsuranceFund(sdkCtx, req.UserAddress)
	if err != nil {
		return nil, err
	}
	return &types.QueryUserInsuranceFundResponse{
		Balance: insuranceFund.Balance,
		Used:    usedInsuranceFund,
	}, nil
}

// UserInsuranceFunds implements types.QueryServer.
func (q Querier) UserInsuranceFunds(goCtx context.Context, req *types.QueryUserInsuranceFundsRequest) (*types.QueryUserInsuranceFundsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(goCtx)

	insuranceFunds, pagination, err := query.CollectionPaginate(sdkCtx, q.insuranceFunds, req.Pagination,
		func(userAddress sdk.AccAddress, insuranceFund types.UserInsuranceFund) (types.UserInsuranceFundData, error) {
			stringAddr, err := q.accountKeeper.AddressCodec().BytesToString(userAddress)
			if err != nil {
				return types.UserInsuranceFundData{}, err
			}

			usedInsuranceFund, err := q.GetUserUsedInsuranceFund(sdkCtx, stringAddr)
			if err != nil {
				return types.UserInsuranceFundData{}, err
			}

			return types.NewUserInsuranceFundData(stringAddr, insuranceFund.Balance, usedInsuranceFund), nil
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
func (q Querier) UserRestakableAssets(goCtx context.Context, req *types.QueryUserRestakableAssetsRequest) (*types.QueryUserRestakableAssetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(goCtx)

	balance, err := q.GetUserInsuranceFundBalance(sdkCtx, req.UserAddress)
	if err != nil {
		return nil, err
	}
	params, err := q.GetParams(sdkCtx)
	if err != nil {
		return nil, err
	}

	// Compute the amount of tokens that the user can restake
	restakableCoins := sdk.NewCoins()
	for _, coin := range balance {
		restakableAmount := math.LegacyNewDecFromInt(coin.Amount).
			MulInt64(100).
			QuoTruncate(params.InsurancePercentage).TruncateInt()
		vestedDenom, err := types.GetVestedRepresentationDenom(coin.Denom)
		if err != nil {
			return nil, err
		}
		restakableCoins = restakableCoins.Add(sdk.NewCoin(vestedDenom, restakableAmount))
	}

	return &types.QueryUserRestakableAssetsResponse{Amount: restakableCoins}, nil
}
