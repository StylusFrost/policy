package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/StylusFrost/policy/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	keeper types.PolicyOpsKeeper
}

func NewMsgServerImpl(k types.PolicyOpsKeeper) types.MsgServer {
	return &msgServer{keeper: k}
}

func (m msgServer) RefundPolicy(goCtx context.Context, msg *types.MsgRefundPolicy) (*types.MsgRefundPolicyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	policyAddr, err := sdk.AccAddressFromBech32(msg.Policy)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "policy")
	}

	if err := m.keeper.Refund(ctx, policyAddr, senderAddr, msg.Refunds); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyPolicyAddr, msg.Policy),
	))

	return &types.MsgRefundPolicyResponse{}, nil
}

func (m msgServer) ExecutePolicy(goCtx context.Context, msg *types.MsgExecutePolicy) (*types.MsgExecutePolicyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	policyAddr, err := sdk.AccAddressFromBech32(msg.Policy)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "policy")
	}

	data, err := m.keeper.Execute(ctx, policyAddr, senderAddr, msg.EntryPoint, msg.Input, msg.Funds)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyPolicyAddr, msg.Policy),
		sdk.NewAttribute(types.AttributeKeyResultDataHex, hex.EncodeToString(data)),
	))

	return &types.MsgExecutePolicyResponse{
		Data: data,
	}, nil
}

func (m msgServer) StoreRego(goCtx context.Context, msg *types.MsgStoreRego) (*types.MsgStoreRegoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}

	regoID, err := m.keeper.Create(ctx, senderAddr, msg.REGOByteCode, msg.Source, msg.EntryPoints, msg.InstantiatePermission)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyRegoID, fmt.Sprintf("%d", regoID)),
	))

	return &types.MsgStoreRegoResponse{
		RegoID: regoID,
	}, nil
}

func (m msgServer) InstantiatePolicy(goCtx context.Context, msg *types.MsgInstantiatePolicy) (*types.MsgInstantiatePolicyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	var adminAddr sdk.AccAddress
	if msg.Admin != "" {
		if adminAddr, err = sdk.AccAddressFromBech32(msg.Admin); err != nil {
			return nil, sdkerrors.Wrap(err, "admin")
		}
	}

	policyAddr, err := m.keeper.Instantiate(ctx, msg.RegoID, senderAddr, adminAddr, msg.EntryPoints, msg.Label, msg.Funds)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyRegoID, fmt.Sprintf("%d", msg.RegoID)),
		sdk.NewAttribute(types.AttributeKeyPolicyAddr, policyAddr.String()),
	))

	return &types.MsgInstantiatePolicyResponse{
		Address: policyAddr.String(),
	}, nil
}

func (m msgServer) UpdateAdmin(goCtx context.Context, msg *types.MsgUpdateAdmin) (*types.MsgUpdateAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	policyAddr, err := sdk.AccAddressFromBech32(msg.Policy)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "policy")
	}
	newAdminAddr, err := sdk.AccAddressFromBech32(msg.NewAdmin)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "new admin")
	}

	if err := m.keeper.UpdatePolicyAdmin(ctx, policyAddr, senderAddr, newAdminAddr); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyPolicy, msg.Policy),
	))

	return &types.MsgUpdateAdminResponse{}, nil
}

func (m msgServer) ClearAdmin(goCtx context.Context, msg *types.MsgClearAdmin) (*types.MsgClearAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	policyAddr, err := sdk.AccAddressFromBech32(msg.Policy)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "policy")
	}

	if err := m.keeper.ClearPolicyAdmin(ctx, policyAddr, senderAddr); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyPolicy, msg.Policy),
	))

	return &types.MsgClearAdminResponse{}, nil
}

func (m msgServer) MigratePolicy(goCtx context.Context, msg *types.MsgMigratePolicy) (*types.MsgMigratePolicyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	policyAddr, err := sdk.AccAddressFromBech32(msg.Policy)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "policy")
	}

	err = m.keeper.Migrate(ctx, policyAddr, senderAddr, msg.RegoID, msg.EntryPoints)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyRegoID, fmt.Sprintf("%d", msg.RegoID)),
		sdk.NewAttribute(types.AttributeKeyPolicyAddr, msg.Policy),
	))

	return &types.MsgMigratePolicyResponse{}, nil
}
