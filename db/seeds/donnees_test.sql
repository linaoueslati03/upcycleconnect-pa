-- Données de démonstration, chargées uniquement à la création de la base locale.
-- Mot de passe de tous les comptes : Test1234!
-- Le salarié est créé en premier (id 1) : le front Salarié envoie encore createur_id = 1 en dur.

INSERT INTO utilisateurs (role_id, nom, prenom, email, mot_de_passe_hash, langue_preferee_id) VALUES
    ((SELECT id FROM roles WHERE code = 'salarie'), 'Martin', 'Sophie', 'salarie@upcycleconnect.fr',
     '$2a$10$HDK2AXzOOGuvVHeeI759R.Xz2oREQp2BAQRL381FMBAvIrXZZGqGO', 1),
    ((SELECT id FROM roles WHERE code = 'administrateur'), 'Admin', 'Compte', 'admin@upcycleconnect.fr',
     '$2a$10$HDK2AXzOOGuvVHeeI759R.Xz2oREQp2BAQRL381FMBAvIrXZZGqGO', 1),
    ((SELECT id FROM roles WHERE code = 'particulier'), 'Durand', 'Paul', 'particulier@upcycleconnect.fr',
     '$2a$10$HDK2AXzOOGuvVHeeI759R.Xz2oREQp2BAQRL381FMBAvIrXZZGqGO', 1),
    ((SELECT id FROM roles WHERE code = 'professionnel'), 'Bernard', 'Luc', 'pro@upcycleconnect.fr',
     '$2a$10$HDK2AXzOOGuvVHeeI759R.Xz2oREQp2BAQRL381FMBAvIrXZZGqGO', 1);

INSERT INTO salaries (utilisateur_id, est_responsable)
    SELECT id, TRUE FROM utilisateurs WHERE email = 'salarie@upcycleconnect.fr';

INSERT INTO categories_materiaux (code) VALUES ('bois'), ('metal'), ('textile'), ('verre'), ('plastique');

INSERT INTO conteneurs (site) VALUES ('Paris 11e'), ('Paris 13e'), ('Montreuil');

INSERT INTO formations (titre, description, date_debut, lieu, nb_places, tarif, statut, createur_id)
    SELECT 'Initiation à la menuiserie de récup', 'Fabriquer un meuble à partir de palettes.',
           now() + interval '7 days', 'Paris 11e', 12, 25.00, 'publie', utilisateur_id
    FROM salaries LIMIT 1;

INSERT INTO ateliers (titre, description, date_debut, lieu, statut, createur_id)
    SELECT 'Réparer ses vêtements', 'Couture de base et techniques de raccommodage.',
           now() + interval '10 days', 'Montreuil', 'publie', utilisateur_id
    FROM salaries LIMIT 1;

INSERT INTO evenements (titre, description, date_debut, lieu, site, statut, createur_id)
    SELECT 'Repair Café', 'Venez réparer vos objets avec nos bénévoles.',
           now() + interval '14 days', '12 rue de la Roquette', 'Paris 11e', 'publie', utilisateur_id
    FROM salaries LIMIT 1;

INSERT INTO conseils (titre, contenu, categorie, statut, auteur_id)
    SELECT 'Bien trier son bois', 'Séparer le bois brut du bois traité avant tout dépôt.', 'bois', 'publie', utilisateur_id
    FROM salaries LIMIT 1;

INSERT INTO prestations (titre, categorie, tarif, statut) VALUES
    ('Audit recyclage', 'Conseil', 450.00, 'publie'),
    ('Collecte de bois', 'Logistique', 250.00, 'publie'),
    ('Formation de tri', 'Formation', 800.00, 'brouillon');
