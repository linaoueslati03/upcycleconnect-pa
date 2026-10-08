-- Un utilisateur ne peut s'inscrire qu'une fois à une même offre : sans cette règle, une
-- double inscription consommait deux places et comptait deux fois les points du score.
-- Les doublons éventuels sont supprimés (on garde la première inscription) avant d'ajouter
-- la contrainte, qui sert de filet de sécurité à la vérification faite par l'API.

DELETE FROM inscriptions i
USING inscriptions plus_ancienne
WHERE i.utilisateur_id = plus_ancienne.utilisateur_id
  AND i.type_offre = plus_ancienne.type_offre
  AND i.offre_id = plus_ancienne.offre_id
  AND i.id > plus_ancienne.id;

ALTER TABLE inscriptions
    ADD CONSTRAINT inscriptions_unique UNIQUE (utilisateur_id, type_offre, offre_id);
