package types

import (
	"net/url"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	MaxRegoSize = 500 * 1024
	// MaxLabelSize is the longest label that can be used when Instantiating a policy
	MaxLabelSize = 128
)

func validateSourceURL(source string) error {
	if source != "" {
		u, err := url.Parse(source)
		if err != nil {
			return sdkerrors.Wrap(ErrInvalid, "not an url")
		}
		if !u.IsAbs() {
			return sdkerrors.Wrap(ErrInvalid, "not an absolute url")
		}
		if u.Scheme != "https" {
			return sdkerrors.Wrap(ErrInvalid, "must use https")
		}
	}
	return nil
}

func validateRegoCode(s []byte) error {
	if len(s) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "is required")
	}
	if len(s) > MaxRegoSize {
		return sdkerrors.Wrapf(ErrLimit, "cannot be longer than %d bytes", MaxRegoSize)
	}
	return nil
}

func validateLabel(label string) error {
	if label == "" {
		return sdkerrors.Wrap(ErrEmpty, "is required")
	}
	if len(label) > MaxLabelSize {
		return sdkerrors.Wrap(ErrLimit, "cannot be longer than 128 characters")
	}
	return nil
}
