PLATFORM=linux/amd64

.PHONY: test
## test: Runs the tests
test:
	go test -v -race ./...

.PHONY: unit-test
## unit-test: Runs the tests with the short flag
unit-test:
	go test -v -short -race ./...

.PHONY: int-test
## int-test: Runs the integration tests
int-test:
	docker-compose run --entrypoint=make dblab test

.PHONY: linter
## linter: Runs the golangci-lint command
linter:
	golangci-lint run --enable=golint --enable=godot ./...

.PHONY: test-all
## test-all: Runs the integration testing bash script with different database docker image versions
test-all:
	@./scripts/test_all.sh

.PHONY: docker-build
## docker-build: Builds de Docker image
docker-build:
	@docker build --target bin --output bin/ --platform ${PLATFORM} -t dblab .

.PHONY: build
## build: Builds the Go program
build:
	go build -o dblab .

.PHONY: run
## run: Runs the application
run: build
	./dblab --host localhost --user postgres --db users --pass password --ssl disable --port 5432 --driver postgres

.PHONY: run-mysql
## run-mysql: Runs the application with a connection to mysql
run-mysql: build
	./dblab --host localhost --user myuser --db mydb --pass 5@klkbN#ABC --ssl enable --port 3306 --driver mysql

.PHONY: run-url
## run-url: Runs the app passing the url as paramenter
run-url: build
	./dblab --url postgres://postgres:password@localhost:5432/users?sslmode=disable

.PHONY: up
## up: Runs all the containers listed in the docker-compose.yml file
up:
	docker-compose up --build -d

.PHONY: down
## down: Shut down all the containers listed in the docker-compose.yml file
down:
	docker-compose down

.PHONY: form
## form: Runs the application with no arguments
form: build
	./dblab

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
