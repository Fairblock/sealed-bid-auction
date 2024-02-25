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

func createNAuction(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Auction {
	items := make([]types.Auction, n)
	for i := range items {
		items[i].Id = keeper.AppendAuction(ctx, items[i])
	}
	return items
}

func TestAuctionGet(t *testing.T) {
	keeper, ctx := keepertest.AuctionKeeper(t)
	items := createNAuction(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetAuction(ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestAuctionRemove(t *testing.T) {
	keeper, ctx := keepertest.AuctionKeeper(t)
	items := createNAuction(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveAuction(ctx, item.Id)
		_, found := keeper.GetAuction(ctx, item.Id)
		require.False(t, found)
	}
}

func TestAuctionGetAll(t *testing.T) {
	keeper, ctx := keepertest.AuctionKeeper(t)
	items := createNAuction(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllAuction(ctx)),
	)
}

func TestAuctionCount(t *testing.T) {
	keeper, ctx := keepertest.AuctionKeeper(t)
	items := createNAuction(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetAuctionCount(ctx))
}
