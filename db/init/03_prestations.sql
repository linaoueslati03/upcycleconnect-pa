CREATE TABLE IF NOT EXISTS prestations (
    id SERIAL PRIMARY KEY,
    titre VARCHAR(255) NOT NULL,
    categorie VARCHAR(100) NOT NULL,
    tarif DECIMAL(10, 2) NOT NULL,
    statut VARCHAR(50) DEFAULT 'Brouillon'
);

INSERT INTO prestations (titre, categorie, tarif, statut) VALUES
('Audit recyclage', 'Conseil', 450.00, 'Publiée'),
('Collecte de bois', 'Logistique', 250.00, 'Publiée'),
('Formation de tri', 'Formation', 800.00, 'Brouillon');