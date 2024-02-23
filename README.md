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

TODO

3. Add the following line to the end of `go.mod`

```
replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
```

4. Run `go mod tidy`