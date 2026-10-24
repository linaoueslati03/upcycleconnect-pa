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

INSERT INTO categories_materiaux (code, libelle) VALUES
    ('bois', 'Bois'), ('metal', 'Métal'), ('textile', 'Textile'), ('verre', 'Verre'), ('plastique', 'Plastique');

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
    SELECT 'Bien trier son bois', 'Séparer le bois brut du bois traité avant tout dépôt.', 'conseils', 'publie', utilisateur_id
    FROM salaries LIMIT 1;

INSERT INTO categories_prestations (libelle) VALUES ('Conseil'), ('Logistique'), ('Formation');

INSERT INTO prestations (titre, categorie_id, tarif, statut)
    SELECT p.titre, c.id, p.tarif, p.statut
    FROM (VALUES
        ('Audit recyclage', 'Conseil', 450.00, 'publie'),
        ('Collecte de bois', 'Logistique', 250.00, 'publie'),
        ('Formation de tri', 'Formation', 800.00, 'brouillon')
    ) AS p(titre, categorie, tarif, statut)
    JOIN categories_prestations c ON c.libelle = p.categorie;

INSERT INTO forum_sujets (titre, auteur_utilisateur_id)
    SELECT 'Où trouver des palettes à Paris ?', id FROM utilisateurs WHERE email = 'particulier@upcycleconnect.fr';

INSERT INTO forum_messages (sujet_id, auteur_utilisateur_id, contenu)
    SELECT s.id, u.id, 'Je cherche des palettes pour fabriquer une table basse, vous avez des adresses ?'
    FROM forum_sujets s, utilisateurs u
    WHERE s.titre = 'Où trouver des palettes à Paris ?' AND u.email = 'particulier@upcycleconnect.fr';

INSERT INTO forum_messages (sujet_id, auteur_utilisateur_id, contenu)
    SELECT s.id, u.id, 'Achetez vos palettes neuves sur mon site, livraison gratuite !!!'
    FROM forum_sujets s, utilisateurs u
    WHERE s.titre = 'Où trouver des palettes à Paris ?' AND u.email = 'pro@upcycleconnect.fr';

INSERT INTO forum_signalements (message_id, signale_par_utilisateur_id, motif)
    SELECT m.id, u.id, 'Publicité sans rapport avec l''upcycling'
    FROM forum_messages m, utilisateurs u
    WHERE m.contenu LIKE 'Achetez vos palettes%' AND u.email = 'particulier@upcycleconnect.fr';
