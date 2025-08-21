#!/bin/bash
set -e

export PGPASSWORD=$POSTGRESQL_POSTGRES_PGPASSWORD
psql -v ON_ERROR_STOP=1 -U "$POSTGRESQL_USER" -d "$POSTGRESQL_DATABASE" <<-EOSQL
    --Set up pg_partman schema and extension
    CREATE SCHEMA IF NOT EXISTS partman;
    CREATE EXTENSION IF NOT EXISTS pg_partman SCHEMA partman;
EOSQL

echo "Partition created successfully"