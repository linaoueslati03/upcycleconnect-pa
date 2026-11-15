package api

import (
	"math"
	"net/http"
	"slices"
)

const (
	filtreActeurs = `FROM utilisateurs u
		JOIN roles r ON r.id = u.role_id
		WHERE r.code IN ('particulier', 'professionnel')`
	filtrePeriodeDemandes = `d.date_demande >= $1::timestamp AND d.date_demande < $2::timestamp + interval '1 day'`
)

var (
	rolesActeurs       = []string{roleParticulier, roleProfessionnel}
	statutsUtilisateur = []string{"actif", "suspendu"}
	abonnementsPro     = []string{"freemium", "premium"}
	tranchesScore      = []string{"0", "1-24", "25-49", "50-99", "100-199", "200+"}
	sitesReporting     = []string{"Paris 10e", "Paris 11e", "Paris 13e", "Paris 16e", "Bourg-la-Reine", "Ivry", "Montreuil"}
)

type Comptage struct {
	Libelle string `json:"libelle"`
	Nombre  int    `json:"nombre"`
}

type TotauxActeurs struct {
	Acteurs         int `json:"acteurs"`
	Particuliers    int `json:"particuliers"`
	Professionnels  int `json:"professionnels"`
	NouveauxPeriode int `json:"nouveaux_periode"`
	ActifsPeriode   int `json:"actifs_periode"`
}

type NouveauxMois struct {
	Mois           string `json:"mois"`
	Particuliers   int    `json:"particuliers"`
	Professionnels int    `json:"professionnels"`
}

type ActiviteSite struct {
	Site     string `json:"site"`
	Demandes int    `json:"demandes"`
	Acteurs  int    `json:"acteurs"`
}

type ReportingActeurs struct {
	Debut             string         `json:"debut"`
	Fin               string         `json:"fin"`
	Totaux            TotauxActeurs  `json:"totaux"`
	ParRole           []Comptage     `json:"par_role"`
	ParStatut         []Comptage     `json:"par_statut"`
	ParAbonnement     []Comptage     `json:"par_abonnement"`
	DistributionScore []Comptage     `json:"distribution_score"`
	NouveauxParMois   []NouveauxMois `json:"nouveaux_par_mois"`
	ActiviteParSite   []ActiviteSite `json:"activite_par_site"`
}

func (s *Serveur) compterParLibelle(libelles []string, requete string, args ...any) ([]Comptage, error) {
	lignes, err := s.db.Query(requete, args...)
	if err != nil {
		return nil, err
	}
	defer lignes.Close()

	nombres := make(map[string]int)
	var autres []Comptage
	for lignes.Next() {
		var c Comptage
		if err := lignes.Scan(&c.Libelle, &c.Nombre); err != nil {
			return nil, err
		}
		nombres[c.Libelle] = c.Nombre
		if !slices.Contains(libelles, c.Libelle) {
			autres = append(autres, c)
		}
	}
	if err := lignes.Err(); err != nil {
		return nil, err
	}

	resultat := make([]Comptage, 0, len(libelles)+len(autres))
	for _, libelle := range libelles {
		resultat = append(resultat, Comptage{Libelle: libelle, Nombre: nombres[libelle]})
	}
	return append(resultat, autres...), nil
}

func (s *Serveur) gererReportingActeurs(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	debut, fin, message := periodeDepuisRequete(r)
	if message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}

	reporting := ReportingActeurs{Debut: debut, Fin: fin}
	var err error

	err = s.db.QueryRow(`
		SELECT count(*),
		       count(*) FILTER (WHERE r.code = 'particulier'),
		       count(*) FILTER (WHERE r.code = 'professionnel'),
		       count(*) FILTER (WHERE u.date_creation >= $1::timestamp
		                          AND u.date_creation < $2::timestamp + interval '1 day'),
		       (SELECT count(DISTINCT d.utilisateur_id)
		        FROM vue_demandes_acteurs d
		        JOIN utilisateurs ua ON ua.id = d.utilisateur_id
		        JOIN roles ra ON ra.id = ua.role_id
		        WHERE ra.code IN ('particulier', 'professionnel') AND `+filtrePeriodeDemandes+`)
		`+filtreActeurs, debut, fin,
	).Scan(&reporting.Totaux.Acteurs, &reporting.Totaux.Particuliers, &reporting.Totaux.Professionnels,
		&reporting.Totaux.NouveauxPeriode, &reporting.Totaux.ActifsPeriode)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	reporting.ParRole, err = s.compterParLibelle(rolesActeurs,
		`SELECT r.code, count(*) `+filtreActeurs+` GROUP BY r.code`)
	if err == nil {
		reporting.ParStatut, err = s.compterParLibelle(statutsUtilisateur,
			`SELECT u.statut, count(*) `+filtreActeurs+` GROUP BY u.statut`)
	}
	if err == nil {
		reporting.ParAbonnement, err = s.compterParLibelle(abonnementsPro,
			`SELECT u.abonnement, count(*) `+filtreActeurs+` AND r.code = 'professionnel' GROUP BY u.abonnement`)
	}
	if err == nil {
		reporting.DistributionScore, err = s.compterParLibelle(tranchesScore, `
			SELECT CASE
			           WHEN u.upcycling_score <= 0 THEN '0'
			           WHEN u.upcycling_score < 25 THEN '1-24'
			           WHEN u.upcycling_score < 50 THEN '25-49'
			           WHEN u.upcycling_score < 100 THEN '50-99'
			           WHEN u.upcycling_score < 200 THEN '100-199'
			           ELSE '200+'
			       END, count(*)
			`+filtreActeurs+` GROUP BY 1`)
	}
	if err == nil {
		reporting.NouveauxParMois, err = s.nouveauxActeursParMois(debut, fin)
	}
	if err == nil {
		reporting.ActiviteParSite, err = s.activiteParSite(debut, fin)
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, reporting)
}

func (s *Serveur) nouveauxActeursParMois(debut, fin string) ([]NouveauxMois, error) {
	lignes, err := s.db.Query(`
		SELECT to_char(m.debut, 'YYYY-MM'),
		       count(u.id) FILTER (WHERE r.code = 'particulier'),
		       count(u.id) FILTER (WHERE r.code = 'professionnel')
		FROM generate_series(date_trunc('month', $1::timestamp), $2::timestamp, interval '1 month') AS m(debut)
		LEFT JOIN (utilisateurs u JOIN roles r ON r.id = u.role_id AND r.code IN ('particulier', 'professionnel'))
		       ON u.date_creation >= greatest(m.debut, $1::timestamp)
		      AND u.date_creation < least(m.debut + interval '1 month', $2::timestamp + interval '1 day')
		GROUP BY m.debut
		ORDER BY m.debut`, debut, fin)
	if err != nil {
		return nil, err
	}
	defer lignes.Close()

	mois := make([]NouveauxMois, 0)
	for lignes.Next() {
		var m NouveauxMois
		if err := lignes.Scan(&m.Mois, &m.Particuliers, &m.Professionnels); err != nil {
			return nil, err
		}
		mois = append(mois, m)
	}
	return mois, lignes.Err()
}

func (s *Serveur) activiteParSite(debut, fin string) ([]ActiviteSite, error) {
	lignes, err := s.db.Query(`
		SELECT d.site, count(*), count(DISTINCT d.utilisateur_id)
		FROM vue_demandes_acteurs d
		JOIN utilisateurs u ON u.id = d.utilisateur_id
		JOIN roles r ON r.id = u.role_id
		WHERE r.code IN ('particulier', 'professionnel')
		  AND d.site IS NOT NULL
		  AND `+filtrePeriodeDemandes+`
		GROUP BY d.site`, debut, fin)
	if err != nil {
		return nil, err
	}
	defer lignes.Close()

	parSite := make(map[string]ActiviteSite)
	for lignes.Next() {
		var a ActiviteSite
		if err := lignes.Scan(&a.Site, &a.Demandes, &a.Acteurs); err != nil {
			return nil, err
		}
		parSite[a.Site] = a
	}
	if err := lignes.Err(); err != nil {
		return nil, err
	}

	sites := make([]ActiviteSite, 0, len(sitesReporting)+1)
	for _, site := range sitesReporting {
		a, trouve := parSite[site]
		if !trouve {
			a = ActiviteSite{Site: site}
		}
		sites = append(sites, a)
	}
	if autre, trouve := parSite["autre"]; trouve {
		sites = append(sites, autre)
	}
	return sites, nil
}

const ctePrestations = `
	WITH o AS (
		SELECT * FROM (
			SELECT 'formation'::text AS type, id, titre, date_debut, nb_places, statut FROM formations
			UNION ALL
			SELECT 'atelier', id, titre, date_debut, NULL::int, statut FROM ateliers
			UNION ALL
			SELECT 'evenement', id, titre, date_debut, NULL::int, statut FROM evenements
		) offres
		WHERE statut IN ('publie', 'termine') AND ($3::text = '' OR type = $3::text)
	),
	ins AS (
		SELECT o.type, o.titre, i.utilisateur_id, i.date_inscription
		FROM inscriptions i
		JOIN o ON o.type = i.type_offre AND o.id = i.offre_id
		WHERE i.date_inscription >= $1::timestamp AND i.date_inscription < $2::timestamp + interval '1 day'
	),
	participants AS (
		SELECT type, titre, utilisateur_id, count(*) AS nb
		FROM ins
		GROUP BY type, titre, utilisateur_id
	)`

var typesPrestation = []string{"formation", "atelier", "evenement"}

type StatsPrestation struct {
	Type              string   `json:"type"`
	Titre             string   `json:"titre,omitempty"`
	Sessions          int      `json:"sessions"`
	Inscriptions      int      `json:"inscriptions"`
	Participants      int      `json:"participants"`
	Reinscrits        int      `json:"reinscrits"`
	TauxReinscription *float64 `json:"taux_reinscription"`
	Places            int      `json:"places"`
	TauxRemplissage   *float64 `json:"taux_remplissage"`
	ChiffreAffaires   float64  `json:"chiffre_affaires"`
	inscritsSessions  int
}

type MoisPrestations struct {
	Mois      string `json:"mois"`
	Formation int    `json:"formation"`
	Atelier   int    `json:"atelier"`
	Evenement int    `json:"evenement"`
}

type ReportingPrestations struct {
	Debut       string            `json:"debut"`
	Fin         string            `json:"fin"`
	Type        string            `json:"type"`
	Totaux      StatsPrestation   `json:"totaux"`
	ParType     []StatsPrestation `json:"par_type"`
	ParMois     []MoisPrestations `json:"par_mois"`
	Prestations []StatsPrestation `json:"prestations"`
}

func pourcentage(partie, total int) *float64 {
	if total == 0 {
		return nil
	}
	valeur := math.Round(float64(partie)*1000/float64(total)) / 10
	return &valeur
}

func (p *StatsPrestation) calculerTaux() {
	p.ChiffreAffaires = arrondirCentimes(p.ChiffreAffaires)
	p.TauxReinscription = pourcentage(p.Reinscrits, p.Participants)
	p.TauxRemplissage = pourcentage(p.inscritsSessions, p.Places)
}

func (s *Serveur) gererReportingPrestations(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	debut, fin, message := periodeDepuisRequete(r)
	if message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}
	typeOffre := r.URL.Query().Get("type")
	if typeOffre != "" && !slices.Contains(typesPrestation, typeOffre) {
		envoyerErreur(w, http.StatusBadRequest, "type invalide (formation, atelier ou evenement attendu)")
		return
	}

	reporting := ReportingPrestations{Debut: debut, Fin: fin, Type: typeOffre, Totaux: StatsPrestation{Type: typeOffre}}
	types := typesPrestation
	if typeOffre != "" {
		types = []string{typeOffre}
	}

	prestations, err := s.statsParPrestation(debut, fin, typeOffre)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	parType := make(map[string]*StatsPrestation)
	for _, t := range types {
		parType[t] = &StatsPrestation{Type: t}
	}
	for _, p := range prestations {
		cumul := parType[p.Type]
		cumul.Sessions += p.Sessions
		cumul.Places += p.Places
		cumul.inscritsSessions += p.inscritsSessions
		cumul.ChiffreAffaires += p.ChiffreAffaires
		reporting.Totaux.Sessions += p.Sessions
		reporting.Totaux.Places += p.Places
		reporting.Totaux.inscritsSessions += p.inscritsSessions
		reporting.Totaux.ChiffreAffaires += p.ChiffreAffaires
	}

	if err := s.frequentationParType(debut, fin, typeOffre, parType, &reporting.Totaux); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	reporting.ParMois, err = s.inscriptionsParMois(debut, fin, typeOffre)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	reporting.Totaux.calculerTaux()
	reporting.ParType = make([]StatsPrestation, 0, len(types))
	for _, t := range types {
		parType[t].calculerTaux()
		reporting.ParType = append(reporting.ParType, *parType[t])
	}
	reporting.Prestations = prestations

	envoyerJSON(w, http.StatusOK, reporting)
}

func (s *Serveur) statsParPrestation(debut, fin, typeOffre string) ([]StatsPrestation, error) {
	lignes, err := s.db.Query(ctePrestations+`,
		sessions AS (
			SELECT o.type, o.titre, count(*) AS sessions,
			       COALESCE(sum(o.nb_places), 0) AS places,
			       COALESCE(sum(i.nb) FILTER (WHERE o.nb_places IS NOT NULL), 0) AS inscrits_sessions
			FROM o
			LEFT JOIN (
				SELECT type_offre, offre_id, count(*) AS nb FROM inscriptions GROUP BY type_offre, offre_id
			) i ON i.type_offre = o.type AND i.offre_id = o.id
			WHERE o.date_debut >= $1::timestamp AND o.date_debut < $2::timestamp + interval '1 day'
			GROUP BY o.type, o.titre
		),
		frequentation AS (
			SELECT type, titre, sum(nb) AS inscriptions, count(*) AS participants,
			       count(*) FILTER (WHERE nb >= 2) AS reinscrits
			FROM participants
			GROUP BY type, titre
		),
		ventes AS (
			SELECT o.type, o.titre, sum(p.montant) AS chiffre_affaires
			FROM paiements p
			JOIN o ON p.source = o.type AND p.reference_id = o.id
			WHERE p.statut = 'reussi'
			  AND p.date_paiement >= $1::timestamp AND p.date_paiement < $2::timestamp + interval '1 day'
			GROUP BY o.type, o.titre
		)
		SELECT t.type, t.titre,
		       COALESCE(s.sessions, 0), COALESCE(f.inscriptions, 0), COALESCE(f.participants, 0),
		       COALESCE(f.reinscrits, 0), COALESCE(s.places, 0), COALESCE(s.inscrits_sessions, 0),
		       COALESCE(v.chiffre_affaires, 0)
		FROM (SELECT DISTINCT type, titre FROM o) t
		LEFT JOIN sessions s USING (type, titre)
		LEFT JOIN frequentation f USING (type, titre)
		LEFT JOIN ventes v USING (type, titre)
		WHERE s.sessions IS NOT NULL OR f.inscriptions IS NOT NULL
		ORDER BY COALESCE(f.inscriptions, 0) DESC, t.titre`, debut, fin, typeOffre)
	if err != nil {
		return nil, err
	}
	defer lignes.Close()

	prestations := make([]StatsPrestation, 0)
	for lignes.Next() {
		var p StatsPrestation
		if err := lignes.Scan(&p.Type, &p.Titre, &p.Sessions, &p.Inscriptions, &p.Participants,
			&p.Reinscrits, &p.Places, &p.inscritsSessions, &p.ChiffreAffaires); err != nil {
			return nil, err
		}
		p.calculerTaux()
		prestations = append(prestations, p)
	}
	return prestations, lignes.Err()
}

func (s *Serveur) frequentationParType(debut, fin, typeOffre string, parType map[string]*StatsPrestation, totaux *StatsPrestation) error {
	lignes, err := s.db.Query(ctePrestations+`
		SELECT COALESCE(type, ''), COALESCE(sum(nb), 0), count(DISTINCT utilisateur_id),
		       count(DISTINCT utilisateur_id) FILTER (WHERE nb >= 2)
		FROM participants
		GROUP BY GROUPING SETS ((type), ())`, debut, fin, typeOffre)
	if err != nil {
		return err
	}
	defer lignes.Close()

	for lignes.Next() {
		var t string
		var inscriptions, participants, reinscrits int
		if err := lignes.Scan(&t, &inscriptions, &participants, &reinscrits); err != nil {
			return err
		}
		cible := totaux
		if t != "" {
			cible = parType[t]
		}
		if cible != nil {
			cible.Inscriptions, cible.Participants, cible.Reinscrits = inscriptions, participants, reinscrits
		}
	}
	return lignes.Err()
}

func (s *Serveur) inscriptionsParMois(debut, fin, typeOffre string) ([]MoisPrestations, error) {
	lignes, err := s.db.Query(ctePrestations+`
		SELECT to_char(m.debut, 'YYYY-MM'),
		       count(ins.type) FILTER (WHERE ins.type = 'formation'),
		       count(ins.type) FILTER (WHERE ins.type = 'atelier'),
		       count(ins.type) FILTER (WHERE ins.type = 'evenement')
		FROM generate_series(date_trunc('month', $1::timestamp), $2::timestamp, interval '1 month') AS m(debut)
		LEFT JOIN ins ON ins.date_inscription >= m.debut AND ins.date_inscription < m.debut + interval '1 month'
		GROUP BY m.debut
		ORDER BY m.debut`, debut, fin, typeOffre)
	if err != nil {
		return nil, err
	}
	defer lignes.Close()

	mois := make([]MoisPrestations, 0)
	for lignes.Next() {
		var m MoisPrestations
		if err := lignes.Scan(&m.Mois, &m.Formation, &m.Atelier, &m.Evenement); err != nil {
			return nil, err
		}
		mois = append(mois, m)
	}
	return mois, lignes.Err()
}
