package types

import (
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/policy module sentinel errors
var (
	DefaultCodespace = ModuleName

	// ErrLimit error for content that exceeds a limit
	ErrLimit = sdkErrors.Register(DefaultCodespace, 13, "exceeds limit")

	// ErrCreateFailed error for rego code that has already been uploaded or failed
	ErrCreateFailed = sdkErrors.Register(DefaultCodespace, 2, "create Rego Policy failed")

	// ErrInvalid error for content that is invalid in this context
	ErrInvalid = sdkErrors.Register(DefaultCodespace, 14, "invalid")

	// ErrEmpty error for empty content
	ErrEmpty = sdkErrors.Register(DefaultCodespace, 12, "empty")

	// ErrNotFound error for an entry not found in the store
	ErrNotFound = sdkErrors.Register(DefaultCodespace, 8, "not found")

	// ErrNotFound error for an entry not found in the store
	ErrCompileFailed = sdkErrors.Register(DefaultCodespace, 9, "compile rego code failed")
)
