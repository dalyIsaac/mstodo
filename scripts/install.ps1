# Downloads and unzips mstodo from GitHub to the given directory

Param(
    $directory = "~\bin",
    $version = "v1.0.0"
)

# Create the directory if it doesn't exist
If (!(Test-Path $DIRECTORY)) {
    New-Item $DIRECTORY -ItemType Directory -Force | Out-Null
}

# # Go to the directory
cd $DIRECTORY

# Download the zip file
$EXE_URL="https://github.com/dalyIsaac/mstodo/releases/download/$version/mstodo-$version-windows-amd64.zip"
curl -sL $EXE_URL -o mstodo.zip

# Remove mstodo.exe if it exists
if (Test-Path "mstodo.exe") {
    Remove-Item mstodo.exe -Force -ErrorAction Ignore
}

# Unzip the file
Expand-Archive mstodo.zip -DestinationPath .
Remove-Item mstodo.zip

# Create .mstodo/config.yaml in the user's home directory
cd ~
If (!(Test-Path .mstodo)) {
    New-Item .mstodo -ItemType Directory -Force | Out-Null

    cd ~/.mstodo

    $CONFIG_URL="https://github.com/dalyIsaac/mstodo/releases/download/$version/config.yaml"
    curl -sL $CONFIG_URL > config.yaml

    # Tell the user they need to populate config.yaml
    echo "Please populate config.yaml with your configuration"
    echo "For more, see https://github.com/dalyIsaac/mstodo"
}
