package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		AuctionList:          []Auction{},
		BidList:              []Bid{},
		FinalizedAuctionList: []FinalizedAuction{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated ID in auction
	auctionIdMap := make(map[uint64]bool)
	auctionCount := gs.GetAuctionCount()
	for _, elem := range gs.AuctionList {
		if _, ok := auctionIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for auction")
		}
		if elem.Id >= auctionCount {
			return fmt.Errorf("auction id should be lower or equal than the last id")
		}
		auctionIdMap[elem.Id] = true
	}
	// Check for duplicated ID in bid
	bidIdMap := make(map[uint64]bool)
	bidCount := gs.GetBidCount()
	for _, elem := range gs.BidList {
		if _, ok := bidIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for bid")
		}
		if elem.Id >= bidCount {
			return fmt.Errorf("bid id should be lower or equal than the last id")
		}
		bidIdMap[elem.Id] = true
	}
	// Check for duplicated ID in finalizedAuction
	finalizedAuctionIdMap := make(map[uint64]bool)
	finalizedAuctionCount := gs.GetFinalizedAuctionCount()
	for _, elem := range gs.FinalizedAuctionList {
		if _, ok := finalizedAuctionIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for finalizedAuction")
		}
		if elem.Id >= finalizedAuctionCount {
			return fmt.Errorf("finalizedAuction id should be lower or equal than the last id")
		}
		finalizedAuctionIdMap[elem.Id] = true
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
