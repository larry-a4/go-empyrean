#!/bin/sh
set +e
dropdb 'shyftdb'
sh ./shyft-cli/resetShyftGeth.sh &&                     # Reset geth data - Remove pg and chain data
sh ./shyft-cli/initShyftGeth.sh                         # Init Shyft Geth
