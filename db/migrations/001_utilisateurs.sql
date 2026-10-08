-- Comptes de la plateforme : un seul compte par personne, son rôle dit à quel espace il accède.

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    libelle VARCHAR(50) NOT NULL
);

INSERT INTO roles (code, libelle) VALUES
    ('particulier', 'Particulier'),
    ('professionnel', 'Professionnel'),
    ('salarie', 'Salarié'),
    ('administrateur', 'Administrateur');


CREATE TABLE langues (
    id SERIAL PRIMARY KEY,
    code VARCHAR(5) UNIQUE NOT NULL,
    libelle VARCHAR(50) NOT NULL,
    actif BOOLEAN DEFAULT TRUE
);

INSERT INTO langues (code, libelle) VALUES
    ('fr', 'Français'),
    ('en', 'English');


CREATE TABLE utilisateurs (
    id SERIAL PRIMARY KEY,
    role_id INT NOT NULL REFERENCES roles(id),
    nom VARCHAR(100) NOT NULL,
    prenom VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    mot_de_passe_hash VARCHAR(255) NOT NULL,
    statut VARCHAR(20) DEFAULT 'actif',
    langue_preferee_id INT REFERENCES langues(id),
    upcycling_score INT DEFAULT 0,
    tutoriel_vu BOOLEAN DEFAULT FALSE,
    date_tutoriel_vu TIMESTAMP,
    date_creation TIMESTAMP DEFAULT now()
);
