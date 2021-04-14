.PHONY: test
## test: Runs the tests
test:
	go test -v -race ./...

.PHONY: unit-test
## unit-test: Runs the tests with the short flag
unit-test:
	go test -v -short -race ./...

.PHONY: linter
## linter: Runs the colangci-lint command
linter:
	golangci-lint run --enable=golint --enable=godot ./...

.PHONY: test-all
## test-all: Runs the integration testing bash script with different database docker image versions
test-all:
	@./scripts/test_all.sh

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
