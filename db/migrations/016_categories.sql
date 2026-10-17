CREATE TABLE categories_prestations (
    id SERIAL PRIMARY KEY,
    libelle VARCHAR(100) UNIQUE NOT NULL,
    actif BOOLEAN NOT NULL DEFAULT TRUE
);

INSERT INTO categories_prestations (libelle)
    SELECT DISTINCT categorie FROM prestations;

ALTER TABLE prestations ADD COLUMN categorie_id INT REFERENCES categories_prestations(id);

UPDATE prestations p SET categorie_id = c.id
    FROM categories_prestations c
    WHERE c.libelle = p.categorie;

ALTER TABLE prestations
    ALTER COLUMN categorie_id SET NOT NULL,
    DROP COLUMN categorie;

ALTER TABLE categories_materiaux
    ADD COLUMN libelle VARCHAR(100),
    ADD COLUMN actif BOOLEAN NOT NULL DEFAULT TRUE;

UPDATE categories_materiaux SET libelle = initcap(code);
UPDATE categories_materiaux SET libelle = 'Métal' WHERE code = 'metal';

ALTER TABLE categories_materiaux ALTER COLUMN libelle SET NOT NULL;
