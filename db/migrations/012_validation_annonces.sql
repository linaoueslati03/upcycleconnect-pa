-- Le sujet impose qu'une annonce soit validée par le service administratif avant d'être
-- publiée : une nouvelle annonce est créée « en_attente », puis validée ou refusée.
-- Les annonces déjà en ligne restent en ligne.

ALTER TABLE annonces ALTER COLUMN statut SET DEFAULT 'en_attente';
ALTER TABLE annonces ADD CONSTRAINT annonces_statut_valide
    CHECK (statut IN ('en_attente', 'refusee', 'en_ligne', 'reservee', 'cedee'));

-- Textes de l'interface liés à la validation des annonces
INSERT INTO traductions (langue_id, cle, texte)
SELECT l.id, t.cle, CASE l.code WHEN 'fr' THEN t.fr ELSE t.en END
FROM langues l
CROSS JOIN (VALUES
    ('annonce.en_attente', 'En attente de validation', 'Awaiting approval'),
    ('annonce.refusee', 'Refusée', 'Rejected'),
    ('annonces.envoyee_validation', 'Annonce envoyée : elle sera visible après validation par notre équipe.', 'Listing sent: it will be visible once approved by our team.'),
    ('admin.annonces.titre', 'Validation des annonces', 'Listing approval'),
    ('admin.annonces.a_valider', 'À valider', 'To approve'),
    ('admin.annonces.refusees', 'Refusées', 'Rejected'),
    ('admin.annonces.aucune', 'Aucune annonce.', 'No listing.'),
    ('admin.annonces.auteur', 'Auteur', 'Author'),
    ('admin.annonces.erreur_chargement', 'Erreur lors du chargement des annonces', 'Error while loading listings'),
    ('commun.refuser', 'Refuser', 'Reject')
) AS t(cle, fr, en)
WHERE l.code IN ('fr', 'en');
