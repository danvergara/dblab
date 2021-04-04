#!/bin/bash
#
# Integartion testing with dockerized database servers
#

set -e

# database credentials
export DBHOST=${DBHOST:-localhost}
export PGUSER="postgres"
export DBPASSWORD="password"
export DATABASE="users"
export PGPORT="15432"

# Test different versions of postgres available on Docker Hub.
pgversions="9.6 10.16 11.11 12.6 13.2"

for i in $pgversions
do
	PGVERSION="$i"
	echo "--------------BEGIN POSTGRES TESTS-------------"
	echo "Running test against PostgreSQL v$PGVERSION"
	  docker rm -f postgres || true
	  docker run -p $PGPORT:5432 --name postgres -e POSTGRES_PASSWORD=$DBPASSWORD -d postgres:$PGVERSION
		sleep 5
		make test
	echo "--------------END POSTGRES TESTS-------------"
done
