# Base build image
FROM golang:1.17-alpine AS build_base
ARG SVC_NAME

RUN apk add make protobuf git gcc g++ libc-dev
WORKDIR /go/src/deblasis.net/space-traffic-control

COPY go.mod .
COPY go.sum .
RUN go mod download

RUN echo $SVC_NAME

FROM build_base AS server_builder

ADD https://releases.hashicorp.com/consul/1.10.2/consul_1.10.2_linux_amd64.zip /tmp/consul.zip
RUN cd /bin && unzip /tmp/consul.zip && chmod +x /bin/consul && rm /tmp/consul.zip

RUN go install github.com/mattn/goreman@v0.3.7 && cp /go/bin/goreman /bin
ARG SVC_NAME


COPY ./Makefile ./Makefile
COPY ./common ./common
COPY ./gen ./gen
COPY ./scripts ./scripts
COPY ./services/auth_dbsvc ./services/auth_dbsvc
COPY ./services/authsvc ./services/authsvc
COPY ./services/centralcommand_dbsvc ./services/centralcommand_dbsvc
COPY ./services/centralcommandsvc ./services/centralcommandsvc
COPY ./services/apigateway ./services/apigateway

RUN cp ./services/$SVC_NAME/scripts/*.sh /
RUN cp ./services/$SVC_NAME/Procfile /
# RUN cp ./services/$SVC_NAME/supervisord.conf /
RUN cp ./scripts/*.sh /



RUN make $SVC_NAME \
    && mv build/deblasis-$SVC_NAME /exe

FROM alpine:3.14.2
RUN apk update && apk upgrade && apk add bash && apk add supervisor
RUN mkdir -p /var/log/supervisor
COPY --from=server_builder /entrypoint.sh /consul-agent.sh /exe /wait-for-it.sh /Procfile  /
COPY --from=server_builder /bin/consul /bin/goreman /bin/
# COPY --from=server_builder /supervisord.conf /

# ENTRYPOINT ["/usr/bin/supervisord", "-c", "/supervisord.conf"]

#ENTRYPOINT [ "/consul-agent.sh" ]
# CMD [ "-f", "/Procfile", "start" ]

ENTRYPOINT [ "goreman" ]
CMD [ "-f", "/Procfile", "start" ]