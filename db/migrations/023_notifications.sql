CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    titre_fr VARCHAR(100) NOT NULL,
    message_fr VARCHAR(500) NOT NULL,
    titre_en VARCHAR(100) NOT NULL,
    message_en VARCHAR(500) NOT NULL,
    cible VARCHAR(20) NOT NULL
        CHECK (cible IN ('tous', 'particuliers', 'professionnels', 'utilisateur', 'prediction')),
    destinataire_id INT REFERENCES utilisateurs(id) ON DELETE SET NULL,
    classe_predite VARCHAR(20),
    nb_destinataires INT NOT NULL DEFAULT 0 CHECK (nb_destinataires >= 0),
    statut VARCHAR(10) NOT NULL CHECK (statut IN ('envoyee', 'echec')),
    erreur TEXT,
    onesignal_id VARCHAR(100),
    auteur_id INT REFERENCES utilisateurs(id) ON DELETE SET NULL,
    date_envoi TIMESTAMP NOT NULL DEFAULT now(),
    CHECK ((cible = 'prediction') = (classe_predite IS NOT NULL))
);

CREATE INDEX notifications_date ON notifications (date_envoi DESC);
