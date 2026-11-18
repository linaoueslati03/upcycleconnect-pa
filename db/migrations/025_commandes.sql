-- Achats : paiement d'une formation et achat d'un objet mis en vente.
-- Une commande garde la trace d'un achat du début (en attente de paiement) à la fin (payée ou annulée).
-- Le paiement réussi est ensuite enregistré dans la table paiements (via enregistrerPaiement, comme convenu).
-- « fournisseur » vaut 'simulation' tant que Stripe n'est pas branché, puis 'stripe'.
CREATE TABLE commandes (
    id SERIAL PRIMARY KEY,
    utilisateur_id INT NOT NULL REFERENCES utilisateurs(id),
    type VARCHAR(20) NOT NULL CHECK (type IN ('formation', 'vente')),
    reference_id INT NOT NULL,
    libelle VARCHAR(255) NOT NULL,
    montant NUMERIC(10, 2) NOT NULL CHECK (montant > 0),
    statut VARCHAR(20) NOT NULL DEFAULT 'en_attente'
        CHECK (statut IN ('en_attente', 'payee', 'annulee', 'expiree')),
    fournisseur VARCHAR(20) NOT NULL DEFAULT 'simulation' CHECK (fournisseur IN ('simulation', 'stripe')),
    reference_externe VARCHAR(255) UNIQUE,
    paiement_id INT REFERENCES paiements(id),
    date_creation TIMESTAMP NOT NULL DEFAULT now(),
    date_expiration TIMESTAMP NOT NULL,
    date_paiement TIMESTAMP
);

CREATE INDEX commandes_utilisateur ON commandes (utilisateur_id, date_creation DESC);
CREATE INDEX commandes_en_attente ON commandes (type, reference_id) WHERE statut = 'en_attente';

-- Qui a acheté l'objet (rempli quand le paiement est confirmé)
ALTER TABLE annonces ADD COLUMN acheteur_id INT REFERENCES utilisateurs(id) ON DELETE SET NULL;

-- Facture PDF archivée automatiquement après chaque achat
ALTER TABLE documents ADD COLUMN commande_id INT REFERENCES commandes(id) ON DELETE SET NULL;
ALTER TABLE documents DROP CONSTRAINT documents_type_valide;
ALTER TABLE documents ADD CONSTRAINT documents_type_valide CHECK (
    type_document IN (
        'fiche_atelier',
        'fiche_evenement',
        'fiche_formation',
        'attestation_inscription',
        'contrat',
        'document_administratif',
        'bon_depot',
        'facture'
    )
);
