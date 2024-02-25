package keeper

import (
	"auction/x/auction/types"
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateAuction(goCtx context.Context, msg *types.MsgCreateAuction) (*types.MsgCreateAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.StartPrice.IsZero() {
		return nil, types.AuctionPriceInvalid
	}

	if msg.Duration < 10 {
		return nil, types.AuctionDurationInvalid
	}

	var auction = types.Auction{
		Creator:    msg.Creator,
		Name:       msg.Name,
		StartPrice: msg.StartPrice,
		Duration:   msg.Duration,
		CreatedAt:  uint64(ctx.BlockHeight()),
		Ended:      false,
	}

	id := k.AppendAuction(
		ctx,
		auction,
	)

	return &types.MsgCreateAuctionResponse{
		Id: id,
	}, nil
}
