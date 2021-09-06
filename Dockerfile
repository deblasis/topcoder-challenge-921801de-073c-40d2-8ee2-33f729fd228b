# Base build image
FROM golang:1.17-alpine AS build_base
ARG SVC_NAME

RUN apk add make protobuf git gcc g++ libc-dev
WORKDIR /go/src/deblasis.net/space-traffic-control

COPY go.mod .
COPY go.sum .
RUN go mod download


RUN go install google.golang.org/protobuf/cmd/protoc-gen-go
RUN cp /go/bin/protoc-gen-go /usr/local/bin/

RUN go install github.com/gogo/protobuf/protoc-gen-gofast
RUN cp /go/bin/protoc-gen-gofast /usr/local/bin/

RUN echo $SVC_NAME

FROM build_base AS server_builder
ARG SVC_NAME

COPY . .

RUN make $SVC_NAME \
    && mv build/deblasis-$SVC_NAME /exe \
    && ls
##    && cp build/app.yaml /

FROM scratch
#COPY --from=server_builder /app.yaml /
COPY --from=server_builder /exe /
ENTRYPOINT ["/exe"]