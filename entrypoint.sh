#!/bin/sh
set -e

echo "Run migration"

/app/server migrate up

echo "Starting Todo application..."

exec /app/server serve