#!/bin/bash


echo ""
echo "###########################################################"
echo "#       Test Submit Encrypted Vote to Auction Chain       #"
echo "###########################################################"
echo ""


ENCRYPTER=encrypter
BINARY1=fairyringd
BINARY=auctiond
CHAIN_DIR=$(pwd)/data
CHAINID_1=fairyring_test_1
CHAIN1_NODE=tcp://localhost:16657
CHAINID_2=auction_test_1
CHAIN2_NODE=tcp://localhost:26657
BLOCK_TIME=5

WALLET_2=$($BINARY keys show wallet2 -a --keyring-backend test --home $CHAIN_DIR/$CHAINID_2)
VALIDATOR_2=$($BINARY keys show val2 -a --keyring-backend test --home $CHAIN_DIR/$CHAINID_2)

check_tx_code () {
  local TX_CODE=$(echo "$1" | jq -r '.code')
  if [ "$TX_CODE" != "0" ]; then
    echo "ERROR: Tx failed with code: $TX_CODE"
    exit 1
  fi
}

wait_for_tx () {
  sleep $BLOCK_TIME
  local TXHASH=$(echo "$1" | jq -r '.txhash')
  RESULT=$($BINARY q tx --type=hash $TXHASH --home $CHAIN_DIR/$CHAINID_2 --chain-id $CHAINID_2 --node $CHAIN2_NODE -o json)
  echo "$RESULT"
}

echo "Query new account pep nonce from pep module on chain $CHAINID_2"
RESULT=$($BINARY query pep show-pep-nonce $VALIDATOR_2 --home $CHAIN_DIR/$CHAINID_2 --chain-id $CHAINID_2 --node $CHAIN2_NODE -o json)
VALIDATOR_PEP_NONCE=$(echo "$RESULT" | jq -r '.pepNonce.nonce')
if [ "$VALIDATOR_PEP_NONCE" != "1" ]; then
  echo "ERROR: Pep module query Pep Nonce error. Expected Pep Nonce to be 1, got '$VALIDATOR_PEP_NONCE'"
  echo "ERROR MESSAGE: $(echo "$RESULT" | jq -r '.raw_log')"
  exit 1
fi

echo "Query master public key from key share module for submitting to pep module on chain $CHAINID_1"
PUB_KEY=$($BINARY1 query keyshare show-active-pub-key --node $CHAIN1_NODE -o json | jq -r '.activePubKey.publicKey')
if [ "$PUB_KEY" == "" ]; then
  echo "ERROR: Query master public key from key share module error, expecting an active public key, got '$PUB_KEY'"
  exit 1
fi

echo "Query master public key expiry height from key share module for submitting to pep module on chain $CHAINID_1"
PUB_KEY_EXPIRY=$($BINARY1 query keyshare show-active-pub-key --node $CHAIN1_NODE -o json | jq -r '.activePubKey.expiry')
if [ "$PUB_KEY_EXPIRY" == "" ]; then
  echo "ERROR: Query master public key expiry height from key share module error, expecting an active public key, got '$PUB_KEY'"
  exit 1
fi

echo "Pub Key expires at: $PUB_KEY_EXPIRY"


echo "Creating auction with starting price 1000token for 10 blocks"
RESULT=$($BINARY tx auction create-auction "Testing Auction 0" 1000token 10 --from $VALIDATOR_2 --home $CHAIN_DIR/$CHAINID_2 --chain-id $CHAINID_2 --node $CHAIN2_NODE --keyring-backend test -o json -y)
check_tx_code $RESULT
RESULT=$(wait_for_tx $RESULT)
CODE=$(echo $RESULT | jq -r '.code')
if [ "$CODE" != "0" ]; then
  echo "ERROR: Create auction tx failed: $RESULT"
  exit 1
fi

echo "Signing place bid tx with pep nonce: '$VALIDATOR_PEP_NONCE'"
echo "Create 1001token bid for auction id 0"
$BINARY tx auction create-bid 0 1001token --from $VALIDATOR_2 --home $CHAIN_DIR/$CHAINID_2 --gas-prices 1token --gas 300000 --chain-id $CHAINID_2 --node $CHAIN2_NODE --keyring-backend test --generate-only -o json -y > unsigned.json
SIGNED_DATA=$($BINARY tx sign unsigned.json --from $VALIDATOR_2 --offline --account-number 0 --sequence $VALIDATOR_PEP_NONCE --gas-prices 1token --home $CHAIN_DIR/$CHAINID_2 --chain-id $CHAINID_2 --node $CHAIN2_NODE  --keyring-backend test -y)

echo "Query aggregated key share from key share module for submitting to pep module on chain $CHAINID_1"
CURRENT_BLOCK=$($BINARY1 query block --home $CHAIN_DIR/$CHAINID_1 --node $CHAIN1_NODE | jq -r '.block.header.height')
RESULT=$($BINARY1 query keyshare list-aggregated-key-share --node $CHAIN1_NODE -o json)
AGG_KEY_HEIGHT=$(echo "$RESULT" | jq -r '.aggregatedKeyShare | last | .height')
AGG_KEY=$(echo "$RESULT" | jq -r '.aggregatedKeyShare | last | .data')
if [ "$CURRENT_BLOCK" -gt "$AGG_KEY_HEIGHT" ]; then
  echo "ERROR: Height of the aggregated key from key share module '$AGG_KEY_HEIGHT' is less than current block height '$CURRENT_BLOCK'"
  exit 1
fi

CURRENT_BLOCK=$($BINARY query block --home $CHAIN_DIR/$CHAINID_2 --node $CHAIN2_NODE | jq -r '.block.header.height')
echo "Chain 2 Current Block: $CURRENT_BLOCK"
echo "Submit valid aggregated key to pep module on chain $CHAINID_2 from address: $VALIDATOR_2"
RESULT=$($BINARY tx pep create-aggregated-key-share $AGG_KEY_HEIGHT $AGG_KEY --from $VALIDATOR_2 --home $CHAIN_DIR/$CHAINID_2 --chain-id $CHAINID_2 --node $CHAIN2_NODE --broadcast-mode sync --keyring-backend test -o json -y)
check_tx_code $RESULT
RESULT=$(wait_for_tx $RESULT)
ACTION=$(echo "$RESULT" | jq -r '.logs[0].events[0].attributes[0].value')
if [ "$ACTION" != "/fairyring.pep.MsgCreateAggregatedKeyShare" ]; then
  echo "ERROR: Pep module submit aggregated key error. Expected tx action to be MsgCreateAggregatedKeyShare,  got '$ACTION'"
  echo "ERROR MESSAGE: $(echo "$RESULT" | jq -r '.raw_log')"
  exit 1
fi

sleep $BLOCK_TIME

echo "Query aggregated key share from key share module for submitting to pep module on chain $CHAINID_1"
CURRENT_BLOCK=$($BINARY1 query block --home $CHAIN_DIR/$CHAINID_1 --node $CHAIN1_NODE | jq -r '.block.header.height')
RESULT=$($BINARY1 query keyshare list-aggregated-key-share --node $CHAIN1_NODE -o json)
AGG_KEY_HEIGHT=$(echo "$RESULT" | jq -r '.aggregatedKeyShare | last | .height')
AGG_KEY=$(echo "$RESULT" | jq -r '.aggregatedKeyShare | last | .data')
if [ "$CURRENT_BLOCK" -gt "$AGG_KEY_HEIGHT" ]; then
  echo "ERROR: Height of the aggregated key from key share module '$AGG_KEY_HEIGHT' is less than current block height '$CURRENT_BLOCK'"
  exit 1
fi

echo "Encrypting signed tx with Pub key: '$PUB_KEY'"
CIPHER=$($ENCRYPTER $AGG_KEY_HEIGHT $PUB_KEY $SIGNED_DATA)

rm -r unsigned.json &> /dev/null

echo "Submit encrypted create bid tx to pep module on chain $CHAINID_2"
RESULT=$($BINARY tx pep submit-encrypted-tx $CIPHER $AGG_KEY_HEIGHT --from $VALIDATOR_2 --gas-prices 1token --gas 300000 --home $CHAIN_DIR/$CHAINID_2 --chain-id $CHAINID_2 --node $CHAIN2_NODE --broadcast-mode sync --keyring-backend test -o json -y)
check_tx_code $RESULT
RESULT=$(wait_for_tx $RESULT)
EVENT_TYPE=$(echo "$RESULT" | jq -r '.logs[0].events[5].type')
TARGET_HEIGHT=$(echo "$RESULT" | jq -r '.logs[0].events[5].attributes[1].value')
if [ "$EVENT_TYPE" != "new-encrypted-tx-submitted" ] && [ "$TARGET_HEIGHT" != "$AGG_KEY_HEIGHT" ]; then
  echo "ERROR: Pep module submit encrypted tx error. Expected tx to submitted without error with target height '$AGG_KEY_HEIGHT', got '$TARGET_HEIGHT' and '$EVENT_TYPE' | '$CURRENT_BLOCK'"
  echo "ERROR MESSAGE: $(echo "$RESULT" | jq -r '.raw_log')"
  echo "ERROR MESSAGE: $(echo "$RESULT" | jq '.')"
  exit 1
fi


CURRENT_BLOCK=$($BINARY query block --home $CHAIN_DIR/$CHAINID_2 --node $CHAIN2_NODE | jq -r '.block.header.height')
echo "Chain 2 Current Block: $CURRENT_BLOCK"
echo "Submit valid aggregated key to pep module on chain $CHAINID_2 from address: $VALIDATOR_2"
RESULT=$($BINARY tx pep create-aggregated-key-share $AGG_KEY_HEIGHT $AGG_KEY --from $VALIDATOR_2 --home $CHAIN_DIR/$CHAINID_2 --chain-id $CHAINID_2 --node $CHAIN2_NODE --broadcast-mode sync --keyring-backend test -o json -y)
check_tx_code $RESULT
RESULT=$(wait_for_tx $RESULT)
ACTION=$(echo "$RESULT" | jq -r '.logs[0].events[0].attributes[0].value')
if [ "$ACTION" != "/fairyring.pep.MsgCreateAggregatedKeyShare" ]; then
  echo "ERROR: Pep module submit aggregated key error. Expected tx action to be MsgCreateAggregatedKeyShare,  got '$ACTION'"
  echo "ERROR MESSAGE: $(echo "$RESULT" | jq -r '.raw_log')"
  exit 1
fi


sleep $BLOCK_TIME

echo "Query latest height from pep module on chain $CHAINID_2"
RESULT=$($BINARY q pep latest-height --node $CHAIN2_NODE -o json | jq -r '.height')
if [ "$RESULT" != "$AGG_KEY_HEIGHT" ]; then
  echo "ERROR: Pep module query latest height error, Expected latest height to be same as aggregated key share height: '$AGG_KEY_HEIGHT', got '$RESULT'"
  exit 1
fi

$BINARY q auction list-bid

echo ""
echo "##############################################################"
echo "# Successfully Tested Submit Encrypted Vote to Auction Chain #"
echo "##############################################################"
echo ""