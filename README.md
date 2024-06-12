dblab ![integration tests](https://github.com/danvergara/dblab/actions/workflows/ci.yaml/badge.svg) ![unit tests](https://github.com/danvergara/dblab/actions/workflows/test.yaml/badge.svg) [![Release](https://img.shields.io/github/release/danvergara/dblab.svg?label=Release)](https://github.com/danvergara/dblab/releases)
===============

<p align="center">
  <img style="float: right;" src="assets/gopher-dblab.png" alt="dblab logo"/  width=200>
</p>

__Interactive client for PostgreSQL, MySQL, SQLite3 and Oracle.__

<img src="screenshots/dblab-cover.png" />

---

**Documentation**: <a href="https://dblab.danvergara.com" target="_blank">https://dblab.danvergara.com</a>

---

## Table of contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
    - [Homebrew](#homebrew)
    - [Binary Release](#binary-release-linuxosxwindows)
    - [Automated installation/update](#automated-installationupdate)
- [Help Command](#help)
- [Usage](#usage)
    - [Configuration](#configuration)
- [Navigation](#navigation)
    - [Key Bindings](#key-bindings)
- [Contribute](#contribute)
- [License](#license)

## Overview

dblab is a fast and lightweight interactive terminal based UI application for PostgreSQL, MySQL and SQLite3,
written in Go and works on OSX, Linux and Windows machines. Main idea behind using Go for backend development
is to utilize ability of the compiler to produce zero-dependency binaries for
multiple platforms. dblab was created as an attempt to build very simple and portable
application to work with local or remote PostgreSQL/MySQL/SQLite3/Oracle databases.

## Features

- Cross-platform support OSX/Linux/Windows 32/64-bit
- Simple installation (distributed as a single binary)
- Zero dependencies.

## Installation

> ~~if you need to work with SQLite3, install the CGO enabled binary using the proper bash script listed below.~~
> The above comment is deprecated and CGO is not needed anymore. There will be a single binary capable to deal with all supported clients.

### Homebrew

It works with Linux, too.

```
$ brew install danvergara/tools/dblab
```

Or

```
$ brew tap danvergara/tools
$ brew install dblab
```

### Binary Release (Linux/OSX/Windows)
You can manually download a binary release from [the release page](https://github.com/danvergara/dblab/releases).

## Automated installation/update
> Don't forget to always verify what you're piping into bash

Install the binarry using our bash script:

```sh
curl https://raw.githubusercontent.com/danvergara/dblab/master/scripts/install_update_linux.sh | bash
```

## Help

```
dblab is a terminal UI based interactive database client for Postgres, MySQL and SQLite3.

Usage:
  dblab [flags]
  dblab [command]

Available Commands:
  help        Help about any command
  version     The version of the project

Flags:
      --cfg-name string      Database config name section
      --config               Get the connection data from a config file (default locations are: current directory, $HOME/.dblab.yaml or $XDG_CONFIG_HOME/.dblab.yaml)
      --db string            Database name
      --driver string        Database driver
  -h, --help                 help for dblab
      --host string          Server host name or IP
      --limit uint           Size of the result set from the table content query (should be greater than zero, otherwise the app will error out) (default 100)
      --pass string          Password for user
      --port string          Server port
      --schema string        Database schema (postgres only)
      --socket string        Path to a Unix socket file
      --ssl string           SSL mode
      --ssl-verify string    [enable|disable] or [true|false] enable ssl verify for the server
      --sslcert string       This parameter specifies the file name of the client SSL certificate, replacing the default ~/.postgresql/postgresql.crt
      --sslkey string        This parameter specifies the location for the secret key used for the client certificate. It can either specify a file name that will be used instead of the default ~/.postgresql/postgresql.key, or it can specify a key obtained from an external “engine”
      --sslpassword string   This parameter specifies the password for the secret key specified in sslkey
      --sslrootcert string   This parameter specifies the name of a file containing SSL certificate authority (CA) certificate(s) The default is ~/.postgresql/root.crt
      --trace-file string    File name for trace log
  -u, --url string           Database connection string
      --user string          Database user
      --wallet string        Path for auto-login oracle wallet

Use "dblab [command] --help" for more information about a command.
```

## Usage

You can start the app passing no flags or parameters, you'll be asked for connection data instead.
![Alt Text](screenshots/dblab-default-form.gif)


```sh
$ dblab --host localhost --user myuser --db users --pass password --ssl disable --port 5432 --driver postgres --limit 50
$ dblab --db path/to/file.sqlite3 --driver sqlite
$ dblab --host localhost --user system --db FREEPDB1 --pass password --port 1521 --driver oracle --limit 50
```

Connection URL scheme is also supported:

```sh
$ dblab --url 'postgres://user:password@host:port/database?sslmode=[mode]'
$ dblab --url 'mysql://user:password@tcp(host:port)/db'
$ dblab --url 'file:test.db?_pragma=foreign_keys(1)&_time_format=sqlite'
$ dblab --url 'oracle://user:password@localhost:1521/db'
```

if you're using PostgreSQL, you have the option to define the schema you want to work with, the default value is `public`.

```sh
$ dblab --host localhost --user myuser --db users --pass password --schema myschema --ssl disable --port 5432 --driver postgres --limit 50
$ dblab --url postgres://user:password@host:port/database?sslmode=[mode] --schema myschema
```

As a request made in [#125](https://github.com/danvergara/dblab/issues/125), support for MySQL/MariaDB sockets was integrated.

```sh
$ dblab --url "mysql://user:password@unix(/path/to/socket/mysql.sock)/dbname?charset=utf8"
$ dblab --socket /path/to/socket/mysql.sock --user user --db dbname --pass password --ssl disable --port 5432 --driver mysql --limit 50
```

Postgres connection through Unix sockets:

```sh
$ dblab --url "postgres://user:password@/dbname?host=/path/to/socket"
$ dblab --socket /path/to/socket --user user --db dbname --pass password --ssl disable --port 5432 --driver postgres --limit 50
```

Now, it is possible to ensure SSL connections with `PostgreSQL` databases. SSL related parameters has been added, such as `--sslcert`, `--sslkey`, `--sslpassword`, `--sslrootcert`. More information on how to use such connection flags can be found [here](https://www.postgresql.org/docs/current/libpq-connect.html).

```{ .sh .copy }
dblab --host  db-postgresql-nyc3-56456-do-user-foo-0.fake.db.ondigitalocean.com --user myuser --db users --pass password --schema myschema --port 5432 --driver postgres --limit 50 --ssl require --sslrootcert ~/Downloads/foo.crt
```

### Configuration

Enter previous flags every time is tedious, so `dblab` provides a couple of flags to help with it: `--config` and `--cfg-name`.

`dblab` is going to look for a file called `.dblab.yaml`. Currently, there are three places where you can drop a config file:

- $XDG_CONFIG_HOME ($XDG_CONFIG_HOME/.dblab.yaml)
- $HOME ($HOME/.dblab.yaml)
- . (the current directory where you run the command line tool)

If you want to use this feature, `--config` is mandatory and `--cfg-name` may be omitted. The config file can store one or multiple database connection sections under the `database` field. `database` is an array, previously was an object only able to store a single connection section at a time. 

We strongly encourgae you to adopt the new format as of `v0.18.0`. `--cfg-name` takes the name of the desired database section to connect with. It can be omitted and its default values will be the first item on the array. 

As of `v0.21.0`, ssl connections options are supported in the config file.

```sh
# default: test
$ dbladb --config

$ dblab --config --cfg-name "prod"
```

`.dblab.yaml` example:

```yaml
database:
  - name: "test"
    host: "localhost"
    port: 5432
    db: "users"
    password: "password"
    user: "postgres"
    driver: "postgres"
    # optional
    # postgres only
    # default value: public
    schema: "myschema"
  - name: "prod"
    # example endpoint
    host: "mydb.123456789012.us-east-1.rds.amazonaws.com"
    port: 5432
    db: "users"
    password: "password"
    user: "postgres"
    schema: "public"
    driver: "postgres"
    ssl: "require"
    sslrootcert: "~/.postgresql/root.crt."
  - name: "oracle"
    host: "localhost"
    port: 1521
    db: "FREEPDB1 "
    password: "password"
    user: "system"
    driver: "oracle"
    ssl: "enable"
    wallet: "path/to/wallet"
    ssl-verify: true
# should be greater than 0, otherwise the app will error out
limit: 50
```

Or for sqlite:

```yaml
database:
  - name: "prod"
    db: "path/to/file.sqlite3"
    driver: "sqlite"
```

Only the `host` and `ssl` fields are optionals. `127.0.0.1` and `disable`, respectively.

## Navigation

If the query panel is active, type the desired query and press <kbd>Ctrl+Space</kbd> to see the results on the rows panel below.
Otherwise, you might me located at the tables panel, then you can navigate by using the arrows <kbd>Up</kbd> and <kbd>Down</kbd> (or the keys <kbd>k</kbd> and <kbd>j</kbd> respectively). If you want to see the rows of a table, press <kbd>Enter</kbd>. To see the the schema of a table, locate yourself on the `rows` panel and press <kbd>Ctrl+S</kbd> to switch to the `structure` panel, then switch <kbd>Ctrl+S</kbd> to switch back.
The same can be achieved for the `constraints` view by pressing <kbd>Ctrl+F</kbd> to go back and forth between the `rows` and the `constraints` panels.

Now, there's a menu to navigate between hidden views by just clicking on the desired options:

<img src="screenshots/rows-view.png" />
<img src="screenshots/structure-view.png" />
<img src="screenshots/constraints-view.png" />
<img src="screenshots/indexes-view.png" />

As you may have noticed, navigation has already been added, so every time you query the content of a listed table, the result set is going to be paginated. This allows to the user dealing with large tables, optimizing resources.
Just hit the `BACK` and `NEXT` buttons to go back and forth.

### Key Bindings
Key                                     | Description
----------------------------------------|---------------------------------------
<kbd>Ctrl+Space</kbd>                   | If the query panel is active, execute the query
<kbd>Enter</kbd>                        | If the tables panel is active, list all the rows as a result set on the rows panel and display the structure of the table on the structure panel
<kbd>Ctrl+S</kbd>                       | If the rows panel is active, switch to the schema panel. The opposite is true
<kbd>Ctrl+F</kbd>                       | If the rows panel is active, switch to the constraints view. The opposite is true
<kbd>Ctrl+I</kbd>                       | If the rows panel is active, switch to the indexes view. The opposite is true
<kbd>Ctrl+H</kbd>                       | Toggle to the panel on the left
<kbd>Ctrl+J</kbd>                       | Toggle to the panel below
<kbd>Ctrl+K</kbd>                       | Toggle to the panel above
<kbd>Ctrl+L</kbd>                       | Toggle to the panel on the right
<kbd>Arrow Up</kbd>                     | Next row of the result set on the panel. Views: rows, table, constraints, structure and indexes
<kbd>k</kbd>                            | Next row of the result set on the panel. Views: rows, table, constraints, structure and indexes
<kbd>Arrow Down</kbd>                   | Previous row of the result set on the panel. Views: rows, table, constraints, structure and indexes
<kbd>j</kbd>                            | Previous row of the result set on the panel. Views: rows, table, constraints, structure and indexes
<kbd>Arrow Right</kbd>                  | Horizontal scrolling on the panel. Views: rows, constraints, structure and indexes
<kbd>l</kbd>                            | Horizontal scrolling on the panel. Views: rows, constraints, structure and indexes
<kbd>Arrow Left</kbd>                   | Horizontal scrolling on the panel. Views: rows, constraints, structure and indexes
<kbd>h</kbd>                            | Horizontal scrolling on the panel. Views: rows, constraints, structure and indexes
<kbd>0</kbd>                            | Move cursor to the start of the current line. Views: rows, constraints, and structure
<kbd>$</kbd>                            | Move cursor to the end of the current line. Views: rows, constraints, and structure
<kbd>g</kbd>                            | Move cursor to the top of the panel's dataset. Views: rows, constraints, and structure
<kbd>G</kbd>                            | Move cursor to the bottom of the panel's dataset. Views: rows, constraints, and structure
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
