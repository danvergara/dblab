PLATFORM=linux/amd64

.PHONY: test
## test: Runs the tests
test: test-postgres test-mysql

.PHONY: test-postgres
## test-postgres: Runs the tests with a connection to the postgres sakila database
test-postgres:
	DB_USER=sakila \
	DB_PASSWORD=sakila \
	DB_NAME=sakila \
	DB_DRIVER=postgres \
	go test -v -race ./...

.PHONY: test-mysql
## test-mysql: Runs the tests with a connection to the mysql sakila database
test-mysql:
	DB_USER=sakila \
	DB_PASSWORD=sakila \
	DB_NAME=sakila \
	DB_DRIVER=mysql \
	go test -v -race ./...

.PHONY: unit-test
## unit-test: Runs the tests with the short flag
unit-test:
	go test -v -short -race ./...

.PHONY: linter
## linter: Runs the golangci-lint command
linter:
	golangci-lint run ./...

.PHONY: build
## build: Builds the Go program
build:
	@CGO_ENABLED=0 \
	go build -o dblab .

.PHONY: connect
## connect: Runs the connect command to list the previous connections available
connect: build
	@./dblab connect

.PHONY: run
## run: Runs the application with a connection to the PostgreSQL sakila database
run: run-sakila-postgres

.PHONY: run-sakila-postgres
## run-sakila-postgres: Runs the application and connects to the PostgreSQL sakila database
run-sakila-postgres: build
	DBLAB_DEBUG="" ./dblab --host localhost --user sakila --db sakila --pass sakila --schema public --ssl disable --port 5432 --driver postgres --limit 50 --save-as sakila -k

.PHONY: run-sakila-mysql
## run-sakila-mysql: Runs the application and connects to the MySQL sakila database
run-sakila-mysql: build
	DBLAB_DEBUG="" ./dblab --host localhost --user sakila --db sakila --pass sakila --ssl disable --port 3306 --driver mysql --limit 50 -k

.PHONY: run-sakila-sqlite
## run-sakila-sqlite: Runs the application with a connection to the Sqlite Sakila database
run-sakila-sqlite: build
	./dblab --db sakila.db --driver sqlite

.PHONY: run-ssh
## run-ssh: Runs the application through a ssh tunnel
run-ssh: build
	./dblab --host postgres --user sakila --pass sakila --schema public --ssl disable --port 5432 --driver postgres --limit 50 --ssh-host 127.0.0.1 --ssh-port 2222 --ssh-user root --ssh-pass root

.PHONY: run-ssh-key
## run-ssh-key: Runs the application through a ssh tunnel using a private key file
run-ssh-key: build
	./dblab --host postgres --user sakila --pass sakila --schema public --ssl disable --port 5432 --driver postgres --limit 50 --ssh-host 127.0.0.1 --ssh-port 2222 --ssh-user root --ssh-key my_ssh_key

.PHONY: run-mysql
## run-mysql: Runs the application with a connection to mysql
run-mysql: build
	./dblab --host localhost --user sakila --db sakila --pass sakila --ssl enable --port 3306 --driver mysql --save-as mysql

.PHONY: run-mysql-ssh
## run-mysql-ssh: Runs the application through a ssh tunnel
run-mysql-ssh: build
	./dblab --host mysql --user sakila --db sakila --pass sakila --ssl enable --port 3306 --driver mysql --limit 50 --ssh-host 127.0.0.1 --ssh-port 2222 --ssh-user root --ssh-pass root

.PHONY: run-mysql-socket
## run-mysql-socket: Runs the application with a connection to mysql through a socket file. In this example the socke file is located in /var/lib/mysql/mysql.sock.
run-mysql-socket: build
	./dblab --socket /var/lib/mysql/mysql.sock --user sakila --pass sakila --db sakila --ssl enable --port 3306 --driver mysql
	
.PHONY: run-postgres-socket
## run-postgres-socket: Runs the application with a connection to postgres through a socket file. In this example the socke file is located in /var/run/postgresql.
run-postgres-socket: build
	./dblab --socket /var/run/postgresql --user  sakila --db sakila --pass sakila --ssl disable --port 5432 --driver postgres --limit 50

.PHONY: run-oracle
## run-oracle: Runs the application making a connection to the Oracle database
run-oracle: build
	./dblab --host localhost --user sys --db FREEPDB1 --pass password --port 1521 --driver oracle --limit 50

.PHONY: run-sql-server
## run-sql-server: Runs the application making a connection to the SQL Server database
run-sql-server: build
	./dblab --host localhost --user SA --db msdb --pass '5@klkbN#ABC' --port 1433 --driver sqlserver --limit 50

.PHONY: run-mysql-socket-url
## run-mysql-socket-url: Runs the application with a connection to mysql through a socket file. In this example the socke file is located in /var/lib/mysql/mysql.sock.
run-mysql-socket-url: build
	./dblab --url "mysql://sakila:sakila@unix(/var/lib/mysql/mysql.sock)/sakila?charset=utf8"

.PHONY: run-sqlite3
## run-sqlite3: Runs the application with a connection to sqlite3
run-sqlite3: build
	./dblab --db sakila.db --driver sqlite

.PHONY: run-sqlite3-url
## run-sqlite3-url: Runs the application with a connection string to sqlite3
run-sqlite3-url: build
	./dblab --url 'file:sakila.db?_pragma=foreign_keys(1)&_time_format=sqlite'

.PHONY: run-url
## run-url: Runs the app passing the url as parameter
run-url: build
	./dblab --url postgres://sakila:sakila@localhost:5432/sakila?sslmode=disable

.PHONY: run-url-ssh
## run-url-ssh: Runs the application through a ssh tunnel providing the url as parameter
run-url-ssh: build
	./dblab --url postgres://sakila:sakila@postgres:5432/sakila?sslmode=disable --schema public --ssh-host 127.0.0.1 --ssh-port 2222 --ssh-user root --ssh-pass root

.PHONY: run-mysql-url
## run-mysql-url: Runs the app passing the url as parameter
run-mysql-url: build
	./dblab --url "mysql://sakila:sakila@tcp(localhost:3306)/sakila" 

.PHONY: run-mysql-url-ssh
## run-mysql-url-ssh: Runs the app passing the url as parameter through a ssh tunnel providing the url as parameter
run-mysql-url-ssh: build
	./dblab --url "mysql://sakila:sakila@mysql+tcp(mysql:3306)/sakila" --driver mysql --ssh-host 127.0.0.1 --ssh-port 2222 --ssh-user root --ssh-pass root

.PHONY: run-config
## run-config: Runs the client using the config file.
run-config: build
	./dblab --config --cfg-name "test"

.PHONY: up
## up: Runs all the containers listed in the docker-compose.yml file
up:
	docker compose up --build -d

.PHONY: up-ssh
## up-ssh: Runs all the containers listed in the docker-compose.ssh.yml file to test the ssh tunnel
up-ssh:
	docker compose -f docker-compose.ssh.yml up -d

.PHONY: down
## down: Shut down all the containers listed in the docker-compose.yml file
down:
	docker compose down

.PHONY: stop-ssh
## stop-ssh: Shut down all the containers listed in the docker-compose.ssh.yml file
stop-ssh:
	docker compose -f docker-compose.ssh.yml down

.PHONY: form
## form: Runs the application with no arguments
form: build
	./dblab

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
