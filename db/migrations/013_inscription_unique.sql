DELETE FROM inscriptions i
USING inscriptions plus_ancienne
WHERE i.utilisateur_id = plus_ancienne.utilisateur_id
  AND i.type_offre = plus_ancienne.type_offre
  AND i.offre_id = plus_ancienne.offre_id
  AND i.id > plus_ancienne.id;

ALTER TABLE inscriptions
    ADD CONSTRAINT inscriptions_unique UNIQUE (utilisateur_id, type_offre, offre_id);
