# Downloads and unzips mstodo from GitHub to the given directory

$DIRECTORY=$args[0]

# Create the directory if it doesn't exist
If (!(Test-Path $DIRECTORY)) {
    New-Item $DIRECTORY -ItemType Directory -Force
}

# Go to the directory
cd $DIRECTORY

# Download the zip file
$EXE_URL=https://github.com/dalyIsaac/mstodo/releases/download/$env:VERSION/mstodo-$env:VERSION-windows-amd64.zip
curl -L $EXE_URL -o mstodo.zip
unzip mstodo.zip
mv mstodo-$env:VERSION-windows-amd64/mstodo.exe mstodo.exe
Remove-Item mstodo.zip

# Create .mstodo/config.yaml in the user's home directory
cd ~
mkdir .mstodo
cd .mstodo
$CONFIG_URL=https://github.com/dalyIsaac/mstodo/releases/download/$env:VERSION/config.yaml
curl -L $CONFIG_URL > config.yaml

# Tell the user they need to populate config.yaml
echo "Please populate config.yaml with your configuration"
echo "For more, see https://github.com/dalyIsaac/mstodo"

# Remove VERSION variable
Remove-Item Env:\VERSION
