package keeper

import (
	"github.com/StylusFrost/policy/x/policy/types"
)

var _ types.QueryServer = Keeper{}
