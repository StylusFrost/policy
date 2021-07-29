package types

import (
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/policy module sentinel errors
var (
	DefaultCodespace = ModuleName

	// ErrExecuteFailed error for  execution policy failure
	ErrExecuteFailed = sdkErrors.Register(DefaultCodespace, 5, "execute rego policy failed")

	// ErrLimit error for content that exceeds a limit
	ErrLimit = sdkErrors.Register(DefaultCodespace, 13, "exceeds limit")

	// ErrEntryPoint error for entry point
	ErrEntryPoint = sdkErrors.Register(DefaultCodespace, 15, "invalid Entry Point")

	// ErrCreateFailed error for rego code that has already been uploaded or failed
	ErrCreateFailed = sdkErrors.Register(DefaultCodespace, 2, "create Rego Policy failed")

	// ErrAccountExists error for a policy account that already exists
	ErrAccountExists = sdkErrors.Register(DefaultCodespace, 3, "policy account already exists")

	// ErrInstantiateFailed error for rust instantiate policy failure
	ErrInstantiateFailed = sdkErrors.Register(DefaultCodespace, 4, "instantiate rego policy failed")

	// ErrInvalid error for content that is invalid in this context
	ErrInvalid = sdkErrors.Register(DefaultCodespace, 14, "invalid")

	// ErrEmpty error for empty content
	ErrEmpty = sdkErrors.Register(DefaultCodespace, 12, "empty")

	// ErrNotFound error for an entry not found in the store
	ErrNotFound = sdkErrors.Register(DefaultCodespace, 8, "not found")

	// ErrNotFound error for an entry not found in the store
	ErrCompileFailed = sdkErrors.Register(DefaultCodespace, 9, "compile rego code failed")

	// ErrInvalidMsg error when we cannot process the error returned from the policy
	ErrInvalidMsg = sdkErrors.Register(DefaultCodespace, 10, "invalid msg from the policy")
)
