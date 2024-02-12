#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER mylky ENCRYPTED PASSWORD 'mylky' LOGIN;
	CREATE DATABASE servord OWNER mylky;
EOSQL

psql -v ON_ERROR_STOP=1 --username "mylky" --dbname "servord" -f /app/sql/init-db.sql
