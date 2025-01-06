# How to contribute to dblab

Thank you for considering contributing to dblab!

## First time setup

- Fork dblab to your GitHub account by clicking the [Fork](https://github.com/danvergara/dblab/fork) button.
- [Clone](https://docs.github.com/en/github/getting-started-with-github/fork-a-repo#step-2-create-a-local-clone-of-your-fork) the main repository locally.

```bash
git clone https://github.com/danvergara/dblab.git
cd dblab
```

- Add your fork as a remote to push your work to. Replace ``{username}`` with your username. This names the remote "fork", the  default dblab remote is "origin".

```bash
git remote add fork https://github.com/{username}/dblab
```

## Start coding

- Create a branch to identify the issue, feature addition or change you would like to work on.

```bash
git fetch origin
git checkout -b your-branch-name origin/main
```

- Using your favorite editor, make your changes.
- Include tests that cover any code changes you make. Make sure the
  test fails without your patch. Run the tests as described below.
- Commit your changes following the [conventional commit](https://www.conventionalcommits.org/en/v1.0.0/) ([commitizen](https://github.com/commitizen-tools/commitizen) is a great tool for this purpose). The reasoning behind it is that it is easier to read, and enforces writing descriptive commits.
- Push your commits to your fork on GitHub and
  [create a pull request](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request). Link to the issue being addressed with
  ``fixes #123`` in the pull request in the case that you are working on an issue.

```bash
git push --set-upstream fork your-branch-name
```

## Running the tests

### Setting up the environment

First of all, the repository includes a docker-compose file and a makefile. The docker file declares 4 services: a postgres container, a mysql one, a dblab container that runs migrations and seeds on the postgres container and other one that does the same for the mysql container. The makefile is used to automate some task. To spin up all the containers type:

```bash
make up
```

Then, you'll have 2 containers up and running in your system for every kind of database supported by dblab.

Run dblab opening a connection with the postgres container.

```bash
make run
```

You can do the same with the mysql container by running:

```bash
make run-mysql
```

Run the unit test suite with make.

```bash
make unit-test
```

Run the integration test with make too, but make sure the database containers are up and running.

```bash
make int-test
```

## SSH Tunnel

There is an special compose file that spins up an ssh server, to test the ssh tunnel and work with it. The compose file also provides postgres and mysql containers but they are not exposed to the localhost. The sshd server is the intermediary between the client and those containers.

Run the command below to spin up the ssh server and the databases containers behind it.

```bash
make up-ssh
```

To connect to the databases, the make file provides a new series of targets addding the ssh related parameters:

```bash
make run-ssh
```

The command above, is the equivalent of this command:

```bash
dblab --host postgres --user postgres --pass password --schema public --ssl disable --port 5432 --driver postgres --limit 50 --ssh-host localhost --ssh-port 2222 --ssh-user root --ssh-pass root
```

You can check all the options with `help` command.

```bash
Usage:
  test                   Runs the tests
  unit-test              Runs the tests with the short flag
  int-test               Runs the integration tests
  linter                 Runs the golangci-lint command
  test-all               Runs the integration testing bash script with different database docker image versions
  docker-build           Builds de Docker image
  build                  Builds the Go program
  run                    Runs the application
  run-ssh                Runs the application through a ssh tunnel
  run-ssh-key            Runs the application through a ssh tunnel using a private key file
  run-mysql              Runs the application with a connection to mysql
  run-mysql-ssh          Runs the application through a ssh tunnel
  run-mysql-socket       Runs the application with a connection to mysql through a socket file. In this example the socke file is located in /var/lib/mysql/mysql.sock.
  run-postgres-socket    Runs the application with a connection to mysql through a socket file. In this example the socke file is located in /var/lib/mysql/mysql.sock.
  run-oracle             Runs the application making a connection to the Oracle database
  run-sql-server         Runs the application making a connection to the SQL Server database
  run-mysql-socket-url   Runs the application with a connection to mysql through a socket file. In this example the socke file is located in /var/lib/mysql/mysql.sock.
  run-sqlite3            Runs the application with a connection to sqlite3
  run-sqlite3-url        Runs the application with a connection string to sqlite3
  run-url                Runs the app passing the url as parameter
  run-url-ssh            Runs the application through a ssh tunnel providing the url as parameter
  run-mysql-url          Runs the app passing the url as parameter
  run-mysql-url-ssh      Runs the app passing the url as parameter through a ssh tunnel providing the url as parameter
  run-config             Runs the client using the config file.
  up                     Runs all the containers listed in the docker-compose.yml file
  up-ssh                 Runs all the containers listed in the docker-compose.ssh.yml file to test the ssh tunnel
  down                   Shut down all the containers listed in the docker-compose.yml file
  stop-ssh               Shut down all the containers listed in the docker-compose.ssh.yml file
  form                   Runs the application with no arguments
  create                 Creates golang-migrate migration files
  help                   Prints this help message
```
