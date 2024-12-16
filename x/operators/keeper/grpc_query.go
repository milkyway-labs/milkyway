package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/milkyway-labs/milkyway/v7/x/operators/types"
)

var _ types.QueryServer = &Keeper{}

// Operator implements the Query/Operator gRPC method
func (k *Keeper) Operator(ctx context.Context, request *types.QueryOperatorRequest) (*types.QueryOperatorResponse, error) {
	operator, err := k.GetOperator(ctx, request.OperatorId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "operator not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOperatorResponse{Operator: operator}, nil
}

// Operators implements the Query/Operators gRPC method
func (k *Keeper) Operators(ctx context.Context, request *types.QueryOperatorsRequest) (*types.QueryOperatorsResponse, error) {
	operators, pageRes, err := query.CollectionPaginate(ctx, k.operators, request.Pagination, func(_ uint32, operator types.Operator) (types.Operator, error) {
		return operator, nil
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
	_, err := k.GetOperator(ctx, request.OperatorId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "operator not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
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
