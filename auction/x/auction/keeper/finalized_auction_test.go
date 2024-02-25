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

func createNFinalizedAuction(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.FinalizedAuction {
	items := make([]types.FinalizedAuction, n)
	for i := range items {
		items[i].Id = keeper.AppendFinalizedAuction(ctx, items[i])
	}
	return items
}

func TestFinalizedAuctionGet(t *testing.T) {
	keeper, ctx := keepertest.AuctionKeeper(t)
	items := createNFinalizedAuction(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetFinalizedAuction(ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestFinalizedAuctionRemove(t *testing.T) {
	keeper, ctx := keepertest.AuctionKeeper(t)
	items := createNFinalizedAuction(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveFinalizedAuction(ctx, item.Id)
		_, found := keeper.GetFinalizedAuction(ctx, item.Id)
		require.False(t, found)
	}
}

func TestFinalizedAuctionGetAll(t *testing.T) {
	keeper, ctx := keepertest.AuctionKeeper(t)
	items := createNFinalizedAuction(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllFinalizedAuction(ctx)),
	)
}

func TestFinalizedAuctionCount(t *testing.T) {
	keeper, ctx := keepertest.AuctionKeeper(t)
	items := createNFinalizedAuction(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetFinalizedAuctionCount(ctx))
}
