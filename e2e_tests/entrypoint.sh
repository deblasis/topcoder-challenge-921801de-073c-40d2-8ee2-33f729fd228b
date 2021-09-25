#!/usr/bin/env bash

APIGATEWAY="${APIGATEWAY:-http://apigateway:8081}"

while [[ true ]]
do
    status="$(curl --connect-timeout 2 -s -o /dev/null -w ''%{http_code}'' $APIGATEWAY/health)"
    if [[ $status != "200" ]]
    then
        echo ⏳ [status:$status] waiting for apigateway... 
        sleep 5
    else
        echo ✅ apigateway ready!
        break
    fi
done
make dockertest