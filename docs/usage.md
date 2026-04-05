You can get started by using connection flags or by using a configuration file with the connection parameters.

```sh
$ dblab [flags] 

```

or

```sh
$ dblab [command]

```

### Available Commands                               

|         `help`         |    Help about any command    | 
|:----------------------:|:----------------------------:|
|       `version`        |  The version of the project  |

### Flags

```
dblab is a terminal UI-based interactive database client for Postgres, MySQL, and SQLite3.

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
      --encrypt string                    [strict|disable|false|true] whether data sent between client and server is encrypted
  -h, --help                              help for dblab
      --host string                       Server host name or IP
      --limit uint                        Size of the result set for the table content query (should be greater than zero, otherwise the app will error out) (default 100)
      --pass string                       Password for user
      --port string                       Server port
      --schema string                     Database schema (postgres and oracle only)
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
      --wallet string                     Path for auto-login oracle wallet

Use "dblab [command] --help" for more information about a command.
```

## Navigation

If the query panel is active, type the desired query and press <kbd>ctrl+e</kbd> to see the results on the rows panel below.

Otherwise, you might be located at the tables panel, where you can navigate using the arrows <kbd>Up</kbd> and <kbd>Down</kbd> (or the keys <kbd>k</kbd> and <kbd>j</kbd> respectively). If you want to see the rows of a table, press <kbd>Enter</kbd>. To see the schema of a table, locate yourself on the `tables` panel and press <kbd>tab</kbd> to switch to the `columns` panel, then use <kbd>shift+tab</kbd> to switch back.

The `--db` flag is now optional (except for Oracle), meaning that the user will be able to see the list of databases they have access to. The regular list of tables will be replaced with a tree structure showing a list of databases and their respective list of tables, branching off each database. Due to the nature of the vast majority of DBMSs that don't allow cross-database queries, dblab has to open an independent connection for each database. The side effect of this decision is that the user has to press `Enter` on the specific database of interest. An indicator showing the current active database will appear at the bottom-right of the screen. To change the focus, just hit enter on another database. Once a database is selected, the usual behavior of inspecting tables remains the same.

<img src="https://raw.githubusercontent.com/danvergara/dblab/main/screenshots/tree-view.png" />

Now, there's a menu to navigate between hidden views by just clicking on the desired options:

![Alt Text](https://raw.githubusercontent.com/danvergara/dblab/main/screenshots/rows-view.png){ width="700" : .center }
![Alt Text](https://raw.githubusercontent.com/danvergara/dblab/main/screenshots/structure-view.png){ width="700" : .center }
![Alt Text](https://raw.githubusercontent.com/danvergara/dblab/main/screenshots/constraints-view.png){ width="700" : .center }
![Alt Text](https://raw.githubusercontent.com/danvergara/dblab/main/screenshots/indexes-view.png){ width="700" : .center }

When navigating query result sets, the cell will be highlighted so the user can see which table cell is selected. This is important because you can press the `Enter` key on a cell of interest to copy its content.

### Key Bindings
| Key                                    | Description                           |
|----------------------------------------|----------------------------------------|
|<kbd>ctrl+e</kbd>                       | If the query editor is active, execute the query |
|<kbd>Ctrl+D</kbd>                       | Clears all text from the query editor when it is selected |
|<kbd>Enter</kbd>                        | If the tables panel is active, list all rows as a result set on the rows panel and display the structure of the table on the structure panel |
|<kbd>tab</kbd>                          | If the result set panel is active, press tab to navigate to the next metadata tab |
|<kbd>shift+tab</kbd>                    | If the result set panel is active, press shift+tab to navigate to the previous metadata tab |
|<kbd>Ctrl+H</kbd>                       | Toggle to the panel on the left |
|<kbd>Ctrl+J</kbd>                       | Toggle to the panel below |
|<kbd>Ctrl+K</kbd>                       | Toggle to the panel above |
|<kbd>Ctrl+L</kbd>                       | Toggle to the panel on the right |
|<kbd>Arrow Up</kbd>                     | Vertical scrolling on the panel. Views: rows, table, constraints, structure, and indexes |
|<kbd>k</kbd>                            | Vertical scrolling on the panel. Views: rows, table, constraints, structure, and indexes |
|<kbd>Arrow Down</kbd>                   | Vertical scrolling on the panel. Views: rows, table, constraints, structure, and indexes |
|<kbd>j</kbd>                            | Vertical scrolling on the panel. Views: rows, table, constraints, structure, and indexes |
|<kbd>Arrow Right</kbd>                  | Horizontal scrolling on the panel. Views: rows, constraints, structure, and indexes |
|<kbd>l</kbd>                            | Horizontal scrolling on the panel. Views: rows, constraints, structure, and indexes |
|<kbd>Arrow Left</kbd>                   | Horizontal scrolling on the panel. Views: rows, constraints, structure, and indexes |
|<kbd>h</kbd>                            | Horizontal scrolling on the panel. Views: rows, constraints, structure, and indexes |
|<kbd>g</kbd>                            | Move cursor to the top of the panel's dataset. Views: rows, constraints, structure, and indexes |
|<kbd>G</kbd>                            | Move cursor to the bottom of the panel's dataset. Views: rows, constraints, structure, and indexes |
|<kbd>Ctrl+c</kbd>                       | Quit |


## Connection Examples

You can start the app without passing flags or parameters; an interactive command prompt will ask for the connection details.  

![Alt Text](https://raw.githubusercontent.com/danvergara/dblab/main/screenshots/dblab-demo.gif){ width="500" : .center }

Otherwise, you can explicitly include the connection details using multiple parameters:

```{ .sh .copy }
dblab --host localhost --user myuser --db users --pass password --ssl disable --port 5432 --driver postgres --limit 50
```
```{ .sh .copy }
dblab --db path/to/file.sqlite3 --driver sqlite
```
```{ .sh .copy }
dblab --host localhost --user system --db FREEPDB1 --pass password --port 1521 --driver oracle --limit 50
```
```{ .sh .copy }
dblab --host localhost --user SA --db msdb --pass '5@klkbN#ABC' --port 1433 --driver sqlserver --limit 50
```

Connection URL scheme is also supported:

```{ .sh .copy }
dblab --url postgres://user:password@host:port/database?sslmode=[mode]
```
```{ .sh .copy }
dblab --url mysql://user:password@tcp(host:port)/db
```
```{ .sh .copy }
dblab --url file:test.db?cache=shared&mode=memory
```
```{ .sh .copy }
dblab --url 'oracle://user:password@localhost:1521/db'
```
```{ .sh .copy }
dblab --url 'sqlserver://SA:myStrong(!)Password@localhost:1433?database=tempdb&encrypt=true&trustservercertificate=false&connection+timeout=30'
```

If you're using PostgreSQL or Oracle, you have the option to define the schema you want to work with; the default value is `public` for Postgres, and empty for Oracle.

**Postgres**

```{ .sh .copy }
dblab --host localhost --user myuser --db users --pass password --schema myschema --ssl disable --port 5432 --driver postgres --limit 50
```
```{ .sh .copy }
dblab --url postgres://user:password@host:port/database?sslmode=[mode] --schema myschema
```

**Oracle**

```{ .sh .copy }
dblab --host localhost --user user2 --db FREEPDB1 --pass password --port 1521 --driver oracle --limit 50 --schema user1
```
```{ .sh .copy }
dblab --url 'oracle://user2:password@localhost:1521/FREEPDB1' --schema user1
```

As requested in [#125](https://github.com/danvergara/dblab/issues/125), support for MySQL/MariaDB sockets was integrated.

```{ .sh .copy }
dblab --url "mysql://user:password@unix(/path/to/socket/mysql.sock)/dbname?charset=utf8"
```
```{ .sh .copy }
dblab --socket /path/to/socket/mysql.sock --user user --db dbname --pass password --ssl disable --port 5432 --driver mysql --limit 50
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

Entering parameters and flags every time you use the tool is tedious, 
so `dblab` provides a couple of flags to help with it: `--config` and `--cfg-name`.

`dblab` is going to look for a file called `.dblab.yaml`. Currently, there are three places where you can drop a config file:

- $XDG_CONFIG_HOME ($XDG_CONFIG_HOME/.dblab.yaml)
- $HOME ($HOME/.dblab.yaml)
- . (the current directory where you run the command line tool)

If you want to use this feature, `--config` is mandatory and `--cfg-name` may be omitted. The config file can store one or multiple database connection sections under the `database` field. `database` is an array; previously it was an object only able to store a single connection section at a time. 

We strongly encourage you to adopt the new format as of `v0.18.0`. `--cfg-name` takes the name of the desired database section to connect with. It can be omitted and its default value will be the first item in the array.

As of `v0.21.0`, SSL connection options are supported in the config file.

```{ .sh .copy } 
dblab --config
```
```{ .sh .copy } 
dblab --config --cfg-name "prod"
```

#### Key bindings configuration

Key bindings can be configured through the `.dblab.yaml` file. There is a field called `keybindings` where key bindings can be modified. By default, the keybindings are not loaded, so you need to use the `--keybindings` or `-k` flag to load them. See the example to see the full list of the key bindings subject to change. The file shows the default values. The list of the available key bindings belongs to the [tcell](https://github.com/gdamore/tcell) library. Specifically, see the [KeyNames map](https://github.com/gdamore/tcell/blob/781586687ddb57c9d44727dc9320340c4d049b11/key.go#L83) for an accurate reference.


#### .dblab.yaml example

```{ .yaml .copy } 
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
  execute-query: 'ctrl+e'
  next-tab: 'tab'
  prev-tab: 'shift+tab'
  page-top: 'g'
  page-bottom: 'G'
  navigation:
    up: 'ctrl+k'
    down: 'ctrl+j'
    left: 'ctrl+h'
    right: 'ctrl+l'
```

Or for SQLite:

```{ .yaml .copy } 
database:
  - name: "prod"
    db: "path/to/file.sqlite3"
    driver: "sqlite"
```

Only the `host` and `ssl` fields are optional. They default to `127.0.0.1` and `disable`, respectively.
