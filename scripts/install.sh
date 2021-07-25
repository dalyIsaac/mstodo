#!/bin/sh

# Downloads and installs mstodo from GitHub to the given directory, or /usr/local/bin (by default)
DIRECTORY=${1:-/usr/local/bin}

cd $DIRECTORY

EXE_URL=https://github.com/dalyIsaac/mstodo/releases/download/$VERSION/mstodo-$VERSION-linux-amd64.tar.gz
curl -L $EXE_URL | tar -xz

# Create .mstodo/config.yaml in the user's home directory
cd ~
mkdir .mstodo
cd .mstodo

CONFIG_URL=https://github.com/dalyIsaac/mstodo/releases/download/$VERSION/config.yaml
curl -L $CONFIG_URL > config.yaml

# Tell the user they need to populate config.yaml
echo "Please populate config.yaml with your configuration."
echo "For more, see https://github.com/dalyIsaac/mstodo"
