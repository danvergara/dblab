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
	docker compose run --entrypoint=make dblab test

.PHONY: linter
## linter: Runs the golangci-lint command
linter:
	golangci-lint run ./...

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
	CGO_ENABLED=0 \
	go build -o dblab .

.PHONY: run
## run: Runs the application
run: build
	./dblab --host localhost --user postgres --db users --pass password --schema public --ssl disable --port 5432 --driver postgres --limit 50

.PHONY: run-mysql
## run-mysql: Runs the application with a connection to mysql
run-mysql: build
	./dblab --host localhost --user myuser --db mydb --pass 5@klkbN#ABC --ssl enable --port 3306 --driver mysql


.PHONY: run-mysql-socket
## run-mysql-socket: Runs the application with a connection to mysql through a socket file. In this example the socke file is located in /var/lib/mysql/mysql.sock.
run-mysql-socket: build
	./dblab --socket /var/lib/mysql/mysql.sock --user myuser --pass password --db mydb --ssl enable --port 3306 --driver mysql
	
.PHONY: run-postgres-socket
## run-postgres-socket: Runs the application with a connection to mysql through a socket file. In this example the socke file is located in /var/lib/mysql/mysql.sock.
run-postgres-socket: build
	./dblab --socket /var/run/postgresql --user  myuser --db my_project --pass postgres --ssl disable --port 5432 --driver postgres --limit 50

.PHONY: run-oracle
## run-oracle: Runs the application making a connection to the Oracle database
run-oracle: build
	./dblab --host localhost --user system --db FREEPDB1 --pass password --port 1521 --driver oracle --limit 50

.PHONY: run-sql-server
## run-sql-server: Runs the application making a connection to the SQL Server database
run-sql-server: build
	./dblab --host localhost --user SA --db msdb --pass '5@klkbN#ABC' --port 1433 --driver sqlserver --limit 50

.PHONY: run-mysql-socket-url
## run-mysql-socket-url: Runs the application with a connection to mysql through a socket file. In this example the socke file is located in /var/lib/mysql/mysql.sock.
run-mysql-socket-url: build
	./dblab --url "mysql://myuser:password@unix(/var/lib/mysql/mysql.sock)/mydb?charset=utf8"

.PHONY: run-sqlite3
## run-sqlite3: Runs the application with a connection to sqlite3
run-sqlite3: build
	./dblab --db db/dblab.db --driver sqlite

.PHONY: run-sqlite3-url
## run-sqlite3-url: Runs the application with a connection string to sqlite3
run-sqlite3-url: build
	./dblab --url 'file:db/dblab.db?_pragma=foreign_keys(1)&_time_format=sqlite'

.PHONY: run-url
## run-url: Runs the app passing the url as parameter
run-url: build
	./dblab --url postgres://postgres:password@localhost:5432/users?sslmode=disable

.PHONY: run-mysql-url
## run-mysql-url: Runs the app passing the url as parameter
run-mysql-url: build
	./dblab --url "mysql://myuser:5@klkbN#ABC@tcp(localhost:3306)/mydb"

.PHONY: run-config
## run-config: Runs the client using the config file.
run-config: build
	./dblab --config --cfg-name "test"

.PHONY: up
## up: Runs all the containers listed in the docker-compose.yml file
up:
	docker compose up --build -d

.PHONY: down
## down: Shut down all the containers listed in the docker-compose.yml file
down:
	docker compose down

.PHONY: form
## form: Runs the application with no arguments
form: build
	./dblab

.PHONY: create-migration
## create: Creates golang-migrate migration files
create-migration:
	 migrate create -ext sql -dir ./db/migrations $(file_name)

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
