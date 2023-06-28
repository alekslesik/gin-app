# Include variables from the .envrc file
# include .env

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
