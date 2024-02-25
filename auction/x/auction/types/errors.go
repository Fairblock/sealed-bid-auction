package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/auction module sentinel errors
var (
	AuctionEnded            = sdkerrors.Register(ModuleName, 1100, "target auction already ended")
	AuctionNotFound         = sdkerrors.Register(ModuleName, 1200, "target auction not found")
	AuctionPriceInvalid     = sdkerrors.Register(ModuleName, 1300, "auction start price must larger than 0")
	AuctionDurationInvalid  = sdkerrors.Register(ModuleName, 1400, "auction duration must be at least 10")
	AuctionDurationPassed   = sdkerrors.Register(ModuleName, 1500, "auction duration passed, not accepting new bid")
	AuctionFinalizeTooEarly = sdkerrors.Register(ModuleName, 1600, "please wait until auction duration passed to finalize the result")

	BidNotFound = sdkerrors.Register(ModuleName, 1700, "target bid not found")
	BidPriceLow = sdkerrors.Register(ModuleName, 1800, "bid price is lower / equals to the highest bid / auction start price")

	NotAuctionOwner = sdkerrors.Register(ModuleName, 1900, "you are not the owner of this auction")

	InsufficientBalance = sdkerrors.Register(ModuleName, 2000, "insufficient balance for the bid price")

	InternalError = sdkerrors.Register(ModuleName, 500, "internal error")
)
