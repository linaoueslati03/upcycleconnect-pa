UPDATE evenements SET statut = 'brouillon' WHERE statut IS NULL;
UPDATE ateliers SET statut = 'brouillon' WHERE statut IS NULL;
UPDATE formations SET statut = 'brouillon' WHERE statut IS NULL;

ALTER TABLE evenements
    ALTER COLUMN statut SET NOT NULL,
    ADD COLUMN motif_refus TEXT,
    ADD CONSTRAINT evenements_statut_valide
        CHECK (statut IN ('brouillon', 'en_attente', 'publie', 'refuse', 'annule', 'termine'));

ALTER TABLE ateliers
    ALTER COLUMN statut SET NOT NULL,
    ADD COLUMN motif_refus TEXT,
    ADD CONSTRAINT ateliers_statut_valide
        CHECK (statut IN ('brouillon', 'en_attente', 'publie', 'refuse', 'annule', 'termine'));

ALTER TABLE formations
    ALTER COLUMN statut SET NOT NULL,
    ADD COLUMN motif_refus TEXT,
    ADD CONSTRAINT formations_statut_valide
        CHECK (statut IN ('brouillon', 'en_attente', 'publie', 'refuse', 'annule', 'termine'));
