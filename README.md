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
    - [Binary Release](#binary-release-linuxmacoswindows)
    - [Automated installation/update](#automated-installationupdate)
- [Help Command](#help)
- [Usage](#usage)
    - [SSH Tunnel](#ssh-tunnel)
    - [Configuration](#configuration)
        - [Key bindings configuration](#key-bindings-configuration)
    - [Connection Profiles](#connection-profiles)
- [Navigation](#navigation)
    - [Panels and the sidebar tree](#panels-and-the-sidebar-tree)
    - [Result sets](#result-sets)
- [Query editor](#query-editor)
    - [Modes](#modes)
    - [Editing and motions](#editing-and-motions)
    - [Executing queries](#executing-queries)
    - [Query history](#query-history)
- [Help modal](#help-modal)
- [Key bindings](#key-bindings)
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
- Single-query execution: press <kbd>ctrl+r</kbd> to execute only the query on the current cursor line, without running other statements in the editor.
- Connection profiles with secure credential storage in the OS keyring.
- Query history: executed queries are persisted across sessions and can be browsed/re-used via a filterable list.
- Read-only mode: use `--readonly` to prevent accidental writes by forcing the database session into read-only mode (supported for PostgreSQL, MySQL, SQLite, Oracle, and SQL Server).
- Built-in help modal: press <kbd>?</kbd> to display a help overlay showing all available key bindings; press <kbd>Esc</kbd> to dismiss it.
- Per-panel key bindings: the query editor, the sidebar tree, the result set panel and the panel navigation each have their own section in `.dblab.yaml`, so every panel can be rebound independently.

## Installation

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

### Automated installation/update
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

Key bindings can be configured through the `.dblab.yaml` file. There is a field called `keybindings` where key bindings can be modified. By default, the keybindings are not loaded, so you need to use the `--keybindings` or `-k` flag to load them.

Bindings are grouped by the part of the UI they belong to, so every panel can be rebound independently of the others:

| Section | What it controls |
|---------|------------------|
| `keybindings` (top level) | `help`, `quit` and `history`, which are global |
| `keybindings.navigation` | moving focus between the three panels |
| `keybindings.editor` | the Vim-style query editor: cursor motion in normal mode, mode switching, and query execution |
| `keybindings.sidebar` | jumping to the top / bottom of the sidebar tree |
| `keybindings.resultset` | tab switching and horizontal motion in the result set panel |

Every field has a default value, so you only need to list the ones you want to change; anything you leave out falls back to the default shown in the [example below](#dblabyaml-example). The list of the available key bindings belongs to the [bubbletea](https://github.com/charmbracelet/bubbletea) library. Specifically, see the [KeyNames map](https://github.com/charmbracelet/bubbletea/blob/1ed724a2d1316ace504f87a2f0bbbcc189d280f6/key.go#L15) for an accurate reference.

##### Migrating from the flat layout

Key bindings used to be partly flat: `next-tab`, `prev-tab`, `page-top`, `page-bottom`, `end-of-line` and `beginning-of-line` sat at the top level of `keybindings` and were shared by more than one panel. They now live under the panel that uses them:

| Old (top level) | New |
|-----------------|-----|
| `next-tab` | `resultset.next-tab` |
| `prev-tab` | `resultset.prev-tab` |
| `beginning-of-line` | `resultset.line-start` and `editor.line-start` |
| `end-of-line` | `resultset.line-end` and `editor.line-end` |
| `page-top` | `sidebar.go-top` (plus `editor.go-top` and `resultset.go-top`) |
| `page-bottom` | `sidebar.go-bottom` (plus `editor.go-bottom` and `resultset.go-bottom`) |

The old top-level fields are no longer read: if your config still sets them, those bindings silently fall back to their defaults. Two things to be aware of while migrating:

- The sidebar's jump-to-top / jump-to-bottom defaults changed from <kbd>g</kbd> / <kbd>G</kbd> to <kbd>alt+k</kbd> / <kbd>alt+j</kbd>, which leaves <kbd>g</kbd> / <kbd>G</kbd> free for the editor and the result set.
- The editor gained `line-start`, `line-end`, `go-top` and `go-bottom`. These motions already worked, but were hardcoded to <kbd>0</kbd>, <kbd>$</kbd>, <kbd>g</kbd> and <kbd>G</kbd>; they are now configurable.

The previously deprecated top-level `execute-query` field is gone as well — use `execute-query` under `keybindings.editor`.

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
  help: '?'
  quit: 'ctrl+c'
  history: 'alt+h'
  # moving focus between the three panels
  navigation:
    up: 'ctrl+k'
    down: 'ctrl+j'
    left: 'ctrl+h'
    right: 'ctrl+l'
  # the query editor, in normal mode
  editor:
    up: 'k'
    down: 'j'
    left: 'h'
    right: 'l'
    line-start: '0'
    line-end: '$'
    go-top: 'g'
    go-bottom: 'G'
    insert: 'i'
    normal: 'esc'
    execute-query: 'ctrl+e'
    execute-single-query: 'ctrl+r'
  # the database tree on the left
  sidebar:
    go-top: 'alt+k'
    go-bottom: 'alt+j'
  # the result set panel and its metadata tabs
  resultset:
    next-tab: 'tab'
    prev-tab: 'shift+tab'
    line-start: '0'
    line-end: '$'
    go-top: 'g'
    go-bottom: 'G'
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

The UI is split into three panels: the **sidebar tree** on the left, the **query editor** on the top right, and the **result set** panel below it. Move focus between them with <kbd>Ctrl+H</kbd>, <kbd>Ctrl+J</kbd>, <kbd>Ctrl+K</kbd> and <kbd>Ctrl+L</kbd>.

Every key binding in this README is the default. All of them can be replaced through the `.dblab.yaml` configuration file — there are no flags for it — so if you've customized a binding, substitute yours for the default shown here. See [Key bindings configuration](#key-bindings-configuration).

### Panels and the sidebar tree

dblab connects to a single database (the `--db` flag is mandatory) and displays its catalog as a tree in the sidebar. For PostgreSQL and Oracle, the tree shows the database, its schemas, and the tables under each schema. For MySQL, SQLite, and SQL Server, the tree shows the database and its tables directly. If the `--schema` flag is provided for PostgreSQL or Oracle, only that schema is shown; otherwise, all accessible schemas are listed.

<img src="screenshots/tree-view.png" />

Navigate the tree with <kbd>Up</kbd> and <kbd>Down</kbd> (or <kbd>k</kbd> and <kbd>j</kbd>), and press <kbd>Enter</kbd> on a table to load its rows into the result set panel. Jump straight to the first or last visible node with <kbd>alt+k</kbd> and <kbd>alt+j</kbd> (`keybindings.sidebar.go-top` and `keybindings.sidebar.go-bottom`), and scroll the tree sideways with <kbd>h</kbd> and <kbd>l</kbd> when a name is wider than the panel.

Press <kbd>/</kbd> to search the tree by name and <kbd>Esc</kbd> to leave the search. While a search is active every character you type — including <kbd>h</kbd> and <kbd>l</kbd> — goes to the search box instead of scrolling the tree.

### Result sets

Selecting a table populates the result set panel, which has one tab per view of the table. Press <kbd>tab</kbd> and <kbd>shift+tab</kbd> to move between them:

- **Data** — the rows of the table, or the result of the query you executed
- **Columns** — the schema of the table
- **Indexes** — the indexes on the table
- **Constraints** — the constraints on the table

<img src="screenshots/rows-view.png" />
<img src="screenshots/structure-view.png" />
<img src="screenshots/indexes-view.png" />
<img src="screenshots/constraints-view.png" />

Move around a result set with the arrow keys or <kbd>h</kbd>/<kbd>j</kbd>/<kbd>k</kbd>/<kbd>l</kbd>. The selected cell is highlighted so you can see where you are; press <kbd>Enter</kbd> on a cell to copy its content.

There are no pagination controls — they proved too slow to page through a table effectively. To work through a large table, write a `SELECT` with explicit `OFFSET` and `LIMIT` instead.

## Query editor

### Modes

The query editor uses **normal** and **insert** modes, similar to Vim. When you focus the editor it starts in **normal** mode. Press <kbd>i</kbd> to enter insert mode and type or edit SQL; press <kbd>Escape</kbd> to return to normal mode (the cursor moves one character to the left, as in Vim).

Cursor movement depends on the mode: in insert mode use the arrow keys, in normal mode use <kbd>h</kbd>, <kbd>j</kbd>, <kbd>k</kbd> and <kbd>l</kbd>.

### Editing and motions

In normal mode:

- <kbd>dd</kbd> deletes the current line, <kbd>yy</kbd> yanks it into an internal register, and <kbd>p</kbd> pastes the yanked or deleted line after the current line
- <kbd>x</kbd> deletes the character under the cursor
- <kbd>0</kbd> and <kbd>$</kbd> move to the beginning and end of the current line (`keybindings.editor.line-start` and `keybindings.editor.line-end`)
- <kbd>g</kbd> and <kbd>G</kbd> jump to the first and last line of the buffer (`keybindings.editor.go-top` and `keybindings.editor.go-bottom`)
- <kbd>Ctrl+D</kbd> clears the entire editor content

The cursor motions, the mode switches and the execute shortcuts are all configurable under `keybindings.editor`; the line-oriented commands (<kbd>dd</kbd>, <kbd>yy</kbd>, <kbd>p</kbd>, <kbd>x</kbd>) and <kbd>Ctrl+D</kbd> are fixed.

### Executing queries

Press <kbd>ctrl+e</kbd> to execute the contents of the editor (`keybindings.editor.execute-query`). Whitespace-only queries are ignored.

Press <kbd>ctrl+r</kbd> to execute only the query on the current cursor line (`keybindings.editor.execute-single-query`), leaving the other statements in the editor untouched. Both bindings work from either mode.

#### Multiple statements

<img src="screenshots/dblab-multi-query.png" />

You can write multiple SQL statements separated by semicolons (`;`) and execute them all at once with <kbd>ctrl+e</kbd>:

```sql
SELECT * FROM users; SELECT * FROM orders; SELECT count(*) FROM products;
```

The statements run concurrently and each result is displayed in its own tab ("query #1", "query #2", and so on) — three tabs, for the example above. If a statement fails, its tab shows the error message while the successful ones still show their results. A maximum of 5 statements can be executed per batch.

While a batch is running, press <kbd>Ctrl+c</kbd> to cancel it; press <kbd>Ctrl+c</kbd> again to quit dblab.

### Query history

<img src="screenshots/query-history.png" />

dblab automatically saves every executed query to a local history file (`$XDG_CONFIG_HOME/dblab/dblab.gob`). Press <kbd>alt+h</kbd> (see [Key bindings configuration](#key-bindings-configuration) to configure it) to open the query history view, which displays past queries sorted newest-first in a filterable list. Use the built-in search to narrow results, press <kbd>Enter</kbd> to load the selected query back into the editor, or press <kbd>Esc</kbd> to return without selecting anything.


## Help modal

<img src="screenshots/dblab-help-modal.png" />

Press <kbd>?</kbd> at any time to open the help modal, which displays all available key bindings in a centered overlay. Press <kbd>Esc</kbd> to dismiss it; focus returns to the query editor.

## Key bindings

These are the defaults; see [Key bindings configuration](#key-bindings-configuration) to change them. The **Config field** column gives the `.dblab.yaml` key for the bindings that can be customized — the rest are fixed.

### Panel navigation

| Key | Description | Config field |
|-----|-------------|--------------|
| <kbd>Ctrl+H</kbd> | Toggle to the panel on the left | `navigation.left` |
| <kbd>Ctrl+J</kbd> | Toggle to the panel below | `navigation.down` |
| <kbd>Ctrl+K</kbd> | Toggle to the panel above | `navigation.up` |
| <kbd>Ctrl+L</kbd> | Toggle to the panel on the right | `navigation.right` |

### Query editor (both modes)

| Key | Description | Config field |
|-----|-------------|--------------|
| <kbd>ctrl+e</kbd> | Execute the contents of the editor | `editor.execute-query` |
| <kbd>ctrl+r</kbd> | Execute only the query on the current cursor line | `editor.execute-single-query` |

### Query editor (normal mode)

| Key | Description | Config field |
|-----|-------------|--------------|
| <kbd>i</kbd> | Enter insert mode | `editor.insert` |
| <kbd>h</kbd> <kbd>j</kbd> <kbd>k</kbd> <kbd>l</kbd> | Move the cursor left, down, up, right | `editor.left` / `editor.down` / `editor.up` / `editor.right` |
| <kbd>dd</kbd> | Delete the current line | — |
| <kbd>yy</kbd> | Yank the current line | — |
| <kbd>p</kbd> | Paste the yanked or deleted line after the current line | — |
| <kbd>x</kbd> | Delete the character under the cursor | — |
| <kbd>0</kbd> / <kbd>$</kbd> | Move to the start / end of the current line | `editor.line-start` / `editor.line-end` |
| <kbd>g</kbd> / <kbd>G</kbd> | Jump to the first / last line of the buffer | `editor.go-top` / `editor.go-bottom` |
| <kbd>Ctrl+D</kbd> | Clear the entire editor content | — |

### Query editor (insert mode)

| Key | Description | Config field |
|-----|-------------|--------------|
| <kbd>Escape</kbd> | Return to normal mode | `editor.normal` |
| <kbd>Arrow keys</kbd> | Move the cursor | — |

### Sidebar tree

| Key | Description | Config field |
|-----|-------------|--------------|
| <kbd>Arrow Up</kbd> / <kbd>k</kbd> | Move up the tree | — |
| <kbd>Arrow Down</kbd> / <kbd>j</kbd> | Move down the tree | — |
| <kbd>alt+k</kbd> | Jump to the first visible node | `sidebar.go-top` |
| <kbd>alt+j</kbd> | Jump to the last visible node | `sidebar.go-bottom` |
| <kbd>h</kbd> / <kbd>l</kbd> | Scroll the tree left / right | — |
| <kbd>/</kbd> | Search the tree by name | — |
| <kbd>Esc</kbd> | Leave the search | — |
| <kbd>Enter</kbd> | List all rows of the selected table and display its structure | — |

### Result set panel

Applies to all tabs of the result set panel.

| Key | Description | Config field |
|-----|-------------|--------------|
| <kbd>tab</kbd> / <kbd>shift+tab</kbd> | Navigate to the next / previous metadata tab | `resultset.next-tab` / `resultset.prev-tab` |
| <kbd>Arrow Up</kbd> / <kbd>k</kbd> | Navigate the table upward | — |
| <kbd>Arrow Down</kbd> / <kbd>j</kbd> | Navigate the table downward | — |
| <kbd>Arrow Left</kbd> / <kbd>h</kbd> | Navigate the table to the left | — |
| <kbd>Arrow Right</kbd> / <kbd>l</kbd> | Navigate the table to the right | — |
| <kbd>g</kbd> / <kbd>G</kbd> | Move to the top / bottom of the dataset | — |
| <kbd>0</kbd> / <kbd>$</kbd> | Move to the left / right edge of the row | `resultset.line-start` / `resultset.line-end` |
| <kbd>Enter</kbd> | Copy the content of the selected cell | — |


### Global

| Key | Description | Config field |
|-----|-------------|--------------|
| <kbd>Alt+h</kbd>  | Open the query history view | `history` |
| <kbd>?</kbd> | Open the help modal showing all key bindings | `help` |
| <kbd>Esc</kbd> | Dismiss the help modal (or return to normal mode in the query editor) | — |
| <kbd>Ctrl+c</kbd> | Cancel running queries if any; otherwise quit the application | `quit` |

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
