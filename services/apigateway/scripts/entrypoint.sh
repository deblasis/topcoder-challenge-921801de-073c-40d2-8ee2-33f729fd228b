#!/usr/bin/env bash

AUTHSVC_ENDPOINT="${DEBLASIS_AUTHSERVICE_GRPCENDPOINT:-deblasis-v1-AuthService.service.consul:9082}"
CCSVC_ENDPOINT="${DEBLASIS_CENTRALCOMMANDERVICE_GRPCENDPOINT:-deblasis-v1-CentralCommandService.service.consul:9482}"

./wait-for-it.sh consul:8500 --timeout=60 -- \
./wait-for-it.sh $AUTHSVC_ENDPOINT --timeout=120 -- \
./wait-for-it.sh $CCSVC_ENDPOINT --timeout=120 -- \
/exe