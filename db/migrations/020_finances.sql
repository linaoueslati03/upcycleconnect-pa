CREATE TABLE parametres_finances (
    id INT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    taux_commission NUMERIC(4, 2) NOT NULL DEFAULT 8.00
        CHECK (taux_commission BETWEEN 5 AND 10),
    prix_premium_mensuel NUMERIC(8, 2) NOT NULL DEFAULT 20.00
        CHECK (prix_premium_mensuel BETWEEN 15 AND 30),
    date_modification TIMESTAMP NOT NULL DEFAULT now()
);

INSERT INTO parametres_finances (id) VALUES (1);

CREATE TABLE abonnements (
    id SERIAL PRIMARY KEY,
    utilisateur_id INT NOT NULL REFERENCES utilisateurs(id),
    montant_mensuel NUMERIC(8, 2) NOT NULL CHECK (montant_mensuel > 0),
    statut VARCHAR(20) NOT NULL DEFAULT 'actif'
        CHECK (statut IN ('actif', 'resilie', 'expire')),
    stripe_abonnement_id VARCHAR(100) UNIQUE,
    date_debut DATE NOT NULL DEFAULT CURRENT_DATE,
    date_fin DATE,
    CHECK (date_fin IS NULL OR date_fin >= date_debut)
);

CREATE UNIQUE INDEX abonnements_un_actif_par_utilisateur
    ON abonnements (utilisateur_id)
    WHERE statut = 'actif';

INSERT INTO abonnements (utilisateur_id, montant_mensuel)
    SELECT u.id, p.prix_premium_mensuel
    FROM utilisateurs u, parametres_finances p
    WHERE u.abonnement = 'premium';

CREATE TABLE campagnes_publicitaires (
    id SERIAL PRIMARY KEY,
    utilisateur_id INT NOT NULL REFERENCES utilisateurs(id),
    cible_type VARCHAR(20) NOT NULL CHECK (cible_type IN ('projet', 'annonce')),
    cible_id INT NOT NULL,
    montant_mensuel NUMERIC(8, 2) NOT NULL CHECK (montant_mensuel BETWEEN 100 AND 500),
    statut VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (statut IN ('active', 'terminee', 'annulee')),
    date_debut DATE NOT NULL,
    date_fin DATE,
    CHECK (date_fin IS NULL OR date_fin >= date_debut)
);

CREATE TABLE paiements (
    id SERIAL PRIMARY KEY,
    utilisateur_id INT NOT NULL REFERENCES utilisateurs(id),
    source VARCHAR(20) NOT NULL
        CHECK (source IN ('abonnement', 'vente', 'publicite', 'formation', 'atelier', 'evenement')),
    reference_id INT,
    montant NUMERIC(10, 2) NOT NULL CHECK (montant > 0),
    taux_commission NUMERIC(4, 2) CHECK (taux_commission BETWEEN 5 AND 10),
    montant_commission NUMERIC(10, 2) CHECK (montant_commission >= 0),
    statut VARCHAR(20) NOT NULL DEFAULT 'reussi'
        CHECK (statut IN ('en_attente', 'reussi', 'echoue', 'rembourse')),
    stripe_id VARCHAR(100) UNIQUE,
    date_paiement TIMESTAMP NOT NULL DEFAULT now(),
    CHECK ((source = 'vente') = (taux_commission IS NOT NULL AND montant_commission IS NOT NULL))
);

CREATE INDEX paiements_date ON paiements (date_paiement DESC);
CREATE INDEX paiements_source_date ON paiements (source, date_paiement DESC);
CREATE INDEX paiements_utilisateur ON paiements (utilisateur_id, date_paiement DESC);

INSERT INTO paiements (utilisateur_id, source, reference_id, montant, date_paiement)
    SELECT utilisateur_id, type_offre, offre_id, montant_paye, date_inscription
    FROM inscriptions
    WHERE montant_paye > 0 AND type_offre IN ('formation', 'atelier', 'evenement');

CREATE TABLE charges (
    id SERIAL PRIMARY KEY,
    libelle VARCHAR(150) NOT NULL,
    categorie VARCHAR(30) NOT NULL
        CHECK (categorie IN ('salaires', 'loyers', 'infrastructure', 'marketing', 'autre')),
    montant NUMERIC(10, 2) NOT NULL CHECK (montant > 0),
    date_charge DATE NOT NULL,
    auteur_id INT REFERENCES utilisateurs(id),
    date_creation TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX charges_date ON charges (date_charge DESC);
