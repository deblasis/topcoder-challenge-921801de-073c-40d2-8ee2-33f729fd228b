#!/usr/bin/env bash

CCDB="${DEBLASIS_CENTRALCOMMANDDB:-centralcommand_db:5432}"

./wait-for-it.sh consul:8500 --timeout=60 -- \
./wait-for-it.sh $DEBLASIS_CENTRALCOMMANDDB --timeout=120 -- \
/exe