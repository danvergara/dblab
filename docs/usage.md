
You can get started by using the connection flags or by using a configuration file with the connection params.

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

## Navigation

If the query panel is active, type the desired query and press <kbd>Ctrl+Space</kbd> to see the results on the rows panel below.
Otherwise, you might me located at the tables panel, then you can navigate by using the arrows <kbd>Up</kbd> and <kbd>Down</kbd> (or the keys <kbd>k</kbd> and <kbd>j</kbd> respectively).  

If you want to see the rows of a table, press <kbd>Enter</kbd>.  

To see the schema of a table, locate yourself on the `rows` panel and press <kbd>Ctrl+S</kbd> to switch to the `structure` panel, then switch <kbd>Ctrl+S</kbd> to switch back.  

The same can be achieved for the `constraints` view by pressing <kbd>Ctrl+F</kbd> to go back and forth between the `rows` and the `constraints` panels.

Now, there's a menu to navigate between hidden views by just clicking on the desired options:

![Alt Text](https://raw.githubusercontent.com/danvergara/dblab/main/screenshots/rows-view.png){ width="700" : .center }
![Alt Text](https://raw.githubusercontent.com/danvergara/dblab/main/screenshots/structure-view.png){ width="700" : .center }
![Alt Text](https://raw.githubusercontent.com/danvergara/dblab/main/screenshots/constraints-view.png){ width="700" : .center }
![Alt Text](https://raw.githubusercontent.com/danvergara/dblab/main/screenshots/indexes-view.png){ width="700" : .center }

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
<kbd>Ctrl+c</kbd>                       | Quit


## Connection Examples

You can start the app without passing flags or parameters, so then an interactive command prompt will ask for the connection details.  

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

if you're using PostgreSQL, you have the option to define the schema you want to work with, the default value is `public`.

```{ .sh .copy }
dblab --host localhost --user myuser --db users --pass password --schema myschema --ssl disable --port 5432 --driver postgres --limit 50
```
```{ .sh .copy }
dblab --url postgres://user:password@host:port/database?sslmode=[mode] --schema myschema
```

As a request made in [#125](https://github.com/danvergara/dblab/issues/125), support for MySQL/MariaDB sockets was integrated.

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

Entering the parameters and flags every time you want to use it is tedious, 
so `dblab` provides a couple of flags to help with it: `--config` and `--cfg-name`.

`dblab` is going to look for a file called `.dblab.yaml`. Currently, there are three places where you can drop a config file:

- $XDG_CONFIG_HOME ($XDG_CONFIG_HOME/.dblab.yaml)
- $HOME ($HOME/.dblab.yaml)
- . (the current directory where you run the command line tool)

If you want to use this feature, `--config` is mandatory and `--cfg-name` may be omitted. The config file can store one or multiple database connection sections under the `database` field. `database` is an array, previously was an object only able to store a single connection section at a time. 

We strongly encourage you to adopt the new format as of `v0.18.0`. `--cfg-name` takes the name of the desired database section to connect with. It can be omitted and its default values will be the first item on the array.

As of `v0.21.0`, ssl connections options are supported in the config file.

```{ .sh .copy } 
dbladb --config
```
```{ .sh .copy } 
dblab --config --cfg-name "prod"
```

`.dblab.yaml` example:

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
```

Or for sqlite:

```{ .yaml .copy } 
database:
  - name: "prod"
    db: "path/to/file.sqlite3"
    driver: "sqlite"
```

only the `host` and `ssl` fields are optionals. `127.0.0.1` and `disable`, respectively.
