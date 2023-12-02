#!/bin/bash

ARCH=$(uname -m)

if [ "$ARCH" == "x86_64" ]; then
  DOWNLOAD_URL="https://github.com/csznet/tgState/releases/latest/download/tgState.zip"
elif [ "$ARCH" == "aarch64" ]; then
  DOWNLOAD_URL="https://github.com/csznet/tgState/releases/latest/download/tgState_arm64.zip"
else
  echo "Unsupported architecture: $ARCH"
  exit 1
fi

# Download and unzip
wget "$DOWNLOAD_URL" && unzip "tgState${ARCH}_latest.zip" && rm "tgState${ARCH}_latest.zip"

# Set permissions
chmod +x tgState

echo "successfully."
