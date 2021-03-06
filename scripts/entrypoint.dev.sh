#!/bin/bash

set -e

echo "Waiting for postgres..."

while ! nc -z postgres 5432; do
  sleep 0.1
done

echo "PostgreSQL started"

echo "Running the migrations against the DB"
go run cmd/dbmigrate/main.go migrate up

echo "Seeding the database"
go run cmd/seeder/main.go seed
