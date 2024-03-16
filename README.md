# Seal Bid Auction Demo

This document provide step by steps guide on integrating your Cosmsos SDK chain with `fairyring`'s  `pep` module.

The following binaries are for testing the integration:

7. Install [Hermes relayer](https://github.com/informalsystems/hermes) by following the [official guide](https://hermes.informal.systems/quick-start/pre-requisites.html)

8. Install [fairyring binary](https://github.com/Fairblock/fairyring) by following [this guide](https://docs.fairblock.network/docs/running-a-node/installation)

9. Install [encrypter](https://github.com/Fairblock/encrypter) by following [this guide](https://docs.fairblock.network/docs/advanced/encrypt_tx#install-encrypter)

10. Install [ShareGenerator](https://github.com/Fairblock/ShareGenerator) by following [this guide](https://docs.fairblock.network/docs/advanced/share_generator)

11. Install [Fairyport](https://github.com/Fairblock/fairyport) by following [this guide](https://docs.fairblock.network/docs/advanced/fairyport)

## Integration of `pep` module

1. Import pep module by adding the following lines to the import section in `app/app.go`

```go
pepmodule "github.com/Fairblock/fairyring/x/pep"
pepmodulekeeper "github.com/Fairblock/fairyring/x/pep/keeper"
pepmoduletypes "github.com/Fairblock/fairyring/x/pep/types"
```

2. Add `pep` modules to `app/app.go`:

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

That's it! Now you should have the `pep` module integrated with your chain.

## Testing the Integration

The scripts in the `tests` directory allows you to put all of the moving components together and test the functionality.
