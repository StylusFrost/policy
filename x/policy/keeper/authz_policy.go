package keeper

import (
	"github.com/StylusFrost/policy/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AuthorizationPolicy interface {
	CanCreateRego(c types.AccessConfig, creator sdk.AccAddress) bool
	CanInstantiatePolicy(c types.AccessConfig, actor sdk.AccAddress) bool
	CanModifyPolicy(admin, actor sdk.AccAddress) bool
	CanRefundPolicy(admin, actor sdk.AccAddress) bool
}

type DefaultAuthorizationPolicy struct {
}

func (p DefaultAuthorizationPolicy) CanCreateRego(config types.AccessConfig, actor sdk.AccAddress) bool {
	return config.Allowed(actor)
}
func (p DefaultAuthorizationPolicy) CanInstantiatePolicy(config types.AccessConfig, actor sdk.AccAddress) bool {
	return config.Allowed(actor)
}
func (p DefaultAuthorizationPolicy) CanRefundPolicy(admin, actor sdk.AccAddress) bool {
	return admin != nil && admin.Equals(actor)
}
func (p DefaultAuthorizationPolicy) CanModifyPolicy(admin, actor sdk.AccAddress) bool {
	return admin != nil && admin.Equals(actor)
}

type GovAuthorizationPolicy struct {
}

func (p GovAuthorizationPolicy) CanCreateRego(types.AccessConfig, sdk.AccAddress) bool {
	return true
}

func (p GovAuthorizationPolicy) CanInstantiatePolicy(types.AccessConfig, sdk.AccAddress) bool {
	return true
}

func (p GovAuthorizationPolicy) CanModifyPolicy(sdk.AccAddress, sdk.AccAddress) bool {
	return true
}
func (p GovAuthorizationPolicy) CanRefundPolicy(sdk.AccAddress, sdk.AccAddress) bool {
	return true
}
