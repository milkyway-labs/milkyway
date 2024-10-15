package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

var _ types.QueryServer = &Keeper{}

// Operator implements the Query/Operator gRPC method
func (k *Keeper) Operator(ctx context.Context, request *types.QueryOperatorRequest) (*types.QueryOperatorResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	operator, found := k.GetOperator(sdkCtx, request.OperatorId)
	if !found {
		return nil, status.Error(codes.NotFound, "operator not found")
	}

	return &types.QueryOperatorResponse{Operator: operator}, nil
}

// Operators implements the Query/Operators gRPC method
func (k *Keeper) Operators(ctx context.Context, request *types.QueryOperatorsRequest) (*types.QueryOperatorsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := sdkCtx.KVStore(k.storeKey)
	operatorsStore := prefix.NewStore(store, types.OperatorPrefix)

	var operators []types.Operator
	pageRes, err := query.Paginate(operatorsStore, request.Pagination, func(key []byte, value []byte) error {
		var operator types.Operator
		if err := k.cdc.Unmarshal(value, &operator); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		operators = append(operators, operator)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOperatorsResponse{
		Operators:  operators,
		Pagination: pageRes,
	}, nil
}

func (k *Keeper) OperatorParams(ctx context.Context, request *types.QueryOperatorParamsRequest) (*types.QueryOperatorParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	params, found, err := k.GetOperatorParams(sdkCtx, request.OperatorId)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, status.Error(codes.NotFound, "operator params not found")
	}

	return &types.QueryOperatorParamsResponse{OperatorParams: params}, nil
}

// Params implements the Query/Params gRPC method
func (k *Keeper) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(sdkCtx)
	return &types.QueryParamsResponse{Params: params}, nil
}
