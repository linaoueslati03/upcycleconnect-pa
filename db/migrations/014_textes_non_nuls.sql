-- L'API lit ces colonnes comme du texte (ou un nombre) toujours présent : une valeur NULL
-- insérée hors de l'API (ex. directement en SQL) faisait échouer toute la liste concernée.
-- On remplace les NULL existants puis on interdit les nouveaux, avec une valeur par défaut.

UPDATE evenements SET description = '' WHERE description IS NULL;
UPDATE evenements SET lieu = '' WHERE lieu IS NULL;
UPDATE evenements SET site = '' WHERE site IS NULL;
ALTER TABLE evenements
    ALTER COLUMN description SET DEFAULT '', ALTER COLUMN description SET NOT NULL,
    ALTER COLUMN lieu SET DEFAULT '', ALTER COLUMN lieu SET NOT NULL,
    ALTER COLUMN site SET DEFAULT '', ALTER COLUMN site SET NOT NULL;

UPDATE annonces SET description = '' WHERE description IS NULL;
UPDATE annonces SET localisation = '' WHERE localisation IS NULL;
ALTER TABLE annonces
    ALTER COLUMN description SET DEFAULT '', ALTER COLUMN description SET NOT NULL,
    ALTER COLUMN localisation SET DEFAULT '', ALTER COLUMN localisation SET NOT NULL;

UPDATE projets_upcycling SET description = '' WHERE description IS NULL;
ALTER TABLE projets_upcycling ALTER COLUMN description SET DEFAULT '', ALTER COLUMN description SET NOT NULL;

UPDATE etapes_projet SET description = '' WHERE description IS NULL;
ALTER TABLE etapes_projet ALTER COLUMN description SET DEFAULT '', ALTER COLUMN description SET NOT NULL;

UPDATE utilisateurs SET statut = 'actif' WHERE statut IS NULL;
UPDATE utilisateurs SET upcycling_score = 0 WHERE upcycling_score IS NULL;
ALTER TABLE utilisateurs ALTER COLUMN statut SET NOT NULL, ALTER COLUMN upcycling_score SET NOT NULL;
