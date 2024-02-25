#!/bin/bash

CHAIN_DIR=$(pwd)/data

fairyport start --config fairyport_config.yml > $CHAIN_DIR/fairyport.log 2>&1 &