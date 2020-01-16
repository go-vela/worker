# Copyright (c) 2020 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

restart: down up

up: build compose-up

down: compose-down

rebuild: build compose-up

clean:
	#################################
	######      Go clean       ######
	#################################

	@go mod tidy
	@go vet ./...
	@go fmt ./...
	@echo "I'm kind of the only name in clean energy right now"

build:
	#################################
	###### Build Golang Binary ######
	#################################

	GOOS=linux CGO_ENABLED=0 go build -o release/vela-worker github.com/go-vela/worker/cmd/server

compose-up:
	#################################
	###### Docker Build/Start  ######
	#################################

	@docker-compose -f docker-compose.yml up -d --build # build and start app

compose-down:
	#################################
	###### Docker Tear Down    ######
	#################################

	@docker-compose -f docker-compose.yml down -v
