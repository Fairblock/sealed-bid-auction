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
