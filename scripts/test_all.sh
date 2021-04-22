#!/bin/bash
#
# Integartion testing with dockerized database servers
#

set -e

# database credentials

export MYSQL_PORT="3306"

# Test different versions of postgres available on Docker Hub.
pgversions="9.6 10.16 11.11 12.6 13.2"

for i in $pgversions
do
	PGVERSION="$i"
	echo "--------------BEGIN POSTGRES TESTS-------------"
	echo "Running test against PostgreSQL v$PGVERSION"

	export DB_HOST=${DBHOST:-localhost}
	export DB_USER="postgres"
	export DB_PASSWORD="password"
	export DB_NAME="users"
	export DB_PORT="5432"
	export DB_DRIVER="postgres"
	docker rm -f postgres || true
	docker run -p $DB_PORT:5432 --name postgres -e POSTGRES_PASSWORD=$DB_PASSWORD -e POSTGRES_DB=$DB_NAME -d postgres:$PGVERSION
	sleep 5
	make test
	echo "--------------END POSTGRES TESTS-------------"
done


# Test different versions of mysql available on Docker Hub.
mysql_versions="5.6 5.7 8.0"

for i in $mysql_versions
do
	MYSQL_VERSION="$i"
	echo "--------------BEGIN MYSQL TESTS-------------"
	echo "Running test against MySQL v$MYSQL_VERSION"

	export DB_HOST=${DBHOST:-localhost}
	export DB_USER="root"
	export DB_PASSWORD="password"
	export DB_NAME="users"
	export DB_PORT="3306"
	export DB_DRIVER="mysql"
	docker rm -f mysql || true
	docker run -p $DB_PORT:3306 --name mysql -e MYSQL_PASSWORD=$DB_PASSWORD -e MYSQL_DATABASE=$DB_NAME -d mysql:$MYSQL_VERSION
	sleep 5
	make test
	echo "--------------END MYSQL TESTS-------------"
done
