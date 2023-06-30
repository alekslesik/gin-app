# Include variables from the .envrc file
# include .env

PROJECT := aklatan
ALLGO := $(wildcard *.go */*.go cmd/*/*.go)
ALLHTML := $(wildcard templates/*/*.html)

#=====================================#
# HELPERS #
#=====================================#

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.DELETE_ON_ERROR:

.PHONY: all
all: lint check $(PROJECT) $(CMDS)

.PHONY: lint
lint: .lint

.lint: $(ALLGO)
	golangci-lint run --timeout 180s
	@touch $@

.coverage:
	@mkdir -p .coverage

.PHONY: check
check: .coverage ./.coverage/$(PROJECT).out
./.coverage/$(PROJECT).out: $(ALLGO) $(ALLHTML) Makefile
	go test $(TESTFLAGS) -coverprofile=./.coverage/$(PROJECT).out ./...

.PHONY: cover
# When running manually, capture just the total percentage (and
# beautify it slightly because the tool output is usually too wide).
cover: .coverage ./.coverage/$(PROJECT).html
	@echo "Checking overall code coverage..."
	@go tool cover -func .coverage/$(PROJECT).out | sed -n -e '/^total/s/:.*statements)[^0-9]*/: /p'

./.coverage/$(PROJECT).html: ./.coverage/$(PROJECT).out
	go tool cover -html=./.coverage/$(PROJECT).out -o ./.coverage/$(PROJECT).html

# XXX *.go isn't quite right here -- it will rebuild when tests are
# touched, but it's good enough.
$(PROJECT): $(ALLGO)
	go build .

$(CMDS): $(ALLGO)
	for cmd in $(CMDS); do go build ./cmd/$$cmd; done

# XXX This only works if go-junit-report is installed. It's not part of go.mod
# because I don't want to force a dependency, but it is part of the ci docker
# image.
report.xml: $(ALLGO) Makefile
	go test $(TESTFLAGS) -v ./... 2>&1 | go-junit-report > $@
	go tool cover -func .coverage/$(PROJECT).out


#=====================================#
# DOCKER #
#=====================================#

# ## run: Exec docker-up
# .PHONY: run
# up: docker-up

## docker-up: Build images before starting containers; Create and start containers
.PHONY: docker-up
docker-up:
	docker-compose up --build -d

## docker-rebuild: Stop, rebuild and start containers
.PHONY: docker-rebuild
docker-rebuild:
	make docker-down
	docker-compose up --build -d

## docker-down: Stop and remove containers, networks; Remove containers for services not defined in the Compose file.
.PHONY: docker-down
docker-down:
	docker-compose down --remove-orphans

## docker-restart: Restart service containers
.PHONY: docker-restart
docker-restart:
	docker compose restart

## docker-stop: Stop service containers
.PHONY: docker-stop
docker-stop:
	docker compose stop

## docker-prune: Delete all stopped containers, images and networks
.PHONY: docker-prune
docker-prune:
	docker system prune

## mysql-up: Build bitrix-mysql image before starting containers; Create and start containers
.PHONY: mysql-up
mysql-up:
	docker-compose up --build -d mysql

## mysql-exec: Start bin/sh in bitrix-mysql
.PHONY: mysql-exec
mysql-exec:
	docker container exec -it ${PROJECT_NAME}-mysql sh

## golang-up: Build php-apache image before starting containers; Create and start containers
.PHONY: golang-up
golang-up:
	docker-compose up --build -d ${PROJECT_NAME}-golang

## golang-exec: Start bin/sh in php-docker
.PHONY: golang-exec
golang-exec:
	docker container exec -it ${PROJECT_NAME}-golang sh



	
