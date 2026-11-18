-- Génération automatique des documents PDF :
-- un bon de dépôt (avec le code d'ouverture et le code-barre) est archivé à la validation d'un dépôt.
ALTER TABLE documents
    ADD COLUMN depot_id INTEGER REFERENCES depots(id) ON DELETE SET NULL;

ALTER TABLE documents DROP CONSTRAINT documents_type_valide;
ALTER TABLE documents ADD CONSTRAINT documents_type_valide CHECK (
    type_document IN (
        'fiche_atelier',
        'fiche_evenement',
        'fiche_formation',
        'attestation_inscription',
        'contrat',
        'document_administratif',
        'bon_depot'
    )
);

-- Retrouver rapidement le bon de dépôt déjà archivé pour un dépôt
CREATE INDEX IF NOT EXISTS idx_documents_depot ON documents(depot_id);
