package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
)

// PolicyOpsKeeper contains mutable operations on a policy.
type PolicyOpsKeeper interface {
	Create(ctx sdk.Context, creator sdk.AccAddress, regoCode []byte, source string, entry_points []byte, instantiateAccess *AccessConfig) (regoID uint64, err error)
}

// ViewKeeper provides read only operations
type ViewKeeper interface {
	GetRegoInfo(ctx types.Context, regoID uint64) *RegoInfo
	GetByteRego(ctx types.Context, regoID uint64) ([]byte, error)
	IterateRegoInfos(ctx types.Context, cb func(uint64, RegoInfo) bool)
}
