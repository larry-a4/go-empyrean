#!/bin/sh
  /bin/geth --config config.toml --gcmode archive --ws --wsaddr="0.0.0.0" --wsorigins "*" --nat=any --minerthreads 4 --targetgaslimit 80000000 

