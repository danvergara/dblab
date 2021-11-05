#!/bin/bash

set -e

echo "Running the migrations against the DB"
go run cmd/dbmigrate/main.go migrate up

echo "Seeding the database"
go run cmd/seeder/main.go seed
