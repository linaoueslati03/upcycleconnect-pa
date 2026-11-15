-- Bloc Sciences : données communes au reporting et au modèle de prédiction.
-- variables_acteurs(p_date) ne regarde que ce qui s'est passé avant p_date :
-- Python l'appelle avec une date de coupure pour l'entraînement,
-- Go lit vue_variables_acteurs (calculée à now()) pour prédire en temps réel.

CREATE OR REPLACE FUNCTION normaliser_site(texte TEXT)
RETURNS TEXT
LANGUAGE sql IMMUTABLE AS $$
    SELECT CASE
        WHEN texte IS NULL OR btrim(texte) = '' THEN NULL
        WHEN lower(texte) LIKE '%montreuil%' THEN 'Montreuil'
        WHEN lower(texte) LIKE '%ivry%' THEN 'Ivry'
        WHEN lower(texte) LIKE '%bourg%' THEN 'Bourg-la-Reine'
        WHEN lower(texte) ~ '(^|[^0-9])(750)?10([^0-9]|$)' THEN 'Paris 10e'
        WHEN lower(texte) ~ '(^|[^0-9])(750)?11([^0-9]|$)' THEN 'Paris 11e'
        WHEN lower(texte) ~ '(^|[^0-9])(750)?13([^0-9]|$)' THEN 'Paris 13e'
        WHEN lower(texte) ~ '(^|[^0-9])(750)?16([^0-9]|$)' THEN 'Paris 16e'
        ELSE 'autre'
    END
$$;

-- Une ligne par demande de prestation. Les 6 types sont les classes du modèle.
CREATE OR REPLACE VIEW vue_demandes_acteurs AS
    SELECT i.utilisateur_id, 'formation'::TEXT AS type_demande, i.date_inscription AS date_demande,
           normaliser_site(f.lieu) AS site, NULL::INT AS categorie_id
    FROM inscriptions i
    JOIN formations f ON f.id = i.offre_id
    WHERE i.type_offre = 'formation'
UNION ALL
    SELECT i.utilisateur_id, 'atelier', i.date_inscription, normaliser_site(a.lieu), NULL
    FROM inscriptions i
    JOIN ateliers a ON a.id = i.offre_id
    WHERE i.type_offre = 'atelier'
UNION ALL
    SELECT i.utilisateur_id, 'evenement', i.date_inscription, normaliser_site(COALESCE(e.site, e.lieu)), NULL
    FROM inscriptions i
    JOIN evenements e ON e.id = i.offre_id
    WHERE i.type_offre = 'evenement'
UNION ALL
    SELECT d.utilisateur_id, 'depot', d.date_demande, normaliser_site(c.site), NULL
    FROM depots d
    LEFT JOIN conteneurs c ON c.id = d.conteneur_id
UNION ALL
    -- Pour une vente, la personne qui paie est l'acheteur.
    SELECT p.utilisateur_id, 'achat_annonce', p.date_paiement, normaliser_site(a.localisation), a.categorie_id
    FROM paiements p
    LEFT JOIN annonces a ON a.id = p.reference_id
    WHERE p.source = 'vente' AND p.statut = 'reussi'
UNION ALL
    SELECT ab.utilisateur_id, 'premium', ab.date_debut::TIMESTAMP, NULL, NULL
    FROM abonnements ab;

CREATE OR REPLACE FUNCTION variables_acteurs(p_date TIMESTAMP)
RETURNS TABLE (
    utilisateur_id INT,
    role TEXT,
    premium INT,
    anciennete_jours INT,
    upcycling_score INT,
    nb_annonces INT,
    nb_annonces_vente INT,
    nb_formations INT,
    nb_ateliers INT,
    nb_evenements INT,
    nb_depots INT,
    nb_achats_annonces INT,
    montant_depense NUMERIC,
    jours_depuis_derniere_demande INT,
    derniere_demande TEXT,
    categorie_preferee INT,
    site_principal TEXT
)
LANGUAGE sql STABLE AS $$
    WITH acteurs AS (
        SELECT u.id, r.code AS role, u.date_creation, u.upcycling_score
        FROM utilisateurs u
        JOIN roles r ON r.id = u.role_id
        WHERE r.code IN ('particulier', 'professionnel')
          AND u.date_creation < p_date
    ),
    demandes AS (
        SELECT d.*
        FROM vue_demandes_acteurs d
        WHERE d.date_demande < p_date
    ),
    compteurs AS (
        SELECT d.utilisateur_id,
               COUNT(*) FILTER (WHERE d.type_demande = 'formation') AS nb_formations,
               COUNT(*) FILTER (WHERE d.type_demande = 'atelier') AS nb_ateliers,
               COUNT(*) FILTER (WHERE d.type_demande = 'evenement') AS nb_evenements,
               COUNT(*) FILTER (WHERE d.type_demande = 'depot') AS nb_depots,
               COUNT(*) FILTER (WHERE d.type_demande = 'achat_annonce') AS nb_achats_annonces,
               MAX(d.date_demande) AS derniere_date,
               mode() WITHIN GROUP (ORDER BY d.site) FILTER (WHERE d.site IS NOT NULL) AS site_principal
        FROM demandes d
        GROUP BY d.utilisateur_id
    ),
    derniere AS (
        SELECT DISTINCT ON (d.utilisateur_id) d.utilisateur_id, d.type_demande
        FROM demandes d
        ORDER BY d.utilisateur_id, d.date_demande DESC, d.type_demande
    ),
    annonces_avant AS (
        SELECT a.utilisateur_id,
               COUNT(*) AS nb_annonces,
               COUNT(*) FILTER (WHERE a.type = 'vente') AS nb_annonces_vente
        FROM annonces a
        WHERE a.date_creation < p_date
        GROUP BY a.utilisateur_id
    ),
    categories AS (
        SELECT c.utilisateur_id, mode() WITHIN GROUP (ORDER BY c.categorie_id) AS categorie_preferee
        FROM (
            SELECT a.utilisateur_id, a.categorie_id
            FROM annonces a
            WHERE a.date_creation < p_date AND a.categorie_id IS NOT NULL
            UNION ALL
            SELECT d.utilisateur_id, d.categorie_id
            FROM demandes d
            WHERE d.categorie_id IS NOT NULL
        ) c
        GROUP BY c.utilisateur_id
    ),
    -- Score à p_date = score actuel moins ce qui a été gagné depuis.
    score_depuis AS (
        SELECT s.utilisateur_id, SUM(s.delta) AS gagne_depuis
        FROM score_historique s
        WHERE s.date >= p_date
        GROUP BY s.utilisateur_id
    ),
    depenses AS (
        SELECT p.utilisateur_id, SUM(p.montant) AS montant_depense
        FROM paiements p
        WHERE p.statut = 'reussi' AND p.date_paiement < p_date
        GROUP BY p.utilisateur_id
    ),
    -- utilisateurs.abonnement ne donne que l'état actuel : on lit l'historique des abonnements.
    premium AS (
        SELECT DISTINCT ab.utilisateur_id
        FROM abonnements ab
        WHERE ab.date_debut <= p_date::DATE
          AND (ab.date_fin > p_date::DATE OR (ab.date_fin IS NULL AND ab.statut = 'actif'))
    )
    SELECT a.id,
           a.role::TEXT,
           (pr.utilisateur_id IS NOT NULL)::INT,
           (p_date::DATE - a.date_creation::DATE),
           (a.upcycling_score - COALESCE(sd.gagne_depuis, 0))::INT,
           COALESCE(an.nb_annonces, 0)::INT,
           COALESCE(an.nb_annonces_vente, 0)::INT,
           COALESCE(c.nb_formations, 0)::INT,
           COALESCE(c.nb_ateliers, 0)::INT,
           COALESCE(c.nb_evenements, 0)::INT,
           COALESCE(c.nb_depots, 0)::INT,
           COALESCE(c.nb_achats_annonces, 0)::INT,
           COALESCE(dp.montant_depense, 0),
           COALESCE(p_date::DATE - c.derniere_date::DATE, 999),
           COALESCE(dn.type_demande, 'aucune'),
           cat.categorie_preferee,
           COALESCE(c.site_principal, 'inconnu')
    FROM acteurs a
    LEFT JOIN compteurs c ON c.utilisateur_id = a.id
    LEFT JOIN derniere dn ON dn.utilisateur_id = a.id
    LEFT JOIN annonces_avant an ON an.utilisateur_id = a.id
    LEFT JOIN categories cat ON cat.utilisateur_id = a.id
    LEFT JOIN score_depuis sd ON sd.utilisateur_id = a.id
    LEFT JOIN depenses dp ON dp.utilisateur_id = a.id
    LEFT JOIN premium pr ON pr.utilisateur_id = a.id
$$;

CREATE OR REPLACE VIEW vue_variables_acteurs AS
    SELECT * FROM variables_acteurs(now()::TIMESTAMP);

CREATE INDEX IF NOT EXISTS inscriptions_utilisateur_date ON inscriptions (utilisateur_id, date_inscription);
CREATE INDEX IF NOT EXISTS depots_utilisateur_date ON depots (utilisateur_id, date_demande);
CREATE INDEX IF NOT EXISTS paiements_utilisateur_date ON paiements (utilisateur_id, date_paiement);
CREATE INDEX IF NOT EXISTS annonces_utilisateur ON annonces (utilisateur_id);
CREATE INDEX IF NOT EXISTS score_historique_utilisateur_date ON score_historique (utilisateur_id, date);
