#!/bin/bash

ARCH=$(uname -m)

if [ "$ARCH" == "x86_64" ]; then
  File="tgState.zip"
elif [ "$ARCH" == "aarch64" ]; then
  File="tgState_arm64.zip"
else
  echo "Unsupported architecture: $ARCH"
  exit 1
fi

# Download and unzip
wget "https://github.com/csznet/tgState/releases/latest/download/$File" && unzip "$File" && rm "$File"

# Set permissions
chmod +x tgState

echo "successfully."
