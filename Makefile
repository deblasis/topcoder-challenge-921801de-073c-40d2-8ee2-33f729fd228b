#!/bin/bash

CGO_ENABLED ?= 0
GOOS ?= linux

#built using Docker version 20.10.2, build 2291f61 on Windows with WSL2
#depending on your setup, you might have to change the below line to DOCKERCOMPOSE = docker-compose 
DOCKERCOMPOSE ?= docker compose

SERVICES = apigateway auth_dbsvc centralcommand_dbsvc authsvc centralcommandsvc shippingstationsvc clessidrasvc
MIGRATORS = auth_dbsvc_migrator centralcommand_dbsvc_migrator
SEEDERS = auth_dbsvc_seeder
DOCKERBUILD = $(addprefix docker_,$(SERVICES))
DOCKERCLEANBUILD = $(addprefix docker_clean_,$(SERVICES))
INJECTPROTOTAGS = inject_prototags_ $(addprefix inject_prototags_,$(SERVICES))
APIGATEWAY ?= http://localhost:8081
APIGATEWAY_NOPROTOCOL=$(shell echo $(APIGATEWAY) | sed -E 's/^\s*.*:\/\///g')
WAIT4IT=./scripts/wait-for-it.sh

PREFIX=/usr/local
BUFVERSION=1.0.0-rc2


define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -ldflags "-s -w" -o ./build/deblasis-$(1) ./services/$(1)/cmd/app/main.go
endef
define compile_migrator
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -ldflags "-s -w" -o ./build/deblasis-$(subst _migrator,,$(1))_migrator ./services/$(subst _migrator,,$(1))/cmd/migrator/main.go
endef
define compile_seeder
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -ldflags "-s -w" -o ./build/deblasis-$(subst _seeder,,$(1))_seeder ./services/$(subst _seeder,,$(1))/cmd/seeder/main.go
endef
define make_docker_cleanbuild
	docker build --no-cache --build-arg SVC_NAME=$(subst docker_clean_,,$(1)) --tag=deblasis/stc_$(subst docker_clean_,,$(1)) .
endef

define make_docker_build
	docker build --build-arg SVC_NAME=$(subst docker_,,$(1)) --tag=deblasis/stc_$(subst docker_,,$(1)) .
endef

define make_inject_prototags ## FIX THIS FOR apigateway
	protoc-go-inject-tag -input="gen/proto/go/$(subst inject_prototags_,,$(1))/v1/*.pb.go" -verbose || true
endef

all: $(SERVICES)

.PHONY: install-buf
install-buf:
ifeq ("", "$(shell which buf)")
		curl -sSL "https://github.com/bufbuild/buf/releases/download/v$(BUFVERSION)/buf-$(shell uname -s)-$(shell uname -m).tar.gz" | \
		tar -xvzf - -C "$(PREFIX)" --strip-components 1
endif

.PHONY: protodeps
protodeps: 
	go install github.com/favadi/protoc-go-inject-tag@v1.3.0
	go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc


.PHONY: proto
proto: protodeps install-buf
	buf generate
	make injectprototags

.PHONY: binaries
binaries:
	make services
	make migrators
	make seeders


.PHONY: migrate-auth_dbsvc
migrate-auth_dbsvc: ## do migration
	cd ./services/auth_dbsvc/cmd/migrator && go run main.go -dir ../../scripts/migrations -init


.PHONY: migrate-centralcommand_dbsvc
migrate-centralcommand_dbsvc: ## do migration
	cd ./services/centralcommand_dbsvc/cmd/migrator && go run main.go -dir ../../scripts/migrations -init	


.PHONY: seed-auth_dbsvc
seed-auth_dbsvc: ## do migration
# the following commented line is one way to seed if the db is accessible from localhost
##	cd ./services/auth_dbsvc/cmd/seeder && go run main.go -file ../../scripts/seeding/users.csv
# this is a better approach that injects the file into the container and uses the current environment
	docker cp ./services/auth_dbsvc/scripts/seeding $(shell docker ps -qf "ancestor=deblasis/stc_auth_dbsvc:latest"):/seeding
	docker exec -t $(shell docker ps -qf "ancestor=deblasis/stc_auth_dbsvc:latest") /bin/bash -c "./seeder -file ./seeding/users.csv"
# the seeding seeds a tmp table, a service restart is required to seed the real table, this is for security reasons	
	docker restart $(shell docker ps -qf "ancestor=deblasis/stc_auth_dbsvc:latest")

.PHONY: certdeps
certdeps:
	wget https://github.com/FiloSottile/mkcert/releases/download/v1.4.3/mkcert-v1.4.3-linux-amd64 \
	&& chmod +x mkcert-v1.4.3-linux-amd64 \
	&& mv mkcert-v1.4.3-linux-amd64 /usr/local/bin/mkcert

.PHONY: host-gencerts
host-gencerts: 
	mkdir -p ./certs \
	&& docker build -f ./common/tools/jose-jwk/Dockerfile ./common/tools/jose-jwk -t jose-jwt \
	&& docker run jose-jwt -c "jose jwk gen -i '{\"alg\": \"RS256\"}'" > ./certs/jwk-private.json \
	&& cat ./certs/jwk-private.json | jq '{kid: "$(shell openssl rand -base64 32)", alg: .alg, kty: .kty , use: "sig", n: .n , e: .e }'  > ./certs/jwk.json \
	&& npx pem-jwk ./certs/jwk-private.json > ./certs/jwt.pem.key \
	&& openssl rsa -in ./certs/jwt.pem.key -pubout -outform PEM -out ./certs/jwt.pem.pub \
	&& cd certs && mkcert -cert-file deblasis-stc.pem -key-file deblasis-stc-key.pem spacetrafficcontrol.127.0.0.1.nip.io localhost 127.0.0.1 ::1 authsvc \
	&& ls


.PHONY: docker-gencerts
docker-gencerts: certdeps
	mkdir -p ./certs \
	&& jose jwk gen -i '{"alg": "RS256"}' > ./certs/jwk-private.json \
	&& cat ./certs/jwk-private.json | jq '{kid: "$(shell openssl rand -base64 32)", alg: .alg, kty: .kty , use: "sig", n: .n , e: .e }'  > ./certs/jwk.json \
	&& npx pem-jwk ./certs/jwk-private.json > ./certs/jwt.pem.key \
	&& openssl rsa -in ./certs/jwt.pem.key -pubout -outform PEM -out ./certs/jwt.pem.pub \
	&& cd certs && mkcert -cert-file deblasis-stc.pem -key-file deblasis-stc-key.pem spacetrafficcontrol.127.0.0.1.nip.io localhost 127.0.0.1 ::1 authsvc \
	&& ls
	

.PHONY: getcerts
getcerts:
	docker cp $(shell docker ps -qf "ancestor=deblasis/stc_authsvc:latest"):/certs/jwt.pem.key ./certs/jwt.pem.key
	docker cp $(shell docker ps -qf "ancestor=deblasis/stc_authsvc:latest"):/certs/jwt.pem.pub ./certs/jwt.pem.pub


.PHONY: builder
builder:
	docker build . --tag=deblasis/stc_builder  

.PHONY: docker-build
docker-build: builder
	$(DOCKERCOMPOSE) -f docker-compose.yml -f docker-compose.build.yml build --parallel
	
.PHONY: host-build
host-build: proto binaries builder
	$(DOCKERCOMPOSE) -f docker-compose.yml -f docker-compose.hostbuild.yml -f docker-compose.prod.yml -f docker-compose.integrationtests.yml build --parallel

.PHONY: integrationtests-build
integrationtests-build: builder
	COMPOSE_PROJECT_NAME=deblasis-stc-e2e_tests $(DOCKERCOMPOSE) -f docker-compose.yml -f docker-compose.build.yml -f docker-compose.integrationtests.yml build --parallel

.PHONY: integrationtests-up
integrationtests-up: integrationtests-build
	COMPOSE_PROJECT_NAME=deblasis-stc-e2e_tests	$(DOCKERCOMPOSE) -f docker-compose.yml -f docker-compose.integrationtests.yml up -d --force-recreate --remove-orphans

.PHONY: integrationtests-run
integrationtests-run: 
	COMPOSE_PROJECT_NAME=deblasis-stc-e2e_tests	$(DOCKERCOMPOSE) -f docker-compose.yml -f docker-compose.integrationtests.yml up integrationtester


.PHONY: dockertest
dockertest:
	go install github.com/onsi/ginkgo/ginkgo@v1.16.4
	ginkgo -race -v -tags integration ./e2e_tests

.PHONY: run-fast
run-fast: host-build
	$(DOCKERCOMPOSE) -f docker-compose.yml -f docker-compose.hostports.yml -f docker-compose.prod.yml up --remove-orphans

.PHONY: run-fast-integrationtests
run-fast-integrationtests: host-build
	$(DOCKERCOMPOSE) -f docker-compose.yml -f docker-compose.hostports.yml -f docker-compose.integrationtests.yml up -d --remove-orphans
	make getcerts
	$(DOCKERCOMPOSE) -f docker-compose.yml -f docker-compose.hostports.yml -f docker-compose.integrationtests.yml up --remove-orphans



services: $(SERVICES)
migrators: $(MIGRATORS)
seeders: $(SEEDERS)
docker-build: $(DOCKERBUILD)
docker-cleanbuild: $(DOCKERCLEANBUILD)
injectprototags: $(INJECTPROTOTAGS)

$(SERVICES):
	$(call compile_service,$(@))

$(MIGRATORS):
	$(call compile_migrator,$(@))	

$(SEEDERS):
	$(call compile_seeder,$(@))	

$(DOCKERBUILD):
	$(call make_docker_build,$(@))

$(DOCKERCLEANBUILD):
	$(call make_docker_cleanbuild,$(@))

$(INJECTPROTOTAGS):
	$(call make_inject_prototags,$(@))
