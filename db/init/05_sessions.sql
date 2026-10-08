
CREATE TABLE sessions (
    token VARCHAR(64) PRIMARY KEY,
    utilisateur_id INT NOT NULL REFERENCES utilisateurs(id) ON DELETE CASCADE,
    date_creation TIMESTAMP DEFAULT now(),
    date_expiration TIMESTAMP NOT NULL
);
