package keeper

import (
	"github.com/StylusFrost/policy/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.PolicyOpsKeeper = PermissionedKeeper{}

// decoratedKeeper contains a subset of the polocy keeper that are already or can be guarded by an authorization policy in the future
type decoratedKeeper interface {
	create(ctx sdk.Context, creator sdk.AccAddress, regoCode []byte, source string, entry_points []byte, instantiateAccess *types.AccessConfig, authZ AuthorizationPolicy) (regoID uint64, err error)
}

type PermissionedKeeper struct {
	authZPolicy AuthorizationPolicy
	nested      decoratedKeeper
}

func NewPermissionedKeeper(nested decoratedKeeper, authZPolicy AuthorizationPolicy) *PermissionedKeeper {
	return &PermissionedKeeper{ authZPolicy: authZPolicy,nested: nested}
}

func NewGovPermissionKeeper(nested decoratedKeeper) *PermissionedKeeper {
	return NewPermissionedKeeper(nested , GovAuthorizationPolicy{})
}

func NewDefaultPermissionKeeper(nested decoratedKeeper) *PermissionedKeeper {
	return NewPermissionedKeeper(nested , DefaultAuthorizationPolicy{})
}

func (p PermissionedKeeper) Create(ctx sdk.Context, creator sdk.AccAddress, regoCode []byte, source string, entry_points []byte, instantiateAccess *types.AccessConfig) (regoID uint64, err error) {
	return p.nested.create(ctx, creator, regoCode, source, entry_points, instantiateAccess, p.authZPolicy)
}
