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

## Help

```
dblab is a terminal UI based interactive database client for Postgres, MySQL and SQLite.

Usage:
  dblab [flags]
  dblab [command]

Available Commands:
  help        Help about any command
  version     The version of the project

Flags:
      --config          get the connection data from a config file (default is $HOME/.dblab.yaml or the current directory)
      --db string       Database name
      --driver string   Database driver
  -h, --help            help for dblab
      --host string     Server host name or IP
      --pass string     Password for user
      --port string     Server port
      --ssl string      SSL mode
  -u, --url string      Database connection string
      --user string     Database user

Use "dblab [command] --help" for more information about a command.
```

## Usage

You can start the app passing no flags or parameters, you'll be asked for connection data instead.
![Alt Text](screenshots/dblab-default-form.gif)

```sh
$ dblab --host localhost --user myuser --db users --pass password --ssl disable --port 5432 --driver postgres
```

Connection URL scheme is also supported:

```sh
$ dblab --url postgres://user:password@host:port/database?sslmode=[mode]
$ dblab --url mysql://user:password@tcp(host:port)/db
```

Now, you can use a configuration file to make a connection to the database.

```sh
$ dbladb --config
```

`dblab` is going to look for a file called `.dblab.yaml`. For now, the only two places where you can drop a config file are $HOME ($HOME/.dblab.yaml) and the current directory where you run the command line tool.

`.dblab.yaml` example:

```yaml
database:
  host: "localhost"
  port: 5432
  db: "users"
  password: "password"
  user: "postgres"
  driver: "postgres"
```

Only the `host` and `ssl` fields are optionals. `127.0.0.1` and `disable`, respectively.

If the query panel is active, type the desired query and press <kbd>Ctrl+Space</kbd> to see the results on the rows panel below.
Otherwise, you might me located at the tables panel, then you can navigate by using the arrows <kbd>Up</kbd> and <kbd>Down</kbd> (or the keys <kbd>k</kbd> and <kbd>j</kbd> respectively). If you want to see the rows of a table, press <kbd>Enter</kbd>. To see the the schema of a table, locate yourself on the `rows` panel and press <kbd>Ctrl+S</kbd> to switch to the `structure` panel, then switch <kbd>Ctrl+S</kbd> to switch back.
The same can be achieved for the `constraints` view by pressing <kbd>Ctrl+F</kbd> to go back and forth between the `rows` and the `constraints` panels.

<img src="screenshots/dblab-screen-shot.png" />


### Key Bindings
Key                                     | Description
----------------------------------------|---------------------------------------
<kbd>Ctrl+Space</kbd>                   | If the query panel is active, execute the query
<kbd>Enter</kbd>                        | If the tables panel is active, list all the rows as a result set on the rows panel and display the structure of the table on the structure panel
<kbd>Ctrl+S</kbd>                       | If the rows panel is active, switch to the schema panel
<kbd>Ctrl+F</kbd>                       | If the rows panel is active, switch to the constraints view
<kbd>Ctrl+h</kbd>                       | Toggle to the panel on the left
<kbd>Ctrl+j</kbd>                       | Toggle to the panel below
<kbd>Ctrl+k</kbd>                       | Toggle to the panel above
<kbd>Ctrl+l</kbd>                       | Toggle to the panel on the right
<kbd>Arrow Up</kbd>                     | Next row of the result set, on the rows panel. Next table on the tables panel
<kbd>k</kbd>                            | Next row of the result set, on the rows panel. Next table on the tables panel
<kbd>Arrow Down</kbd>                   | Previous row of the result set, on the rows panel. Previous table on the tables panel
<kbd>j</kbd>                            | Previous row of the result set, on the rows panel. Previous table on the tables panel
<kbd>Arrow Right</kbd>                  | Horizontal scrolling on the rows panel
<kbd>l</kbd>                            | Horizontal scrolling on the rows panel
<kbd>Arrow Left</kbd>                   | Horizontal scrolling on the rows panel
<kbd>h</kbd>                            | Horizontal scrolling on the rows panel
<kbd>Ctrl+c</kbd>                       | Quit

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
