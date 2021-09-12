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
ARG SVC_NAME

COPY ./Makefile ./Makefile
COPY ./common ./common
COPY ./gen ./gen
COPY ./scripts ./scripts
COPY ./services/auth_dbsvc ./services/auth_dbsvc


RUN make ${SVC_NAME} \
    && mv build/deblasis-$SVC_NAME /exe


RUN make ${SVC_NAME}_migrator \
    && mv build/deblasis-${SVC_NAME}_migrator /migrator

RUN cp ./services/auth_dbsvc/scripts/entrypoint.sh / \
    && mkdir -p /scripts/migrations \
    && cp ./services/auth_dbsvc/scripts/migrations/* /scripts/migrations

RUN make ${SVC_NAME}_seeder \
    && mv build/deblasis-${SVC_NAME}_seeder /seeder

RUN cp ./scripts/*.sh /

FROM alpine:3.14.2
RUN apk update && apk upgrade && apk add bash
COPY --from=server_builder /entrypoint.sh /exe /migrator /seeder /wait-for-it.sh  /
COPY --from=server_builder /scripts /scripts
ENTRYPOINT ["/entrypoint.sh"]