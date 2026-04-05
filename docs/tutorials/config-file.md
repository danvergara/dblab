
We know that writing parameters and flags every time you want to connect to a database can be very tedious, 
so `dblab` provides the option to create a YAML file with the configuration to connect to the database.
You can define multiple database configurations under the `database` field; just be sure to use different names for each of them.

For this example, we'll be using a PostgreSQL database, so the driver to use would be `postgres`. Remember that every driver has different args that we can include in the config YAML.

### Single database

In order to connect to a local database hosted at `0.0.0.0:5432`, we can just copy and paste the following configuration to a file named `.dblab.yaml` stored either in the root of your current directory, in your $HOME path ($HOME/.dblab.yaml), or in your $XDG_CONFIG_HOME path ($XDG_CONFIG_HOME/.dblab.yaml).

```{ .yaml .copy }

database:
  - name: "local"
    host: "0.0.0.0"
    port: 5432
    db: "postgres"
    password: "postgres"
    user: "postgres"
    driver: "postgres"
    # optional
    # postgres only, schema default value: public
    # schema: "public"
limit: 50

```

Once created, we can launch `dblab` with the command:

```{ .sh .copy }

dblab --config

```
If you don't specify the name of the database configuration with `--cfg-name`, then `dblab` will use the first configuration defined under the `database` field.

### Multiple Databases

But, as we all know, on a daily basis we tend to access multiple databases or the "same" database in different environments.
So using the `--cfg-name` flag can be very handy in these cases.

In the following case we have 3 environments: `local`, `staging`, and `prod`. The YAML file would look like this (but with your own credentials):

```{ .yaml .copy }

database:
  - name: "local"
    host: "<LOCAL HOST ADDRESS>"
    port: 5432
    db: "<DB NAME>"
    password: "<PASSWORD>"
    user: "<USERNAME>"
    schema: "public"
    driver: "postgres"
  - name: "staging"
    host: "<STAGING HOST ADDRESS>"
    port: 5432
    db: "<DB NAME>"
    password: "<PASSWORD>"
    user: "<USERNAME>"
    schema: "public"
    driver: "postgres"
  - name: "prod"
    host: "<PROD HOST ADDRESS>"
    port: 5432
    db: "<DB NAME>"
    password: "<PASSWORD>"
    user: "<USERNAM>"
    schema: "public"
    driver: "postgres"
# should be greater than 0, otherwise the app will error out
limit: 50

```
And in order to launch a specific environment/configuration we have to use the `--cfg-name` flag, followed by the name of the database configuration.

```{ .sh .copy }

dblab --config --cfg-name "prod"

```
