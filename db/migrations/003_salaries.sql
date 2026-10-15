CREATE TABLE salaries (
    utilisateur_id INT PRIMARY KEY REFERENCES utilisateurs(id) ON DELETE CASCADE,
    est_responsable BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now()
);
