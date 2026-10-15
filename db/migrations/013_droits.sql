UPDATE utilisateurs SET statut = 'actif' WHERE statut IS NULL;

ALTER TABLE utilisateurs
    ALTER COLUMN statut SET NOT NULL,
    ADD CONSTRAINT utilisateurs_statut_valide CHECK (statut IN ('actif', 'suspendu')),
    ADD COLUMN abonnement VARCHAR(20) NOT NULL DEFAULT 'freemium',
    ADD CONSTRAINT utilisateurs_abonnement_valide CHECK (abonnement IN ('freemium', 'premium'));

ALTER TABLE evenements
    DROP CONSTRAINT evenements_createur_id_fkey,
    DROP CONSTRAINT evenements_valide_par_id_fkey,
    ADD CONSTRAINT evenements_createur_id_fkey FOREIGN KEY (createur_id) REFERENCES utilisateurs(id),
    ADD CONSTRAINT evenements_valide_par_id_fkey FOREIGN KEY (valide_par_id) REFERENCES utilisateurs(id);

ALTER TABLE formations
    DROP CONSTRAINT formations_createur_id_fkey,
    DROP CONSTRAINT formations_valide_par_id_fkey,
    ADD CONSTRAINT formations_createur_id_fkey FOREIGN KEY (createur_id) REFERENCES utilisateurs(id),
    ADD CONSTRAINT formations_valide_par_id_fkey FOREIGN KEY (valide_par_id) REFERENCES utilisateurs(id);

ALTER TABLE ateliers
    DROP CONSTRAINT ateliers_createur_id_fkey,
    DROP CONSTRAINT ateliers_responsable_id_fkey,
    ADD CONSTRAINT ateliers_createur_id_fkey FOREIGN KEY (createur_id) REFERENCES utilisateurs(id),
    ADD CONSTRAINT ateliers_responsable_id_fkey FOREIGN KEY (responsable_id) REFERENCES utilisateurs(id);

ALTER TABLE conseils
    DROP CONSTRAINT conseils_auteur_id_fkey,
    ADD CONSTRAINT conseils_auteur_id_fkey FOREIGN KEY (auteur_id) REFERENCES utilisateurs(id);

ALTER TABLE moderation_actions RENAME COLUMN salarie_id TO moderateur_id;
ALTER TABLE moderation_actions
    DROP CONSTRAINT moderation_actions_salarie_id_fkey,
    ADD CONSTRAINT moderation_actions_moderateur_id_fkey FOREIGN KEY (moderateur_id) REFERENCES utilisateurs(id);
