package keeper_test

import (
	"testing"

	keepertest "auction/testutil/keeper"
	"auction/testutil/nullify"
	"auction/x/auction/keeper"
	"auction/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func createNBid(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Bid {
	items := make([]types.Bid, n)
	for i := range items {
		items[i].Id = keeper.AppendBid(ctx, items[i])
	}
	return items
}

func TestBidGet(t *testing.T) {
	keeper, ctx := keepertest.AuctionKeeper(t)
	items := createNBid(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetBid(ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestBidRemove(t *testing.T) {
	keeper, ctx := keepertest.AuctionKeeper(t)
	items := createNBid(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveBid(ctx, item.Id)
		_, found := keeper.GetBid(ctx, item.Id)
		require.False(t, found)
	}
}

func TestBidGetAll(t *testing.T) {
	keeper, ctx := keepertest.AuctionKeeper(t)
	items := createNBid(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllBid(ctx)),
	)
}

func TestBidCount(t *testing.T) {
	keeper, ctx := keepertest.AuctionKeeper(t)
	items := createNBid(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetBidCount(ctx))
}
