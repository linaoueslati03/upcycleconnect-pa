-- Offres proposées par les salariés : événements, formations et ateliers.
-- Statuts : 'brouillon' -> 'en_attente' (soumis) -> 'publie' (validé par un responsable).

CREATE TABLE evenements (
    id SERIAL PRIMARY KEY,
    titre VARCHAR(150) NOT NULL,
    description TEXT,
    date_debut TIMESTAMP NOT NULL,
    date_fin TIMESTAMP,
    lieu VARCHAR(150),
    site VARCHAR(50),
    statut VARCHAR(20) DEFAULT 'brouillon',
    createur_id INT NOT NULL REFERENCES salaries(utilisateur_id),
    valide_par_id INT REFERENCES salaries(utilisateur_id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE evenement_salarie (
    evenement_id INT REFERENCES evenements(id) ON DELETE CASCADE,
    salarie_id INT REFERENCES salaries(utilisateur_id) ON DELETE CASCADE,
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
    tarif NUMERIC(8, 2) DEFAULT 0,
    materiel_necessaire TEXT,
    statut VARCHAR(20) DEFAULT 'brouillon',
    createur_id INT NOT NULL REFERENCES salaries(utilisateur_id),
    valide_par_id INT REFERENCES salaries(utilisateur_id),
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
    createur_id INT NOT NULL REFERENCES salaries(utilisateur_id),
    responsable_id INT REFERENCES salaries(utilisateur_id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE atelier_intervenant (
    atelier_id INT REFERENCES ateliers(id) ON DELETE CASCADE,
    salarie_id INT REFERENCES salaries(utilisateur_id) ON DELETE CASCADE,
    PRIMARY KEY (atelier_id, salarie_id)
);
