#!/bin/bash
set -e
echo $POSTGRES_DB

cd ./shyft-cli/postgres_setup
psql -U postgres --set=pgdb="$POSTGRES_DB" -f create_shyftdb.psql
