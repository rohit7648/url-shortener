#!/bin/bash

# Get the PostgreSQL container ID
PG_CONTAINER=$(docker ps | grep postgres | awk '{print $1}')

# Execute the schema SQL
docker exec -i $PG_CONTAINER psql -U postgres -d urlshortener < internal/data/schema.sql 