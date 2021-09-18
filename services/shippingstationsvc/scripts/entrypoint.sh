#!/usr/bin/env bash

./wait-for-it.sh centralcommand_dbsvc:9382 --timeout=60 -- /exe
