-- Les prestations utilisaient 'Brouillon' / 'Publiée' alors que toutes les autres tables
-- utilisent 'brouillon' / 'publie' : on aligne les valeurs et on les contraint.

UPDATE prestations SET statut = 'publie' WHERE statut = 'Publiée';
UPDATE prestations SET statut = 'brouillon' WHERE statut <> 'publie';

ALTER TABLE prestations ALTER COLUMN statut SET DEFAULT 'brouillon';
ALTER TABLE prestations ADD CONSTRAINT prestations_statut_valide CHECK (statut IN ('brouillon', 'publie'));
