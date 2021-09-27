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

RUN go install github.com/mattn/goreman@v0.3.7 && cp /go/bin/goreman /bin
ADD https://releases.hashicorp.com/consul/1.10.2/consul_1.10.2_linux_amd64.zip /tmp/consul.zip
RUN cd /bin && unzip /tmp/consul.zip && chmod +x /bin/consul && rm /tmp/consul.zip

FROM build_base AS server_builder

ARG SVC_NAME

COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/scripts ./scripts
COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/services/$SVC_NAME/scripts/*.sh /
COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/services/$SVC_NAME/Procfile /
COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/scripts/*.sh /

COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/build/deblasis-$SVC_NAME /exe


FROM alpine:3.14.2
RUN apk update && apk upgrade && apk add bash && apk add supervisor
RUN mkdir -p /var/log/supervisor
COPY --from=build_base /bin/consul /bin/goreman /bin/
COPY --from=server_builder /entrypoint.sh /consul-agent.sh /exe /wait-for-it.sh /Procfile  /
COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/certs/jwt.pem.pub /certs/jwt.pem.pub

ENTRYPOINT [ "goreman" ]
CMD [ "-f", "/Procfile", "start" ]