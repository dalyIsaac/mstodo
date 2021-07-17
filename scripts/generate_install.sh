#!/bin/sh

# This script generates install-$VERSION.sh
VERSION=$1

# Go to the location of this script
cd `dirname $0`

echo "Creating $VERSION install script"

# Replace VERSION=VERSION in install.sh with the given version, and save to install-$VERSION.sh
cp install.sh install-$VERSION.sh
sed -i "s/VERSION=.*/VERSION=$VERSION/g" install-$VERSION.sh
