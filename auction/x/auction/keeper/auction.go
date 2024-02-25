package keeper

import (
	"encoding/binary"

	"auction/x/auction/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAuctionCount get the total number of auction
func (k Keeper) GetAuctionCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.AuctionCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// UpdateAuctionHighestBidId set the highest bid id of auction
func (k Keeper) UpdateAuctionHighestBidId(ctx sdk.Context, id uint64, bidId uint64) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))

	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)

	b := store.Get(bz)
	if b == nil {
		return types.AuctionNotFound
	}

	var auction types.Auction
	k.cdc.MustUnmarshal(b, &auction)

	auction.HighestBidExists = true
	auction.CurrentHighestBidId = bidId

	appendedValue := k.cdc.MustMarshal(&auction)

	store.Set(bz, appendedValue)
	return nil
}

// EndAuction set the auction status to ended
func (k Keeper) EndAuction(ctx sdk.Context, id uint64) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))

	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)

	b := store.Get(bz)
	if b == nil {
		return types.AuctionNotFound
	}

	var auction types.Auction
	k.cdc.MustUnmarshal(b, &auction)
	auction.Ended = true

	appendedValue := k.cdc.MustMarshal(&auction)

	store.Set(bz, appendedValue)
	return nil
}

// SetAuctionCount set the total number of auction
func (k Keeper) SetAuctionCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.AuctionCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendAuction appends a auction in the store with a new id and update the count
func (k Keeper) AppendAuction(
	ctx sdk.Context,
	auction types.Auction,
) uint64 {
	// Create the auction
	count := k.GetAuctionCount(ctx)

	// Set the ID of the appended value
	auction.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))
	appendedValue := k.cdc.MustMarshal(&auction)
	store.Set(GetAuctionIDBytes(auction.Id), appendedValue)

	// Update auction count
	k.SetAuctionCount(ctx, count+1)

	return count
}

// SetAuction set a specific auction in the store
func (k Keeper) SetAuction(ctx sdk.Context, auction types.Auction) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))
	b := k.cdc.MustMarshal(&auction)
	store.Set(GetAuctionIDBytes(auction.Id), b)
}

// GetAuction returns a auction from its id
func (k Keeper) GetAuction(ctx sdk.Context, id uint64) (val types.Auction, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))
	b := store.Get(GetAuctionIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveAuction removes a auction from the store
func (k Keeper) RemoveAuction(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))
	store.Delete(GetAuctionIDBytes(id))
}

// GetAllAuction returns all auction
func (k Keeper) GetAllAuction(ctx sdk.Context) (list []types.Auction) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Auction
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAuctionIDBytes returns the byte representation of the ID
func GetAuctionIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetAuctionIDFromBytes returns ID in uint64 format from a byte array
func GetAuctionIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
