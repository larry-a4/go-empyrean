#!/bin/sh
if [ -z "${DBENV}" ]; then
  ./build/bin/geth --config config.toml --nodiscover --ws --wsaddr="0.0.0.0" --wsorigins "*" --nat=any --mine --minerthreads 4 --targetgaslimit 80000000
else
  /bin/geth --config config.toml --nodiscover --gcmode archive --ws --wsaddr="0.0.0.0" --wsorigins "*" --nat=any --mine --minerthreads 4 --targetgaslimit 80000000
fi