-- Catalogue des prestations, géré depuis le back-office administrateur.

CREATE TABLE prestations (
    id SERIAL PRIMARY KEY,
    titre VARCHAR(255) NOT NULL,
    categorie VARCHAR(100) NOT NULL,
    tarif DECIMAL(10, 2) NOT NULL,
    statut VARCHAR(50) DEFAULT 'Brouillon'
);
