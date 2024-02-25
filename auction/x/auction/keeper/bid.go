package keeper

import (
	"encoding/binary"

	"auction/x/auction/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetBidCount get the total number of bid
func (k Keeper) GetBidCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.BidCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetBidCount set the total number of bid
func (k Keeper) SetBidCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.BidCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendBid appends a bid in the store with a new id and update the count
func (k Keeper) AppendBid(
	ctx sdk.Context,
	bid types.Bid,
) uint64 {
	// Create the bid
	count := k.GetBidCount(ctx)

	// Set the ID of the appended value
	bid.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BidKey))
	appendedValue := k.cdc.MustMarshal(&bid)
	store.Set(GetBidIDBytes(bid.Id), appendedValue)

	// Update bid count
	k.SetBidCount(ctx, count+1)

	return count
}

// SetBid set a specific bid in the store
func (k Keeper) SetBid(ctx sdk.Context, bid types.Bid) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BidKey))
	b := k.cdc.MustMarshal(&bid)
	store.Set(GetBidIDBytes(bid.Id), b)
}

// GetBid returns a bid from its id
func (k Keeper) GetBid(ctx sdk.Context, id uint64) (val types.Bid, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BidKey))
	b := store.Get(GetBidIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveBid removes a bid from the store
func (k Keeper) RemoveBid(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BidKey))
	store.Delete(GetBidIDBytes(id))
}

// GetAllBid returns all bid
func (k Keeper) GetAllBid(ctx sdk.Context) (list []types.Bid) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BidKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Bid
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetBidIDBytes returns the byte representation of the ID
func GetBidIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetBidIDFromBytes returns ID in uint64 format from a byte array
func GetBidIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
