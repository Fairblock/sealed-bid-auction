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

func NewMsgCreateAuction(creator string, name string, startPrice sdk.Coin, duration uint64) *MsgCreateAuction {
	return &MsgCreateAuction{
		Creator:    creator,
		Name:       name,
		StartPrice: startPrice,
		Duration:   duration,
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
