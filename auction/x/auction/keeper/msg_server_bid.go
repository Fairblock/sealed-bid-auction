package keeper

import (
	"context"
	"fmt"

	"auction/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateBid(goCtx context.Context, msg *types.MsgCreateBid) (*types.MsgCreateBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var bid = types.Bid{
		Creator:   msg.Creator,
		AuctionId: msg.AuctionId,
		BidPrice:  msg.BidPrice,
	}

	id := k.AppendBid(
		ctx,
		bid,
	)

	return &types.MsgCreateBidResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdateBid(goCtx context.Context, msg *types.MsgUpdateBid) (*types.MsgUpdateBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var bid = types.Bid{
		Creator:   msg.Creator,
		Id:        msg.Id,
		AuctionId: msg.AuctionId,
		BidPrice:  msg.BidPrice,
	}

	// Checks that the element exists
	val, found := k.GetBid(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetBid(ctx, bid)

	return &types.MsgUpdateBidResponse{}, nil
}

func (k msgServer) DeleteBid(goCtx context.Context, msg *types.MsgDeleteBid) (*types.MsgDeleteBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	val, found := k.GetBid(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveBid(ctx, msg.Id)

	return &types.MsgDeleteBidResponse{}, nil
}
