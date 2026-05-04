#!/usr/bin/env sh

set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname "$0")" && pwd)
MIGRATIONS_DIR="$SCRIPT_DIR/migrations"

if [ "$#" -lt 1 ]; then
    echo "Usage: $0 <psql command...>" >&2
    exit 1
fi

run_psql() {
    "$@"
}

run_psql "$@" -v ON_ERROR_STOP=1 <<'SQL'
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
SQL

for file in $(find "$MIGRATIONS_DIR" -maxdepth 1 -type f -name '*.sql' | sort); do
    version=$(basename "$file")

    already_applied=$(
        run_psql "$@" -t -A -v ON_ERROR_STOP=1 \
            -c "SELECT 1 FROM schema_migrations WHERE version = '$version' LIMIT 1;"
    )

    if [ "$already_applied" = "1" ]; then
        echo "Skipping already applied migration: $version"
        continue
    fi

    echo "Applying migration: $version"
    run_psql "$@" -v ON_ERROR_STOP=1 -f /dev/stdin < "$file"
    run_psql "$@" -v ON_ERROR_STOP=1 \
        -c "INSERT INTO schema_migrations (version) VALUES ('$version');"
done
