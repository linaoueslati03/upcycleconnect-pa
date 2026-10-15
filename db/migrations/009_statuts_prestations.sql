UPDATE prestations SET statut = 'publie' WHERE statut = 'Publiée';
UPDATE prestations SET statut = 'brouillon' WHERE statut <> 'publie';

ALTER TABLE prestations ALTER COLUMN statut SET DEFAULT 'brouillon';
ALTER TABLE prestations ADD CONSTRAINT prestations_statut_valide CHECK (statut IN ('brouillon', 'publie'));
