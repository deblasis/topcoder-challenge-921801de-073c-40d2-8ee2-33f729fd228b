# Base build image
FROM golang:1.17-alpine AS build_base
ARG SVC_NAME

RUN apk add make
WORKDIR /go/src/deblasis.net/space-traffic-control

RUN echo $SVC_NAME

FROM build_base AS server_builder
ARG SVC_NAME

COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/scripts ./scripts
COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/services/${SVC_NAME}/scripts/entrypoint.sh /
COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/scripts/*.sh /   
COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/build/deblasis-$SVC_NAME /exe

FROM alpine:3.14.2
RUN apk update && apk upgrade && apk add bash
COPY --from=server_builder /entrypoint.sh /exe /wait-for-it.sh  /

COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/certs/jwt.pem.key /certs/jwt.pem.key 
COPY --from=deblasis/stc_builder:latest /go/src/deblasis.net/space-traffic-control/certs/jwt.pem.pub /certs/jwt.pem.pub

ENTRYPOINT ["/entrypoint.sh"]