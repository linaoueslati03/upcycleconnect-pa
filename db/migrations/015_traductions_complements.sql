-- Textes ajoutés lors de la relecture (messages d'erreur et compte professionnel).

INSERT INTO traductions (langue_id, cle, texte)
SELECT l.id, t.cle, CASE l.code WHEN 'fr' THEN t.fr ELSE t.en END
FROM langues l
CROSS JOIN (VALUES
    ('commun.serveur_injoignable', 'Le serveur est injoignable, réessayez dans un instant.', 'The server cannot be reached, please try again shortly.'),
    ('inscription.compte_cree_espace_indisponible', 'Votre compte est créé. L''espace professionnel sera bientôt disponible.', 'Your account has been created. The professional area will be available soon.')
) AS t(cle, fr, en)
WHERE l.code IN ('fr', 'en');
