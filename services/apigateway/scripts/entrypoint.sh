#!/usr/bin/env bash

AUTHSVC_ENDPOINT="${DEBLASIS_APIGATEWAY_AUTHSERVICE_GRPCENDPOINT:-deblasis-v1-AuthService.service.consul:9082}"
CCSVC_ENDPOINT="${DEBLASIS_APIGATEWAY_CENTRALCOMMANDSERVICE_GRPCENDPOINT:-deblasis-v1-CentralCommandService.service.consul:9482}"
SSSVC_ENDPOINT="${DEBLASIS_APIGATEWAY_SHIPPINGSTATIONSERVICE_GRPCENDPOINT:-deblasis-v1-ShippingStationService.service.consul:9282}"

./wait-for-it.sh consul:8500 --timeout=60 -- \
./wait-for-it.sh $AUTHSVC_ENDPOINT --timeout=120 -- \
./wait-for-it.sh $CCSVC_ENDPOINT --timeout=120 -- \
./wait-for-it.sh $SSSVC_ENDPOINT --timeout=120 -- \
/exe