#!/usr/bin/env bash

./wait-for-it.sh centralcommandsvc:9482 --timeout=60 -- /exe
