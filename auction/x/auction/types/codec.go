package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateAuction{}, "auction/CreateAuction", nil)
	cdc.RegisterConcrete(&MsgCreateBid{}, "auction/CreateBid", nil)
	cdc.RegisterConcrete(&MsgFinalizeAuction{}, "auction/FinalizeAuction", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateAuction{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateBid{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgFinalizeAuction{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
