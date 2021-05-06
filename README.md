# dblab

Interactive client for PostgreSQL and MySQL.

[![Release](https://img.shields.io/github/release/danvergara/dblab.svg?label=Release)](https://github.com/danvergara/dblab/releases)
![integration tests](https://github.com/danvergara/dblab/actions/workflows/ci.yaml/badge.svg)
![unit tests](https://github.com/danvergara/dblab/actions/workflows/test.yaml/badge.svg)

## Overview

dblab is a fast and lightweight interactive terminal based UI application for PostgreSQL and MySQL,
written in Go and works on OSX, Linux and Windows machines. Main idea behind using Go for backend development
is to utilize ability of the compiler to produce zero-dependency binaries for
multiple platforms. dblab was created as an attempt to build very simple and portable
application to work with local or remote PostgreSQL/MySQL databases.

## Features

- Cross-platform support OSX/Linux/Windows 32/64-bit
- Simple installation (distributed as a single binary)
- Zero dependencies

## Installation

- [Precompiled binaries](https://github.com/danvergara/dblab/releases) for supported
operating systems are available.

## Usage


```
dblab --host localhost --user myuser --db users --pass password --ssl disable --port 5432 --driver postgres
```

Connection URL scheme is also supported:

```
db --url postgres://user:password@host:port/database?sslmode=[mode]
db --url mysql://user:password@tcp(host:port)/db
```

## Contribute

- Fork this repository
- Create a new feature branch for a new functionality or bugfix
- Commit your changes
- Execute test suite
- Push your code and open a new pull request
- Use [issues](https://github.com/danvergara/dblab/issues) for any questions
- Check [wiki](https://github.com/danvergara/dblab/wiki) for extra documentation

## License
The MIT License (MIT). See [LICENSE](LICENSE) file for more details.
