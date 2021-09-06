#!/bin/bash

CGO_ENABLED ?= 0
GOOS ?= linux

SERVICES = auth_dbsvc authsvc
DOCKERBUILD = $(addprefix docker_,$(SERVICES))
DOCKERCLEANBUILD = $(addprefix docker_clean_,$(SERVICES))


define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -ldflags "-s -w" -o build/deblasis-$(1) services/$(1)/cmd/app/main.go
endef
define make_docker_cleanbuild
	docker build --no-cache --build-arg SVC_NAME=$(subst docker_clean_,,$(1)) --tag=deblasis/stc-$(subst docker_clean_,,$(1)) .
endef

define make_docker_build
	docker build --build-arg SVC_NAME=$(subst docker_,,$(1)) --tag=deblasis/stc-$(subst docker_,,$(1)) .
endef

all: $(SERVICES)

PHONY: proto
proto:
#	@./scripts/protobuf-gen.sh
	buf generate

.PHONY: migrate-auth_dbsvc
migrate-auth_dbsvc: ## do migration
	cd ./services/auth_dbsvc/cmd/migration && go run main.go -dir ../../scripts/migrations -init

.PHONY: seed-auth_dbsvc
seed-auth_dbsvc: ## do migration
	cd ./services/auth_dbsvc/cmd/seeder && go run main.go -file ../../scripts/seeding/users.csv

.PHONY: gencert
gencert:
	cd certs && mkcert -cert-file deblasis-stc.pem -key-file deblasis-stc-key.pem spacetrafficcontrol.127.0.0.1.nip.io localhost 127.0.0.1 ::1 auth_dbsvc 

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


	