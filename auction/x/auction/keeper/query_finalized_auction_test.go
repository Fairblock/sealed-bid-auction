package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "auction/testutil/keeper"
	"auction/testutil/nullify"
	"auction/x/auction/types"
)

func TestFinalizedAuctionQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.AuctionKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNFinalizedAuction(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetFinalizedAuctionRequest
		response *types.QueryGetFinalizedAuctionResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetFinalizedAuctionRequest{Id: msgs[0].Id},
			response: &types.QueryGetFinalizedAuctionResponse{FinalizedAuction: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetFinalizedAuctionRequest{Id: msgs[1].Id},
			response: &types.QueryGetFinalizedAuctionResponse{FinalizedAuction: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetFinalizedAuctionRequest{Id: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.FinalizedAuction(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestFinalizedAuctionQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.AuctionKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNFinalizedAuction(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllFinalizedAuctionRequest {
		return &types.QueryAllFinalizedAuctionRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.FinalizedAuctionAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.FinalizedAuction), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.FinalizedAuction),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.FinalizedAuctionAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.FinalizedAuction), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.FinalizedAuction),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.FinalizedAuctionAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.FinalizedAuction),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.FinalizedAuctionAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
