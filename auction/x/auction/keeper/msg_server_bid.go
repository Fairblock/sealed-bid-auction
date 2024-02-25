package keeper

import (
	"auction/x/auction/types"
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateBid(goCtx context.Context, msg *types.MsgCreateBid) (*types.MsgCreateBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the auction ID in the Create Bid Message exists
	auction, found := k.GetAuction(ctx, msg.AuctionId)
	if !found {
		return nil, types.AuctionNotFound
	}

	// If the auction exists, check if it is ended
	if auction.Ended {
		return nil, types.AuctionEnded
	}

	// Check If the auction already passed the duration
	if auction.CreatedAt+auction.Duration <= uint64(ctx.BlockHeight()) {
		return nil, types.AuctionDurationPassed
	}

	bidPriceInCoins := sdk.NewCoins(msg.BidPrice)

	// Check if the bid price is greater than the auction start price / current highest bid
	if auction.StartPrice.IsGTE(msg.BidPrice) {
		return nil, types.BidPriceLow
	}

	sender, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	// Check if bid creator have sufficient coin for placing the bid
	senderBalance := k.bankKeeper.SpendableCoins(ctx, sender)

	if bidPriceInCoins.IsAllGT(senderBalance) {
		return nil, types.InsufficientBalance
	}

	if auction.HighestBidExists {
		// If highest bid exists, get the current highest bid
		currentHighestBid, found := k.GetBid(ctx, auction.CurrentHighestBidId)
		if !found {
			return nil, types.InternalError
		}

		// Check if the current highest bid price is greater or equals to the bid price
		if currentHighestBid.BidPrice.IsGTE(msg.BidPrice) {
			return nil, types.BidPriceLow
		}

		// Current bid price is higher than the highest bid price, returning the coins to the highest bid creator
		receiver, err := sdk.AccAddressFromBech32(currentHighestBid.Creator)
		if err != nil {
			return nil, err
		}

		if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(currentHighestBid.BidPrice)); err != nil {
			return nil, err
		}
	}

	var bid = types.Bid{
		Creator:   msg.Creator,
		AuctionId: msg.AuctionId,
		BidPrice:  msg.BidPrice,
	}

	id := k.AppendBid(ctx, bid)

	// Update the highest bid to current bid
	if err := k.UpdateAuctionHighestBidId(ctx, msg.AuctionId, id); err != nil {
		return nil, err
	}

	// Send the coin from creator to the module
	if err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		sender,
		types.ModuleName,
		bidPriceInCoins,
	); err != nil {
		return nil, err
	}

	return &types.MsgCreateBidResponse{
		Id: id,
	}, nil
}
