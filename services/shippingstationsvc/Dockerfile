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
COPY ./services/centralcommand_dbsvc ./services/centralcommand_dbsvc
COPY ./services/centralcommandsvc ./services/centralcommandsvc
RUN cp ./services/${SVC_NAME}/scripts/entrypoint.sh /
RUN cp ./scripts/*.sh /   

RUN make $SVC_NAME \
    && mv build/deblasis-$SVC_NAME /exe

FROM alpine:3.14.2
RUN apk update && apk upgrade && apk add bash
COPY --from=server_builder /entrypoint.sh /exe /wait-for-it.sh  /
ENTRYPOINT ["/entrypoint.sh"]