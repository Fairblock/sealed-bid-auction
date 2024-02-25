# Seal Bid Auction Demo

This document provide step by steps guide on creating a seal bid auction chain with fairyring pep module.

## Install dependencies

This guide assumes that you are using Ubuntu.

1. Upgrade your operating system:

```bash
sudo apt update && sudo apt upgrade -y
 ```

2. Install essential packages:

```bash
sudo apt install git curl tar wget libssl-dev jq build-essential gcc make
```

3. Download and install Go:

```bash
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install golang-go
```

4. Add `/usr/local/go/bin` & `$HOME/go/bin` directories to your `$PATH`:

```bash
echo "export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin" >> $HOME/.profile
source $HOME/.profile
```

5. Verify Go was installed correctly. Note that the `fairyring` binary requires at least Go `v1.21`:

```bash
go version
```

6. Install Ignite cli v0.27.2

```bash
git clone https://github.com/ignite/cli
cd cli
git fetch --all --tags
git checkout tags/v0.27.2
make install
```

The following binaries are for testing the seal bid auction chain:

7. Install [Hermes relayer](https://github.com/informalsystems/hermes) by following the [official guide](https://hermes.informal.systems/quick-start/pre-requisites.html)

8. Install [fairyring binary](https://github.com/Fairblock/fairyring) by following [this guide](https://docs.fairblock.network/docs/running-a-node/installation)

9. Install [encrypter](https://github.com/Fairblock/encrypter) by following [this guide](https://docs.fairblock.network/docs/advanced/encrypt_tx#install-encrypter)

10. Install [ShareGenerator](https://github.com/Fairblock/ShareGenerator) by following [this guide](https://docs.fairblock.network/docs/advanced/share_generator) 

## Scaffold the seal bid auction chain

### 1. Scaffold the chain with Ignite CLI

```bash
ignite scaffold chain auction
```

### 2. Integrate pep module

1. Import pep module by adding the following lines to the import section in `app/app.go`

```go
pepmodule "github.com/Fairblock/fairyring/x/pep"
pepmodulekeeper "github.com/Fairblock/fairyring/x/pep/keeper"
pepmoduletypes "github.com/Fairblock/fairyring/x/pep/types"
```

It will look something like this:

```go
package app

import (
    pepmodule "github.com/Fairblock/fairyring/x/pep"
    pepmodulekeeper "github.com/Fairblock/fairyring/x/pep/keeper"
    pepmoduletypes "github.com/Fairblock/fairyring/x/pep/types"

 "encoding/json"
 "fmt"
 ...
)
```

2. Add pep modules to `app/app.go`:

- Add the module to `ModuleBasics`

```go
ModuleBasics = module.NewBasicManager(
  // ... 
  pepmodule.AppModuleBasic{},
 )
```

- Update module account permissions

```go
maccPerms = map[string][]string{
    // ...
    pepmoduletypes.ModuleName:  {authtypes.Minter, authtypes.Burner, authtypes.Staking},
 }
```

- Add keepers to app

```go
type App struct {
 // ...
 ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
 ScopedTransferKeeper capabilitykeeper.ScopedKeeper
 ScopedICAHostKeeper  capabilitykeeper.ScopedKeeper
 ScopedPepKeeper      capabilitykeeper.ScopedKeeper

 PepKeeper     pepmodulekeeper.Keeper
 // ...
 }
```

- update kv store keys

```go
 keys := sdk.NewKVStoreKeys(
    // ... 
  auctionmoduletypes.StoreKey, pepmoduletypes.StoreKey,
 )
```

- configure keepers and module

```go
 scopedPepKeeper := app.CapabilityKeeper.ScopeToModule(pepmoduletypes.ModuleName)
 app.PepKeeper = *pepmodulekeeper.NewKeeper(
  appCodec,
  keys[pepmoduletypes.StoreKey],
  keys[pepmoduletypes.MemStoreKey],
  app.GetSubspace(pepmoduletypes.ModuleName),
  app.IBCKeeper.ChannelKeeper,
  &app.IBCKeeper.PortKeeper,
  scopedPepKeeper,
  app.IBCKeeper.ConnectionKeeper,
  app.BankKeeper,
 )
 pepModule := pepmodule.NewAppModule(
  appCodec,
  app.PepKeeper,
  app.AccountKeeper,
  app.BankKeeper,
  app.MsgServiceRouter(),
  encodingConfig.TxConfig,
  app.SimCheck,
 )

 pepIBCModule := pepmodule.NewIBCModule(app.PepKeeper)
```

- Add IBC route

```go
ibcRouter.AddRoute(icahosttypes.SubModuleName, icaHostIBCModule).
  AddRoute(ibctransfertypes.ModuleName, transferIBCModule).
  AddRoute(pepmoduletypes.ModuleName, pepIBCModule)
```

- Add to module manager

```go
app.mm = module.NewManager(
 // ...  
  icaModule,
  auctionModule,
  pepModule,
 // ... 
)
```

- Set begin and end blockers

```go
app.mm.SetOrderBeginBlockers(
  // ... 
  pepmoduletypes.ModuleName,
 )

app.mm.SetOrderEndBlockers(
  // ... 
  pepmoduletypes.ModuleName,
 )
```

- Modify genesis modules

```go
genesisModuleOrder := []string{
  // ...  
  pepmoduletypes.ModuleName,
 }
```

- Scoped keeper  

```go
app.ScopedPepKeeper = scopedPepKeeper
```

- Init params keeper

```go
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
 paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

 // ... 
 paramsKeeper.Subspace(pepmoduletypes.ModuleName)

 return paramsKeeper
}
```

3. Add the following line to the end of `go.mod`

```
replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
```

4. Run `go mod tidy`

### 3. Scaffold the types and messages for the auction chain

1. Create aution type & messages
   
```bash
ignite scaffold list auction name startPrice:coin duration:uint createdAt:uint currentHighestBidId:uint highestBidExists:bool ended:bool --module auction --no-simulation
```

1. Create bid types & message for user to place bid

```bash
ignite scaffold list bid auctionId:uint bidPrice:coin --module auction --no-simulation
```

3. Create finalize-auction message for auction creator to end the auction

```bash
ignite scaffold message finalize-auction auctionId:uint --module auction
```

4. Create `FinalizedAuction` types

```bash
ignite scaffold list finalizedAuction auctionId:uint bidId:uint finalPrice:coin bidder creator --module auction --no-simulation --no-message
```

### 4. Implement logic for all the messages

1. Implement all the errors for all the messages

- Replace line 10 - 12 in `x/auction/types/errors/go` with the code below

```go
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
```

2. Implement logic for create auction message

- `CreateAuction()` in `x/auction/keeper/msg_server_auction.go`

```go
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
```

3. Add `UpdateAuctionHighestBidId()` & `EndAuction()` to `/x/auction/keeper/auction.go`

```go
// UpdateAuctionHighestBidId set the highest bid id of auction
func (k Keeper) UpdateAuctionHighestBidId(ctx sdk.Context, id uint64, bidId uint64) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))

	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)

	b := store.Get(bz)
	if b == nil {
		return types.AuctionNotFound
	}

	var auction types.Auction
	k.cdc.MustUnmarshal(b, &auction)

	auction.HighestBidExists = true
	auction.CurrentHighestBidId = bidId

	appendedValue := k.cdc.MustMarshal(&auction)

	store.Set(bz, appendedValue)
	return nil
}

// EndAuction set the auction status to ended
func (k Keeper) EndAuction(ctx sdk.Context, id uint64) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuctionKey))

	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)

	b := store.Get(bz)
	if b == nil {
		return types.AuctionNotFound
	}

	var auction types.Auction
	k.cdc.MustUnmarshal(b, &auction)
	auction.Ended = true

	appendedValue := k.cdc.MustMarshal(&auction)

	store.Set(bz, appendedValue)
	return nil
}
```

4. Remove `UpdateAuction` & `DeleteAuction` message handler

- Remove `UpdateAuction()` and `DeleteAuction()` in `x/auction/keeper/msg_server_aution.go`

- Remove `"fmt"` import, the import will look something like this:

```go
import (
	"auction/x/auction/types"
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)
```

5. Implement logic for create bid message

- Add `BankKeeper` to auction module keeper in `x/auction/keeper/keeper.go`

```go
type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
		bankKeeper types.BankKeeper
	}
)

func NewKeeper(
    cdc codec.BinaryCodec,
    storeKey,
    memKey storetypes.StoreKey,
    ps paramtypes.Subspace,
    bankKeeper types.BankKeeper,
) *Keeper {
    // set KeyTable if it has not already been set
    if !ps.HasKeyTable() {
        ps = ps.WithKeyTable(types.ParamKeyTable())
    }

    return &Keeper{
        cdc:        cdc,
        storeKey:   storeKey,
        memKey:     memKey,
        paramstore: ps,
        bankKeeper: bankKeeper,
    }
}
```

- Add all the `SendCoins()` function to `BankKeeper` interface under `x/auction/types/expected_keepers.go`

```go
// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}
```

- Add `BankKeeper` to `AuctionKeeper` in `app/app.go`

```go
app.AuctionKeeper = *auctionmodulekeeper.NewKeeper(
    appCodec,
    keys[auctionmoduletypes.StoreKey],
    keys[auctionmoduletypes.MemStoreKey],
    app.GetSubspace(auctionmoduletypes.ModuleName),
    app.BankKeeper,
)
```

- Implement logic for `CreateBid()` in `x/auction/keeper/msg_server_bid.go`

```go
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
```

6. Remove `UpdateBid` & `DeleteBid` message handler

- Remove `UpdateBid()` and `DeleteBid()` in `x/auction/keeper/msg_server_bid.go`

- Remove `"fmt"` import, the import will look something like this:

```go
import (
	"auction/x/auction/types"
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)
```

7. Implement finalize auction message logic

- `FinalizeAuction()` in `x/auction/keeper/msg_server_finalize_auction.go`

```go
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
```

8. Update tx proto

- Remove Update Bid & Auction, Delete Bid & Auction TX Msgs in `proto/auction/auction/tx.proto`

```go
service Msg {
  rpc CreateAuction   (MsgCreateAuction  ) returns (MsgCreateAuctionResponse  );
  rpc CreateBid       (MsgCreateBid      ) returns (MsgCreateBidResponse      );
  rpc FinalizeAuction (MsgFinalizeAuction) returns (MsgFinalizeAuctionResponse);
}
```

- Remove all the related messages

```go
message MsgUpdateAuction {
  string                   creator             = 1;
  uint64                   id                  = 2;
  string                   name                = 3;
  cosmos.base.v1beta1.Coin startPrice          = 4 [(gogoproto.nullable) = false];
  uint64                   duration            = 5;
  uint64                   createdAt           = 6;
  uint64                   currentHighestBidId = 7;
  bool                     highestBidExists    = 8;
  bool                     ended               = 9;
}

message MsgUpdateAuctionResponse {}

message MsgDeleteAuction {
  string creator = 1;
  uint64 id      = 2;
}

message MsgDeleteAuctionResponse {}

message MsgUpdateBid {
string                   creator   = 1;
uint64                   id        = 2;
uint64                   auctionId = 3;
cosmos.base.v1beta1.Coin bidPrice  = 4 [(gogoproto.nullable) = false];
}

message MsgUpdateBidResponse {}

message MsgDeleteBid {
string creator = 1;
uint64 id      = 2;
}

message MsgDeleteBidResponse {}
```

9. Remove Update Bid & Auction, Delete Bid & Auction TX in cli client tx command

- Remove these lines in `GetTxCmd()` in `x/auction/client/cli/tx.go`

```go
	cmd.AddCommand(CmdUpdateAuction())
	cmd.AddCommand(CmdDeleteAuction())
	cmd.AddCommand(CmdUpdateBid())
	cmd.AddCommand(CmdDeleteBid())
```

- Remove `CmdUpdateAuction()` and `CmdDeleteAuction()` and `"strconv"` import in `x/auction/client/cli/tx_auction.go`

- Remove `CmdUpdateBid()` and `CmdDeleteBid()` and `"strconv"` import in `x/auction/client/cli/tx_bid.go`

10. Remove Update Bid & Auction, Delete Bid & Auction in messages type

- Remove all the following lines in `x/auction/types/messages_auction.go`

```go
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
```

- Remove all the following lines in `x/auction/types/messages_bid.go`

```go

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
```

11. Remove Update Bid & Auction, Delete Bid & Auction in the codec

- Remove all update bid & auction, delete bid & auction code in `x/auction/types/codec.go`, The file should look like this:

```go
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
```

12. Update `MsgFinalizeAuctionResponse` to include Bidder address, ID and final price

- Update `MsgFinalizeAuctionResponse{}` in `proto/auction/auction/tx.proto` to:

```protobuf
message MsgFinalizeAuctionResponse {
  uint64 id = 1;
  cosmos.base.v1beta1.Coin finalPrice = 2 [(gogoproto.nullable) = false];
  string bidder = 3;
}
```