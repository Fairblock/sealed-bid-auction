package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgFinalizeAuction = "finalize_auction"

var _ sdk.Msg = &MsgFinalizeAuction{}

func NewMsgFinalizeAuction(creator string, auctionId uint64) *MsgFinalizeAuction {
	return &MsgFinalizeAuction{
		Creator:   creator,
		AuctionId: auctionId,
	}
}

func (msg *MsgFinalizeAuction) Route() string {
	return RouterKey
}

func (msg *MsgFinalizeAuction) Type() string {
	return TypeMsgFinalizeAuction
}

func (msg *MsgFinalizeAuction) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgFinalizeAuction) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgFinalizeAuction) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
