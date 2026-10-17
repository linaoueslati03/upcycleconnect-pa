ALTER TABLE forum_sujets
    ALTER COLUMN statut SET NOT NULL,
    ADD CONSTRAINT forum_sujets_statut_valide CHECK (statut IN ('en_attente', 'valide', 'rejete'));

ALTER TABLE forum_messages
    ALTER COLUMN statut SET NOT NULL,
    ADD CONSTRAINT forum_messages_statut_valide CHECK (statut IN ('visible', 'masque'));

ALTER TABLE forum_signalements
    ALTER COLUMN statut SET NOT NULL,
    ADD CONSTRAINT forum_signalements_statut_valide CHECK (statut IN ('a_traiter', 'traite', 'rejete'));
