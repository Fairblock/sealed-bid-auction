#!/bin/bash

./start.sh

./relayer.sh

./fairyport.sh

./keyshare.sh

sleep 15

./auction.sh

./stop.sh