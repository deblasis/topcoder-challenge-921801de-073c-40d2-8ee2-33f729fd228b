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
# Base build image
FROM golang:1.17-alpine AS build_base
ARG SVC_NAME

RUN apk add make
WORKDIR /go/src/deblasis.net/space-traffic-control

RUN echo $SVC_NAME

FROM build_base AS server_builder
ARG SVC_NAME

COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/scripts ./scripts

COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/build/deblasis-${SVC_NAME}_migrator /migrator

COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/services/${SVC_NAME}/scripts/entrypoint.sh /
RUN mkdir -p /scripts/migrations
COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/services/${SVC_NAME}/scripts/migrations/* /scripts/migrations/

COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/scripts/*.sh /
COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/build/deblasis-$SVC_NAME /exe

FROM alpine:3.14.2
RUN apk update && apk upgrade && apk add bash
COPY --from=server_builder /entrypoint.sh /exe /migrator /wait-for-it.sh /
COPY --from=server_builder /scripts /scripts
ENTRYPOINT ["/entrypoint.sh"]
