package keeper

import (
	"fmt"

	"github.com/StylusFrost/policy/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// NewPolicyProposalHandler creates a new governance Handler for policy proposals
func NewPolicyProposalHandler(k decoratedKeeper, enabledProposalTypes []types.ProposalType) govtypes.Handler {
	return NewPolicyProposalHandlerX(NewGovPermissionKeeper(k), enabledProposalTypes)
}

// NewPolicyProposalHandlerX creates a new governance Handler for policy proposals
func NewPolicyProposalHandlerX(k types.PolicyOpsKeeper, enabledProposalTypes []types.ProposalType) govtypes.Handler {
	enabledTypes := make(map[string]struct{}, len(enabledProposalTypes))
	for i := range enabledProposalTypes {
		enabledTypes[string(enabledProposalTypes[i])] = struct{}{}
	}
	return func(ctx sdk.Context, content govtypes.Content) error {
		if content == nil {
			return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "content must not be empty")
		}
		if _, ok := enabledTypes[content.ProposalType()]; !ok {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unsupported policy proposal content type: %q", content.ProposalType())
		}
		switch c := content.(type) {
		case *types.StoreRegoProposal:
			return handleStoreCodeRego(ctx, k, *c)
		case *types.InstantiatePolicyProposal:
			return handleInstantiateProposal(ctx, k, *c)
		case *types.MigratePolicyProposal:
			return handleMigrateProposal(ctx, k, *c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized policy proposal content type: %T", c)
		}
	}
}

func handleStoreCodeRego(ctx sdk.Context, k types.PolicyOpsKeeper, p types.StoreRegoProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	runAsAddr, err := sdk.AccAddressFromBech32(p.RunAs)
	if err != nil {
		return sdkerrors.Wrap(err, "run as address")
	}
	regoID, err := k.Create(ctx, runAsAddr, p.REGOByteCode, p.Source, p.EntryPoints, p.InstantiatePermission)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyRegoID, fmt.Sprintf("%d", regoID)),
	))

	return err
}

func handleInstantiateProposal(ctx sdk.Context, k types.PolicyOpsKeeper, p types.InstantiatePolicyProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}
	runAsAddr, err := sdk.AccAddressFromBech32(p.RunAs)
	if err != nil {
		return sdkerrors.Wrap(err, "run as address")
	}
	adminAddr, err := sdk.AccAddressFromBech32(p.Admin)
	if err != nil {
		return sdkerrors.Wrap(err, "admin")
	}

	policyAddress, err := k.Instantiate(ctx, p.RegoID, runAsAddr, adminAddr, p.EntryPoints, p.Label, p.Funds)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeGovPolicyResult,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyRegoID, fmt.Sprintf("%d", p.RegoID)),
		sdk.NewAttribute(types.AttributeKeyPolicyAddr, policyAddress.String()),
	))
	return nil
}

func handleMigrateProposal(ctx sdk.Context, k types.PolicyOpsKeeper, p types.MigratePolicyProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	policyAddr, err := sdk.AccAddressFromBech32(p.Policy)
	if err != nil {
		return sdkerrors.Wrap(err, "policy")
	}
	runAsAddr, err := sdk.AccAddressFromBech32(p.RunAs)
	if err != nil {
		return sdkerrors.Wrap(err, "run as address")
	}
	err = k.Migrate(ctx, policyAddr, runAsAddr, p.RegoID, p.EntryPoints)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeGovPolicyResult,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyRegoID, fmt.Sprintf("%d", p.RegoID)),
		sdk.NewAttribute(types.AttributeKeyPolicyAddr, policyAddr.String()),
	))
	return err
}
