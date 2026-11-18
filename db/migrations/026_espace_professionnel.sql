-- Espace professionnel : entreprise, abonnement Premium, publicité, statistiques d'impact,
-- alertes de collecte et notifications dans le site.

-- 1. Entreprise du professionnel (Mon compte)
ALTER TABLE utilisateurs
    ADD COLUMN siret CHAR(14) CHECK (siret ~ '^[0-9]{14}$'),
    ADD COLUMN adresse_entreprise VARCHAR(255);

-- 2. Matériau et poids des dépôts (calcul du CO2 évité, alertes par matériau)
ALTER TABLE depots
    ADD COLUMN categorie_id INT REFERENCES categories_materiaux(id),
    ADD COLUMN poids_kg NUMERIC(8, 2) CHECK (poids_kg > 0),
    -- vrai quand le poids n'a pas été saisi mais déduit (dépôts existants avant cette migration)
    ADD COLUMN poids_estime BOOLEAN NOT NULL DEFAULT false;

-- Facteur d'impact par matériau, modifiable par l'administrateur (page Catégories)
ALTER TABLE categories_materiaux
    ADD COLUMN co2_evite_par_kg NUMERIC(6, 3) CHECK (co2_evite_par_kg >= 0),
    ADD COLUMN source_facteur TEXT;

-- Dépôts existants : leur description permet de retrouver exactement le matériau ;
-- le poids est un poids type par objet, marqué comme estimé.
UPDATE depots d
SET categorie_id = c.id, poids_kg = o.poids, poids_estime = true
FROM (VALUES
    ('Lot de planches en chêne', 'bois', 25), ('Chaises bistrot à restaurer', 'bois', 10), ('Palettes', 'bois', 20),
    ('Portes anciennes', 'bois', 25), ('Caisses de vin', 'bois', 4),
    ('Radiateur en fonte', 'metal', 40), ('Tôles de zinc', 'metal', 10), ('Pieds de table en acier', 'metal', 8),
    ('Jantes de vélo', 'metal', 3), ('Chutes de cuivre', 'metal', 5),
    ('Rouleaux de tissu ancien', 'textile', 6), ('Chutes de cuir', 'textile', 3), ('Jeans usés', 'textile', 3),
    ('Rideaux en lin', 'textile', 3), ('Sacs de laine', 'textile', 5),
    ('Bocaux en verre', 'verre', 5), ('Portes vitrées', 'verre', 20), ('Bouteilles colorées', 'verre', 6),
    ('Miroirs anciens', 'verre', 8), ('Vitraux', 'verre', 10),
    ('Bidons propres', 'plastique', 3), ('Caisses en plastique', 'plastique', 4), ('Bouchons triés', 'plastique', 2),
    ('Bâches', 'plastique', 4), ('Jouets à réparer', 'plastique', 3)
) AS o(description, code, poids)
JOIN categories_materiaux c ON c.code = o.code
WHERE d.description_objet = o.description AND d.categorie_id IS NULL;

-- 3. Alertes de collecte (Premium) : être prévenu quand un objet correspondant est déposé
CREATE TABLE alertes_collecte (
    id SERIAL PRIMARY KEY,
    utilisateur_id INT NOT NULL REFERENCES utilisateurs(id) ON DELETE CASCADE,
    categorie_id INT REFERENCES categories_materiaux(id),
    site VARCHAR(50),
    mot_cle VARCHAR(100),
    date_creation TIMESTAMP NOT NULL DEFAULT now(),
    CHECK (categorie_id IS NOT NULL OR site IS NOT NULL OR mot_cle IS NOT NULL)
);
CREATE INDEX alertes_collecte_utilisateur ON alertes_collecte (utilisateur_id);

-- 4. Notifications affichées dans le site (cloche), en plus du push OneSignal
CREATE TABLE notifications_utilisateurs (
    id SERIAL PRIMARY KEY,
    utilisateur_id INT NOT NULL REFERENCES utilisateurs(id) ON DELETE CASCADE,
    titre VARCHAR(150) NOT NULL,
    message TEXT NOT NULL,
    lien VARCHAR(255),
    lue BOOLEAN NOT NULL DEFAULT false,
    date_creation TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX notifications_utilisateurs_liste ON notifications_utilisateurs (utilisateur_id, date_creation DESC);

-- 5. Abonnement Premium et campagnes publicitaires achetés comme les autres commandes
ALTER TABLE commandes DROP CONSTRAINT commandes_type_check;
ALTER TABLE commandes ADD CONSTRAINT commandes_type_check
    CHECK (type IN ('formation', 'vente', 'abonnement', 'publicite'));

-- Une campagne est « en attente » le temps de son paiement
ALTER TABLE campagnes_publicitaires DROP CONSTRAINT campagnes_publicitaires_statut_check;
ALTER TABLE campagnes_publicitaires ADD CONSTRAINT campagnes_publicitaires_statut_check
    CHECK (statut IN ('en_attente', 'active', 'terminee', 'annulee'));
