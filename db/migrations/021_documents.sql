CREATE TABLE IF NOT EXISTS documents (
    id BIGSERIAL PRIMARY KEY,
    type_document VARCHAR(40) NOT NULL,
    titre VARCHAR(255) NOT NULL,
    utilisateur_id INTEGER REFERENCES utilisateurs(id) ON DELETE SET NULL,
    type_offre VARCHAR(20),
    offre_id INTEGER,
    inscription_id INTEGER REFERENCES inscriptions(id) ON DELETE SET NULL,
    cree_par INTEGER REFERENCES utilisateurs(id) ON DELETE SET NULL,
    contenu_pdf BYTEA NOT NULL,
    date_creation TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT documents_type_valide CHECK (
        type_document IN (
            'fiche_atelier',
            'fiche_evenement',
            'fiche_formation',
            'attestation_inscription',
            'contrat',
            'document_administratif'
        )
    )
);

CREATE INDEX IF NOT EXISTS idx_documents_utilisateur
    ON documents(utilisateur_id);

CREATE INDEX IF NOT EXISTS idx_documents_inscription
    ON documents(inscription_id);

CREATE INDEX IF NOT EXISTS idx_documents_offre
    ON documents(type_offre, offre_id);

CREATE INDEX IF NOT EXISTS idx_documents_date
    ON documents(date_creation DESC);