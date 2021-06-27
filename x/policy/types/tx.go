package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (msg MsgStoreRego) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return err
	}

	if err := validateRegoCode(msg.REGOByteCode); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "code bytes %s", err.Error())
	}

	if err := validateSourceURL(msg.Source); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "source %s", err.Error())
	}

	if !json.Valid(msg.EntryPoints) {
		return sdkerrors.Wrap(ErrInvalid, "entry points json")
	}

	if msg.InstantiatePermission != nil {
		if err := msg.InstantiatePermission.ValidateBasic(); err != nil {
			return sdkerrors.Wrap(err, "instantiate permission")
		}
	}

	return nil
}

func (msg MsgStoreRego) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))

}

func (msg MsgStoreRego) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err.Error())
	}
	return []sdk.AccAddress{senderAddr}
}

func (msg MsgStoreRego) Route() string {
	return RouterKey
}

func (msg MsgStoreRego) Type() string {
	return "store-rego"
}
