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
	RegoKeyPrefix                                = []byte{0x01}
	PolicyKeyPrefix                              = []byte{0x02}
	SequenceKeyPrefix                            = []byte{0x04}
	PolicyRegoHistoryElementPrefix               = []byte{0x05}
	PolicyByRegoIDAndCreatedSecondaryIndexPrefix = []byte{0x06}
	RegoHashKeyPrefix                            = []byte{0x07}

	KeyLastRegoID     = append(SequenceKeyPrefix, []byte("lastRegoId")...)
	KeyLastInstanceID = append(SequenceKeyPrefix, []byte("lastPolicyId")...)
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// GetPolicyAddressKey returns the key for the REGO policy instance
func GetPolicyAddressKey(addr sdk.AccAddress) []byte {
	return append(PolicyKeyPrefix, addr...)
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

// GetPolicyByCreatedSecondaryIndexKey returns the key for the secondary index:
// `<prefix><regoID><created/last-migrated><policyAddr>`
func GetPolicyByCreatedSecondaryIndexKey(policyAddr sdk.AccAddress, c PolicyRegoHistoryEntry) []byte {
	prefix := GetPolicyByRegoIDSecondaryIndexPrefix(c.RegoID)
	prefixLen := len(prefix)
	policyAddrInvr := sdk.CopyBytes(policyAddr)
	r := make([]byte, prefixLen+AbsoluteTxPositionLen+len(policyAddrInvr))
	copy(r[0:], prefix)
	copy(r[prefixLen:], c.Updated.Bytes())
	copy(r[prefixLen+AbsoluteTxPositionLen:], policyAddr)
	return r
}

// GetPolicyByRegoIDSecondaryIndexPrefix returns the prefix for the second index: `<prefix><regoID>`
func GetPolicyByRegoIDSecondaryIndexPrefix(regoID uint64) []byte {
	prefixLen := len(PolicyByRegoIDAndCreatedSecondaryIndexPrefix)
	const regoIDLen = 8
	r := make([]byte, prefixLen+regoIDLen)
	copy(r[0:], PolicyByRegoIDAndCreatedSecondaryIndexPrefix)
	copy(r[prefixLen:], sdk.Uint64ToBigEndian(regoID))
	return r
}

// GetPolicyRegoHistoryElementPrefix returns the key prefix for a policy rego history entry: `<prefix><policyAddr>`
func GetPolicyRegoHistoryElementPrefix(policyAddr sdk.AccAddress) []byte {
	prefixLen := len(PolicyRegoHistoryElementPrefix)
	policyAddrInvr := sdk.CopyBytes(policyAddr)
	r := make([]byte, prefixLen+len(policyAddrInvr))
	copy(r[0:], PolicyRegoHistoryElementPrefix)
	copy(r[prefixLen:], policyAddr)
	return r
}

// GetPolicyRegoHistoryElementKey returns the key a policy rego history entry: `<prefix><policyAddr><position>`
func GetPolicyRegoHistoryElementKey(policyAddr sdk.AccAddress, pos uint64) []byte {
	prefix := GetPolicyRegoHistoryElementPrefix(policyAddr)
	prefixLen := len(prefix)
	r := make([]byte, prefixLen+8)
	copy(r[0:], prefix)
	copy(r[prefixLen:], sdk.Uint64ToBigEndian(pos))
	return r
}
