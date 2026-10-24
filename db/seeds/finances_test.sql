DO $$ BEGIN PERFORM setseed(0.42); END $$;

INSERT INTO utilisateurs (role_id, nom, prenom, email, mot_de_passe_hash, langue_preferee_id, abonnement, date_creation)
    SELECT r.id, u.nom, u.prenom, u.email,
        '$2a$10$HDK2AXzOOGuvVHeeI759R.Xz2oREQp2BAQRL381FMBAvIrXZZGqGO', 1, u.abonnement,
        now() - make_interval(months => u.anciennete)
    FROM (VALUES
        ('professionnel', 'Lefebvre', 'Claire', 'atelier.lefebvre@upcycleconnect.fr', 'premium', 13),
        ('professionnel', 'Moreau', 'Julien', 'ebenisterie.moreau@upcycleconnect.fr', 'premium', 11),
        ('professionnel', 'Garnier', 'Inès', 'garnier.design@upcycleconnect.fr', 'premium', 8),
        ('professionnel', 'Roux', 'Thomas', 'roux.metal@upcycleconnect.fr', 'freemium', 10),
        ('professionnel', 'Fontaine', 'Léa', 'fontaine.textile@upcycleconnect.fr', 'freemium', 6),
        ('particulier', 'Petit', 'Emma', 'emma.petit@upcycleconnect.fr', 'freemium', 12),
        ('particulier', 'Robert', 'Hugo', 'hugo.robert@upcycleconnect.fr', 'freemium', 12),
        ('particulier', 'Richard', 'Chloé', 'chloe.richard@upcycleconnect.fr', 'freemium', 11),
        ('particulier', 'Dubois', 'Nathan', 'nathan.dubois@upcycleconnect.fr', 'freemium', 10),
        ('particulier', 'Laurent', 'Manon', 'manon.laurent@upcycleconnect.fr', 'freemium', 9),
        ('particulier', 'Simon', 'Louis', 'louis.simon@upcycleconnect.fr', 'freemium', 8),
        ('particulier', 'Michel', 'Camille', 'camille.michel@upcycleconnect.fr', 'freemium', 7),
        ('particulier', 'Lemoine', 'Arthur', 'arthur.lemoine@upcycleconnect.fr', 'freemium', 5),
        ('particulier', 'Mercier', 'Jade', 'jade.mercier@upcycleconnect.fr', 'freemium', 4),
        ('particulier', 'Blanc', 'Raphaël', 'raphael.blanc@upcycleconnect.fr', 'freemium', 2)
    ) AS u(role, nom, prenom, email, abonnement, anciennete)
    JOIN roles r ON r.code = u.role;

INSERT INTO abonnements (utilisateur_id, montant_mensuel, statut, stripe_abonnement_id, date_debut, date_fin)
    SELECT u.id, a.montant, a.statut, 'sub_test_' || u.id || '_' || a.debut,
        (date_trunc('month', now()) - make_interval(months => a.debut))::date,
        CASE WHEN a.fin IS NULL THEN NULL
             ELSE (date_trunc('month', now()) - make_interval(months => a.fin) - interval '1 day')::date END
    FROM (VALUES
        ('atelier.lefebvre@upcycleconnect.fr', 20.00, 'actif', 11, NULL),
        ('ebenisterie.moreau@upcycleconnect.fr', 15.00, 'resilie', 10, 6),
        ('ebenisterie.moreau@upcycleconnect.fr', 20.00, 'actif', 4, NULL),
        ('garnier.design@upcycleconnect.fr', 20.00, 'actif', 7, NULL),
        ('roux.metal@upcycleconnect.fr', 20.00, 'resilie', 9, 5)
    ) AS a(email, montant, statut, debut, fin)
    JOIN utilisateurs u ON u.email = a.email;

INSERT INTO paiements (utilisateur_id, source, reference_id, montant, stripe_id, date_paiement)
    SELECT a.utilisateur_id, 'abonnement', a.id, a.montant_mensuel,
        'pi_test_abo_' || a.id || '_' || to_char(echeance, 'YYYYMM'),
        echeance + interval '9 hours'
    FROM abonnements a
    CROSS JOIN LATERAL generate_series(
        a.date_debut::timestamp,
        COALESCE(a.date_fin, CURRENT_DATE)::timestamp,
        interval '1 month'
    ) AS echeance
    WHERE a.stripe_abonnement_id LIKE 'sub_test_%';

INSERT INTO formations (titre, description, date_debut, lieu, nb_places, tarif, statut, createur_id)
    SELECT f.titre, 'Session passée utilisée pour les données financières de test.',
        date_trunc('month', now()) - make_interval(months => f.mois) + interval '14 days 10 hours',
        f.lieu, 15, f.tarif, 'termine', s.utilisateur_id
    FROM (VALUES
        ('Upcycling textile : les bases', 'Paris 11e', 25.00, 11),
        ('Restaurer un meuble ancien', 'Montreuil', 60.00, 10),
        ('Fabriquer avec des palettes', 'Paris 13e', 35.00, 9),
        ('Soudure et récupération de métal', 'Montreuil', 90.00, 8),
        ('Couture zéro déchet', 'Paris 11e', 30.00, 7),
        ('Luminaires en matériaux récupérés', 'Paris 16e', 45.00, 6),
        ('Réparer son petit électroménager', 'Ivry', 40.00, 5),
        ('Tapisserie de sièges', 'Bourg-la-Reine', 75.00, 4),
        ('Mosaïque en verre recyclé', 'Paris 13e', 50.00, 3),
        ('Initiation à l''ébénisterie durable', 'Montreuil', 100.00, 2),
        ('Créer sa marque upcyclée', 'Paris 10e', 80.00, 1)
    ) AS f(titre, lieu, tarif, mois)
    CROSS JOIN (SELECT utilisateur_id FROM salaries LIMIT 1) s;

INSERT INTO inscriptions (utilisateur_id, type_offre, offre_id, date_inscription, statut_paiement, montant_paye)
    SELECT u.id, 'formation', f.id,
        f.date_debut - make_interval(days => 3 + floor(random() * 20)::int),
        'paye', f.tarif
    FROM formations f
    CROSS JOIN utilisateurs u
    JOIN roles r ON r.id = u.role_id
    WHERE f.statut = 'termine'
      AND f.description = 'Session passée utilisée pour les données financières de test.'
      AND r.code IN ('particulier', 'professionnel')
      AND u.date_creation < f.date_debut
      AND random() < 0.55;

INSERT INTO paiements (utilisateur_id, source, reference_id, montant, stripe_id, date_paiement)
    SELECT i.utilisateur_id, 'formation', i.offre_id, i.montant_paye,
        'pi_test_ins_' || i.id, i.date_inscription
    FROM inscriptions i
    WHERE i.type_offre = 'formation' AND i.statut_paiement = 'paye';

INSERT INTO annonces (utilisateur_id, titre, description, type, prix, categorie_id, localisation, statut, date_creation)
    SELECT vendeur.id,
        (ARRAY['Lot de planches en chêne', 'Chaises bistrot à restaurer', 'Chutes de cuir', 'Radiateur en fonte',
               'Rouleaux de tissu ancien', 'Portes vitrées', 'Tôles de zinc', 'Bocaux en verre'])[1 + (n % 8)],
        'Annonce vendue, utilisée pour les données financières de test.',
        'vente',
        round((15 + random() * 285)::numeric, 0),
        c.id,
        (ARRAY['Paris 11e', 'Paris 13e', 'Montreuil', 'Ivry', 'Bourg-la-Reine'])[1 + (n % 5)],
        'cedee',
        LEAST(now() - interval '5 days',
            date_trunc('month', now()) - make_interval(months => n % 12) + make_interval(days => floor(random() * 25)::int))
    FROM generate_series(1, 140) AS n
    CROSS JOIN LATERAL (
        SELECT u.id FROM utilisateurs u JOIN roles r ON r.id = u.role_id
        WHERE r.code = 'particulier' AND n > 0
        ORDER BY random() LIMIT 1
    ) AS vendeur
    CROSS JOIN LATERAL (
        SELECT id FROM categories_materiaux WHERE n > 0 ORDER BY random() LIMIT 1
    ) AS c;

INSERT INTO paiements (utilisateur_id, source, reference_id, montant, taux_commission, montant_commission,
                       statut, stripe_id, date_paiement)
    SELECT acheteur.id, 'vente', a.id, a.prix, taux.valeur,
        round(a.prix * taux.valeur / 100, 2),
        CASE WHEN random() < 0.04 THEN 'rembourse' ELSE 'reussi' END,
        'pi_test_vente_' || a.id,
        LEAST(now() - interval '1 hour',
            a.date_creation + make_interval(days => 1 + floor(random() * 4)::int, hours => floor(random() * 10)::int))
    FROM annonces a
    CROSS JOIN LATERAL (
        SELECT CASE WHEN a.date_creation < date_trunc('month', now()) - interval '5 months' THEN 7.00 ELSE 8.00 END AS valeur
    ) AS taux
    CROSS JOIN LATERAL (
        SELECT u.id FROM utilisateurs u JOIN roles r ON r.id = u.role_id
        WHERE r.code = 'professionnel' AND a.id > 0
        ORDER BY random() LIMIT 1
    ) AS acheteur
    WHERE a.description = 'Annonce vendue, utilisée pour les données financières de test.';

INSERT INTO paiements (utilisateur_id, source, reference_id, montant, statut, stripe_id, date_paiement)
    SELECT u.id, 'abonnement', NULL, 20.00, 'echoue', 'pi_test_echec_' || n,
        date_trunc('month', now()) - make_interval(months => n * 3) + interval '9 hours'
    FROM generate_series(1, 3) AS n
    JOIN utilisateurs u ON u.email = 'fontaine.textile@upcycleconnect.fr';

INSERT INTO projets_upcycling (utilisateur_id, titre, description, partage_public, sponsorise, date_creation)
    SELECT u.id, p.titre, p.description, TRUE, p.sponsorise, now() - make_interval(months => p.mois)
    FROM (VALUES
        ('atelier.lefebvre@upcycleconnect.fr', 'Bibliothèque en portes anciennes',
            'Trois portes de récupération transformées en bibliothèque.', TRUE, 10),
        ('ebenisterie.moreau@upcycleconnect.fr', 'Table basse en chêne de charpente',
            'Poutres de charpente rabotées et assemblées.', FALSE, 9),
        ('garnier.design@upcycleconnect.fr', 'Suspensions en bocaux',
            'Luminaires réalisés à partir de bocaux en verre.', TRUE, 5)
    ) AS p(email, titre, description, sponsorise, mois)
    JOIN utilisateurs u ON u.email = p.email;

INSERT INTO campagnes_publicitaires (utilisateur_id, cible_type, cible_id, montant_mensuel, statut, date_debut, date_fin)
    SELECT p.utilisateur_id, 'projet', p.id, c.montant, c.statut,
        (date_trunc('month', now()) - make_interval(months => c.debut))::date,
        CASE WHEN c.fin IS NULL THEN NULL
             ELSE (date_trunc('month', now()) - make_interval(months => c.fin) - interval '1 day')::date END
    FROM (VALUES
        ('Bibliothèque en portes anciennes', 250.00, 'active', 9, NULL),
        ('Table basse en chêne de charpente', 150.00, 'terminee', 8, 5),
        ('Suspensions en bocaux', 400.00, 'active', 3, NULL)
    ) AS c(titre, montant, statut, debut, fin)
    JOIN projets_upcycling p ON p.titre = c.titre;

INSERT INTO paiements (utilisateur_id, source, reference_id, montant, stripe_id, date_paiement)
    SELECT c.utilisateur_id, 'publicite', c.id, c.montant_mensuel,
        'pi_test_pub_' || c.id || '_' || to_char(echeance, 'YYYYMM'),
        echeance + interval '10 hours'
    FROM campagnes_publicitaires c
    CROSS JOIN LATERAL generate_series(
        c.date_debut::timestamp,
        COALESCE(c.date_fin, CURRENT_DATE)::timestamp,
        interval '1 month'
    ) AS echeance;

INSERT INTO charges (libelle, categorie, montant, date_charge, auteur_id)
    SELECT ch.libelle, ch.categorie, round((ch.montant * (0.95 + random() * 0.1))::numeric, 2),
        (date_trunc('month', now()) - make_interval(months => m))::date + ch.jour,
        (SELECT id FROM utilisateurs WHERE email = 'admin@upcycleconnect.fr')
    FROM generate_series(0, 11) AS m
    CROSS JOIN (VALUES
        ('Salaires des animateurs et formateurs', 'salaires', 520.00, 27),
        ('Loyer des locaux', 'loyers', 300.00, 4),
        ('Hébergement et serveurs', 'infrastructure', 60.00, 9),
        ('Licences logicielles', 'infrastructure', 25.00, 9)
    ) AS ch(libelle, categorie, montant, jour)
    WHERE (date_trunc('month', now()) - make_interval(months => m))::date + ch.jour <= CURRENT_DATE;

INSERT INTO charges (libelle, categorie, montant, date_charge, auteur_id)
    SELECT ch.libelle, ch.categorie, ch.montant,
        (date_trunc('month', now()) - make_interval(months => ch.mois))::date + 15,
        (SELECT id FROM utilisateurs WHERE email = 'admin@upcycleconnect.fr')
    FROM (VALUES
        ('Campagne réseaux sociaux', 'marketing', 180.00, 9),
        ('Salon de l''économie circulaire', 'marketing', 350.00, 6),
        ('Achat de conteneurs', 'autre', 900.00, 8),
        ('Campagne réseaux sociaux', 'marketing', 200.00, 2)
    ) AS ch(libelle, categorie, montant, mois);
