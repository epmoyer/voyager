#!/bin/bash
echo "NOTE: This script must be run with sudo"
cp voyager.linux voyager
chmod 755 voyager
chown root:root voyager
mv voyager /usr/local/bin/