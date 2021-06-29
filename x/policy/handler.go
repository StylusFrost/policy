package policy

import (
	"fmt"

	"github.com/StylusFrost/policy/x/policy/keeper"
	"github.com/StylusFrost/policy/x/policy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewHandler returns a handler for "policy" type messages.
func NewHandler(k types.PolicyOpsKeeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		var (
			res proto.Message
			err error
		)
		switch msg := msg.(type) {
		case *types.MsgStoreRego:
			res, err = msgServer.StoreRego(sdk.WrapSDKContext(ctx), msg)

		case *types.MsgInstantiatePolicy:
			res, err = msgServer.InstantiatePolicy(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgUpdateAdmin:
			res, err = msgServer.UpdateAdmin(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgClearAdmin:
			res, err = msgServer.ClearAdmin(sdk.WrapSDKContext(ctx), msg)

		default:
			errMsg := fmt.Sprintf("unrecognized policy message type: %T", msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}

		ctx = ctx.WithEventManager(filterMessageEvents(ctx))
		return sdk.WrapServiceResult(ctx, res, err)
	}
}

// filterMessageEvents returns the same events with all of type == EventTypeMessage removed except
// for policy` message types.
// this is so only our top-level message event comes through
func filterMessageEvents(ctx sdk.Context) *sdk.EventManager {
	m := sdk.NewEventManager()
	for _, e := range ctx.EventManager().Events() {
		if e.Type == sdk.EventTypeMessage &&
			!hasPolicyModuleAttribute(e.Attributes) {
			continue
		}
		m.EmitEvent(e)
	}
	return m
}

func hasPolicyModuleAttribute(attrs []abci.EventAttribute) bool {
	for _, a := range attrs {
		if sdk.AttributeKeyModule == string(a.Key) &&
			types.ModuleName == string(a.Value) {
			return true
		}
	}
	return false
}
