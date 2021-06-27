package keeper

import (
	"context"
	"encoding/binary"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/StylusFrost/policy/x/policy/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
)

var _ types.QueryServer = &grpcQuerier{}

type grpcQuerier struct {
	cdc           codec.Marshaler
	storeKey      sdk.StoreKey
	keeper        types.ViewKeeper
	queryGasLimit sdk.Gas
}

func NewGrpcQuerier(cdc codec.Marshaler, storeKey sdk.StoreKey, keeper types.ViewKeeper, queryGasLimit sdk.Gas) *grpcQuerier {
	return &grpcQuerier{cdc: cdc, storeKey: storeKey, keeper: keeper, queryGasLimit: queryGasLimit}
}

func (q grpcQuerier) Rego(c context.Context, req *types.QueryRegoRequest) (*types.QueryRegoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.RegoId == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "code id")
	}
	rsp, err := queryRego(sdk.UnwrapSDKContext(c), req.RegoId, q.keeper)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}
	return &types.QueryRegoResponse{
		RegoInfoResponse: rsp.RegoInfoResponse,
		Data:             rsp.Data,
	}, nil
}

func (q grpcQuerier) Regos(c context.Context, req *types.QueryRegosRequest) (*types.QueryRegosResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	r := make([]types.RegoInfoResponse, 0)
	prefixStore := prefix.NewStore(ctx.KVStore(q.storeKey), types.RegoKeyPrefix)
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			var c types.RegoInfo
			if err := q.cdc.UnmarshalBinaryBare(value, &c); err != nil {
				return false, err
			}

			r = append(r, types.RegoInfoResponse{
				RegoID:      binary.BigEndian.Uint64(key),
				Creator:     c.Creator,
				RegoHash:    c.RegoHash,
				Source:      c.Source,
				EntryPoints: c.EntryPoints,
			})
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryRegosResponse{RegoInfos: r, Pagination: pageRes}, nil
}

func queryRego(ctx sdk.Context, regoID uint64, keeper types.ViewKeeper) (*types.QueryRegoResponse, error) {
	if regoID == 0 {
		return nil, nil
	}
	res := keeper.GetRegoInfo(ctx, regoID)
	if res == nil {
		// nil, nil leads to 404 in rest handler
		return nil, nil
	}
	info := types.RegoInfoResponse{
		RegoID:      regoID,
		Creator:     res.Creator,
		RegoHash:    res.RegoHash,
		Source:      res.Source,
		EntryPoints: res.EntryPoints,
	}

	rego, err := keeper.GetByteRego(ctx, regoID)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "loading rego code")
	}

	return &types.QueryRegoResponse{RegoInfoResponse: &info, Data: rego}, nil
}
