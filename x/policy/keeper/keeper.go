package keeper

import (
	"bytes"
	"context"
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
	"github.com/open-policy-agent/opa/rego"
)

// CompileCost is how much SDK gas we charge *per byte* for compiling REGO code.
const CompileCost uint64 = 2

type CoinTransferrer interface {
	// TransferCoins sends the coin amounts from the source to the destination with rules applied.
	TransferCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) int
	GetSupply(ctx sdk.Context, denom string) int
}

type (
	Keeper struct {
		cdc           codec.BinaryCodec
		storeKey      sdk.StoreKey
		memKey        sdk.StoreKey
		paramSpace    paramtypes.Subspace
		accountKeeper types.AccountKeeper
		bank          CoinTransferrer
		queryGasLimit uint64
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) *Keeper {

	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	RegistryCosmosBuiltins()

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

func (k Keeper) execute(ctx sdk.Context, policyAddress sdk.AccAddress, caller sdk.AccAddress, entry_point string, input []byte, coins sdk.Coins) ([]byte, error) {

	// TODO: GAS COST

	// Get Policy info
	policyInfo := k.GetPolicyInfo(ctx, policyAddress)
	if policyInfo == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "unknown policy")
	}

	// Verify execute entry point is valid Policy Entry Points

	var entryPointsArr []types.EntryPoint
	err := json.Unmarshal(policyInfo.EntryPoints, &entryPointsArr)

	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrExecuteFailed, err.Error())
	}

	var found bool

	// Verify entrypoints
	for _, entry := range entryPointsArr {
		if entry.Entry == entry_point {
			found = true
		}

	}
	if !found {
		return nil, sdkerrors.Wrap(types.ErrExecuteFailed, "Invalid Entry Point")
	}
	// get rego code

	regoCode, err := k.GetByteRego(ctx, policyInfo.RegoID)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrExecuteFailed, err.Error())
	}

	// add more funds
	if !coins.IsZero() {
		if err := k.bank.TransferCoins(ctx, caller, policyAddress, coins); err != nil {
			return nil, err
		}
	}

	// Execute Policy
	data, err := k.regoExecute(ctx, regoCode, entry_point, input, caller, policyAddress)

	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrExecuteFailed, err.Error())
	}

	return data, nil
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

	// Verify if it is valid REGO code and entry_points are valid
	regoHash, err := k.regoValidateAndCompile(ctx, regoCode, entryPointsArr)

	if err != nil {
		return 0, sdkerrors.Wrap(types.ErrCompileFailed, err.Error())
	}

	// Store Rego Code
	k.storeRegoHash(ctx, regoHash, regoCode)

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
	store.Set(types.GetRegoKey(regoID), k.cdc.MustMarshal(&regoInfo))
}

func (k Keeper) regoValidateAndCompile(ctx sdk.Context, regoCode []byte, entry_points []string) ([]byte, error) {

	// Load Rego Policy
	mod, err := ast.ParseModule("policy", string(regoCode))

	if err != nil {
		// Load Policy parse problem
		return nil, err
	}

	// Create a new compiler instance and compile the module.
	compiler := ast.NewCompiler()

	mods := map[string]*ast.Module{
		"policy": mod,
	}

	if compiler.Compile(mods); compiler.Failed() {
		return nil, compiler.Errors
	}

	// Verify entrypoints
	for _, entry := range entry_points {
		if len(compiler.GetRulesWithPrefix(ast.MustParseRef("data."+entry))) == 0 {
			return nil, sdkerrors.Wrap(types.ErrEntryPoint, entry)
		}
	}

	// Get SHA256 regoCode
	h := sha256.New()
	h.Write(regoCode)
	return h.Sum(nil), nil
}

func (k Keeper) regoExecute(ctx sdk.Context, regoCode []byte, entry_point string, input []byte, caller sdk.AccAddress, policyAddress sdk.AccAddress) ([]byte, error) {

	ctxRego := context.Background()

	// Load Rego Policy
	mod, err := ast.ParseModule("policy", string(regoCode))

	if err != nil {
		// Load Policy parse problem
		return nil, err
	}

	// Create a new compiler instance and compile the module.
	compiler := ast.NewCompiler()

	mods := map[string]*ast.Module{
		"policy": mod,
	}

	// Load Cosmos Builtins
	k.LoadCosmosBuiltins(ctx, caller, policyAddress)

	if compiler.Compile(mods); compiler.Failed() {
		return nil, compiler.Errors
	}

	// Verify entrypoint
	if len(compiler.GetRulesWithPrefix(ast.MustParseRef("data."+entry_point))) == 0 {
		return nil, sdkerrors.Wrap(types.ErrEntryPoint, entry_point)
	}

	// Raw input data that will be used in evaluation.
	raw := string(input)
	d := json.NewDecoder(bytes.NewBufferString(raw))

	var inputDecode interface{}

	if err := d.Decode(&inputDecode); err != nil {
		return nil, sdkerrors.Wrap(types.ErrExecuteFailed, err.Error())
	}

	// Create a new query that uses the compiled policy from above.
	rego := rego.New(
		rego.Query("data."+entry_point),
		rego.Compiler(compiler),
		rego.Input(inputDecode),
	)
	// Run evaluation.
	rs, err := rego.Eval(ctxRego)

	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrExecuteFailed, err.Error())
	}

	if len(rs) == 0 {
		return nil, types.ErrExecuteFailed
	}

	data, err := json.Marshal(rs[0].Expressions[0].Value)

	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrExecuteFailed, err.Error())
	}

	return data, nil
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
	k.cdc.MustUnmarshal(regoInfoBz, &regoInfo)

	return k.GetRegoCode(ctx, regoInfo.RegoHash), nil

}
func (k Keeper) GetRegoInfo(ctx sdk.Context, regoID uint64) *types.RegoInfo {
	store := ctx.KVStore(k.storeKey)
	var regoInfo types.RegoInfo
	regoInfoBz := store.Get(types.GetRegoKey(regoID))
	if regoInfoBz == nil {
		return nil
	}
	k.cdc.MustUnmarshal(regoInfoBz, &regoInfo)
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
		k.cdc.MustUnmarshal(iter.Value(), &c)
		// cb returns true to stop early
		if cb(binary.BigEndian.Uint64(iter.Key()), c) {
			return
		}
	}
}

func (k Keeper) instantiate(ctx sdk.Context, regoID uint64, creator, admin sdk.AccAddress, entry_points []byte, label string, deposit sdk.Coins, authZ AuthorizationPolicy) (sdk.AccAddress, error) {

	var entryPointsArr []types.EntryPoint
	err := json.Unmarshal(entry_points, &entryPointsArr)

	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, err.Error())
	}

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

	regoInfo := k.GetRegoInfo(ctx, regoID)

	if !authZ.CanInstantiatePolicy(regoInfo.InstantiateConfig, creator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "can not instantiate")
	}

	// Verify if all msg Entry points exist into rego id
	for _, v := range entryPointsArr {
		if !contains(regoInfo.EntryPoints, v.Entry) {
			return nil, sdkerrors.Wrap(types.ErrEntryPoint, "entry point")
		}
	}

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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
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
	if err := c.keeper.IsSendEnabledCoins(ctx, amt...); err != nil {
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

// GetTotalSupplyAmount get total Supply from input denom
func (c BankCoinTransferrer) HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool {

	return c.keeper.HasBalance(ctx, addr, amt)

}

// GetBalance
func (c BankCoinTransferrer) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) int {

	balanceCoin := c.keeper.GetBalance(ctx, addr, denom)
	return int(balanceCoin.Amount.Int64())

}

// GetSupply
func (c BankCoinTransferrer) GetSupply(ctx sdk.Context, denom string) int {

	coin := c.keeper.GetSupply(ctx, denom)

	if coin.Denom == denom {
		return int(coin.Amount.Int64())
	}
	return 0

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
		store.Set(key, k.cdc.MustMarshal(&e))
	}
}

// storePolicyInfo persists the PolicyInfo. No secondary index updated here.
func (k Keeper) storePolicyInfo(ctx sdk.Context, policyAddress sdk.AccAddress, policy *types.PolicyInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetPolicyAddressKey(policyAddress), k.cdc.MustMarshal(policy))
}

func (k Keeper) GetPolicyInfo(ctx sdk.Context, policyAddress sdk.AccAddress) *types.PolicyInfo {
	store := ctx.KVStore(k.storeKey)
	var policy types.PolicyInfo
	policyBz := store.Get(types.GetPolicyAddressKey(policyAddress))
	if policyBz == nil {
		return nil
	}
	k.cdc.MustUnmarshal(policyBz, &policy)
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

func (k Keeper) migrate(ctx sdk.Context, policyAddress sdk.AccAddress, caller sdk.AccAddress, newRegoID uint64, entry_points []byte, authZ AuthorizationPolicy) error {

	// TODO: GAS COST

	policyInfo := k.GetPolicyInfo(ctx, policyAddress)
	if policyInfo == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "unknown policy")
	}
	if !authZ.CanModifyPolicy(policyInfo.AdminAddr(), caller) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "can not migrate")
	}

	newCodeInfo := k.GetRegoInfo(ctx, newRegoID)
	if newCodeInfo == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "unknown rego")
	}

	//TODO: verify entry_points vs newCodeInfo.EntryPoints

	policyInfo.RegoID = newRegoID
	policyInfo.EntryPoints = entry_points

	// delete old secondary index entry
	k.removeFromPolicyRegoSecondaryIndex(ctx, policyAddress, k.getLastPolicyHistoryEntry(ctx, policyAddress))
	// persist migration updates
	historyEntry := policyInfo.AddMigration(ctx, newRegoID, entry_points)
	k.appendToPolicyHistory(ctx, policyAddress, historyEntry)
	k.addToPolicyRegoSecondaryIndex(ctx, policyAddress, historyEntry)
	k.storePolicyInfo(ctx, policyAddress, policyInfo)

	return nil
}
func (k Keeper) refund(ctx sdk.Context, policyAddress sdk.AccAddress, caller sdk.AccAddress, coins sdk.Coins, authZ AuthorizationPolicy) error {

	// TODO: GAS COST

	policyInfo := k.GetPolicyInfo(ctx, policyAddress)
	if policyInfo == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "unknown policy")
	}
	if !authZ.CanRefundPolicy(policyInfo.AdminAddr(), caller) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "can not refund")
	}

	// refund coins
	if !coins.IsZero() {
		if err := k.bank.TransferCoins(ctx, policyAddress, caller, coins); err != nil {
			return err
		}
	}

	return nil
}

// removeFromPolicyCodeSecondaryIndex removes element to the index for policies-by-regoid queries
func (k Keeper) removeFromPolicyRegoSecondaryIndex(ctx sdk.Context, policyAddress sdk.AccAddress, entry types.PolicyRegoHistoryEntry) {
	ctx.KVStore(k.storeKey).Delete(types.GetPolicyByCreatedSecondaryIndexKey(policyAddress, entry))
}

func (k Keeper) getLastPolicyHistoryEntry(ctx sdk.Context, policyAddr sdk.AccAddress) types.PolicyRegoHistoryEntry {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetPolicyRegoHistoryElementPrefix(policyAddr))
	iter := prefixStore.ReverseIterator(nil, nil)
	var r types.PolicyRegoHistoryEntry
	if !iter.Valid() {
		// all policys have a history
		panic(fmt.Sprintf("no history for %s", policyAddr.String()))
	}
	k.cdc.MustUnmarshal(iter.Value(), &r)
	return r
}

func (k Keeper) GetPolicyHistory(ctx sdk.Context, policyAddr sdk.AccAddress) []types.PolicyRegoHistoryEntry {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetPolicyRegoHistoryElementPrefix(policyAddr))
	r := make([]types.PolicyRegoHistoryEntry, 0)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var e types.PolicyRegoHistoryEntry
		k.cdc.MustUnmarshal(iter.Value(), &e)
		r = append(r, e)
	}
	return r
}
