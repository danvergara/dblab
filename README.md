# dblab 
![integration tests](https://github.com/danvergara/dblab/actions/workflows/ci.yaml/badge.svg)
![unit tests](https://github.com/danvergara/dblab/actions/workflows/test.yaml/badge.svg)
[![Release](https://img.shields.io/github/release/danvergara/dblab.svg?label=Release)](https://github.com/danvergara/dblab/releases)

<p align="center">
  <img style="float: right;" src="assets/gopher-dblab.png" alt="dblab logo"/  width=200>
</p>

__Interactive client for PostgreSQL, MySQL, SQLite3, Oracle and SQL Server.__

<img src="screenshots/dblab-cover.png" />

---

**Documentation**: <a href="https://dblab.app" target="_blank">https://dblab.app</a>

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
    - [Connection Profiles](#connection-profiles)
- [Navigation](#navigation)
    - [Query editor](#query-editor)
    - [Query History](#query-history)
    - [Key Bindings](#key-bindings)
- [Contribute](#contribute)
- [License](#license)

## Overview

dblab is a fast and lightweight interactive terminal-based UI application for PostgreSQL, MySQL, and SQLite3,
written in Go and works on macOS, Linux, and Windows machines. The main idea behind using Go for backend development
is to utilize the ability of the compiler to produce zero-dependency binaries for
multiple platforms. dblab was created as an attempt to build a very simple and portable
application to work with local or remote PostgreSQL/MySQL/SQLite3/Oracle/SQL Server databases.

## Features

- Cross-platform support for macOS/Linux/Windows (32/64-bit)
- Simple installation (distributed as a single binary)
- Zero dependencies.
- Vim-style query editor (normal and insert modes, line-oriented editing commands).
- Multi-query execution: write multiple SQL statements separated by `;` and run them concurrently with results displayed in separate tabs.
- Connection profiles with secure credential storage in the OS keyring.
- Query history: executed queries are persisted across sessions and can be browsed/re-used via a filterable list.
- Read-only mode: use `--readonly` to prevent accidental writes by forcing the database session into read-only mode (supported for PostgreSQL, MySQL, SQLite, Oracle, and SQL Server).

## Installation

> The above comment is deprecated and CGO is not needed anymore. There will be a single binary capable of dealing with all supported clients.

### Homebrew

It works with Linux, too.

```
brew install --cask danvergara/tools/dblab
```

Or

```
brew tap danvergara/tools
brew install --cask dblab
```

### Binary Release (Linux/macOS/Windows)
You can manually download a binary release from [the release page](https://github.com/danvergara/dblab/releases).

## Automated installation/update
> Don't forget to always verify what you're piping into bash

Install the binary using our bash script:

```sh
curl https://raw.githubusercontent.com/danvergara/dblab/master/scripts/install_update_linux.sh | bash
```

## Help

```
dblab is a terminal UI-based interactive database client

Usage:
  dblab [flags]
  dblab [command]

Available Commands:
  connect     Re-use saved connection profiles
  help        Help about any command
  version     The version of the project

Flags:
      --cfg-name string                   Database config name section
      --config                            Get the connection data from a config file (default locations are: current directory, $HOME/.dblab.yaml or $XDG_CONFIG_HOME/.dblab.yaml)
      --keybindings, -k                   Get the keybindings configuration from the config file (default locations are: current directory, $HOME/.dblab.yaml or $XDG_CONFIG_HOME/.dblab.yaml)
      --db string                         Database name
      --driver string                     Database driver
      --encrypt string                    [strict|disable|false|true] whether data sent between client and server is encrypted
  -h, --help                              help for dblab
      --host string                       Server host name or IP
      --limit uint                        Size of the result set for the table content query (should be greater than zero, otherwise the app will error out) (default 100)
      --pass string                       Password for user
      --port string                       Server port
      --save-as string                    Save the connection as a named profile for later reuse
      --schema string                     Database schema (optional for postgres and oracle only)
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
      --sslrootcert string                This parameter specifies the name of a file containing SSL certificate authority (CA) certificate(s). The default is ~/.postgresql/root.crt
      --timeout string                    in seconds (default is 0 for no timeout), set to 0 for no timeout. Recommended to set to 0 and use context to manage query and connection timeouts
      --trace-file string                 File name for trace log
      --trust-server-certificate string   [false|true] whether the server certificate is checked
  -u, --url string                        Database connection string
      --user string                       Database user
  -v, --version                           version for dblab
      --readonly                            Forces a read-only connection with the target database
      --wallet string                     Path for auto-login oracle wallet

Use "dblab [command] --help" for more information about a command.
```

## Usage

You can start the app without passing flags or parameters; you'll be asked for connection data instead.
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

If you're using PostgreSQL or Oracle, you have the option to define the schema you want to work with. The `--schema` flag is optional: if omitted, dblab will display all schemas the connected user has access to in the sidebar tree. If provided, only that specific schema will be shown.

```sh
# Postgres
$ dblab --host localhost --user myuser --db users --pass password --schema myschema --ssl disable --port 5432 --driver postgres --limit 50
$ dblab --url postgres://user:password@host:port/database?sslmode=[mode] --schema myschema

# Oracle
$ dblab --host localhost --user user2 --db FREEPDB1 --pass password --port 1521 --driver oracle --limit 50 --schema user1
$ dblab --url 'oracle://user2:password@localhost:1521/FREEPDB1' --schema user1
```

You can use the `--readonly` flag to open a connection in read-only mode. This prevents any write operations (INSERT, UPDATE, DELETE, etc.) from being executed, which is useful when you want to safely browse a production database. The same can be achieved via the configuration file by setting `readonly: true` on a database profile (see [Configuration](#configuration)).

```sh
# Postgres
$ dblab --host localhost --user myuser --db users --pass password --ssl disable --port 5432 --driver postgres --limit 50 --readonly

# MySQL
$ dblab --host localhost --user myuser --db mydb --pass password --ssl disable --port 3306 --driver mysql --limit 50 --readonly

# SQLite
$ dblab --db path/to/file.sqlite3 --driver sqlite --readonly

# Oracle
$ dblab --host localhost --user system --db FREEPDB1 --pass password --port 1521 --driver oracle --limit 50 --readonly

# SQL Server
$ dblab --host localhost --user SA --db msdb --pass '5@klkbN#ABC' --port 1433 --driver sqlserver --limit 50 --readonly
```

<img src="screenshots/dblab-read-only.png" />

As requested in [#125](https://github.com/danvergara/dblab/issues/125), support for MySQL/MariaDB sockets was integrated.

```sh
$ dblab --url "mysql://user:password@unix(/path/to/socket/mysql.sock)/dbname?charset=utf8"
$ dblab --socket /path/to/socket/mysql.sock --user user --db dbname --pass password --ssl disable --port 5432 --driver mysql --limit 50
```

Postgres connection through Unix sockets:

```sh
$ dblab --url "postgres://user:password@/dbname?host=/path/to/socket"
$ dblab --socket /path/to/socket --user user --db dbname --pass password --ssl disable --port 5432 --driver postgres --limit 50
```

Now, it is possible to ensure SSL connections with `PostgreSQL` databases. SSL-related parameters have been added, such as `--sslcert`, `--sslkey`, `--sslpassword`, and `--sslrootcert`. More information on how to use such connection flags can be found [here](https://www.postgresql.org/docs/current/libpq-connect.html).

```{ .sh .copy }
dblab --host  db-postgresql-nyc3-56456-do-user-foo-0.fake.db.ondigitalocean.com --user myuser --db users --pass password --schema myschema --port 5432 --driver postgres --limit 50 --ssl require --sslrootcert ~/Downloads/foo.crt
```

### SSH Tunnel

Now, it's possible to connect to Postgres or MySQL (more to come later) databases on a server via SSH using a password or SSH key files.

To do so, 6 new flags have been added to the dblab command:

| Flag                 | Description                                                       |
|----------------------|-------------------------------------------------------------------|
|  --ssh-host          |  SSH Server Hostname/IP                                           |
|  --ssh-port          |  SSH Port                                                         |
|  --ssh-user          |  SSH User                                                         |
|  --ssh-pass          |  SSH Password (Empty string for no password)                      |
|  --ssh-key           |  File with private key for SSH authentication                     |
|  --ssh-key-pass      | Passphrase for protected private key files                        |

#### Examples

Postgres connection via SSH tunnel using a password:

```{ .sh .copy }
dblab --host localhost --user postgres --pass password --schema public --ssl disable --port 5432 --driver postgres --limit 50 --ssh-host example.com --ssh-port 22 --ssh-user root --ssh-pass root
```

Postgres connection via SSH tunnel using an SSH private key file:

```{ .sh .copy }
dblab --host localhost --user postgres --pass password --schema public --ssl disable --port 5432 --driver postgres --limit 50 --ssh-host example.com --ssh-port 22 --ssh-user root --ssh-key my_ssh_key --ssh-key-pass password
```

Postgres connection using the url parameter via SSH tunnel using a password:

```{ .sh .copy }
dblab --url postgres://postgres:password@localhost:5432/users?sslmode=disable --schema public --ssh-host example.com --ssh-port 22 --ssh-user root --ssh-pass root
```

MySQL connection via SSH tunnel using a password:

```{ .sh .copy }
dblab --host localhost --user myuser --db mydb --pass 5@klkbN#ABC --ssl enable --port 3306 --driver mysql --limit 50 --ssh-host example.com --ssh-port 22 --ssh-user root --ssh-pass root
```

MySQL connection via SSH tunnel using an SSH private key file:

```{ .sh .copy }
dblab --host localhost --user postgres --pass password --ssl enable --port 3306 --driver mysql --limit 50 --ssh-host example.com --ssh-port 22 --ssh-user root --ssh-key my_ssh_key --ssh-key-pass passphrase
```

MySQL connection using the url parameter via SSH tunnel using a password:

```{ .sh .copy }
dblab --url "mysql://myuser:5@klkbN#ABC@mysql+tcp(localhost:3306)/mydb" --driver mysql --ssh-host example.com --ssh-port 22 --ssh-user root --ssh-pass root
```

### Configuration

Entering these flags every time is tedious, so `dblab` provides a couple of flags to help with it: `--config` and `--cfg-name`.

`dblab` is going to look for a file called `.dblab.yaml`. Currently, there are three places where you can drop a config file:

- $XDG_CONFIG_HOME ($XDG_CONFIG_HOME/.dblab.yaml)
- $HOME ($HOME/.dblab.yaml)
- . (the current directory where you run the command line tool)

If you want to use this feature, `--config` is mandatory and `--cfg-name` may be omitted. The config file can store one or multiple database connection sections under the `database` field. `database` is an array; previously it was an object only able to store a single connection section at a time. 

We strongly encourage you to adopt the new format as of `v0.18.0`. `--cfg-name` takes the name of the desired database section to connect with. It can be omitted and its default value will be the first item in the array. 

As of `v0.21.0`, SSL connection options are supported in the config file.

```sh
# default: test
$ dblab --config

$ dblab --config --cfg-name "prod"
```

#### Key bindings configuration

Key bindings can be configured through the `.dblab.yaml` file. There is a field called `keybindings` where key bindings can be modified. Under `keybindings`, an `editor` section configures the Vim-style query editor (movement between normal and insert mode, cursor motion in normal mode, and the editor’s execute-query shortcut). By default, the keybindings are not loaded, so you need to use the `--keybindings` or `-k` flag to load them. See the example to see the full list of the key bindings subject to change. The file shows the default values. The list of the available key bindings belongs to the [bubbletea](https://github.com/charmbracelet/bubbletea) library. Specifically, see the [KeyNames map](https://github.com/charmbracelet/bubbletea/blob/1ed724a2d1316ace504f87a2f0bbbcc189d280f6/key.go#L15) for an accurate reference.

**Deprecated:** the top-level `execute-query` field under `keybindings`. Use `execute-query` under `keybindings.editor` instead.

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
    # optional for postgres and oracle
    # if omitted, all accessible schemas are shown
    schema: "myschema"
    # optional: set to true to force a read-only session
    readonly: true
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
    db: "FREEPDB1"
    schema: "user1"
    password: "password"
    user: "user2"
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
  - name: "realistic-ssh-example"
    host: "rds-endpoint.region.rds.amazonaws.com"
    port: 5432
    db: "database_name"
    user: "db_user"
    password: "password"
    schema: "schema_name"
    driver: "postgres"
    ssl: "require"
    ssh-host: "bastion.host.ip"
    ssh-port: 22
    ssh-user: "ec2-user"
    ssh-key-file: "/path/to/ssh/key.pem"
    ssh-key-pass: "hiuwiewnc092"
# should be greater than 0, otherwise the app will error out
limit: 50
keybindings:
  next-tab: 'tab'
  prev-tab: 'shift+tab'
  page-top: 'g'
  page-bottom: 'G'
  end-of-line: '$'
  beginning-of-line: '0'
  navigation:
    up: 'ctrl+k'
    down: 'ctrl+j'
    left: 'ctrl+h'
    right: 'ctrl+l'
  editor:
    up: 'k'
    down: 'j'
    left: 'h'
    right: 'l'
    insert: 'i'
    normal: 'esc'
    execute-query: 'ctrl+e'
```

Or for SQLite:

```yaml
database:
  - name: "prod"
    db: "path/to/file.sqlite3"
    driver: "sqlite"
```

Only the `host`, `ssl`, and `schema` fields are optional. `host` defaults to `127.0.0.1`, `ssl` defaults to `disable`. The `schema` field is only applicable to PostgreSQL and Oracle; if omitted, all accessible schemas are shown.

### Connection Profiles

<img src="screenshots/dblab-connect.png" />

dblab supports saving and reusing database connection profiles. When you successfully connect to a database, you can store the connection parameters as a named profile using the `--save-as` flag. Both the database password and the SSH tunnel password (when using SSH connections) are stored securely in your operating system's keyring (e.g., GNOME Keyring, macOS Keychain, or Windows Credential Manager) rather than in plain text.

#### Saving a profile

Use the `--save-as` flag with any connection to save it as a named profile:

```sh
$ dblab --host localhost --user myuser --db users --pass password --ssl disable --port 5432 --driver postgres --limit 50 --save-as myprofile
```

The connection parameters are saved to `$XDG_CONFIG_HOME/dblab/dblab.json` (excluding passwords), while the database password and SSH password (if provided) are stored in the OS keyring.

#### Using saved profiles

Use the `connect` command to launch an interactive menu that lists all saved profiles:

```sh
$ dblab connect
```

This opens a TUI selector where you can:

- Browse saved database profiles
- Press <kbd>Enter</kbd> to connect to the selected profile
- Press <kbd>Ctrl+D</kbd> to delete a profile
- Press <kbd>Ctrl+C</kbd> to quit

The password is automatically retrieved from the OS keyring when connecting.

#### Profile storage format

Profiles are stored in `$XDG_CONFIG_HOME/dblab/dblab.json`:

```json
{
  "profiles": {
    "myprofile": {
      "host": "localhost",
      "port": "5432",
      "db": "users",
      "user": "postgres",
      "schema": "public",
      "driver": "postgres"
    }
  }
}
```

## Navigation

Key bindings are now configurable; see [Key bindings configuration](#key-bindings-configuration) to learn how to replace existing key bindings. It's worth noting that key bindings are only configurable through the configuration file; there are no flags to do so. If you don't replace them through the configuration file, the information below remains the same; otherwise, just replace the new key binding with the existing information for the default one.

### Query editor

The query editor uses **normal** and **insert** modes (similar to Vim). When you focus the query editor, it starts in **normal** mode. Press <kbd>i</kbd> to enter insert mode and type or edit SQL; press <kbd>Escape</kbd> to return to normal mode (the cursor moves one character to the left, as in Vim). In insert mode, use the arrow keys to move the cursor; in normal mode, use <kbd>h</kbd>, <kbd>j</kbd>, <kbd>k</kbd>, and <kbd>l</kbd> instead (configurable in `.dblab.yaml` with `--keybindings` or `-k`; see [Key bindings configuration](#key-bindings-configuration)). In normal mode, <kbd>dd</kbd> deletes the current line, <kbd>yy</kbd> yanks the current line into an internal register, <kbd>p</kbd> pastes that line after the current line, and <kbd>x</kbd> deletes the character under the cursor. <kbd>0</kbd> and <kbd>$</kbd> move to the beginning or end of the current line in the query buffer. <kbd>g</kbd> jumps to the first line and <kbd>G</kbd> jumps to the last line of the editor buffer. Press <kbd>Ctrl+D</kbd> to clear the entire editor content. Press <kbd>ctrl+e</kbd> to execute the query (this uses the `keybindings.editor.execute-query` binding); whitespace-only queries are ignored.

#### Multi-query execution

<img src="screenshots/dblab-multi-query.png" />

You can write multiple SQL statements in the editor separated by semicolons (`;`) and execute them all at once with <kbd>ctrl+e</kbd>. The queries are run concurrently and each result is displayed in its own tab (e.g., "query #1", "query #2", etc.). If a query fails, its tab will display the error message while other successful queries still show their results. A maximum of 5 queries can be executed per batch. Pressing <kbd>Ctrl+c</kbd> quits the application. If queries are currently running, they are cancelled before exiting.

#### Query History

<img src="screenshots/query-history.png" />

dblab automatically saves every executed query to a local history file (`$XDG_CONFIG_HOME/dblab/dblab.gob`). Press <kbd>F8</kbd> to open the query history view, which displays past queries sorted newest-first in a filterable list. Use the built-in search to narrow results, press <kbd>Enter</kbd> to load the selected query back into the editor, or press <kbd>Esc</kbd> to return without selecting anything.

Otherwise, you might be located at the tables panel, where you can navigate using the arrows <kbd>Up</kbd> and <kbd>Down</kbd> (or the keys <kbd>k</kbd> and <kbd>j</kbd> respectively). If you want to see the rows of a table, press <kbd>Enter</kbd>. To see the schema of a table, locate yourself on the `tables` panel and press <kbd>tab</kbd> to switch to the `columns` panel, then use <kbd>shift+tab</kbd> to switch back.

<img src="screenshots/rows-view.png" />
<img src="screenshots/structure-view.png" />
<img src="screenshots/indexes-view.png" />
<img src="screenshots/constraints-view.png" />

The navigation buttons were removed since they are too slow to navigate the content of a table effectively. The user is better off typing a `SELECT` statement with proper `OFFSET` and `LIMIT`.

The `--db` flag is mandatory. dblab connects to a single database and displays its catalog as a tree in the sidebar. For PostgreSQL and Oracle, the tree shows the database, its schemas, and the tables under each schema. For MySQL, SQLite, and SQL Server, the tree shows the database and its tables directly. If the `--schema` flag is provided for PostgreSQL or Oracle, only that schema is shown; otherwise, all accessible schemas are listed.

<img src="screenshots/tree-view.png" />

When navigating query result sets, the cell will be highlighted so the user can see which table cell is selected. This is important because you can press the `Enter` key on a cell of interest to copy its content.

### Key Bindings
| Key                                    | Description                           |
|----------------------------------------|----------------------------------------|
|<kbd>ctrl+e</kbd>                       | If the query editor is focused, execute the query (also works in insert and normal mode) |
|<kbd>i</kbd>                            | If the query editor is focused in normal mode, enter insert mode |
|<kbd>Escape</kbd>                       | If the query editor is focused in insert mode, return to normal mode |
|<kbd>dd</kbd>                           | If the query editor is focused in normal mode, delete the current line |
|<kbd>yy</kbd>                           | If the query editor is focused in normal mode, yank the current line |
|<kbd>p</kbd>                            | If the query editor is focused in normal mode, paste the yanked or deleted line after the current line |
|<kbd>x</kbd>                            | If the query editor is focused in normal mode, delete the character under the cursor |
|<kbd>Enter</kbd>                        | If the tables panel is focused, list all rows as a result set on the rows panel and display the structure of the table on the structure panel |
|<kbd>tab</kbd>                          | If the result set panel is focused, press tab to navigate to the next metadata tab |
|<kbd>shift+tab</kbd>                    | If the result set panel is focused, press shift+tab to navigate to the previous metadata tab |
|<kbd>Ctrl+H</kbd>                       | Toggle to the panel on the left |
|<kbd>Ctrl+J</kbd>                       | Toggle to the panel below |
|<kbd>Ctrl+K</kbd>                       | Toggle to the panel above |
|<kbd>Ctrl+L</kbd>                       | Toggle to the panel on the right |
|<kbd>Arrow Up</kbd>                     | If the query editor is focused in insert mode, move the cursor up. If the results panel is focused, navigate the table upward (all tabs on the results panel). |
|<kbd>k</kbd>                            | If the query editor is focused in normal mode, move the cursor up. If the results panel is focused, navigate the table upward (all tabs on the results panel). |
|<kbd>Arrow Down</kbd>                   | If the query editor is focused in insert mode, move the cursor down. If the results panel is focused, navigate the table downward (all tabs on the results panel). |
|<kbd>j</kbd>                            | If the query editor is focused in normal mode, move the cursor down. If the results panel is focused, navigate the table downward (all tabs on the results panel). |
|<kbd>Arrow Right</kbd>                  | If the query editor is focused in insert mode, move the cursor right. If the results panel is focused, navigate the table to the right (all tabs on the results panel). |
|<kbd>l</kbd>                            | If the query editor is focused in normal mode, move the cursor right. If the results panel is focused, navigate the table to the right (all tabs on the results panel). |
|<kbd>Arrow Left</kbd>                   | If the query editor is focused in insert mode, move the cursor left. If the results panel is focused, navigate the table to the left (all tabs on the results panel). |
|<kbd>h</kbd>                            | If the query editor is focused in normal mode, move the cursor left. If the results panel is focused, navigate the table to the left (all tabs on the results panel). |
|<kbd>g</kbd>                            | If the query editor is focused in normal mode, jump to the first line of the buffer. If the results panel is focused, move to the top of the dataset (all tabs on the results panel). |
|<kbd>G</kbd>                            | If the query editor is focused in normal mode, jump to the last line of the buffer. If the results panel is focused, move to the bottom of the dataset (all tabs on the results panel). |
|<kbd>0</kbd>                            | If the query editor is focused in normal mode, move to the start of the current line. If the results panel is focused, move to the left edge of the row (all tabs on the results panel). |
|<kbd>$</kbd>                            | If the query editor is focused in normal mode, move to the end of the current line. If the results panel is focused, move to the right edge of the row (all tabs on the results panel). |
|<kbd>Ctrl+D</kbd>                       | If the query editor is focused in normal mode, clear the entire editor content |
|<kbd>F8</kbd>                           | Open the query history view |
|<kbd>Ctrl+c</kbd>                       | Quit the application (cancels in-flight queries if any) |

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
