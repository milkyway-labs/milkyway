package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/milkyway-labs/milkyway/v3/x/operators/types"
)

var _ types.QueryServer = &Keeper{}

// Operator implements the Query/Operator gRPC method
func (k *Keeper) Operator(ctx context.Context, request *types.QueryOperatorRequest) (*types.QueryOperatorResponse, error) {
	operator, found, err := k.GetOperator(ctx, request.OperatorId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !found {
		return nil, status.Error(codes.NotFound, "operator not found")
	}

	return &types.QueryOperatorResponse{Operator: operator}, nil
}

// Operators implements the Query/Operators gRPC method
func (k *Keeper) Operators(ctx context.Context, request *types.QueryOperatorsRequest) (*types.QueryOperatorsResponse, error) {
	store := k.storeService.OpenKVStore(ctx)
	operatorsStore := prefix.NewStore(runtime.KVStoreAdapter(store), types.OperatorPrefix)

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
	_, found, err := k.GetOperator(ctx, request.OperatorId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !found {
		return nil, types.ErrOperatorNotFound
	}

	params, err := k.GetOperatorParams(ctx, request.OperatorId)
	if err != nil {
		return nil, err
	}

	return &types.QueryOperatorParamsResponse{OperatorParams: params}, nil
}

// Params implements the Query/Params gRPC method
func (k *Keeper) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryParamsResponse{Params: params}, nil
}
