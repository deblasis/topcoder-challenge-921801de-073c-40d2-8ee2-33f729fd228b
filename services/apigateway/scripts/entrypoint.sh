#!/usr/bin/env bash
#
# Copyright 2021 Alessandro De Blasis <alex@deblasis.net>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

AUTHSVC_ENDPOINT="${DEBLASIS_APIGATEWAY_AUTHSERVICE_GRPCENDPOINT:-deblasis-v1-AuthService.service.consul:9082}"
CCSVC_ENDPOINT="${DEBLASIS_APIGATEWAY_CENTRALCOMMANDSERVICE_GRPCENDPOINT:-deblasis-v1-CentralCommandService.service.consul:9482}"
SSSVC_ENDPOINT="${DEBLASIS_APIGATEWAY_SHIPPINGSTATIONSERVICE_GRPCENDPOINT:-deblasis-v1-ShippingStationService.service.consul:9282}"

./wait-for-it.sh consul:8500 --timeout=60 -- \
./wait-for-it.sh $AUTHSVC_ENDPOINT --timeout=120 -- \
./wait-for-it.sh $CCSVC_ENDPOINT --timeout=120 -- \
./wait-for-it.sh $SSSVC_ENDPOINT --timeout=120 -- \
/exe