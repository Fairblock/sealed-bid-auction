package keeper

import (
	"encoding/binary"

	"auction/x/auction/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetFinalizedAuctionCount get the total number of finalizedAuction
func (k Keeper) GetFinalizedAuctionCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.FinalizedAuctionCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetFinalizedAuctionCount set the total number of finalizedAuction
func (k Keeper) SetFinalizedAuctionCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.FinalizedAuctionCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendFinalizedAuction appends a finalizedAuction in the store with a new id and update the count
func (k Keeper) AppendFinalizedAuction(
	ctx sdk.Context,
	finalizedAuction types.FinalizedAuction,
) uint64 {
	// Create the finalizedAuction
	count := k.GetFinalizedAuctionCount(ctx)

	// Set the ID of the appended value
	finalizedAuction.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FinalizedAuctionKey))
	appendedValue := k.cdc.MustMarshal(&finalizedAuction)
	store.Set(GetFinalizedAuctionIDBytes(finalizedAuction.Id), appendedValue)

	// Update finalizedAuction count
	k.SetFinalizedAuctionCount(ctx, count+1)

	return count
}

// SetFinalizedAuction set a specific finalizedAuction in the store
func (k Keeper) SetFinalizedAuction(ctx sdk.Context, finalizedAuction types.FinalizedAuction) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FinalizedAuctionKey))
	b := k.cdc.MustMarshal(&finalizedAuction)
	store.Set(GetFinalizedAuctionIDBytes(finalizedAuction.Id), b)
}

// GetFinalizedAuction returns a finalizedAuction from its id
func (k Keeper) GetFinalizedAuction(ctx sdk.Context, id uint64) (val types.FinalizedAuction, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FinalizedAuctionKey))
	b := store.Get(GetFinalizedAuctionIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveFinalizedAuction removes a finalizedAuction from the store
func (k Keeper) RemoveFinalizedAuction(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FinalizedAuctionKey))
	store.Delete(GetFinalizedAuctionIDBytes(id))
}

// GetAllFinalizedAuction returns all finalizedAuction
func (k Keeper) GetAllFinalizedAuction(ctx sdk.Context) (list []types.FinalizedAuction) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FinalizedAuctionKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.FinalizedAuction
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetFinalizedAuctionIDBytes returns the byte representation of the ID
func GetFinalizedAuctionIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetFinalizedAuctionIDFromBytes returns ID in uint64 format from a byte array
func GetFinalizedAuctionIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
