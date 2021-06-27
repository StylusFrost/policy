package keeper

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"

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

type (
	Keeper struct {
		cdc           codec.Marshaler
		storeKey      sdk.StoreKey
		memKey        sdk.StoreKey
		paramSpace    paramtypes.Subspace
		queryGasLimit uint64
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	storeKey,
	memKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,
) *Keeper {

	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramSpace:    paramSpace,
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
