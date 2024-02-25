package keeper

import (
	"context"

	"auction/x/auction/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) FinalizedAuctionAll(goCtx context.Context, req *types.QueryAllFinalizedAuctionRequest) (*types.QueryAllFinalizedAuctionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var finalizedAuctions []types.FinalizedAuction
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	finalizedAuctionStore := prefix.NewStore(store, types.KeyPrefix(types.FinalizedAuctionKey))

	pageRes, err := query.Paginate(finalizedAuctionStore, req.Pagination, func(key []byte, value []byte) error {
		var finalizedAuction types.FinalizedAuction
		if err := k.cdc.Unmarshal(value, &finalizedAuction); err != nil {
			return err
		}

		finalizedAuctions = append(finalizedAuctions, finalizedAuction)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllFinalizedAuctionResponse{FinalizedAuction: finalizedAuctions, Pagination: pageRes}, nil
}

func (k Keeper) FinalizedAuction(goCtx context.Context, req *types.QueryGetFinalizedAuctionRequest) (*types.QueryGetFinalizedAuctionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	finalizedAuction, found := k.GetFinalizedAuction(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetFinalizedAuctionResponse{FinalizedAuction: finalizedAuction}, nil
}
