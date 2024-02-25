package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateAuction = "create_auction"
	TypeMsgUpdateAuction = "update_auction"
	TypeMsgDeleteAuction = "delete_auction"
)

var _ sdk.Msg = &MsgCreateAuction{}

func NewMsgCreateAuction(creator string, name string, startPrice sdk.Coin, duration uint64, createdAt uint64, currentHighestBidId uint64, highestBidExists bool, ended bool) *MsgCreateAuction {
	return &MsgCreateAuction{
		Creator:             creator,
		Name:                name,
		StartPrice:          startPrice,
		Duration:            duration,
		CreatedAt:           createdAt,
		CurrentHighestBidId: currentHighestBidId,
		HighestBidExists:    highestBidExists,
		Ended:               ended,
	}
}

func (msg *MsgCreateAuction) Route() string {
	return RouterKey
}

func (msg *MsgCreateAuction) Type() string {
	return TypeMsgCreateAuction
}

func (msg *MsgCreateAuction) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateAuction) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateAuction) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateAuction{}

func NewMsgUpdateAuction(creator string, id uint64, name string, startPrice sdk.Coin, duration uint64, createdAt uint64, currentHighestBidId uint64, highestBidExists bool, ended bool) *MsgUpdateAuction {
	return &MsgUpdateAuction{
		Id:                  id,
		Creator:             creator,
		Name:                name,
		StartPrice:          startPrice,
		Duration:            duration,
		CreatedAt:           createdAt,
		CurrentHighestBidId: currentHighestBidId,
		HighestBidExists:    highestBidExists,
		Ended:               ended,
	}
}

func (msg *MsgUpdateAuction) Route() string {
	return RouterKey
}

func (msg *MsgUpdateAuction) Type() string {
	return TypeMsgUpdateAuction
}

func (msg *MsgUpdateAuction) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateAuction) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateAuction) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteAuction{}

func NewMsgDeleteAuction(creator string, id uint64) *MsgDeleteAuction {
	return &MsgDeleteAuction{
		Id:      id,
		Creator: creator,
	}
}
func (msg *MsgDeleteAuction) Route() string {
	return RouterKey
}

func (msg *MsgDeleteAuction) Type() string {
	return TypeMsgDeleteAuction
}

func (msg *MsgDeleteAuction) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteAuction) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteAuction) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
