
You can start the app passing no flags or parameters, you'll be asked for connection data instead.
![Alt Text](screenshots/dblab-default-form.gif)


```sh
$ dblab --host localhost --user myuser --db users --pass password --ssl disable --port 5432 --driver postgres --limit 50
$ dblab --db path/to/file.sqlite3 --driver sqlite
```

Connection URL scheme is also supported:

```sh
$ dblab --url postgres://user:password@host:port/database?sslmode=[mode]
$ dblab --url mysql://user:password@tcp(host:port)/db
$ dblab --url file:test.db?cache=shared&mode=memory
```

if you're using PostgreSQL, you have the option to define the schema you want to work with, the default value is `public`.

```sh
$ dblab --host localhost --user myuser --db users --pass password --schema myschema --ssl disable --port 5432 --driver postgres --limit 50
$ dblab --url postgres://user:password@host:port/database?sslmode=[mode] --schema myschema
```

As a request made in [#125](https://github.com/danvergara/dblab/issues/125), support for MySQL/MariaDB sockets was integrated.

```sh
$ dblab --url "mysql://user:pasword@unix(/path/to/socket/mysql.sock)/dbname?charset=utf8"
$ dblab --socket /path/to/socket/mysql.sock --user user --db dbname --pass password --ssl disable --port 5432 --driver mysql --limit 50
```

### Configuration

Enter previous flags every time is tedious, so `dblab` provides a couple of flags to help with it: `--config` and `--cfg-name`.

`dblab` is going to look for a file called `.dblab.yaml`. For now, the only two places where you can drop a config file are $HOME ($HOME/.dblab.yaml) and the current directory where you run the command line tool.
If you want to use this feature, `--config` is mandatory and `--cfg-name` may be omitted. The config file can store one or multiple database connection sections under the `database` field. `database` is an array, previously was an object only able to store a single connection section at a time. We strongly encourgae you to adopt the new format as of `v0.18.0`. `--cfg-name` takes the name of the desired database section to connect with. It can be omitted and its default values will be the first item on the array.

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
<kbd>Ctrl+c</kbd>                       | Quit

Usage:
  dblab [flags]
  dblab [command]

Available Commands:
  help        Help about any command
  version     The version of the project

Flags:
      --cfg-name string   Database config name section
      --config          get the connection data from a config file (default is $HOME/.dblab.yaml or the current directory)
      --db string       Database name
      --driver string   Database driver
  -h, --help            help for dblab
      --host string     Server host name or IP
      --limit int       Size of the result set from the table content query (should be greater than zero, otherwise the app will error out) (default 100)
      --pass string     Password for user
      --port string     Server port
      --schema string   Database schema (postgres only)
      --socket string     Path to a Unix socket file
      --ssl string      SSL mode
  -u, --url string      Database connection string
      --user string     Database user

Use "dblab [command] --help" for more information about a command.