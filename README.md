dblab ![integration tests](https://github.com/danvergara/dblab/actions/workflows/ci.yaml/badge.svg) ![unit tests](https://github.com/danvergara/dblab/actions/workflows/test.yaml/badge.svg) [![Release](https://img.shields.io/github/release/danvergara/dblab.svg?label=Release)](https://github.com/danvergara/dblab/releases)
===============

<p align="center">
  <img style="float: right;" src="assets/gopher-dblab.png" alt="dblab logo"/  width=200>
</p>

__Interactive client for PostgreSQL, MySQL, SQLite3, Oracle and SQL Server.__

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
    - [SSH Tunnel](#ssh-tunnel)
    - [Configuration](#configuration)
        - [Key bindings configuration](#key-bindings-configuration) 
- [Navigation](#navigation)
    - [Key Bindings](#key-bindings)
- [Contribute](#contribute)
- [License](#license)

## Overview

dblab is a fast and lightweight interactive terminal based UI application for PostgreSQL, MySQL and SQLite3,
written in Go and works on OSX, Linux and Windows machines. Main idea behind using Go for backend development
is to utilize ability of the compiler to produce zero-dependency binaries for
multiple platforms. dblab was created as an attempt to build very simple and portable
application to work with local or remote PostgreSQL/MySQL/SQLite3/Oracle/SQL Server databases.

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
      --cfg-name string                   Database config name section
      --config                            Get the connection data from a config file (default locations are: current directory, $HOME/.dblab.yaml or $XDG_CONFIG_HOME/.dblab.yaml)
      --keybindings, -k                   Get the keybindings configuration from the config file (default locations are: current directory, $HOME/.dblab.yaml or $XDG_CONFIG_HOME/.dblab.yaml)
      --db string                         Database name (optional)
      --driver string                     Database driver
      --encrypt string                    [strict|disable|false|true] data sent between client and server is encrypted or not
  -h, --help                              help for dblab
      --host string                       Server host name or IP
      --limit uint                        Size of the result set from the table content query (should be greater than zero, otherwise the app will error out) (default 100)
      --pass string                       Password for user
      --port string                       Server port
      --schema string                     Database schema (postgres only)
      --socket string                     Path to a Unix socket file
      --ssh-host string                   SSH Server Hostname/IP
      --ssh-key string                    File with private key for SSH authentication
      --ssh-key-pass string               Supports connections with protected private keys with passphrase
      --ssh-pass string                   SSH Password (Empty string for no password)
      --ssh-port string                   SSH Port
      --ssh-user string                   SSH User
      --ssl string                        SSL mode
      --ssl-verify string                 [enable|disable] or [true|false] enable ssl verify for the server
      --sslcert string                    This parameter specifies the file name of the client SSL certificate, replacing the default ~/.postgresql/postgresql.crt
      --sslkey string                     This parameter specifies the location for the secret key used for the client certificate. It can either specify a file name that will be used instead of the default ~/.postgresql/postgresql.key, or it can specify a key obtained from an external “engine”
      --sslpassword string                This parameter specifies the password for the secret key specified in sslkey
      --sslrootcert string                This parameter specifies the name of a file containing SSL certificate authority (CA) certificate(s) The default is ~/.postgresql/root.crt
      --timeout string                    in seconds (default is 0 for no timeout), set to 0 for no timeout. Recommended to set to 0 and use context to manage query and connection timeouts
      --trace-file string                 File name for trace log
      --trust-server-certificate string   [false|true] server certificate is checked or not
  -u, --url string                        Database connection string
      --user string                       Database user
  -v, --version                           version for dblab
      --wallet string                     Path for auto-login oracle wallet

Use "dblab [command] --help" for more information about a command.
```

## Usage

You can start the app passing no flags or parameters, you'll be asked for connection data instead.
![dblab-demo](screenshots/dblab-demo.gif)


```sh
$ dblab --host localhost --user myuser --db users --pass password --ssl disable --port 5432 --driver postgres --limit 50
$ dblab --db path/to/file.sqlite3 --driver sqlite
$ dblab --host localhost --user system --db FREEPDB1 --pass password --port 1521 --driver oracle --limit 50
$ dblab --host localhost --user SA --db msdb --pass '5@klkbN#ABC' --port 1433 --driver sqlserver --limit 50
```

Connection URL scheme is also supported:

```sh
$ dblab --url 'postgres://user:password@host:port/database?sslmode=[mode]'
$ dblab --url 'mysql://user:password@tcp(host:port)/db'
$ dblab --url 'file:test.db?_pragma=foreign_keys(1)&_time_format=sqlite'
$ dblab --url 'oracle://user:password@localhost:1521/db'
$ dblab --url 'sqlserver://SA:myStrong(!)Password@localhost:1433?database=tempdb&encrypt=true&trustservercertificate=false&connection+timeout=30'
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

### SSH Tunnel

Now, it's possible to connect to Postgres or MySQL (more to come later) databases on a server via SSH using password or a ssh key files.

To do so, 6 new flags has been added to the dblab command:

| Flag                 | Description                                                       |
|----------------------|-------------------------------------------------------------------|
|  --ssh-host          |  SSH Server Hostname/IP                                           |
|  --ssh-port          |  SSH Port                                                         |
|  --ssh-user          |  SSH User                                                         |
|  --ssh-pass          |  SSH Password (Empty string for no password)                      |
|  --ssh-key           |  File with private key for SSH authentication                     |
|  --ssh-key-pass      | Passphrase for protected private key files                        |

#### Examples

Postgres connection via ssh tunnel using password:

```{ .sh .copy }
dblab --host localhost --user postgres --pass password --schema public --ssl disable --port 5432 --driver postgres --limit 50 --ssh-host example.com --ssh-port 22 --ssh-user root --ssh-pass root
```

Postgres connection via ssh tunnel using ssh private key file:

```{ .sh .copy }
dblab --host localhost --user postgres --pass password --schema public --ssl disable --port 5432 --driver postgres --limit 50 --ssh-host example.com --ssh-port 22 --ssh-user root --ssh-key my_ssh_key --ssh-key-pass password
```

Postgres connection using the url parameter via ssh tunnel using password:

```{ .sh .copy }
dblab --url postgres://postgres:password@localhost:5432/users?sslmode=disable --schema public --ssh-host example.com --ssh-port 22 --ssh-user root --ssh-pass root
```

MySQL connection via ssh tunnel using password:

```{ .sh .copy }
dblab --host localhost --user myuser --db mydb --pass 5@klkbN#ABC --ssl enable --port 3306 --driver mysql --limit 50 --ssh-host example.com --ssh-port 22 --ssh-user root --ssh-pass root
```

MySQL connection via ssh tunnel using ssh private key file:

```{ .sh .copy }
dblab --host localhost --user postgres --pass password --ssl enable --port 3306 --driver mysql --limit 50 --ssh-host example.com --ssh-port 22 --ssh-user root --ssh-key my_ssh_key --ssh-key-pass passphrase
```

MySQL connection using the url parameter via ssh tunnel using password:

```{ .sh .copy }
dblab --url "mysql://myuser:5@klkbN#ABC@mysql+tcp(localhost:3306)/mydb" --driver mysql --ssh-host example.com --ssh-port 22 --ssh-user root --ssh-pass root
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
$ dblab --config

$ dblab --config --cfg-name "prod"
```

#### Key bindings configuration

Key bindings can be configured through the `.dblab.yaml` file. There is a field called `keybindings` where key bindings can be modified. By default, the keybindings are not loaded, so you need to use the `--keybindings` or `-k` flag to load them. See the example to see the full list of the key bindings subject to change. The file shows the default values. The list of the available key bindings belongs to the [tcell](https://github.com/gdamore/tcell) library. Specifically, see the [KeyNames map](https://github.com/gdamore/tcell/blob/781586687ddb57c9d44727dc9320340c4d049b11/key.go#L83), for an accurate reference.

#### .dblab.yaml example

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
  - name: "sqlserver"
    driver: "sqlserver"
    host: "localhost"
    port: 1433
    db: "msdb"
    password: "5@klkbN#ABC"
    user: "SA"
  - name: "ssh-tunnel"
    host: "localhost"
    port: 5432
    db: "users"
    password: "password"
    user: "postgres"
    schema: "public"
    driver: "postgres"
    ssh-host: "example.com"
    ssh-port: 22
    ssh-user: "ssh-user"
    ssh-pass: "password"
# should be greater than 0, otherwise the app will error out
limit: 50
keybindings:
  run-query: 'Ctrl-Space'
  structure: 'Ctrl-S'
  indexes: 'Ctrl-I'
  constraints: 'Ctrl-T'
  clear-editor: 'Ctrl-D'
  navigation:
    up: 'Ctrl-K'
    down: 'Ctrl-J'
    left: 'Ctrl-H'
    right: 'Ctrl-L'
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

The key bindings are now configurable, see [Key bindings configuration](#key-bindings-configuration) to learn how to replace existing key bindings. It's worth noting that key bindings are only configurable through the configuration file, there is no flags to do so. If you don't replace them through the configuration file, the information below remains the same, otherwise, just replace the new key binding with the existing information for the default one.

If the query panel is active, type the desired query and press <kbd>Ctrl+Space</kbd> to see the results on the rows panel below.
Otherwise, you might me located at the tables panel, then you can navigate by using the arrows <kbd>Up</kbd> and <kbd>Down</kbd> (or the keys <kbd>k</kbd> and <kbd>j</kbd> respectively). If you want to see the rows of a table, press <kbd>Enter</kbd>. To see the the schema of a table, locate yourself on the `rows` panel and press <kbd>Ctrl+S</kbd> to switch to the `structure` panel, then switch <kbd>Ctrl+S</kbd> to switch back.
The same can be achieved for the `constraints` view by pressing <kbd>Ctrl+F</kbd> to go back and forth between the `rows` and the `constraints` panels.

Now, there's a menu to navigate between hidden views by just clicking on the desired options:

<img src="screenshots/rows-view.png" />
<img src="screenshots/structure-view.png" />
<img src="screenshots/constraints-view.png" />
<img src="screenshots/indexes-view.png" />

~~As you may have noticed, navigation has already been added, so every time you query the content of a listed table, the result set is going to be paginated. This allows to the user dealing with large tables, optimizing resources.
Just hit the `BACK` and `NEXT` buttons to go back and forth.~~

The navigation buttons were removed since they are too slow to really navigate the content of a table. The user is better off typing a `SELECT` statement with proper `OFFSET` and `LIMIT`.

The `--db` flag is now optional (except for Oracle), meaning that the user will be able to see the list of databases they have access to. The regular list of tables will be replaced with a tree structure showing a list of databases and their respective list of tables, branching off each database. Due to the nature of the vast majority of DBMSs that don't allow cross-database queries, dblab has to open an independent connection for each database. The side effect of this decision, is that the user has to press `Enter` on the specific database of interest. An indicator showing the current active database will appear at the bottom-right of the screen. To change the focus, just hit enter on another database. Once a database is selected, the usual behavior of inspecting tables remains the same.

<img src="screenshots/tree-view.png" />

At the time to navigate queries' resultsets, the cell will be highlighted, so the user can know which table cell is on. This is important because you can press the `Enter` key on a cell of interest to copy its content.

### Key Bindings
| Key                                     | Description                           |
|----------------------------------------|----------------------------------------|
|<kbd>Ctrl+Space</kbd>                   | If the query editor is active, execute the query |
|<kbd>Ctrl+D</kbd>                       | Cleans the whole text from the query editor, when the editor is selected |
|<kbd>Enter</kbd>                        | If the tables panel is active, list all the rows as a result set on the rows panel and display the structure of the table on the structure panel |
|<kbd>Ctrl+S</kbd>                       | If the rows panel is active, switch to the schema panel. The opposite is true |
|<kbd>Ctrl+T</kbd>                       | If the rows panel is active, switch to the constraints view. The opposite is true |
|<kbd>Ctrl+I</kbd>                       | If the rows panel is active, switch to the indexes view. The opposite is true |
|<kbd>Ctrl+H</kbd>                       | Toggle to the panel on the left |
|<kbd>Ctrl+J</kbd>                       | Toggle to the panel below |
|<kbd>Ctrl+K</kbd>                       | Toggle to the panel above |
|<kbd>Ctrl+L</kbd>                       | Toggle to the panel on the right |
|<kbd>Arrow Up</kbd>                     | Vertical scrolling on the panel. Views: rows, table, constraints, structure and indexes |
|<kbd>k</kbd>                            | Vertical scrolling on the panel. Views: rows, table, constraints, structure and indexes |
|<kbd>Arrow Down</kbd>                   | Vertical on the panel. Views: rows, table, constraints, structure and indexes |
|<kbd>j</kbd>                            | Vertical on the panel. Views: rows, table, constraints, structure and indexes |
|<kbd>Arrow Right</kbd>                  | Horizontal scrolling on the panel. Views: rows, constraints, structure and indexes |
|<kbd>l</kbd>                            | Horizontal scrolling on the panel. Views: rows, constraints, structure and indexes |
|<kbd>Arrow Left</kbd>                   | Horizontal scrolling on the panel. Views: rows, constraints, structure and indexes |
|<kbd>h</kbd>                            | Horizontal scrolling on the panel. Views: rows, constraints, structure and indexes |
|<kbd>g</kbd>                            | Move cursor to the top of the panel's dataset. Views: rows, constraints, structure and indexes |
|<kbd>G</kbd>                            | Move cursor to the bottom of the panel's dataset. Views: rows, constraints, structure and indexes |
|<kbd>Ctrl-F</kbd>                       | Move down by one page. Views: rows, constraints, structure and indexes |
|<kbd>Ctrl-B</kbd>                       | Move up by one page. Views: rows, constraints, structure and indexes |
|<kbd>Enter</kbd>                        | Press the Enter key on a resultset table cell to copy the content |
|<kbd>Ctrl+c</kbd>                       | Quit |

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
