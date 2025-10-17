### Using flags

You can start the app without passing flags or parameters, so then an interactive command prompt will ask for the connection details.  

![Alt Text](https://raw.githubusercontent.com/danvergara/dblab/main/screenshots/dblab-demo.gif){ width="500" : .center }

Otherwise, you can explicitly include the connection details using multiple parameters:

```{ .sh .copy }
dblab --host localhost --user myuser --db users --pass password --ssl disable --port 5432 --driver postgres --limit 50
```

```{ .sh .copy }
dblab --db path/to/file.sqlite3 --driver sqlite
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

if you're using PostgreSQL or Oracle, you have the option to define the schema you want to work with, the default value is `public` for Postgres, empty for Oracle.

```{ .sh .copy }
dblab --host localhost --user myuser --db users --pass password --schema myschema --ssl disable --port 5432 --driver postgres --limit 50
```
```{ .sh .copy }
dblab --url postgres://user:password@host:port/database?sslmode=[mode] --schema myschema
```
```{ .sh .copy }
dblab --host localhost --user user2 --db FREEPDB1 --pass password --port 1521 --driver oracle --limit 50 --schema user1
```
```{ .sh .copy }
dblab --url 'oracle://user2:password@localhost:1521/FREEPDB1' --schema user1
```

As a request made in [#125](https://github.com/danvergara/dblab/issues/125), support for MySQL/MariaDB sockets was integrated.

```{ .sh .copy }
dblab --url "mysql://user:password@unix(/path/to/socket/mysql.sock)/dbname?charset=utf8"
```
```{ .sh .copy }
dblab --socket /path/to/socket/mysql.sock --user user --db dbname --pass password --ssl disable --port 5432 --driver mysql --limit 50
```

For more information about the available flags check the [Usage section](https://dblab.danvergara.com/usage/#usage).

### Using a config file

default: the first configuration after the `database` field.
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

Only the `host` and `ssl` fields are optionals. `127.0.0.1` and `disable`, respectively.
