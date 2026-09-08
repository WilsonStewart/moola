#!/bin/bash
set -e

# Uses passed argument as binary name, or defaults to the current directory name
BINARY_NAME="${1:-$(basename "$PWD")}"
DEST_DIR="/usr/local/bin"

echo "Building $BINARY_NAME..."
go build -o "$BINARY_NAME" .

echo "Installing to $DEST_DIR (may prompt for password)..."
sudo mv "$BINARY_NAME" "$DEST_DIR/$BINARY_NAME"

echo "Successfully installed $BINARY_NAME!"