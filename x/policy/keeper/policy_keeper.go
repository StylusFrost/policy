package keeper

import (
	"github.com/StylusFrost/policy/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.PolicyOpsKeeper = PermissionedKeeper{}

// decoratedKeeper contains a subset of the polocy keeper that are already or can be guarded by an authorization policy in the future
type decoratedKeeper interface {
	create(ctx sdk.Context, creator sdk.AccAddress, regoCode []byte, source string, entry_points []byte, instantiateAccess *types.AccessConfig, authZ AuthorizationPolicy) (regoID uint64, err error)
	instantiate(ctx sdk.Context, regoID uint64, creator, admin sdk.AccAddress, ntry_points []byte, label string, deposit sdk.Coins, authZ AuthorizationPolicy) (sdk.AccAddress, error)
	setPolicyAdmin(ctx sdk.Context, policyAddress, caller, newAdmin sdk.AccAddress, authZ AuthorizationPolicy) error
	migrate(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, newRegoID uint64, entry_points []byte, authZ AuthorizationPolicy) error
}

type PermissionedKeeper struct {
	authZPolicy AuthorizationPolicy
	nested      decoratedKeeper
}

func NewPermissionedKeeper(nested decoratedKeeper, authZPolicy AuthorizationPolicy) *PermissionedKeeper {
	return &PermissionedKeeper{authZPolicy: authZPolicy, nested: nested}
}

func NewGovPermissionKeeper(nested decoratedKeeper) *PermissionedKeeper {
	return NewPermissionedKeeper(nested, GovAuthorizationPolicy{})
}

func NewDefaultPermissionKeeper(nested decoratedKeeper) *PermissionedKeeper {
	return NewPermissionedKeeper(nested, DefaultAuthorizationPolicy{})
}

func (p PermissionedKeeper) Create(ctx sdk.Context, creator sdk.AccAddress, regoCode []byte, source string, entry_points []byte, instantiateAccess *types.AccessConfig) (regoID uint64, err error) {
	return p.nested.create(ctx, creator, regoCode, source, entry_points, instantiateAccess, p.authZPolicy)
}

func (p PermissionedKeeper) Instantiate(ctx sdk.Context, regoID uint64, creator, admin sdk.AccAddress, entry_points []byte, label string, deposit sdk.Coins) (sdk.AccAddress, error) {
	return p.nested.instantiate(ctx, regoID, creator, admin, entry_points, label, deposit, p.authZPolicy)
}

func (p PermissionedKeeper) UpdatePolicyAdmin(ctx sdk.Context, policyAddress sdk.AccAddress, caller sdk.AccAddress, newAdmin sdk.AccAddress) error {
	return p.nested.setPolicyAdmin(ctx, policyAddress, caller, newAdmin, p.authZPolicy)
}

func (p PermissionedKeeper) ClearPolicyAdmin(ctx sdk.Context, policyAddress sdk.AccAddress, caller sdk.AccAddress) error {
	return p.nested.setPolicyAdmin(ctx, policyAddress, caller, nil, p.authZPolicy)
}

func (p PermissionedKeeper) Migrate(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, newRegoID uint64, entry_points []byte) error {
	return p.nested.migrate(ctx, contractAddress, caller, newRegoID, entry_points, p.authZPolicy)
}
