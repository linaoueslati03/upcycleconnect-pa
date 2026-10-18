ALTER TABLE annonces ALTER COLUMN statut SET DEFAULT 'en_attente';
ALTER TABLE annonces ADD CONSTRAINT annonces_statut_valide
    CHECK (statut IN ('en_attente', 'refusee', 'en_ligne', 'reservee', 'cedee'));
