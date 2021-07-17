# This script generates install-$VERSION.ps1
$VERSION=$args[0]

# Go to the location of this script
cd $PSScriptRoot

echo "Creating $VERSION install script"

# Replace $VERSION=VERSION in install.sh with the given version, and save to install-$VERSION.sh
(Get-Content -path install.ps1) -replace "VERSION=VERSION", "VERSION=$VERSION" | Set-Content -Path "install-$VERSION.ps1"
