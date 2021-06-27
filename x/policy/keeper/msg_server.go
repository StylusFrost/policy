package keeper

import (
	"context"
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

func (m msgServer) StoreRego(goCtx context.Context, msg *types.MsgStoreRego) (*types.MsgStoreRegoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}

	
	regoID, err := m.keeper.Create(ctx, senderAddr, msg.REGOByteCode, msg.Source, msg.EntryPoints,msg.InstantiatePermission)
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
