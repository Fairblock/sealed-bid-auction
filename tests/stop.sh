#!/bin/bash

BINARY=fairyringd
BINARY2=auctiond

if pgrep -x "$BINARY" >/dev/null; then
    echo "Terminating $BINARY..."
    killall $BINARY
fi

if pgrep -x "$BINARY2" >/dev/null; then
    echo "Terminating $BINARY2..."
    killall $BINARY2
fi

if pgrep -x "hermes" >/dev/null; then
    echo "Terminating Hermes Relayer..."
    killall hermes
fi

if pgrep -x "fairyport" >/dev/null; then
    echo "Terminating fairyport..."
    killall fairyport
fi
