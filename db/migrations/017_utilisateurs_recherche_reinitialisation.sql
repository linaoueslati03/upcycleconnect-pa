ALTER TABLE utilisateurs ADD COLUMN doit_changer_mot_de_passe BOOLEAN NOT NULL DEFAULT FALSE;

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX utilisateurs_recherche_index ON utilisateurs
    USING gin ((lower(nom || ' ' || prenom || ' ' || email)) gin_trgm_ops);

CREATE INDEX utilisateurs_role_statut_index ON utilisateurs (role_id, statut);
