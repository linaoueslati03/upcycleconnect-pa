#!/bin/sh
# Crée une base UpcycleConnect VIDE, prête pour un vrai déploiement :
#   - toutes les migrations (structure + référentiels : rôles, langues, matériaux et facteurs CO2,
#     paramètres financiers) ;
#   - un conteneur par site de collecte ;
#   - un seul compte administrateur, dont le mot de passe est demandé au clavier (jamais écrit dans le code).
# Aucun utilisateur, aucune annonce, aucun paiement : tout le reste se crée depuis l'application.
#
# Utilisation, depuis la racine du projet (conteneur postgres démarré) :
#   sh db/creer_base_vide.sh                         # crée la base « upcycleconnect_vide »
#   sh db/creer_base_vide.sh ma_base                 # autre nom
#   sh db/creer_base_vide.sh ma_base --remplacer     # supprime puis recrée la base si elle existe
# Puis, pour que l'API l'utilise : DB_NAME=upcycleconnect_vide dans .env et « docker compose up -d api ».

set -e
cd "$(dirname "$0")/.."

BASE="${1:-upcycleconnect_vide}"
REMPLACER="$2"
SITES="Paris 10e;Paris 11e;Paris 13e;Paris 16e;Bourg-la-Reine;Ivry;Montreuil"

if ! printf '%s' "$BASE" | grep -Eq '^[a-z][a-z0-9_]{2,62}$'; then
    echo "Nom de base invalide : lettres minuscules, chiffres et _ uniquement." >&2
    exit 1
fi

# Les commandes Docker ne doivent pas lire l'entrée standard : elle sert aux réponses du script.
dans_postgres() {
    docker compose exec -T "$@" </dev/null
}

existe=$(dans_postgres postgres sh -c "psql -U \"\$POSTGRES_USER\" -d postgres -tAc \"SELECT 1 FROM pg_database WHERE datname = '$BASE'\"")
if [ -n "$existe" ] && [ "$REMPLACER" != "--remplacer" ]; then
    echo "La base « $BASE » existe déjà. Relancez avec --remplacer pour la supprimer et la recréer." >&2
    exit 1
fi

# Compte administrateur
printf "E-mail de l'administrateur : "
read -r EMAIL || true
printf "Prénom : "
read -r PRENOM || true
printf "Nom : "
read -r NOM || true
if ! printf '%s' "$EMAIL" | grep -Eq '^[^@ ]+@[^@ ]+\.[^@ ]+$' || [ -z "$PRENOM" ] || [ -z "$NOM" ]; then
    echo "E-mail, prénom et nom sont obligatoires." >&2
    exit 1
fi

stty -echo 2>/dev/null || true
printf "Mot de passe (12 caractères minimum, majuscule, minuscule, chiffre, caractère spécial) : "
read -r MOT_DE_PASSE || true
printf "\nConfirmation : "
read -r CONFIRMATION || true
stty echo 2>/dev/null || true
printf "\n"

if [ "$MOT_DE_PASSE" != "$CONFIRMATION" ]; then
    echo "Les deux mots de passe ne correspondent pas." >&2
    exit 1
fi
if [ "${#MOT_DE_PASSE}" -lt 12 ] \
    || ! printf '%s' "$MOT_DE_PASSE" | grep -q '[A-Z]' \
    || ! printf '%s' "$MOT_DE_PASSE" | grep -q '[a-z]' \
    || ! printf '%s' "$MOT_DE_PASSE" | grep -q '[0-9]' \
    || ! printf '%s' "$MOT_DE_PASSE" | grep -q '[^A-Za-z0-9]'; then
    echo "Mot de passe trop faible (même règle que l'application)." >&2
    exit 1
fi

if [ -n "$existe" ]; then
    echo "Suppression de l'ancienne base « $BASE »…"
    dans_postgres postgres sh -c "psql -U \"\$POSTGRES_USER\" -d postgres -qc 'DROP DATABASE \"$BASE\" WITH (FORCE)'"
fi

echo "Création de la base « $BASE »…"
dans_postgres postgres sh -c "psql -U \"\$POSTGRES_USER\" -d postgres -qc 'CREATE DATABASE \"$BASE\"'"

echo "Migrations…"
dans_postgres -e POSTGRES_DB="$BASE" postgres sh /db/migrer.sh

echo "Conteneurs de collecte et compte administrateur…"
# Le mot de passe passe par une variable d'environnement : il n'apparaît ni dans la commande ni dans un fichier.
# Il est haché en bcrypt (même format que l'API) avec pgcrypto, retiré juste après.
docker compose exec -T -e MOT_DE_PASSE_ADMIN="$MOT_DE_PASSE" -e EMAIL_ADMIN="$EMAIL" -e PRENOM_ADMIN="$PRENOM" -e NOM_ADMIN="$NOM" \
    -e SITES="$SITES" postgres sh -c "psql -v ON_ERROR_STOP=1 -U \"\$POSTGRES_USER\" -d \"$BASE\" -q" <<'SQL'
\getenv mdp MOT_DE_PASSE_ADMIN
\getenv email EMAIL_ADMIN
\getenv prenom PRENOM_ADMIN
\getenv nom NOM_ADMIN
\getenv sites SITES
BEGIN;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

INSERT INTO conteneurs (site)
    SELECT unnest(string_to_array(:'sites', ';'));

INSERT INTO utilisateurs (role_id, nom, prenom, email, mot_de_passe_hash, langue_preferee_id, tutoriel_vu)
    SELECT r.id, :'nom', :'prenom', lower(:'email'), crypt(:'mdp', gen_salt('bf', 10)), l.id, false
    FROM roles r, langues l
    WHERE r.code = 'administrateur' AND l.code = 'fr';

DROP EXTENSION pgcrypto;
COMMIT;
SQL

dans_postgres postgres sh -c "psql -U \"\$POSTGRES_USER\" -d \"$BASE\" -tAc \"
    SELECT 'Base « $BASE » prête : '
        || (SELECT count(*) FROM utilisateurs) || ' utilisateur (administrateur), '
        || (SELECT count(*) FROM conteneurs) || ' conteneurs, '
        || (SELECT count(*) FROM categories_materiaux WHERE co2_evite_par_kg IS NOT NULL) || ' matériaux avec facteur CO2, '
        || (SELECT count(*) FROM schema_migrations) || ' migrations.'\""
echo "Pour l'utiliser : DB_NAME=$BASE dans .env, puis « docker compose up -d api »."
