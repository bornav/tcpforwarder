#!/usr/bin/env bash
BINARY_PATH="/usr/bin/tcpforwarder"
LISTEN_HOST=${LISTEN_HOST:-}
LISTEN_PORT=${LISTEN_PORT:-}
REMOTE_HOST=${REMOTE_HOST:-}
REMOTE_PORT=${REMOTE_PORT:-}
DIAL_TIMEOUT=${DIAL_TIMEOUT:-}

args=()

args+=(-lHost "${LISTEN_HOST:-localhost}")
args+=(-lPort "${LISTEN_PORT:-80}")
args+=(-rHost "${REMOTE_HOST:-example.com}")
args+=(-rPort "${REMOTE_PORT:-8080}")

STRING_TO_FORWARD="-lHost ${LISTEN_HOST:-localhost}"
STRING_TO_FORWARD+=" -lPort ${LISTEN_PORT:-80}"
STRING_TO_FORWARD+=" -rHost ${REMOTE_HOST:-1.1.1.1}"
STRING_TO_FORWARD+=" -rPort ${REMOTE_PORT:-8080}"

echo ""
echo "Tcp forwarder command: $BINARY_PATH $STRING_TO_FORWARD"
echo ""

if ! $BINARY_PATH "${args[@]}"; then
    echo "Error running tcp forwarder command:"
    echo "$?"
fi