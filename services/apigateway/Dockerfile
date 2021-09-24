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


ENTRYPOINT [ "goreman" ]
CMD [ "-f", "/Procfile", "start" ]