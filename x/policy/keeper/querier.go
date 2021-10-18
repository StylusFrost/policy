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
	cdc           codec.BinaryCodec
	storeKey      sdk.StoreKey
	keeper        types.ViewKeeper
	queryGasLimit sdk.Gas
}

func NewGrpcQuerier(cdc codec.BinaryCodec, storeKey sdk.StoreKey, keeper types.ViewKeeper, queryGasLimit sdk.Gas) *grpcQuerier {
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
			if err := q.cdc.Unmarshal(value, &c); err != nil {
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

func (q grpcQuerier) PolicyInfo(c context.Context, req *types.QueryPolicyInfoRequest) (*types.QueryPolicyInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	policyAddr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	rsp, err := queryPolicyInfo(sdk.UnwrapSDKContext(c), policyAddr, q.keeper)
	switch {
	case err != nil:
		return nil, err
	case rsp == nil:
		return nil, types.ErrNotFound
	}
	return rsp, nil
}
func queryPolicyInfo(ctx sdk.Context, addr sdk.AccAddress, keeper types.ViewKeeper) (*types.QueryPolicyInfoResponse, error) {
	info := keeper.GetPolicyInfo(ctx, addr)
	if info == nil {
		return nil, types.ErrNotFound
	}
	// redact the Created field (just used for sorting, not part of public API)
	info.Created = nil
	return &types.QueryPolicyInfoResponse{
		Address:    addr.String(),
		PolicyInfo: *info,
	}, nil
}

func (q grpcQuerier) PoliciesByRegoCode(c context.Context, req *types.QueryPoliciesByRegoCodeRequest) (*types.QueryPoliciesByRegoCodeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.RegoId == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "code id")
	}
	ctx := sdk.UnwrapSDKContext(c)
	r := make([]string, 0)

	prefixStore := prefix.NewStore(ctx.KVStore(q.storeKey), types.GetPolicyByRegoIDSecondaryIndexPrefix(req.RegoId))
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			var policyAddr sdk.AccAddress = key[types.AbsoluteTxPositionLen:]
			r = append(r, policyAddr.String())
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryPoliciesByRegoCodeResponse{
		Policies:   r,
		Pagination: pageRes,
	}, nil
}

func (q grpcQuerier) PolicyHistory(c context.Context, req *types.QueryPolicyHistoryRequest) (*types.QueryPolicyHistoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	policyAddr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	r := make([]types.PolicyRegoHistoryEntry, 0)

	prefixStore := prefix.NewStore(ctx.KVStore(q.storeKey), types.GetPolicyRegoHistoryElementPrefix(policyAddr))
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			var e types.PolicyRegoHistoryEntry
			if err := q.cdc.Unmarshal(value, &e); err != nil {
				return false, err
			}
			e.Updated = nil // redact
			r = append(r, e)
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryPolicyHistoryResponse{
		Entries:    r,
		Pagination: pageRes,
	}, nil
}
