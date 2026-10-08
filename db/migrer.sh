#!/bin/sh
# Applique dans l'ordre les fichiers de db/migrations/ qui n'ont pas encore été appliqués.
# La table schema_migrations garde la liste des fichiers déjà passés.
#
# Utilisation (base déjà créée) : docker compose exec postgres sh /db/migrer.sh

set -e

# Masque les messages d'information de PostgreSQL (ex. « table déjà existante »)
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
        # --single-transaction : si le fichier échoue, rien n'est appliqué ni enregistré
        $PSQL --single-transaction \
            -f "$fichier" \
            -c "INSERT INTO schema_migrations (version) VALUES ('$version')"
    fi
done
