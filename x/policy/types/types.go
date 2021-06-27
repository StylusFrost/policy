package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewRegoInfo(regoHash []byte, creator sdk.AccAddress, source string, entry_points []string, instantiatePermission AccessConfig) RegoInfo {
	return RegoInfo{
		RegoHash:          regoHash,
		Creator:           creator.String(),
		Source:            source,
		EntryPoints:       entry_points,
		InstantiateConfig: instantiatePermission,
	}
}
