package keeper

import (
	"context"

	"auction/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) FinalizeAuction(goCtx context.Context, msg *types.MsgFinalizeAuction) (*types.MsgFinalizeAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgFinalizeAuctionResponse{}, nil
}
