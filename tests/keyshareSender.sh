#!/bin/bash

BINARY=$1
HOME=$2
NODE=$3
FROM=$4
CHAINID=$5
GENERATOR=$6
GENERATED_SHARE=$7

BLOCK_TIME=5

check_tx_code () {
  local TX_CODE=$(echo "$1" | jq -r '.code')
  if [ "$TX_CODE" != 0 ]; then
    echo "ERROR: Tx failed with code: $TX_CODE"
    exit 1
  fi
}

wait_for_tx () {
  sleep $BLOCK_TIME
  local TXHASH=$(echo "$1" | jq -r '.txhash')
  RESULT=$($BINARY q tx --type=hash $TXHASH --home $HOME --chain-id $CHAINID --node $NODE -o json)
  echo "$RESULT"
}


wait_for_tx_2 () {
  sleep $BLOCK_TIME
  local TXHASH=$(echo "$1" | jq -r '.txhash')
  RESULT=$($BINARY q tx --type=hash $TXHASH --home /Users/beast/Work-FairBlock/seal-bid-auction-demo/tests/data/auction_test_1/ --chain-id auction_test_1 --node http://localhost:36657 -o json)
  echo "$RESULT"
}


while true
do
  CURRENT_BLOCK=$($BINARY query block --home $HOME --node $NODE | jq -r '.block.header.height')
  TARGET_HEIGHT=$((CURRENT_BLOCK+1))
  EXTRACTED_RESULT=$($GENERATOR derive $GENERATED_SHARE 1 $TARGET_HEIGHT)
  EXTRACTED_SHARE=$(echo "$EXTRACTED_RESULT" | jq -r '.KeyShare')
  RESULT=$($BINARY tx keyshare send-keyshare $EXTRACTED_SHARE 1 $TARGET_HEIGHT --from $FROM --gas-prices 1ufairy --home $HOME --chain-id $CHAINID --node $NODE --broadcast-mode sync --keyring-backend test -o json -y)
  check_tx_code $RESULT
  RESULT=$(wait_for_tx $RESULT)
  RESULT_EVENT=$(echo "$RESULT" | jq -r '.logs[0].events[2].type')
  if [ "$RESULT_EVENT" != "keyshare-aggregated" ]; then
    echo "ERROR: KeyShare module submit invalid key share from registered validator error. Expected the key to be aggregated, got '$RESULT_EVENT'"
    echo "ERROR MESSAGE: $(echo "$RESULT" | jq -r '.raw_log')"
  fi
  echo "Submitted KeyShare for height: $TARGET_HEIGHT"
  echo $EXTRACTED_SHARE
  PUBKEY=$($BINARY q keyshare show-active-pub-key --home $HOME --node $NODE -o json | jq -r '.activePubKey.publicKey')
  echo "PUBKEY: $PUBKEY"
  RESULT=$(auctiond tx pep create-aggregated-key-share $TARGET_HEIGHT $EXTRACTED_SHARE $PUBKEY --from rly2 --home /Users/beast/Work-FairBlock/seal-bid-auction-demo/tests/data/auction_test_1/ --chain-id auction_test_1 --node http://localhost:36657 --broadcast-mode sync --keyring-backend test -o json -y)
  check_tx_code $RESULT
  RESULT=$(wait_for_tx_2 $RESULT)
  echo $RESULT
done