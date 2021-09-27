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
FROM golang:1.17-alpine AS build_base


RUN apk add make protobuf git gcc g++ libc-dev curl jq npm openssl nss-tools jose
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
COPY ./gen/proto/go/centralcommand_dbsvc/v1/extensions.go ./gen/proto/go/centralcommand_dbsvc/v1/extensions.go
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
