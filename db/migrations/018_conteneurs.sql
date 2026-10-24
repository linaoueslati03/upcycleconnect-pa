UPDATE conteneurs SET statut = 'vide' WHERE statut IS NULL;

UPDATE conteneurs c SET statut = 'occupe'
    WHERE EXISTS (SELECT 1 FROM depots d WHERE d.conteneur_id = c.id AND d.statut = 'deposee');

ALTER TABLE conteneurs
    ALTER COLUMN statut SET NOT NULL,
    ADD CONSTRAINT conteneurs_statut_valide CHECK (statut IN ('vide', 'occupe', 'maintenance'));

CREATE UNIQUE INDEX depots_un_objet_par_conteneur
    ON depots (conteneur_id)
    WHERE statut IN ('validee', 'deposee');

CREATE TABLE historique_conteneurs (
    id SERIAL PRIMARY KEY,
    conteneur_id INT NOT NULL REFERENCES conteneurs(id),
    ancien_statut VARCHAR(20),
    nouveau_statut VARCHAR(20) NOT NULL,
    depot_id INT REFERENCES depots(id),
    auteur_id INT REFERENCES utilisateurs(id),
    motif TEXT,
    date_changement TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX historique_conteneurs_conteneur ON historique_conteneurs (conteneur_id, date_changement DESC);
