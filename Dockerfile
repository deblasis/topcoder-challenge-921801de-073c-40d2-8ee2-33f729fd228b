FROM golang:1.17-alpine AS build_base


RUN apk add make protobuf git gcc g++ libc-dev curl jq npm libnss3-tools jose
WORKDIR /go/src/deblasis.net/space-traffic-control


COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./Makefile ./Makefile
COPY ./buf.*.yaml  ./
COPY ./common/proto ./common/proto


COPY ./services/auth_dbsvc/proto ./services/auth_dbsvc/proto
COPY ./services/centralcommand_dbsvc/proto ./services/centralcommand_dbsvc/proto
COPY ./services/authsvc/proto ./services/authsvc/proto
COPY ./services/centralcommandsvc/proto ./services/centralcommandsvc/proto
COPY ./services/shippingstationsvc/proto ./services/shippingstationsvc/proto

#TODO refactor this
COPY ./gen/proto/go/v1/extensions.go ./gen/proto/go/v1/extensions.go
COPY ./gen/proto/go/authsvc/v1/extensions.go ./gen/proto/go/authsvc/v1/extensions.go
COPY ./gen/proto/go/centralcommandsvc/v1/extensions.go ./gen/proto/go/centralcommandsvc/v1/extensions.go
COPY ./gen/proto/go/shippingstationsvc/v1/extensions.go ./gen/proto/go/shippingstationsvc/v1/extensions.go
###
RUN make proto
#Provision certificates... the wrong way, the right one is out of scope here, I'd use Vault
RUN make docker-gencerts

FROM build_base AS server_builder

WORKDIR /go/src/deblasis.net/space-traffic-control

COPY --from=build_base /go/src/deblasis.net/space-traffic-control/gen ./gen
COPY --from=build_base /go/src/deblasis.net/space-traffic-control/certs ./certs

COPY ./Makefile ./Makefile
COPY ./common ./common
COPY ./services ./services
COPY ./scripts ./scripts
COPY ./scripts/*.sh /

RUN make binaries
