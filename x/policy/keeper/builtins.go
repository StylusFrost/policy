package keeper

import (
	"github.com/StylusFrost/policy/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	asttypes "github.com/open-policy-agent/opa/types"
)

type OptionRego func(r *rego.Rego)

func OptionsRego(opts ...OptionRego) OptionRego {
	return func(r *rego.Rego) {
		for _, opt := range opts {
			opt(r)
		}
	}
}

func (k Keeper) LoadCosmosBuiltins(ctx sdk.Context, caller sdk.AccAddress, policyAddres sdk.AccAddress) {

	rego.RegisterBuiltin1(&rego.Function{
		Name: "cosmos.bank.doIHaveBalance",
		Decl: asttypes.NewFunction(
			asttypes.Args(asttypes.S),
			asttypes.B,
		)}, func(_ rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {

		denom, ok := a.Value.(ast.String)

		if !ok {
			return nil, sdkerrors.Wrap(types.ErrExecuteFailed, "Invalid Denom")
		}
		hash := k.bank.HasBalance(ctx, caller, sdk.Coin{Denom: string(denom), Amount: sdk.NewInt(1)})

		return ast.BooleanTerm(hash), nil
	},
	)
	rego.RegisterBuiltin1(&rego.Function{
		Name: "cosmos.bank.hasPolicyBalance",
		Decl: asttypes.NewFunction(
			asttypes.Args(asttypes.S),
			asttypes.B,
		)}, func(_ rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {

		denom, ok := a.Value.(ast.String)

		if !ok {
			return nil, sdkerrors.Wrap(types.ErrExecuteFailed, "Invalid Denom")
		}
		hash := k.bank.HasBalance(ctx, policyAddres, sdk.Coin{Denom: string(denom), Amount: sdk.NewInt(1)})

		return ast.BooleanTerm(hash), nil
	},
	)

	rego.RegisterBuiltin1(&rego.Function{
		Name: "cosmos.bank.getMyBalance",
		Decl: asttypes.NewFunction(
			asttypes.Args(asttypes.S),
			asttypes.N,
		)}, func(_ rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {

		denom, ok := a.Value.(ast.String)

		if !ok {
			return nil, sdkerrors.Wrap(types.ErrExecuteFailed, "Invalid Denom")
		}
		balance := k.bank.GetBalance(ctx, caller, string(denom))

		return ast.IntNumberTerm(balance), nil
	},
	)
	rego.RegisterBuiltin1(&rego.Function{
		Name: "cosmos.bank.getPolicyBalance",
		Decl: asttypes.NewFunction(
			asttypes.Args(asttypes.S),
			asttypes.N,
		)}, func(_ rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {

		denom, ok := a.Value.(ast.String)

		if !ok {
			return nil, sdkerrors.Wrap(types.ErrExecuteFailed, "Invalid Denom")
		}
		balance := k.bank.GetBalance(ctx, policyAddres, string(denom))

		return ast.IntNumberTerm(balance), nil
	},
	)

	rego.RegisterBuiltin1(&rego.Function{
		Name: "cosmos.bank.getSupply",
		Decl: asttypes.NewFunction(
			asttypes.Args(asttypes.S),
			asttypes.N,
		)}, func(_ rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {

		denom, ok := a.Value.(ast.String)

		if !ok {
			return nil, sdkerrors.Wrap(types.ErrExecuteFailed, "Invalid Denom")
		}
		balance := k.bank.GetSupply(ctx, string(denom))

		return ast.IntNumberTerm(balance), nil
	},
	)
}
func RegistryCosmosBuiltins() {

	// Verify if i have balance from input denom
	ast.RegisterBuiltin(&ast.Builtin{
		Name: "cosmos.bank.doIHaveBalance",
		Decl: asttypes.NewFunction(
			asttypes.Args(asttypes.S),
			asttypes.B,
		),
	})
	// Verify if policy balance from input denom
	ast.RegisterBuiltin(&ast.Builtin{
		Name: "cosmos.bank.hasPolicyBalance",
		Decl: asttypes.NewFunction(
			asttypes.Args(asttypes.S),
			asttypes.B,
		),
	})
	// Get my balance from input denom
	ast.RegisterBuiltin(&ast.Builtin{
		Name: "cosmos.bank.getMyBalance",
		Decl: asttypes.NewFunction(
			asttypes.Args(asttypes.S),
			asttypes.N,
		),
	})
	// Get policy balance from input denom
	ast.RegisterBuiltin(&ast.Builtin{
		Name: "cosmos.bank.getPolicyBalance",
		Decl: asttypes.NewFunction(
			asttypes.Args(asttypes.S),
			asttypes.N,
		),
	})
	// Get total Supply from input denom
	ast.RegisterBuiltin(&ast.Builtin{
		Name: "cosmos.bank.getSupply",
		Decl: asttypes.NewFunction(
			asttypes.Args(asttypes.S),
			asttypes.N,
		),
	})

}
