#!/bin/sh
if [[ -z "${DBENV}" ]]; then
  ./build/bin/geth --identity "ShyftTestnetNode" --keystore ./ --datadir "./shyftData" init ./ShyftNetwork.json
else
  if [ -d /go/src/ShyftNetwork/go-empyrean/shyftData/geth/chaindata ]; then
    echo "Skipping Genesis Initialization as already completed"
    :
  else
    echo "Initializing Custom Genesis Block"
    /bin/geth --identity "ShyftTestnetNode" --keystore ./ --datadir "./shyftData" init ShyftNetwork.json
  fi
fi
