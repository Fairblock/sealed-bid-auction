package keeper

import (
	"context"

	"auction/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) FinalizeAuction(goCtx context.Context, msg *types.MsgFinalizeAuction) (*types.MsgFinalizeAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the provided auction exists
	auction, found := k.GetAuction(ctx, msg.AuctionId)
	if !found {
		return nil, types.AuctionNotFound
	}

	// Check if the provided auction is ended
	if auction.Ended {
		return nil, types.AuctionEnded
	}

	// Make sure the message sender is the creator of the auction
	if auction.Creator != msg.Creator {
		return nil, types.NotAuctionOwner
	}

	// Make sure the auction passed its duration
	if auction.CreatedAt+auction.Duration > uint64(ctx.BlockHeight()) {
		return nil, types.AuctionFinalizeTooEarly
	}

	// Get the highest bid price
	finalBidPrice := sdk.NewCoin("stake", sdk.NewInt(0))
	bidCreator := ""
	if auction.HighestBidExists {
		bid, found := k.GetBid(ctx, auction.CurrentHighestBidId)
		if found {
			finalBidPrice = bid.BidPrice
			bidCreator = bid.Creator
		}
	}

	finalizedAuction := types.FinalizedAuction{
		AuctionId:  msg.AuctionId,
		BidId:      auction.CurrentHighestBidId,
		FinalPrice: finalBidPrice,
		Bidder:     bidCreator,
		Creator:    msg.Creator,
	}

	id := k.AppendFinalizedAuction(ctx, finalizedAuction)

	// End the auction
	if err := k.EndAuction(ctx, msg.AuctionId); err != nil {
		return nil, err
	}

	// If there is a bid, send the coins to auction creator
	if !finalBidPrice.IsZero() {
		receiver, err := sdk.AccAddressFromBech32(msg.Creator)
		if err != nil {
			return nil, err
		}

		if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(finalBidPrice)); err != nil {
			return nil, err
		}
	}

	return &types.MsgFinalizeAuctionResponse{
		Id:         id,
		FinalPrice: finalBidPrice,
		Bidder:     bidCreator,
	}, nil
}
