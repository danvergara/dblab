.PHONY: test
## test: Runs the tests
test:
	go test -v -race ./...

.PHONY: test-all
## test-all: Runs the integration testing bash script with different database docker image versions
test-all:
	@./scripts/test_all.sh

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
