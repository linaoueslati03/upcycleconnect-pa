-- Espace particuliers : annonces, dépôts en conteneur, inscriptions aux offres,
-- Upcycling Score et projets d'upcycling.

CREATE TABLE categories_materiaux (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL
);


CREATE TABLE annonces (
    id SERIAL PRIMARY KEY,
    utilisateur_id INT NOT NULL REFERENCES utilisateurs(id),
    titre VARCHAR(150) NOT NULL,
    description TEXT,
    type VARCHAR(10) NOT NULL CHECK (type IN ('don', 'vente')),
    prix NUMERIC(10, 2),
    categorie_id INT REFERENCES categories_materiaux(id),
    localisation VARCHAR(150),
    statut VARCHAR(20) DEFAULT 'en_ligne',
    date_creation TIMESTAMP DEFAULT now()
);


CREATE TABLE conteneurs (
    id SERIAL PRIMARY KEY,
    site VARCHAR(50) NOT NULL,
    statut VARCHAR(20) DEFAULT 'vide'
);


CREATE TABLE depots (
    id SERIAL PRIMARY KEY,
    utilisateur_id INT NOT NULL REFERENCES utilisateurs(id),
    conteneur_id INT REFERENCES conteneurs(id),
    description_objet TEXT,
    code_ouverture VARCHAR(20) UNIQUE,
    code_barre VARCHAR(30) UNIQUE,
    statut VARCHAR(20) DEFAULT 'demande',
    professionnel_recuperateur_id INT REFERENCES utilisateurs(id),
    date_demande TIMESTAMP DEFAULT now(),
    date_validation TIMESTAMP,
    date_recuperation TIMESTAMP
);


-- Pas de FK stricte sur (type_offre, offre_id) : selon type_offre, offre_id pointe vers
-- formations, ateliers ou evenements (3 tables distinctes), Postgres ne supporte pas de
-- FK polymorphe native.
CREATE TABLE inscriptions (
    id SERIAL PRIMARY KEY,
    utilisateur_id INT NOT NULL REFERENCES utilisateurs(id),
    type_offre VARCHAR(20) NOT NULL,
    offre_id INT NOT NULL,
    date_inscription TIMESTAMP DEFAULT now(),
    statut_paiement VARCHAR(20) DEFAULT 'gratuit',
    montant_paye NUMERIC(10, 2) DEFAULT 0
);


CREATE TABLE score_historique (
    id SERIAL PRIMARY KEY,
    utilisateur_id INT NOT NULL REFERENCES utilisateurs(id),
    delta INT NOT NULL,
    motif VARCHAR(150),
    date TIMESTAMP DEFAULT now()
);


CREATE TABLE projets_upcycling (
    id SERIAL PRIMARY KEY,
    utilisateur_id INT NOT NULL REFERENCES utilisateurs(id),
    titre VARCHAR(150) NOT NULL,
    description TEXT,
    partage_public BOOLEAN DEFAULT FALSE,
    sponsorise BOOLEAN DEFAULT FALSE,
    date_creation TIMESTAMP DEFAULT now()
);


CREATE TABLE etapes_projet (
    id SERIAL PRIMARY KEY,
    projet_id INT NOT NULL REFERENCES projets_upcycling(id) ON DELETE CASCADE,
    description TEXT,
    photo_url VARCHAR(255),
    ordre INT NOT NULL,
    date TIMESTAMP DEFAULT now()
);
