#!/bin/bash



PHONY: proto
proto:
	@./scripts/protobuf-gen.sh



.PHONY: migrate-auth_dbsvc
migrate-auth_dbsvc: ## do migration
	cd ./services/auth_dbsvc/cmd/migration && go run main.go -dir ../../scripts/migrations -init


.PHONY: seed-auth_dbsvc
seed-auth_dbsvc: ## do migration
	cd ./services/auth_dbsvc/cmd/seeder && go run main.go -file ../../scripts/seeding/users.csv
