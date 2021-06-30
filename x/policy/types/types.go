package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AbsoluteTxPositionLen number of elements in byte representation
const AbsoluteTxPositionLen = 16

func NewRegoInfo(regoHash []byte, creator sdk.AccAddress, source string, entry_points []string, instantiatePermission AccessConfig) RegoInfo {
	return RegoInfo{
		RegoHash:          regoHash,
		Creator:           creator.String(),
		Source:            source,
		EntryPoints:       entry_points,
		InstantiateConfig: instantiatePermission,
	}
}
func NewPolicyInfo(regoID uint64, creator, admin sdk.AccAddress, label string, createdAt *AbsoluteTxPosition, entry_points []byte) PolicyInfo {
	var adminAddr string
	if !admin.Empty() {
		adminAddr = admin.String()
	}
	return PolicyInfo{
		RegoID:      regoID,
		Creator:     creator.String(),
		Admin:       adminAddr,
		Label:       label,
		Created:     createdAt,
		EntryPoints: entry_points,
	}
}

func NewAbsoluteTxPosition(ctx sdk.Context) *AbsoluteTxPosition {
	// we must safely handle nil gas meters
	var index uint64
	meter := ctx.BlockGasMeter()
	if meter != nil {
		index = meter.GasConsumed()
	}
	height := ctx.BlockHeight()
	if height < 0 {
		panic(fmt.Sprintf("unsupported height: %d", height))
	}
	return &AbsoluteTxPosition{
		BlockHeight: uint64(height),
		TxIndex:     index,
	}
}

func (c PolicyInfo) InitialHistory(entry_points []byte) PolicyRegoHistoryEntry {
	return PolicyRegoHistoryEntry{
		Operation:   PolicyRegoHistoryOperationTypeInit,
		RegoID:      c.RegoID,
		Updated:     c.Created,
		EntryPoints: entry_points,
	}
}

func (a *AbsoluteTxPosition) Bytes() []byte {
	if a == nil {
		panic("object must not be nil")
	}
	r := make([]byte, AbsoluteTxPositionLen)
	copy(r[0:], sdk.Uint64ToBigEndian(a.BlockHeight))
	copy(r[8:], sdk.Uint64ToBigEndian(a.TxIndex))
	return r
}

func (c *PolicyInfo) AdminAddr() sdk.AccAddress {
	if c.Admin == "" {
		return nil
	}
	admin, err := sdk.AccAddressFromBech32(c.Admin)
	if err != nil { // should never happen
		panic(err.Error())
	}
	return admin
}

func (c *PolicyInfo) AddMigration(ctx sdk.Context, regoID uint64, entry_points []byte) PolicyRegoHistoryEntry {
	h := PolicyRegoHistoryEntry{
		Operation:   PolicyRegoHistoryOperationTypeMigrate,
		RegoID:      regoID,
		Updated:     NewAbsoluteTxPosition(ctx),
		EntryPoints: entry_points,
	}
	c.RegoID = regoID
	return h
}
