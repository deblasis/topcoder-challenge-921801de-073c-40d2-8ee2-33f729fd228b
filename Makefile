#!/bin/bash

CGO_ENABLED ?= 0
GOOS ?= linux
#apigateway
SERVICES = auth_dbsvc centralcommand_dbsvc authsvc centralcommandsvc
DOCKERBUILD = $(addprefix docker_,$(SERVICES))
DOCKERCLEANBUILD = $(addprefix docker_clean_,$(SERVICES))


define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -ldflags "-s -w" -o build/deblasis-$(1) services/$(1)/cmd/app/main.go
endef
define make_docker_cleanbuild
	docker build --no-cache --build-arg SVC_NAME=$(subst docker_clean_,,$(1)) --tag=deblasis/stc_$(subst docker_clean_,,$(1)) .
endef

define make_docker_build
	docker build --build-arg SVC_NAME=$(subst docker_,,$(1)) --tag=deblasis/stc_$(subst docker_,,$(1)) .
endef

all: $(SERVICES)

deps:
	go install github.com/favadi/protoc-go-inject-tag@v1.3.0


PHONY: proto
proto: deps
	buf generate
#docker run --volume "$(shell pwd):/workspace" --workdir /workspace bufbuild/buf generate
	protoc-go-inject-tag -input="gen/proto/go/centralcommand_dbsvc/v1/*.pb.go" -verbose
	protoc-go-inject-tag -input="gen/proto/go/centralcommandsvc/v1/*.pb.go" -verbose	


.PHONY: migrate-auth_dbsvc
migrate-auth_dbsvc: ## do migration
	cd ./services/auth_dbsvc/cmd/migration && go run main.go -dir ../../scripts/migrations -init

.PHONY: migrate-centralcommand_dbsvc
migrate-centralcommand_dbsvc: ## do migration
	cd ./services/centralcommand_dbsvc/cmd/migration && go run main.go -dir ../../scripts/migrations -init	

.PHONY: seed-auth_dbsvc
seed-auth_dbsvc: ## do migration
	cd ./services/auth_dbsvc/cmd/seeder && go run main.go -file ../../scripts/seeding/users.csv

.PHONY: gencert
gencert:
	cd certs && mkcert -cert-file deblasis-stc.pem -key-file deblasis-stc-key.pem spacetrafficcontrol.127.0.0.1.nip.io localhost 127.0.0.1 ::1 auth_dbsvc 

.PHONY: build-parallel
build-parallel:
	docker-compose -f docker-compose.yml -f docker-compose-build.yml build --parallel

.PHONY: run-parallel
run-parallel: build-parallel
	docker-compose up --force-recreate --remove-orphans
# auth_dbsvc:
# 	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -ldflags "-s -w" -o build/deblasis-auth_dbsvc services/auth_dbsvc/cmd/app/main.go
# && cp services/auth_dbsvc/app.yaml build


services: $(SERVICES)
docker-build: $(DOCKERBUILD)
docker-cleanbuild: $(DOCKERCLEANBUILD)

$(SERVICES):
	$(call compile_service,$(@))


$(DOCKERBUILD):
	$(call make_docker_build,$(@))

$(DOCKERCLEANBUILD):
	$(call make_docker_cleanbuild,$(@))


	