package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateBid = "create_bid"
	TypeMsgUpdateBid = "update_bid"
	TypeMsgDeleteBid = "delete_bid"
)

var _ sdk.Msg = &MsgCreateBid{}

func NewMsgCreateBid(creator string, auctionId uint64, bidPrice sdk.Coin) *MsgCreateBid {
	return &MsgCreateBid{
		Creator:   creator,
		AuctionId: auctionId,
		BidPrice:  bidPrice,
	}
}

func (msg *MsgCreateBid) Route() string {
	return RouterKey
}

func (msg *MsgCreateBid) Type() string {
	return TypeMsgCreateBid
}

func (msg *MsgCreateBid) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateBid) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateBid) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateBid{}

func NewMsgUpdateBid(creator string, id uint64, auctionId uint64, bidPrice sdk.Coin) *MsgUpdateBid {
	return &MsgUpdateBid{
		Id:        id,
		Creator:   creator,
		AuctionId: auctionId,
		BidPrice:  bidPrice,
	}
}

func (msg *MsgUpdateBid) Route() string {
	return RouterKey
}

func (msg *MsgUpdateBid) Type() string {
	return TypeMsgUpdateBid
}

func (msg *MsgUpdateBid) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateBid) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateBid) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteBid{}

func NewMsgDeleteBid(creator string, id uint64) *MsgDeleteBid {
	return &MsgDeleteBid{
		Id:      id,
		Creator: creator,
	}
}
func (msg *MsgDeleteBid) Route() string {
	return RouterKey
}

func (msg *MsgDeleteBid) Type() string {
	return TypeMsgDeleteBid
}

func (msg *MsgDeleteBid) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteBid) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteBid) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
