#!/bin/bash

set -e

echo "Waiting for mssql..."

while ! nc -z mssql 1433; do
  sleep 0.1
done

echo "SQL Server started"

echo "Running the migrations against the DB"
go run cmd/dbmigrate/main.go migrate up

echo "Seeding the database"
go run cmd/seeder/main.go seed
