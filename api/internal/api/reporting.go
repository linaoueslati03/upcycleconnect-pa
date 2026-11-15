package api

import (
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
