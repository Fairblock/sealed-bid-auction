package auction_test

import (
	"testing"

	keepertest "auction/testutil/keeper"
	"auction/testutil/nullify"
	"auction/x/auction"
	"auction/x/auction/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		AuctionList: []types.Auction{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		AuctionCount: 2,
		BidList: []types.Bid{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		BidCount: 2,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.AuctionKeeper(t)
	auction.InitGenesis(ctx, *k, genesisState)
	got := auction.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.AuctionList, got.AuctionList)
	require.Equal(t, genesisState.AuctionCount, got.AuctionCount)
	require.ElementsMatch(t, genesisState.BidList, got.BidList)
	require.Equal(t, genesisState.BidCount, got.BidCount)
	// this line is used by starport scaffolding # genesis/test/assert
}
