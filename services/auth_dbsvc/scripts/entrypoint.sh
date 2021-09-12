#!/usr/bin/env bash

./wait-for-it.sh ${DEBLASIS_DB_ADDRESS:-auth_db:5432} --timeout=60 -- /migrator -dir /scripts/migrations -init && /exe
