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

.PHONY: run
run: clean test cover undeploy
	PKG_DEBUG=true PKG_SQL_ENGINE=sqlite go run ./cmd/server/main.go serve

test: format
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
	sudo $$(command -v docker-compose) --file deployments/docker-compose.yml build --compress --force-rm
	sudo docker image prune --force
install: pull deploy
pull:
	sudo $$(command -v docker-compose) --file deployments/docker-compose.yml pull
push:
	@docker buildx build --platform linux/amd64,linux/arm64 -t electrocucaracha/pkg_mgr --push .
deploy: undeploy
	sudo $$(command -v docker-compose) --file deployments/docker-compose.yml --env-file deployments/.env up --force-recreate --detach --no-build
	until sudo $$(command -v docker-compose) --file deployments/docker-compose.yml logs api | grep "Serving pkg mgr at"; do \
		sleep 2; \
	done
logs:
	sudo $$(command -v docker-compose) --file deployments/docker-compose.yml logs --follow
undeploy:
	sudo $$(command -v docker-compose) --file deployments/docker-compose.yml down --remove-orphans
pre-deploy-dev: build
	sudo $$(command -v docker-compose) --file deployments/docker-compose.yml pull init
	sudo $$(command -v docker-compose) --file deployments/docker-compose.yml pull mariadb
deploy-dev: pre-deploy-dev deploy
