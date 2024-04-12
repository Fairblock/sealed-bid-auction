#!/bin/sh

# set variables for the chain
VALIDATOR_NAME=validator1
CHAIN_ID=auction
KEY_NAME=auction-key
KEY_2_NAME=auction-key-2
KEY_RELAY=auction-relay
TOKEN_AMOUNT="10000000000000000000000000stake"
STAKING_AMOUNT="1000000000stake"
NAMESPACE=$(openssl rand -hex 8)

# query the DA Layer start height, in this case we are querying
# our local devnet at port 26657, the RPC. The RPC endpoint is
# to allow users to interact with Celestia's nodes by querying
# the node's state and broadcasting transactions on the Celestia
# network. The default port is 26657.
DA_BLOCK_HEIGHT=$(curl http://0.0.0.0:26657/block | jq -r '.result.block.header.height')

# echo variables for the chain
echo -e "\n Your DA_BLOCK_HEIGHT is $DA_BLOCK_HEIGHT \n"

# build the auction chain with Rollkit
# ignite chain build

# reset any existing genesis/chain data
rm -rf $HOME/.auction
auctiond tendermint unsafe-reset-all

# initialize the validator with the chain ID you set
auctiond init $VALIDATOR_NAME --chain-id $CHAIN_ID

# add keys for key 1 and key 2 to keyring-backend test
auctiond keys add $KEY_NAME --keyring-backend test
auctiond keys add $KEY_2_NAME --keyring-backend test
echo "milk verify alley price trust come maple will suit hood clay exotic" | auctiond keys add $KEY_RELAY --keyring-backend test  --recover

# add these as genesis accounts
auctiond add-genesis-account $KEY_NAME $TOKEN_AMOUNT --keyring-backend test
auctiond add-genesis-account $KEY_2_NAME $TOKEN_AMOUNT --keyring-backend test
auctiond add-genesis-account $KEY_RELAY $TOKEN_AMOUNT --keyring-backend test

# set the staking amounts in the genesis transaction
auctiond gentx $KEY_NAME $STAKING_AMOUNT --chain-id $CHAIN_ID --keyring-backend test

# collect genesis transactions
auctiond collect-gentxs

# copy centralized sequencer address into genesis.json
# Note: validator and sequencer are used interchangeably here
ADDRESS=$(jq -r '.address' ~/.auction/config/priv_validator_key.json)
PUB_KEY=$(jq -r '.pub_key' ~/.auction/config/priv_validator_key.json)
jq --argjson pubKey "$PUB_KEY" '.consensus["validators"]=[{"address": "'$ADDRESS'", "pub_key": $pubKey, "power": "1000", "name": "Rollkit Sequencer"}]' ~/.auction/config/genesis.json > temp.json && mv temp.json ~/.auction/config/genesis.json

AUTH_TOKEN=$(docker exec $(docker ps -q)  celestia bridge --node.store /home/celestia/bridge/ auth admin)

# create a restart-local.sh file to restart the chain later
[ -f restart-local.sh ] && rm restart-local.sh
echo "DA_BLOCK_HEIGHT=$DA_BLOCK_HEIGHT" >> restart-local.sh
echo "NAMESPACE=$NAMESPACE" >> restart-local.sh
echo "AUTH_TOKEN=$AUTH_TOKEN" >> restart-local.sh

echo "auctiond start --rollkit.aggregator true --rollkit.da_layer celestia --rollkit.da_config='{\"base_url\":\"http://localhost:26658\",\"timeout\":60000000000,\"fee\":600000,\"gas_limit\":6000000,\"auth_token\":\"'\$AUTH_TOKEN'\"}' --rollkit.namespace_id \$NAMESPACE --rollkit.da_start_height \$DA_BLOCK_HEIGHT --rpc.laddr tcp://127.0.0.1:36657 --p2p.laddr \"0.0.0.0:36656\"" >> restart-local.sh

# start the chain
auctiond start --rollkit.aggregator true --rollkit.da_layer celestia --rollkit.da_config='{"base_url":"http://localhost:26658","timeout":60000000000,"fee":600000,"gas_limit":6000000,"auth_token":"'$AUTH_TOKEN'"}' --rollkit.namespace_id $NAMESPACE --rollkit.da_start_height $DA_BLOCK_HEIGHT --rpc.laddr tcp://127.0.0.1:36657 --p2p.laddr "0.0.0.0:36656"
