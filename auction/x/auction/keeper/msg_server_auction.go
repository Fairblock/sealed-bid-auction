package keeper

import (
	"context"
	"fmt"

	"auction/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateAuction(goCtx context.Context, msg *types.MsgCreateAuction) (*types.MsgCreateAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var auction = types.Auction{
		Creator:             msg.Creator,
		Name:                msg.Name,
		StartPrice:          msg.StartPrice,
		Duration:            msg.Duration,
		CreatedAt:           msg.CreatedAt,
		CurrentHighestBidId: msg.CurrentHighestBidId,
		HighestBidExists:    msg.HighestBidExists,
		Ended:               msg.Ended,
	}

	id := k.AppendAuction(
		ctx,
		auction,
	)

	return &types.MsgCreateAuctionResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdateAuction(goCtx context.Context, msg *types.MsgUpdateAuction) (*types.MsgUpdateAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var auction = types.Auction{
		Creator:             msg.Creator,
		Id:                  msg.Id,
		Name:                msg.Name,
		StartPrice:          msg.StartPrice,
		Duration:            msg.Duration,
		CreatedAt:           msg.CreatedAt,
		CurrentHighestBidId: msg.CurrentHighestBidId,
		HighestBidExists:    msg.HighestBidExists,
		Ended:               msg.Ended,
	}

	// Checks that the element exists
	val, found := k.GetAuction(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetAuction(ctx, auction)

	return &types.MsgUpdateAuctionResponse{}, nil
}

func (k msgServer) DeleteAuction(goCtx context.Context, msg *types.MsgDeleteAuction) (*types.MsgDeleteAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	val, found := k.GetAuction(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveAuction(ctx, msg.Id)

	return &types.MsgDeleteAuctionResponse{}, nil
}
