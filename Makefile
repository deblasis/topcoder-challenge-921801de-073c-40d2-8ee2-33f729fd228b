#!/bin/bash

CGO_ENABLED ?= 0
GOOS ?= linux

SERVICES = apigateway auth_dbsvc centralcommand_dbsvc authsvc centralcommandsvc
MIGRATORS = auth_dbsvc_migrator centralcommand_dbsvc_migrator
SEEDERS = auth_dbsvc_seeder
DOCKERBUILD = $(addprefix docker_,$(SERVICES))
DOCKERCLEANBUILD = $(addprefix docker_clean_,$(SERVICES))
INJECTPROTOTAGS = $(addprefix inject_prototags_,$(SERVICES))


define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -ldflags "-s -w" -o build/deblasis-$(1) services/$(1)/cmd/app/main.go
endef
define compile_migrator
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -ldflags "-s -w" -o build/deblasis-$(subst _migrator,,$(1))_migrator services/$(subst _migrator,,$(1))/cmd/migrator/main.go
endef
define compile_seeder
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -ldflags "-s -w" -o build/deblasis-$(subst _seeder,,$(1))_seeder services/$(subst _seeder,,$(1))/cmd/seeder/main.go
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

protodeps:
	go install github.com/favadi/protoc-go-inject-tag@v1.3.0
	go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc


PHONY: proto
proto: protodeps
	buf generate
	make injectprototags

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
# thie is a better approach that injects the file into the container and uses the current environment
	docker cp ./services/auth_dbsvc/scripts/seeding $(shell docker ps -qf "name=^deblasis-stc_auth_dbsvc"):/seeding
	docker exec -it $(shell docker ps -qf "name=^deblasis-stc_auth_dbsvc") /bin/bash -c "./seeder -file ./seeding/users.csv"
# the seeding seeds a tmp table, a service restart is required to seed the real table, this is for security reasons	
	docker restart $(shell docker ps -qf "name=^deblasis-stc_auth_dbsvc")

.PHONY: gencert
gencert:
	docker build -f ./common/tools/jose-jwk/Dockerfile ./common/tools/jose-jwk -t jose-jwt \
	&& docker run jose-jwt -c "jose jwk gen -i '{\"alg\": \"RS256\"}'" > ./certs/jwk-private.json \
	&& cat ./certs/jwk-private.json | jq '{kid: "$(shell openssl rand -base64 32)", alg: .alg, kty: .kty , use: "sig", n: .n , e: .e }'  > ./certs/jwk.json \
	&& npx pem-jwk ./certs/jwk-private.json > ./certs/jwt-key.pem \
	&& openssl rsa -in ./certs/jwt-key.pem -pubout -outform PEM -out ./certs/jwt-pubout.pem \
	&& cd certs && mkcert -cert-file deblasis-stc.pem -key-file deblasis-stc-key.pem spacetrafficcontrol.127.0.0.1.nip.io localhost 127.0.0.1 ::1 authsvc \

	 

.PHONY: build-parallel
build-parallel: proto
	docker-compose -f docker-compose.yml -f docker-compose-build.yml build --parallel
.PHONY: run-parallel
run-parallel: build-parallel
	docker-compose -f docker-compose.yml up --force-recreate --remove-orphans

.PHONY: build-on-host
build-on-host: proto
	make services
	make migrators
	make seeders
	docker-compose -f docker-compose.yml -f docker-compose-hostbuild.yml build --parallel

.PHONY: run-fast
run-fast: build-on-host
	docker-compose -f docker-compose.yml up --force-recreate --remove-orphans
	


# auth_dbsvc:
# 	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -ldflags "-s -w" -o build/deblasis-auth_dbsvc services/auth_dbsvc/cmd/app/main.go
# && cp services/auth_dbsvc/app.yaml build




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
