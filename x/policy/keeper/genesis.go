package keeper

import (
	"github.com/StylusFrost/policy/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k *Keeper, genState types.GenesisState) {
	k.setParams(ctx, genState.Params)

}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k *Keeper) *types.GenesisState {
	var genState types.GenesisState

	genState.Params = k.GetParams(ctx)
	return &genState
}
