include .env

compose-up:
	docker-compose up --build -d

compose-down:
	docker-compose down

build-api:
	docker build -t ${API_IMAGE}

docs:
	swag init -g internal/app/app.go --pd

mocks:
	mockgen -source=internal/service/service.go -destination=internal/mocks/servicemocks/service.go -package=servicemocks

pg-tests:
	docker run --name postgres --rm -d \
		-p 6000:6000 \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=1234567890 \
		-e POSTGRES_DB=postgres postgres:15 -p 6000

redis-tests:
	docker run --name redis --rm -d -p 6379:6379 redis:latest

init-test-containers: pg-tests redis-tests

stop-test-containers:
	docker stop postgres && docker stop redis

init-tests:
	go test -v ./...

tests: init-test-containers init-tests stop-test-containers

proto:
	protoc -I proto proto/auth/*.proto --go_out=./proto/ --go_opt=paths=source_relative --go-grpc_out=./proto/ --go-grpc_opt=paths=source_relative
.PHONY: proto