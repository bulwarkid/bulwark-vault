#!/bin/bash

if [[ $1 == "--dry-run=0" ]];
then
    ssh ubuntu@bulwark.id '
        cd /opt/bulwarkid/bulwark-vault
        git pull origin master
        make build
        sudo systemctl restart vault.service
    '
else
    echo "Use --dry-run=0 to actually run this command."
fi