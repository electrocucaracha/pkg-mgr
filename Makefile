# SPDX-license-identifier: Apache-2.0
##############################################################################
# Copyright (c) 2020
# All rights reserved. This program and the accompanying materials
# are made available under the terms of the Apache License, Version 2.0
# which accompanies this distribution, and is available at
# http://www.apache.org/licenses/LICENSE-2.0
##############################################################################

PWD := $(shell pwd)
BINARY := pkg_mgr

export GO111MODULE=on

format:
	@go fmt ./...

swagger:
	@rm -rf gen/*
	@swagger generate server -t gen -f ./api/openapi-spec/swagger.yaml --exclude-main -A pkg-mgr

run: clean test cover
	PKG_SQL_ENGINE=sqlite go run ./cmd/main.go

test:
	@go test -v ./...

.PHONY: cover
cover:
	@go test -race ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

clean:
	@rm -f *.db
	@rm -f coverage.*
	@rm -f $(BINARY)

build: clean
	@go build -o $(PWD)/$(BINARY) cmd/main.go

docker: clean
	@docker-compose --file deployments/docker-compose.yml build --compress --force-rm
	@docker image prune --force

deploy: undeploy
	@docker-compose --file deployments/docker-compose.yml --env-file deployments/.env up --force-recreate --detach
undeploy:
	@docker-compose --file deployments/docker-compose.yml down --remove-orphans