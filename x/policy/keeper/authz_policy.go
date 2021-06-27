package keeper

import (
	"github.com/StylusFrost/policy/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AuthorizationPolicy interface {
	CanCreateRego(c types.AccessConfig, creator sdk.AccAddress) bool
}

type DefaultAuthorizationPolicy struct {
}

func (p DefaultAuthorizationPolicy) CanCreateRego(config types.AccessConfig, actor sdk.AccAddress) bool {
	return config.Allowed(actor)
}

type GovAuthorizationPolicy struct {
}

func (p GovAuthorizationPolicy) CanCreateRego(types.AccessConfig, sdk.AccAddress) bool {
	return true
}
