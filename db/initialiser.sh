#!/bin/sh

set -e

sh /db/migrer.sh

for fichier in /db/seeds/*.sql; do
    echo "Données $(basename "$fichier")"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" --quiet -f "$fichier"
done
