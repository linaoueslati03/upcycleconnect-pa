package api

import (
	"net/http"
	"time"
)

// Tableau de bord du back-office : les KPI d'activité, de revenus et de charges sur une période,
// comparés à la période de même durée juste avant. Les chiffres financiers viennent du même calcul
// que la page Finances (resumeFinances), l'activité suit les mêmes règles que la page Reporting.

type IndicateursPeriode struct {
	Revenus         float64 `json:"revenus"`
	Charges         float64 `json:"charges"`
	Resultat        float64 `json:"resultat"`
	NouveauxActeurs int     `json:"nouveaux_acteurs"`
	ActeursActifs   int     `json:"acteurs_actifs"`
	Annonces        int     `json:"annonces"`
	Depots          int     `json:"depots"`
	Inscriptions    int     `json:"inscriptions"`
}

type ActiviteMois struct {
	Mois            string `json:"mois"`
	NouveauxActeurs int    `json:"nouveaux_acteurs"`
	Annonces        int    `json:"annonces"`
	Depots          int    `json:"depots"`
	Inscriptions    int    `json:"inscriptions"`
}

type ElementsATraiter struct {
	OffresAValider        int `json:"offres_a_valider"`
	AnnoncesAValider      int `json:"annonces_a_valider"`
	SujetsAModerer        int `json:"sujets_a_moderer"`
	Signalements          int `json:"signalements"`
	ConteneursMaintenance int `json:"conteneurs_maintenance"`
}

type TableauDeBord struct {
	Debut                  string             `json:"debut"`
	Fin                    string             `json:"fin"`
	DebutSeries            string             `json:"debut_series"`
	DebutPrecedent         string             `json:"debut_precedent"`
	FinPrecedent           string             `json:"fin_precedent"`
	Actuel                 IndicateursPeriode `json:"actuel"`
	Precedent              IndicateursPeriode `json:"precedent"`
	RevenuMensuelRecurrent float64            `json:"revenu_mensuel_recurrent"`
	AbonnesActifs          int                `json:"abonnes_actifs"`
	CampagnesActives       int                `json:"campagnes_actives"`
	RevenusParSource       []RevenuSource     `json:"revenus_par_source"`
	FinancesParMois        []BilanMois        `json:"finances_par_mois"`
	ActiviteParMois        []ActiviteMois     `json:"activite_par_mois"`
	ATraiter               ElementsATraiter   `json:"a_traiter"`
}

// Les graphiques mois par mois couvrent au moins les 12 mois qui finissent à la fin de la période,
// pour garder une tendance lisible même quand on regarde seulement les 30 derniers jours.
func debutDesSeries(debut, fin string) string {
	f, _ := time.Parse(formatDate, fin)
	douzeMois := time.Date(f.Year(), f.Month()-11, 1, 0, 0, 0, 0, time.UTC).Format(formatDate)
	return min(debut, douzeMois)
}

// Période de même durée qui se termine la veille du début.
func periodePrecedente(debut, fin string) (string, string) {
	d, _ := time.Parse(formatDate, debut)
	f, _ := time.Parse(formatDate, fin)
	finPrecedente := d.AddDate(0, 0, -1)
	debutPrecedent := finPrecedente.Add(-f.Sub(d))
	return debutPrecedent.Format(formatDate), finPrecedente.Format(formatDate)
}

func (s *Serveur) indicateursPeriode(debut, fin string) (IndicateursPeriode, *ResumeFinances, error) {
	var ind IndicateursPeriode
	resume, err := s.resumeFinances(debut, fin)
	if err != nil {
		return ind, nil, err
	}
	ind.Revenus = resume.TotalRevenus
	ind.Charges = resume.TotalCharges
	ind.Resultat = resume.Resultat

	err = s.db.QueryRow(`
		SELECT
		    (SELECT count(*) `+filtreActeurs+`
		       AND u.date_creation >= $1::timestamp AND u.date_creation < $2::timestamp + interval '1 day'),
		    (SELECT count(DISTINCT d.utilisateur_id)
		       FROM vue_demandes_acteurs d
		       JOIN utilisateurs ua ON ua.id = d.utilisateur_id
		       JOIN roles ra ON ra.id = ua.role_id
		      WHERE ra.code IN ('particulier', 'professionnel') AND `+filtrePeriodeDemandes+`),
		    (SELECT count(*) FROM annonces
		      WHERE date_creation >= $1::timestamp AND date_creation < $2::timestamp + interval '1 day'),
		    (SELECT count(*) FROM depots
		      WHERE date_demande >= $1::timestamp AND date_demande < $2::timestamp + interval '1 day'),
		    (SELECT count(*) FROM inscriptions
		      WHERE date_inscription >= $1::timestamp AND date_inscription < $2::timestamp + interval '1 day')`,
		debut, fin,
	).Scan(&ind.NouveauxActeurs, &ind.ActeursActifs, &ind.Annonces, &ind.Depots, &ind.Inscriptions)
	return ind, &resume, err
}

func (s *Serveur) activiteParMois(debut, fin string) ([]ActiviteMois, error) {
	lignes, err := s.db.Query(`
		WITH m AS (
		    SELECT greatest(m.debut, $1::timestamp) AS de,
		           least(m.debut + interval '1 month', $2::timestamp + interval '1 day') AS a,
		           to_char(m.debut, 'YYYY-MM') AS mois
		    FROM generate_series(date_trunc('month', $1::timestamp), $2::timestamp, interval '1 month') AS m(debut)
		)
		SELECT m.mois,
		       (SELECT count(*) `+filtreActeurs+` AND u.date_creation >= m.de AND u.date_creation < m.a),
		       (SELECT count(*) FROM annonces WHERE date_creation >= m.de AND date_creation < m.a),
		       (SELECT count(*) FROM depots WHERE date_demande >= m.de AND date_demande < m.a),
		       (SELECT count(*) FROM inscriptions WHERE date_inscription >= m.de AND date_inscription < m.a)
		FROM m
		ORDER BY m.mois`, debut, fin)
	if err != nil {
		return nil, err
	}
	defer lignes.Close()

	resultat := make([]ActiviteMois, 0)
	for lignes.Next() {
		var a ActiviteMois
		if err := lignes.Scan(&a.Mois, &a.NouveauxActeurs, &a.Annonces, &a.Depots, &a.Inscriptions); err != nil {
			return nil, err
		}
		resultat = append(resultat, a)
	}
	return resultat, lignes.Err()
}

func (s *Serveur) elementsATraiter() (ElementsATraiter, error) {
	var t ElementsATraiter
	err := s.db.QueryRow(`
		SELECT (SELECT count(*) FROM formations WHERE statut = 'en_attente')
		     + (SELECT count(*) FROM ateliers WHERE statut = 'en_attente')
		     + (SELECT count(*) FROM evenements WHERE statut = 'en_attente'),
		       (SELECT count(*) FROM annonces WHERE statut = 'en_attente'),
		       (SELECT count(*) FROM forum_sujets WHERE statut = 'en_attente'),
		       (SELECT count(*) FROM forum_signalements WHERE statut = 'a_traiter'),
		       (SELECT count(*) FROM conteneurs WHERE statut = 'maintenance')`,
	).Scan(&t.OffresAValider, &t.AnnoncesAValider, &t.SujetsAModerer, &t.Signalements, &t.ConteneursMaintenance)
	return t, err
}

func (s *Serveur) gererTableauDeBordAdmin(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	debut, fin, message := periodeDepuisRequete(r)
	if message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}

	tableau := TableauDeBord{Debut: debut, Fin: fin, DebutSeries: debutDesSeries(debut, fin)}
	tableau.DebutPrecedent, tableau.FinPrecedent = periodePrecedente(debut, fin)

	actuel, resume, err := s.indicateursPeriode(debut, fin)
	if err == nil {
		tableau.Actuel = actuel
		tableau.RevenuMensuelRecurrent = resume.RevenuMensuelRecurrent
		tableau.AbonnesActifs = resume.AbonnesActifs
		tableau.CampagnesActives = resume.CampagnesActives
		tableau.RevenusParSource = resume.RevenusParSource
		tableau.Precedent, _, err = s.indicateursPeriode(tableau.DebutPrecedent, tableau.FinPrecedent)
	}
	if err == nil {
		var series ResumeFinances
		series, err = s.resumeFinances(tableau.DebutSeries, fin)
		tableau.FinancesParMois = series.ParMois
	}
	if err == nil {
		tableau.ActiviteParMois, err = s.activiteParMois(tableau.DebutSeries, fin)
	}
	if err == nil {
		tableau.ATraiter, err = s.elementsATraiter()
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, tableau)
}
