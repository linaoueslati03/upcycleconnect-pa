-- Forum : tout utilisateur peut écrire, la modération est faite par les salariés.

CREATE TABLE forum_sujets (
    id SERIAL PRIMARY KEY,
    titre VARCHAR(150) NOT NULL,
    auteur_utilisateur_id INT NOT NULL REFERENCES utilisateurs(id),
    statut VARCHAR(20) DEFAULT 'en_attente',
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE forum_messages (
    id SERIAL PRIMARY KEY,
    sujet_id INT NOT NULL REFERENCES forum_sujets(id) ON DELETE CASCADE,
    auteur_utilisateur_id INT NOT NULL REFERENCES utilisateurs(id),
    contenu TEXT NOT NULL,
    statut VARCHAR(20) DEFAULT 'visible',
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE forum_signalements (
    id SERIAL PRIMARY KEY,
    message_id INT NOT NULL REFERENCES forum_messages(id) ON DELETE CASCADE,
    signale_par_utilisateur_id INT NOT NULL REFERENCES utilisateurs(id),
    motif TEXT,
    statut VARCHAR(20) DEFAULT 'a_traiter',
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE moderation_actions (
    id SERIAL PRIMARY KEY,
    salarie_id INT NOT NULL REFERENCES salaries(utilisateur_id),
    cible_type VARCHAR(20) NOT NULL,
    cible_id INT NOT NULL,
    action VARCHAR(30) NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);
