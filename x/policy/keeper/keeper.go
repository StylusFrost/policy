package keeper

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/StylusFrost/policy/x/policy/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/open-policy-agent/opa/ast"
)

// CompileCost is how much SDK gas we charge *per byte* for compiling REGO code.
const CompileCost uint64 = 2

type CoinTransferrer interface {
	// TransferCoins sends the coin amounts from the source to the destination with rules applied.
	TransferCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

type (
	Keeper struct {
		cdc           codec.Marshaler
		storeKey      sdk.StoreKey
		memKey        sdk.StoreKey
		paramSpace    paramtypes.Subspace
		accountKeeper types.AccountKeeper
		bank          CoinTransferrer
		queryGasLimit uint64
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	storeKey,
	memKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) *Keeper {

	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramSpace:    paramSpace,
		accountKeeper: accountKeeper,
		bank:          NewBankCoinTransferrer(bankKeeper),
		queryGasLimit: 10000000000, // TODO: Gas Limit
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) create(ctx sdk.Context, creator sdk.AccAddress, regoCode []byte, source string, entry_points []byte, instantiateAccess *types.AccessConfig, authZ AuthorizationPolicy) (regoID uint64, err error) {

	var entryPointsArr []string
	err = json.Unmarshal(entry_points, &entryPointsArr)

	if err != nil {
		return 0, sdkerrors.Wrap(types.ErrCreateFailed, err.Error())
	}

	if !authZ.CanCreateRego(k.getUploadAccessConfig(ctx), creator) {
		return 0, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "can not create rego")
	}

	regoCode, err = uncompress(regoCode, k.GetMaxRegoCodeSize(ctx))
	if err != nil {
		return 0, sdkerrors.Wrap(types.ErrCreateFailed, err.Error())
	}

	// Gas Consume Rego Compiling
	ctx.GasMeter().ConsumeGas(CompileCost*uint64(len(regoCode)), "Compiling REGO Bytecode")

	// Verify if it is valid REGO code
	regoHash, err := k.regoCompile(regoCode)

	if err != nil {
		return 0, sdkerrors.Wrap(types.ErrCompileFailed, err.Error())
	}

	// Store Rego Code
	k.storeRegoHash(ctx, regoHash, regoCode)

	if err != nil {
		return 0, sdkerrors.Wrap(types.ErrCreateFailed, err.Error())
	}
	regoID = k.autoIncrementID(ctx, types.KeyLastRegoID)

	if instantiateAccess == nil {
		defaultAccessConfig := k.getInstantiateAccessConfig(ctx).With(creator)
		instantiateAccess = &defaultAccessConfig
	}

	regoInfo := types.NewRegoInfo(regoHash, creator, source, entryPointsArr, *instantiateAccess)
	k.storeRegoInfo(ctx, regoID, regoInfo)
	return regoID, nil
}

func (k Keeper) storeRegoInfo(ctx sdk.Context, regoID uint64, regoInfo types.RegoInfo) {
	store := ctx.KVStore(k.storeKey)
	// 0x01 | regoID (uint64) -> PolicyInfo
	store.Set(types.GetRegoKey(regoID), k.cdc.MustMarshalBinaryBare(&regoInfo))
}

func (k Keeper) regoCompile(regoCode []byte) ([]byte, error) {

	// Load Rego Policy
	_, err := ast.ParseModule("policy", string(regoCode))

	if err != nil {
		// Load Policy problem
		return nil, err
	}

	// Get SHA256 regoCode
	h := sha256.New()
	h.Write(regoCode)
	return h.Sum(nil), nil
}

func (k Keeper) storeRegoHash(ctx sdk.Context, regoHash []byte, regoCode []byte) {
	store := ctx.KVStore(k.storeKey)
	// Save if not exist
	if !k.hasRegoHash(ctx, regoHash) {
		// 0x03 | regoHash ([]byte) -> regoCode
		store.Set(types.GetRegoHashKey(regoHash), regoCode)
	}
}

// HasRegoHash checks if the Rego Code exists in the store
func (k Keeper) hasRegoHash(ctx sdk.Context, regoHash []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetRegoHashKey(regoHash))
}

func (k Keeper) autoIncrementID(ctx sdk.Context, lastIDKey []byte) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(lastIDKey)
	id := uint64(1)
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	bz = sdk.Uint64ToBigEndian(id + 1)
	store.Set(lastIDKey, bz)
	return id
}

// Querier creates a new grpc querier instance
func Querier(k *Keeper) *grpcQuerier {
	return NewGrpcQuerier(k.cdc, k.storeKey, k, k.queryGasLimit)
}
func (k Keeper) GetByteRego(ctx sdk.Context, regoID uint64) ([]byte, error) {
	store := ctx.KVStore(k.storeKey)
	var regoInfo types.RegoInfo
	regoInfoBz := store.Get(types.GetRegoKey(regoID))
	if regoInfoBz == nil {
		return nil, nil
	}
	k.cdc.MustUnmarshalBinaryBare(regoInfoBz, &regoInfo)

	return k.GetRegoCode(ctx, regoInfo.RegoHash), nil

}
func (k Keeper) GetRegoInfo(ctx sdk.Context, regoID uint64) *types.RegoInfo {
	store := ctx.KVStore(k.storeKey)
	var regoInfo types.RegoInfo
	regoInfoBz := store.Get(types.GetRegoKey(regoID))
	if regoInfoBz == nil {
		return nil
	}
	k.cdc.MustUnmarshalBinaryBare(regoInfoBz, &regoInfo)
	return &regoInfo
}

// QueryGasLimit returns the gas limit for  queries.
func (k Keeper) QueryGasLimit() sdk.Gas {
	return k.queryGasLimit
}

func (k Keeper) GetRegoCode(ctx sdk.Context, regoHash []byte) []byte {
	store := ctx.KVStore(k.storeKey)
	regoCode := store.Get(types.GetRegoHashKey(regoHash))
	if regoCode == nil {
		return nil
	}
	return regoCode
}

func (k Keeper) getUploadAccessConfig(ctx sdk.Context) types.AccessConfig {
	var a types.AccessConfig
	k.paramSpace.Get(ctx, types.ParamStoreKeyUploadAccess, &a)
	return a
}

func (k Keeper) setParams(ctx sdk.Context, ps types.Params) {
	k.paramSpace.SetParamSet(ctx, &ps)
}

func (k Keeper) getInstantiateAccessConfig(ctx sdk.Context) types.AccessType {
	var a types.AccessType
	k.paramSpace.Get(ctx, types.ParamStoreKeyInstantiateAccess, &a)
	return a
}

func (k Keeper) GetMaxRegoCodeSize(ctx sdk.Context) uint64 {
	var a uint64
	k.paramSpace.Get(ctx, types.ParamStoreKeyMaxRegoCodeSize, &a)
	return a
}

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var params types.Params
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

func (k Keeper) IterateRegoInfos(ctx sdk.Context, cb func(uint64, types.RegoInfo) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.RegoKeyPrefix)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var c types.RegoInfo
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &c)
		// cb returns true to stop early
		if cb(binary.BigEndian.Uint64(iter.Key()), c) {
			return
		}
	}
}

func (k Keeper) instantiate(ctx sdk.Context, regoID uint64, creator, admin sdk.AccAddress, entry_points []byte, label string, deposit sdk.Coins, authZ AuthorizationPolicy) (sdk.AccAddress, error) {

	// create policy address
	policyAddress := k.generatePolicyAddress(ctx, regoID)
	existingAcct := k.accountKeeper.GetAccount(ctx, policyAddress)
	if existingAcct != nil {
		return nil, sdkerrors.Wrap(types.ErrAccountExists, existingAcct.GetAddress().String())
	}

	// deposit initial policy funds
	if !deposit.IsZero() {
		if err := k.bank.TransferCoins(ctx, creator, policyAddress, deposit); err != nil {
			return nil, err
		}

	} else {
		// create an empty account (so we don't have issues later)
		// TODO: can we remove this?
		policyAccount := k.accountKeeper.NewAccountWithAddress(ctx, policyAddress)
		k.accountKeeper.SetAccount(ctx, policyAccount)
	}

	// get rego info
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetRegoKey(regoID))
	if bz == nil {
		return nil, sdkerrors.Wrap(types.ErrNotFound, "rego")
	}
	var regoInfo types.RegoInfo
	k.cdc.MustUnmarshalBinaryBare(bz, &regoInfo)

	if !authZ.CanInstantiatePolicy(regoInfo.InstantiateConfig, creator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "can not instantiate")
	}

	//TODO: verify entry_points vs codeInfo.EntryPoints

	// TODO: GAS CONSUME

	// persist instance first
	createdAt := types.NewAbsoluteTxPosition(ctx)
	policyInfo := types.NewPolicyInfo(regoID, creator, admin, label, createdAt, entry_points)

	// store policy before dispatch so that policy could be called back
	historyEntry := policyInfo.InitialHistory(entry_points)
	k.addToPolicyRegoSecondaryIndex(ctx, policyAddress, historyEntry)
	k.appendToPolicyHistory(ctx, policyAddress, historyEntry)
	k.storePolicyInfo(ctx, policyAddress, &policyInfo)

	return policyAddress, nil
}

func (k Keeper) generatePolicyAddress(ctx sdk.Context, regoID uint64) sdk.AccAddress {
	instanceID := k.autoIncrementID(ctx, types.KeyLastInstanceID)
	return BuildPolicyAddress(regoID, instanceID)
}

func BuildPolicyAddress(regoID, instanceID uint64) sdk.AccAddress {
	if regoID > math.MaxUint32 || instanceID > math.MaxUint32 {
		// NOTE: It is possible to get a duplicate address if either regoID or instanceID
		// overflow 32 bits. This is highly improbable, but something that could be refactored.
		panic(fmt.Sprintf("address uint32 reached: regoID: %d, instanceID: %d", regoID, instanceID))
	}
	policyID := regoID<<32 + instanceID
	return addrFromUint64(policyID)
}
func addrFromUint64(id uint64) sdk.AccAddress {
	addr := make([]byte, 20)
	addr[0] = 'C'
	binary.PutUvarint(addr[1:], id)
	return sdk.AccAddress(crypto.AddressHash(addr))
}

// BankCoinTransferrer replicates the cosmos-sdk behaviour as in
// https://github.com/cosmos/cosmos-sdk/blob/v0.41.4/x/bank/keeper/msg_server.go#L26
type BankCoinTransferrer struct {
	keeper types.BankKeeper
}

func NewBankCoinTransferrer(keeper types.BankKeeper) BankCoinTransferrer {
	return BankCoinTransferrer{
		keeper: keeper,
	}
}

// TransferCoins transfers coins from source to destination account when coin send was enabled for them and the recipient
// is not in the blocked address list.
func (c BankCoinTransferrer) TransferCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	if err := c.keeper.SendEnabledCoins(ctx, amt...); err != nil {
		return err
	}
	if c.keeper.BlockedAddr(fromAddr) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "blocked address can not be used")
	}
	sdkerr := c.keeper.SendCoins(ctx, fromAddr, toAddr, amt)
	if sdkerr != nil {
		return sdkerr
	}
	return nil
}

func (k Keeper) addToPolicyRegoSecondaryIndex(ctx sdk.Context, policyAddress sdk.AccAddress, entry types.PolicyRegoHistoryEntry) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetPolicyByCreatedSecondaryIndexKey(policyAddress, entry), []byte{})
}

func (k Keeper) appendToPolicyHistory(ctx sdk.Context, policyAddr sdk.AccAddress, newEntries ...types.PolicyRegoHistoryEntry) {
	store := ctx.KVStore(k.storeKey)
	// find last element position
	var pos uint64
	prefixStore := prefix.NewStore(store, types.GetPolicyRegoHistoryElementPrefix(policyAddr))
	if iter := prefixStore.ReverseIterator(nil, nil); iter.Valid() {
		pos = sdk.BigEndianToUint64(iter.Value())
	}
	// then store with incrementing position
	for _, e := range newEntries {
		pos++
		key := types.GetPolicyRegoHistoryElementKey(policyAddr, pos)
		store.Set(key, k.cdc.MustMarshalBinaryBare(&e))
	}
}

// storePolicyInfo persists the PolicyInfo. No secondary index updated here.
func (k Keeper) storePolicyInfo(ctx sdk.Context, policyAddress sdk.AccAddress, policy *types.PolicyInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetPolicyAddressKey(policyAddress), k.cdc.MustMarshalBinaryBare(policy))
}

func (k Keeper) GetPolicyInfo(ctx sdk.Context, policyAddress sdk.AccAddress) *types.PolicyInfo {
	store := ctx.KVStore(k.storeKey)
	var policy types.PolicyInfo
	policyBz := store.Get(types.GetPolicyAddressKey(policyAddress))
	if policyBz == nil {
		return nil
	}
	k.cdc.MustUnmarshalBinaryBare(policyBz, &policy)
	return &policy
}

func (k Keeper) IteratePoliciesByRegoCode(ctx sdk.Context, regoID uint64, cb func(address sdk.AccAddress) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetPolicyByRegoIDSecondaryIndexPrefix(regoID))
	for iter := prefixStore.Iterator(nil, nil); iter.Valid(); iter.Next() {
		key := iter.Key()
		if cb(key[types.AbsoluteTxPositionLen:]) {
			return
		}
	}
}

func (k Keeper) setPolicyAdmin(ctx sdk.Context, policyAddress, caller, newAdmin sdk.AccAddress, authZ AuthorizationPolicy) error {
	policyInfo := k.GetPolicyInfo(ctx, policyAddress)
	if policyInfo == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "unknown policy")
	}
	if !authZ.CanModifyPolicy(policyInfo.AdminAddr(), caller) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "can not modify policy")
	}
	policyInfo.Admin = newAdmin.String()
	k.storePolicyInfo(ctx, policyAddress, policyInfo)
	return nil
}