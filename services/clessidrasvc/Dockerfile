# The MIT License (MIT)
#
# Copyright (c) 2021 Alessandro De Blasis <alex@deblasis.net>  
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE. 
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
COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/services/$SVC_NAME/service.json /consul/service.json
COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/scripts/*.sh /
COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/build/deblasis-$SVC_NAME /exe


FROM alpine:3.14.2
RUN apk update && apk upgrade && apk add bash && apk add supervisor
COPY --from=build_base /bin/consul /bin/goreman /bin/
COPY --from=server_builder /entrypoint.sh /consul-agent.sh /exe /wait-for-it.sh /Procfile  /
COPY --from=server_builder /consul/service.json /consul/service.json

ENTRYPOINT [ "goreman" ]
CMD [ "-f", "/Procfile", "start" ]