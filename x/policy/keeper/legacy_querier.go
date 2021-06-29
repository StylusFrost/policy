package keeper

import (
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/StylusFrost/policy/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryListPolicytByRegoCode = "list-policies-by-rego"
	QueryGetRego               = "rego"
	QueryListRego              = "list-rego"
)

// NewLegacyQuerier creates a new querier
func NewLegacyQuerier(keeper types.ViewKeeper, gasLimit sdk.Gas) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		var (
			rsp interface{}
			err error
		)
		switch path[0] {
		//TODO: Other cases
		case QueryGetRego:
			regoID, err := strconv.ParseUint(path[1], 10, 64)
			if err != nil {
				return nil, sdkerrors.Wrapf(types.ErrInvalid, "rego id: %s", err.Error())
			}
			rsp, err = queryRego(ctx, regoID, keeper)
		case QueryListRego:
			rsp, err = queryRegoList(ctx, keeper)
		case QueryListPolicytByRegoCode:
			regoID, err := strconv.ParseUint(path[1], 10, 64)
			if err != nil {
				return nil, sdkerrors.Wrapf(types.ErrInvalid, "rego id: %s", err.Error())
			}
			rsp = queryPolicyListByRegoCode(ctx, regoID, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown data query endpoint")
		}
		if err != nil {
			return nil, err
		}
		if rsp == nil || reflect.ValueOf(rsp).IsNil() {
			return nil, nil
		}
		bz, err := json.MarshalIndent(rsp, "", "  ")
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil
	}
}

func queryRegoList(ctx sdk.Context, k types.ViewKeeper) ([]types.RegoInfoResponse, error) {
	var info []types.RegoInfoResponse
	k.IterateRegoInfos(ctx, func(i uint64, res types.RegoInfo) bool {
		info = append(info, types.RegoInfoResponse{
			RegoID:      i,
			Creator:     res.Creator,
			RegoHash:    res.RegoHash,
			Source:      res.Source,
			EntryPoints: res.EntryPoints,
		})
		return false
	})
	return info, nil
}

func queryPolicyListByRegoCode(ctx sdk.Context, regoID uint64, keeper types.ViewKeeper) []string {
	var policies []string
	keeper.IteratePoliciesByRegoCode(ctx, regoID, func(addr sdk.AccAddress) bool {
		policies = append(policies, addr.String())
		return false
	})
	return policies
}
