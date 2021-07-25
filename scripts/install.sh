#!/bin/sh

# Downloads and installs mstodo from GitHub to the given directory, or /usr/local/bin (by default)
DIRECTORY=${1:-/usr/local/bin}

# Download from URL
EXE_URL=https://github.com/dalyIsaac/mstodo/releases/download/$VERSION/mstodo-$VERSION-linux-amd64.tar.gz
curl --silent --location --output mstodo.tar.gz $EXE_URL
sudo tar -xzf mstodo.tar.gz -C /usr/local/bin
rm mstodo.tar.gz

# Create .mstodo/config.yaml in the user's home directory, if it doesn't exist
cd ~
if [ ! -f .mstodo/config.yaml ]; then
    mkdir -p .mstodo
    cd .mstodo

    CONFIG_URL=https://github.com/dalyIsaac/mstodo/releases/download/$VERSION/config.yaml
    curl -s -L $CONFIG_URL > config.yaml

    # Tell the user they need to populate config.yaml
    echo "Please populate config.yaml with your configuration."
    echo "For more, see https://github.com/dalyIsaac/mstodo"
fi
