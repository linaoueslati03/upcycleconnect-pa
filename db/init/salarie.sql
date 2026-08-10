
CREATE TABLE salaries (
    id SERIAL PRIMARY KEY,
    utilisateur_id INT NOT NULL,
    est_responsable BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now()
);


CREATE TABLE evenements (
    id SERIAL PRIMARY KEY,
    titre VARCHAR(150) NOT NULL,
    description TEXT,
    date_debut TIMESTAMP NOT NULL,
    date_fin TIMESTAMP,
    lieu VARCHAR(150),
    site VARCHAR(50),
    statut VARCHAR(20) DEFAULT 'brouillon',
    createur_id INT NOT NULL REFERENCES salaries(id),
    valide_par_id INT REFERENCES salaries(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE evenement_salarie (
    evenement_id INT REFERENCES evenements(id) ON DELETE CASCADE,
    salarie_id INT REFERENCES salaries(id) ON DELETE CASCADE,
    PRIMARY KEY (evenement_id, salarie_id)
);


CREATE TABLE formations (
    id SERIAL PRIMARY KEY,
    titre VARCHAR(150) NOT NULL,
    description TEXT,
    date_debut TIMESTAMP NOT NULL,
    date_fin TIMESTAMP,
    lieu VARCHAR(150),
    nb_places INT,
    tarif NUMERIC(8,2) DEFAULT 0,
    materiel_necessaire TEXT,
    statut VARCHAR(20) DEFAULT 'brouillon',
    createur_id INT NOT NULL REFERENCES salaries(id),
    valide_par_id INT REFERENCES salaries(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);


CREATE TABLE ateliers (
    id SERIAL PRIMARY KEY,
    titre VARCHAR(150) NOT NULL,
    description TEXT,
    date_debut TIMESTAMP NOT NULL,
    date_fin TIMESTAMP,
    lieu VARCHAR(150),
    statut VARCHAR(20) DEFAULT 'brouillon',
    createur_id INT NOT NULL REFERENCES salaries(id),
    responsable_id INT REFERENCES salaries(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE atelier_intervenant (
    atelier_id INT REFERENCES ateliers(id) ON DELETE CASCADE,
    salarie_id INT REFERENCES salaries(id) ON DELETE CASCADE,
    PRIMARY KEY (atelier_id, salarie_id)
);


CREATE TABLE conseils (
    id SERIAL PRIMARY KEY,
    titre VARCHAR(150) NOT NULL,
    contenu TEXT NOT NULL,
    categorie VARCHAR(50),
    statut VARCHAR(20) DEFAULT 'brouillon',
    auteur_id INT NOT NULL REFERENCES salaries(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE conseils_historique (
    id SERIAL PRIMARY KEY,
    conseil_id INT NOT NULL REFERENCES conseils(id) ON DELETE CASCADE,
    contenu TEXT NOT NULL,
    version_date TIMESTAMP DEFAULT now()
);


CREATE TABLE forum_sujets (
    id SERIAL PRIMARY KEY,
    titre VARCHAR(150) NOT NULL,
    auteur_utilisateur_id INT NOT NULL,
    statut VARCHAR(20) DEFAULT 'en_attente',
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE forum_messages (
    id SERIAL PRIMARY KEY,
    sujet_id INT NOT NULL REFERENCES forum_sujets(id) ON DELETE CASCADE,
    auteur_utilisateur_id INT NOT NULL,
    contenu TEXT NOT NULL,
    statut VARCHAR(20) DEFAULT 'visible',
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE forum_signalements (
    id SERIAL PRIMARY KEY,
    message_id INT NOT NULL REFERENCES forum_messages(id) ON DELETE CASCADE,
    signale_par_utilisateur_id INT NOT NULL,
    motif TEXT,
    statut VARCHAR(20) DEFAULT 'a_traiter',
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE moderation_actions (
    id SERIAL PRIMARY KEY,
    salarie_id INT NOT NULL REFERENCES salaries(id),
    cible_type VARCHAR(20) NOT NULL,
    cible_id INT NOT NULL,
    action VARCHAR(30) NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);