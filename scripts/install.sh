#!/bin/sh

# Downloads and installs mstodo from GitHub to the given directory, or /usr/local/bin (by default)
DIRECTORY=${1:-/usr/local/bin}
VERSION=VERSION
curl -L https://github.com/dalyIsaac/mstodo/releases/download/v1.0.0/mstodo-v1.0.0-linux-amd64.tar.gz | tar -xz

# Create .mstodo/config.yaml in the user's home directory
cd ~
mkdir .mstodo
cd .mstodo
curl -L https://raw.githubusercontent.com/dalyIsaac/mstodo/main/config.yaml > config.yaml

# Tell the user they need to populate config.yaml
echo "Please populate config.yaml with your configuration."
echo "For more, see https://github.com/dalyIsaac/mstodo"
