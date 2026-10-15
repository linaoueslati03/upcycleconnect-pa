#!/bin/sh

set -e

export PGOPTIONS="--client-min-messages=warning"

PSQL="psql -v ON_ERROR_STOP=1 --username $POSTGRES_USER --dbname $POSTGRES_DB --quiet"

$PSQL -c "CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    date_application TIMESTAMP DEFAULT now()
);"

for fichier in /db/migrations/*.sql; do
    version=$(basename "$fichier")
    deja_appliquee=$($PSQL --tuples-only --no-align \
        -c "SELECT 1 FROM schema_migrations WHERE version = '$version'")

    if [ -z "$deja_appliquee" ]; then
        echo "Migration $version"
        $PSQL --single-transaction \
            -f "$fichier" \
            -c "INSERT INTO schema_migrations (version) VALUES ('$version')"
    fi
done
