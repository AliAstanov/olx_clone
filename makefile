#!make

include .env
export $(shell sed 's/=.*//' .env)

CURRENT_DIR=$(shell pwd)

APP=$(shell basename ${CURRENT_DIR})
APP_CMD_DIR=${CURRENT_DIR}/cmd

TAG=latest
ENV_TAG=latest

gen_proto:
	./scripts/gen_proto.sh

pull-proto-module:
	git submodule update --init --recursive

update-proto-module:
	git submodule update --remote --merge

proto:
	make pull-proto-module
	make update-proto-module

POSTGRESQL_URL='postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable'

println:
	@echo ${POSTGRESQL_URL}

migrate-create:
	@read -p "Enter a migration name: " USER_INPUT; \
	migrate create -ext sql -dir migrations -seq $$USER_INPUT

migrate-local-up:
	migrate -database ${POSTGRESQL_URL} -path migrations up

migrate-local-down:
	migrate -database ${POSTGRESQL_URL} -path migrations down 1

migration-up:
	migrate -path ./migrations/postgres -database 'postgres://postgres:123@0.0.0.0:5432/ss_go_auth_service?sslmode=disable' up

migration-down:
	migrate -path ./migrations/postgres -database 'postgres://postgres:123@0.0.0.0:5432/ss_go_auth_service?sslmode=disable' down

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

build-image:
	docker build --rm -t ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} .
	docker tag ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}

push-image:
	docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG}
	docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}

swag-init:
	swag init -g api/api.go -o api/docs

run:
	go run cmd/main.go

linter:
	golangci-lint run
