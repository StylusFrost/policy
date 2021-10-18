package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoRegoc registers the account types and interface
func RegisterLegacyAminoRegoc(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgStoreRego{}, "policy/MsgStoreRego", nil)
	cdc.RegisterConcrete(&MsgInstantiatePolicy{}, "policy/MsgInstantiatePolicy", nil)
	cdc.RegisterConcrete(&MsgUpdateAdmin{}, "policy/MsgUpdateAdmin", nil)
	cdc.RegisterConcrete(&MsgClearAdmin{}, "policy/MsgClearAdmin", nil)
	cdc.RegisterConcrete(&MsgMigratePolicy{}, "policy/MsgMigratePolicy", nil)
	cdc.RegisterConcrete(&MsgExecutePolicy{}, "policy/MsgExecutePolicy", nil)
	cdc.RegisterConcrete(&MsgRefundPolicy{}, "policy/MsgRefundPolicy", nil)

	cdc.RegisterConcrete(&StoreRegoProposal{}, "policy/StoreRegoProposal", nil)
	cdc.RegisterConcrete(&InstantiatePolicyProposal{}, "policy/InstantiatePolicyProposal", nil)
	cdc.RegisterConcrete(&MigratePolicyProposal{}, "policy/MigratePolicyProposal", nil)
}

func RegisterCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgStoreRego{},
		&MsgInstantiatePolicy{},
		&MsgUpdateAdmin{},
		&MsgClearAdmin{},
		&MsgMigratePolicy{},
		&MsgExecutePolicy{},
		&MsgRefundPolicy{},
	)
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&StoreRegoProposal{},
		&InstantiatePolicyProposal{},
		&MigratePolicyProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterLegacyAminoRegoc(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
