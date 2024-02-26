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

WALLET_4=$($BINARY keys show wallet4 -a --keyring-backend test --home $CHAIN_DIR/$CHAINID_2)
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
  echo "ERROR: Auction Chain Pep module query Pep Nonce error. Expected Pep Nonce to be 1, got '$VALIDATOR_PEP_NONCE'"
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


echo "Creating auction with starting price 1000token for 5 blocks"
RESULT=$($BINARY tx auction create-auction "Testing Auction 0" 1000token 5 --from $VALIDATOR_2 --home $CHAIN_DIR/$CHAINID_2 --chain-id $CHAINID_2 --node $CHAIN2_NODE --keyring-backend test -o json -y)
check_tx_code $RESULT
RESULT=$(wait_for_tx $RESULT)
CODE=$(echo $RESULT | jq -r '.code')
if [ "$CODE" != "0" ]; then
  echo "ERROR: Create auction tx failed: $RESULT"
  exit 1
fi

RESULT=$($BINARY query bank balances $VALIDATOR_2 --node $CHAIN2_NODE -o json)
TARGET_BAL_DENOM=$(echo "$RESULT" | jq -r '.balances[1].denom')
TARGET_BAL=$(echo "$RESULT" | jq -r '.balances[1].amount')
echo "Auction creator balance after creating auction: $TARGET_BAL $TARGET_BAL_DENOM"


echo "Signing place bid tx with pep nonce: '$VALIDATOR_PEP_NONCE'"
echo "Create 1001token bid for auction id 0"
$BINARY tx auction create-bid 0 1001token --from $WALLET_4 --home $CHAIN_DIR/$CHAINID_2 --gas-prices 1token --gas 300000 --chain-id $CHAINID_2 --node $CHAIN2_NODE --keyring-backend test --generate-only -o json -y > unsigned.json
SIGNED_DATA=$($BINARY tx sign unsigned.json --from $WALLET_4 --offline --account-number 2 --sequence $VALIDATOR_PEP_NONCE --gas-prices 1token --home $CHAIN_DIR/$CHAINID_2 --chain-id $CHAINID_2 --node $CHAIN2_NODE  --keyring-backend test -y)


echo "Query latest height from pep module on chain $CHAINID_2"
RESULT=$($BINARY q pep latest-height --node $CHAIN2_NODE -o json | jq -r '.height')
echo "$CHAINID_2 Pep module Latest height: $RESULT"

NEW_TARGET_HEIGHT=$(($RESULT+2))

echo "Encrypting signed tx with Pub key: '$PUB_KEY'"
CIPHER=$($ENCRYPTER $NEW_TARGET_HEIGHT $PUB_KEY $SIGNED_DATA)

rm -r unsigned.json &> /dev/null

echo "Submit encrypted create bid tx to pep module on chain $CHAINID_2"
RESULT=$($BINARY tx pep submit-encrypted-tx $CIPHER $NEW_TARGET_HEIGHT --from $WALLET_4 --gas-prices 1token --gas 300000 --home $CHAIN_DIR/$CHAINID_2 --chain-id $CHAINID_2 --node $CHAIN2_NODE --broadcast-mode sync --keyring-backend test -o json -y)
check_tx_code $RESULT
RESULT=$(wait_for_tx $RESULT)
EVENT_TYPE=$(echo "$RESULT" | jq -r '.logs[0].events[5].type')
TARGET_HEIGHT=$(echo "$RESULT" | jq -r '.logs[0].events[5].attributes[1].value')
if [ "$EVENT_TYPE" != "new-encrypted-tx-submitted" ] && [ "$TARGET_HEIGHT" != "$NEW_TARGET_HEIGHT" ]; then
  CURRENT_BLOCK=$($BINARY query block --home $CHAIN_DIR/$CHAINID_2 --node $CHAIN2_NODE | jq -r '.block.header.height')
  echo "ERROR: Auction Chain Pep module submit encrypted tx error. Expected tx to submitted without error with target height '$NEW_TARGET_HEIGHT', got '$TARGET_HEIGHT' and '$EVENT_TYPE' | '$CURRENT_BLOCK'"
  echo "ERROR MESSAGE: $(echo "$RESULT" | jq -r '.raw_log')"
  echo "ERROR MESSAGE: $(echo "$RESULT" | jq '.')"
  exit 1
fi

echo "Encrypted TXs: $($BINARY query pep list-encrypted-tx --node $CHAIN2_NODE -o json | jq '.encryptedTxArray')"

echo "Bids after submitting the encrypted bid: $($BINARY q auction list-bid --node $CHAIN2_NODE -o json | jq '.Bid')"
RESULT=$($BINARY query bank balances $WALLET_4 --node $CHAIN2_NODE -o json)
TARGET_BAL_DENOM=$(echo "$RESULT" | jq -r '.balances[1].denom')
TARGET_BAL=$(echo "$RESULT" | jq -r '.balances[1].amount')
echo "Bidder balance: $TARGET_BAL $TARGET_BAL_DENOM"

echo "Wait for the encrypted vote to be executed..."

sleep $(($BLOCK_TIME * 3))

RESULT=$($BINARY q auction list-bid --node $CHAIN2_NODE -o json | jq -r '.Bid[0].bidPrice.amount')
if [ "$RESULT" != "1001" ]; then
  $BINARY q auction list-bid --node $CHAIN2_NODE -o json | jq
  echo "ERROR: bid not found"
  exit 1
fi

echo "Found bid with price: $RESULT"
RESULT=$($BINARY query bank balances $WALLET_4 --node $CHAIN2_NODE -o json)
TARGET_BAL_DENOM=$(echo "$RESULT" | jq -r '.balances[1].denom')
TARGET_BAL=$(echo "$RESULT" | jq -r '.balances[1].amount')
echo "Bidder balance: $TARGET_BAL $TARGET_BAL_DENOM"

echo "Encrypted TXs After Execution: $($BINARY query pep list-encrypted-tx --node $CHAIN2_NODE -o json | jq '.encryptedTxArray')"

echo "Wait for the auction to end..."

sleep $(($BLOCK_TIME * 3))

echo "Finalizing Auction:"
RESULT=$($BINARY tx auction finalize-auction 0 --from $VALIDATOR_2 --home $CHAIN_DIR/$CHAINID_2 --chain-id $CHAINID_2 --node $CHAIN2_NODE --keyring-backend test -o json -y)
check_tx_code $RESULT
RESULT=$(wait_for_tx $RESULT)

echo "Auction Finalized"

RESULT=$($BINARY query bank balances $VALIDATOR_2 --node $CHAIN2_NODE -o json)
TARGET_BAL_DENOM=$(echo "$RESULT" | jq -r '.balances[1].denom')
TARGET_BAL=$(echo "$RESULT" | jq -r '.balances[1].amount')
echo "Auction creator balance now: $TARGET_BAL $TARGET_BAL_DENOM"

echo ""
echo "##############################################################"
echo "# Successfully Tested Submit Encrypted Vote to Auction Chain #"
echo "##############################################################"
echo ""