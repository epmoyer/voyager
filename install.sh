#/bin/bash

RED=$'\033[31m'
YELLOW=$'\033[33m'
GREEN=$'\033[32m'
ENDCOLOR=$'\033[0m'

echo "Installing to /usr/local/bin.  You may be prompted for sudo permissions..."
sudo cp voyager /usr/local/bin
sudo chmod 755 /usr/local/bin/voyager
echo "${GREEN}   Copied.${ENDCOLOR}"