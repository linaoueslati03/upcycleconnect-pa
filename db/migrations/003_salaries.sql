-- Profil salarié : complète le compte utilisateur (rôle « salarie ») avec les infos propres
-- aux salariés. La clé est l'id du compte : un salarié = un utilisateur, pas de doublon.

CREATE TABLE salaries (
    utilisateur_id INT PRIMARY KEY REFERENCES utilisateurs(id) ON DELETE CASCADE,
    est_responsable BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now()
);
