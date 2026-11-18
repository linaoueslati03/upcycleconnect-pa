-- Les matériaux font partie du référentiel de l'application (dépôts, annonces, facteurs CO2).
-- Jusqu'ici ils n'étaient créés que par les données de test : sur une base neuve (installation,
-- base vide), la migration 027 ne trouvait aucun matériau à qui donner son facteur.
-- Cette migration les crée s'ils manquent et complète les facteurs encore vides (mêmes valeurs et sources que 027).
INSERT INTO categories_materiaux (code, libelle) VALUES
    ('bois', 'Bois'), ('metal', 'Métal'), ('textile', 'Textile'), ('verre', 'Verre'), ('plastique', 'Plastique')
ON CONFLICT (code) DO NOTHING;

UPDATE categories_materiaux c
SET co2_evite_par_kg = f.facteur, source_facteur = f.source
FROM (VALUES
    ('bois', 1.833, 'ADEME Base Carbone n° 20907 « Mobilier – fabrication » : 1 833 kgCO2e/t (2014, incertitude 50 %). Fabrication d''un meuble neuf évitée.'),
    ('metal', 2.211, 'ADEME Base Carbone n° 26729 « Acier ou fer blanc – neuf » : 2 211 kgCO2e/t (FEDEREC/ADEME 2018, incertitude 10 %). Production d''acier neuf évitée.'),
    ('textile', 4.9, 'Refashion, Bilan environnemental – synthèse, juillet 2026 (données 2025), p. 12 : bénéfice net de la seconde main en France -4,9 kgCO2e/kg.'),
    ('verre', 0.923, 'ADEME Base Carbone n° 26741 « Verre – neuf » : 923 kgCO2e/t (FEDEREC/ADEME 2018, incertitude 50 %). Production de verre neuf évitée.'),
    ('plastique', 2.383, 'ADEME Base Carbone n° 20836 « Plastique – moyenne – neuf » : 2 383 kgCO2e/t (2014, incertitude 20 %). Production de plastique neuf évitée.')
) AS f(code, facteur, source)
WHERE c.code = f.code AND c.co2_evite_par_kg IS NULL;
