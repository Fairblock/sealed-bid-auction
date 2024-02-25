package types

const (
	// ModuleName defines the module name
	ModuleName = "auction"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_auction"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	AuctionKey      = "Auction/value/"
	AuctionCountKey = "Auction/count/"
)

const (
	BidKey      = "Bid/value/"
	BidCountKey = "Bid/count/"
)

const (
	FinalizedAuctionKey      = "FinalizedAuction/value/"
	FinalizedAuctionCountKey = "FinalizedAuction/count/"
)
