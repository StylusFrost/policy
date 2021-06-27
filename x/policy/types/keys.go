package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "policy"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// TStoreKey is the string transient store representation
	TStoreKey = "transient_" + ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_policy"
)

// nolint
var (
	RegoKeyPrefix     = []byte{0x01}
	SequenceKeyPrefix = []byte{0x02}
	RegoHashKeyPrefix     = []byte{0x03}

	KeyLastRegoID = append(SequenceKeyPrefix, []byte("lastRegoId")...)
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// GetRegoKey constructs the key for retreiving the ID for the REGO code
func GetRegoKey(regoID uint64) []byte {
	policyIDBz := sdk.Uint64ToBigEndian(regoID)
	return append(RegoKeyPrefix, policyIDBz...)
}

// GetRegoHashKey constructs the key for retreiving the ID for the REGO Hash code
func GetRegoHashKey(regoHash []byte) []byte {
	return append(RegoHashKeyPrefix, regoHash...)
}

