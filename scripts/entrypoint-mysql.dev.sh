#!/bin/bash

set -e

echo "Waiting for postgres..."

while ! nc -z mysql 3306; do
  sleep 0.1
done

echo "MySQL started"

echo "Running the migrations against the DB"
go run cmd/dbmigrate/main.go

echo "Seeding the database"
go run cmd/seeder/main.go seed
