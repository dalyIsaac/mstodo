# Downloads and unzips mstodo from GitHub to the given directory

$DIRECTORY=$args[0]
$VERSION=VERSION
curl -L https://github.com/dalyIsaac/mstodo/releases/download/$VERSION/mstodo-$VERSION-windows-amd64.zip -o mstodo.zip
unzip mstodo.zip
mv mstodo-$VERSION-windows-amd64/mstodo.exe ~/$DIRECTORY/mstodo.exe

# Create .mstodo/config.yaml in the user's home directory
cd ~
mkdir .mstodo
cd .mstodo
curl -L https://raw.githubusercontent.com/dalyIsaac/mstodo/main/config.yaml > config.yaml

# Tell the user they need to populate config.yaml
echo "Please populate config.yaml with your configuration"
echo "For more, see https://github.com/dalyIsaac/mstodo"
