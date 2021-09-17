#!/usr/bin/env bash

./wait-for-it.sh auth_dbsvc:9182 --timeout=60 -- /exe
