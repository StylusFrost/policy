package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
)

// PolicyOpsKeeper contains mutable operations on a policy.
type PolicyOpsKeeper interface {
	// Create uploads and compiles a REGO policy, returning a short identifier for the policy
	Create(ctx sdk.Context, creator sdk.AccAddress, regoCode []byte, source string, entry_points []byte, instantiateAccess *AccessConfig) (regoID uint64, err error)
	// Instantiate creates an instance of a REGO policy
	Instantiate(ctx sdk.Context, regoID uint64, creator, admin sdk.AccAddress, entry_points []byte, label string, deposit sdk.Coins) (sdk.AccAddress, error)
	// UpdatePolicyAdmin sets the admin value on the PolicyInfo. It must be a valid address (use ClearPolicyAdmin to remove it)
	UpdatePolicyAdmin(ctx sdk.Context, policyAddress sdk.AccAddress, caller sdk.AccAddress, newAdmin sdk.AccAddress) error
	// ClearPolicyAdmin sets the admin value on the PolicyInfo to nil, to disable further migrations/ updates.
	ClearPolicyAdmin(ctx sdk.Context, policyAddress sdk.AccAddress, caller sdk.AccAddress) error
	// Migrate allows to upgrade a policy to a new rego
	Migrate(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, newRegoID uint64, msg []byte) error
}

// ViewKeeper provides read only operations
type ViewKeeper interface {
	GetPolicyHistory(ctx types.Context, contractAddr types.AccAddress) []PolicyRegoHistoryEntry
	GetRegoInfo(ctx types.Context, regoID uint64) *RegoInfo
	GetByteRego(ctx types.Context, regoID uint64) ([]byte, error)
	IterateRegoInfos(ctx types.Context, cb func(uint64, RegoInfo) bool)
	GetPolicyInfo(ctx types.Context, policyAddress types.AccAddress) *PolicyInfo
	IteratePoliciesByRegoCode(ctx sdk.Context, regoID uint64, cb func(address sdk.AccAddress) bool)
}
