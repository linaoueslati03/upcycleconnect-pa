-- Articles de conseils rédigés par les salariés, avec l'historique de leurs versions.

CREATE TABLE conseils (
    id SERIAL PRIMARY KEY,
    titre VARCHAR(150) NOT NULL,
    contenu TEXT NOT NULL,
    categorie VARCHAR(50),
    statut VARCHAR(20) DEFAULT 'brouillon',
    auteur_id INT NOT NULL REFERENCES salaries(utilisateur_id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE conseils_historique (
    id SERIAL PRIMARY KEY,
    conseil_id INT NOT NULL REFERENCES conseils(id) ON DELETE CASCADE,
    contenu TEXT NOT NULL,
    version_date TIMESTAMP DEFAULT now()
);
