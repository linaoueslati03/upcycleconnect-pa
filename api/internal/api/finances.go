package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	sourceAbonnement = "abonnement"
	sourceVente      = "vente"
	sourcePublicite  = "publicite"
	sourceFormation  = "formation"
	sourceAtelier    = "atelier"
	sourceEvenement  = "evenement"
)

var sourcesPaiement = []string{sourceAbonnement, sourceVente, sourcePublicite, sourceFormation, sourceAtelier, sourceEvenement}

var statutsPaiement = []string{"en_attente", "reussi", "echoue", "rembourse"}

var categoriesCharge = []string{"salaires", "loyers", "infrastructure", "marketing", "autre"}

const (
	formatDate        = "2006-01-02"
	paiementsParPage  = 20
	tauxCommissionMin = 5
	tauxCommissionMax = 10
	prixPremiumMin    = 15
	prixPremiumMax    = 30
)

const revenuPlateforme = "CASE WHEN p.source = 'vente' THEN p.montant_commission ELSE p.montant END"

type ParametresFinances struct {
	TauxCommission     float64   `json:"taux_commission"`
	PrixPremiumMensuel float64   `json:"prix_premium_mensuel"`
	DateModification   time.Time `json:"date_modification"`
}

type RevenuSource struct {
	Source          string  `json:"source"`
	NbPaiements     int     `json:"nb_paiements"`
	Revenu          float64 `json:"revenu"`
	MontantEncaisse float64 `json:"montant_encaisse"`
}

type BilanMois struct {
	Mois     string  `json:"mois"`
	Revenus  float64 `json:"revenus"`
	Charges  float64 `json:"charges"`
	Resultat float64 `json:"resultat"`
}

type ResumeFinances struct {
	Debut                  string         `json:"debut"`
	Fin                    string         `json:"fin"`
	RevenusParSource       []RevenuSource `json:"revenus_par_source"`
	TotalRevenus           float64        `json:"total_revenus"`
	TotalCharges           float64        `json:"total_charges"`
	Resultat               float64        `json:"resultat"`
	VolumeVentes           float64        `json:"volume_ventes"`
	AbonnesActifs          int            `json:"abonnes_actifs"`
	CampagnesActives       int            `json:"campagnes_actives"`
	RevenuMensuelRecurrent float64        `json:"revenu_mensuel_recurrent"`
	ParMois                []BilanMois    `json:"par_mois"`
}

type Paiement struct {
	ID                int       `json:"id"`
	UtilisateurID     int       `json:"utilisateur_id"`
	Utilisateur       string    `json:"utilisateur"`
	Source            string    `json:"source"`
	ReferenceID       *int      `json:"reference_id"`
	Montant           float64   `json:"montant"`
	TauxCommission    *float64  `json:"taux_commission"`
	MontantCommission *float64  `json:"montant_commission"`
	Statut            string    `json:"statut"`
	StripeID          *string   `json:"stripe_id"`
	DatePaiement      time.Time `json:"date_paiement"`
}

type PagePaiements struct {
	Paiements []Paiement `json:"paiements"`
	Total     int        `json:"total"`
	Page      int        `json:"page"`
	ParPage   int        `json:"par_page"`
}

type AbonnementAdmin struct {
	ID             int        `json:"id"`
	Utilisateur    string     `json:"utilisateur"`
	Email          string     `json:"email"`
	MontantMensuel float64    `json:"montant_mensuel"`
	Statut         string     `json:"statut"`
	DateDebut      time.Time  `json:"date_debut"`
	DateFin        *time.Time `json:"date_fin"`
	TotalPaye      float64    `json:"total_paye"`
}

type CampagneAdmin struct {
	ID             int        `json:"id"`
	Annonceur      string     `json:"annonceur"`
	CibleType      string     `json:"cible_type"`
	CibleID        int        `json:"cible_id"`
	CibleTitre     *string    `json:"cible_titre"`
	MontantMensuel float64    `json:"montant_mensuel"`
	Statut         string     `json:"statut"`
	DateDebut      time.Time  `json:"date_debut"`
	DateFin        *time.Time `json:"date_fin"`
	TotalPaye      float64    `json:"total_paye"`
}

type Charge struct {
	ID         int       `json:"id"`
	Libelle    string    `json:"libelle"`
	Categorie  string    `json:"categorie"`
	Montant    float64   `json:"montant"`
	DateCharge time.Time `json:"date_charge"`
	Auteur     *string   `json:"auteur"`
}

type ChargeEntree struct {
	Libelle    string  `json:"libelle"`
	Categorie  string  `json:"categorie"`
	Montant    float64 `json:"montant"`
	DateCharge string  `json:"date_charge"`
}

func arrondirCentimes(montant float64) float64 {
	return math.Round(montant*100) / 100
}

func periodeDepuisRequete(r *http.Request) (string, string, string) {
	aujourdhui := time.Now()
	fin := aujourdhui.Format(formatDate)
	debut := time.Date(aujourdhui.Year(), aujourdhui.Month()-11, 1, 0, 0, 0, 0, time.Local).Format(formatDate)

	if valeur := r.URL.Query().Get("debut"); valeur != "" {
		if _, err := time.Parse(formatDate, valeur); err != nil {
			return "", "", "date de début invalide (AAAA-MM-JJ attendu)"
		}
		debut = valeur
	}
	if valeur := r.URL.Query().Get("fin"); valeur != "" {
		if _, err := time.Parse(formatDate, valeur); err != nil {
			return "", "", "date de fin invalide (AAAA-MM-JJ attendu)"
		}
		fin = valeur
	}
	if debut > fin {
		return "", "", "la date de début doit précéder la date de fin"
	}
	return debut, fin, ""
}

func enregistrerPaiement(tx *sql.Tx, utilisateurID int, source string, referenceID *int, montant float64, stripeID string) (int, error) {
	if !slices.Contains(sourcesPaiement, source) {
		return 0, fmt.Errorf("source de paiement inconnue : %s", source)
	}
	var id int
	err := tx.QueryRow(`
		INSERT INTO paiements (utilisateur_id, source, reference_id, montant, taux_commission, montant_commission, stripe_id)
		SELECT $1, $2::varchar, $3, $4::numeric,
		       CASE WHEN $2::varchar = 'vente' THEN p.taux_commission END,
		       CASE WHEN $2::varchar = 'vente' THEN round($4::numeric * p.taux_commission / 100, 2) END,
		       NULLIF($5::varchar, '')
		FROM parametres_finances p
		WHERE p.id = 1
		RETURNING id`,
		utilisateurID, source, referenceID, montant, stripeID,
	).Scan(&id)
	return id, err
}

func (s *Serveur) gererResumeFinances(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	debut, fin, message := periodeDepuisRequete(r)
	if message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}

	resume := ResumeFinances{Debut: debut, Fin: fin, RevenusParSource: make([]RevenuSource, 0), ParMois: make([]BilanMois, 0)}

	lignes, err := s.db.Query(`
		SELECT p.source, count(*), COALESCE(sum(`+revenuPlateforme+`), 0), COALESCE(sum(p.montant), 0)
		FROM paiements p
		WHERE p.statut = 'reussi'
		  AND p.date_paiement >= $1::timestamp AND p.date_paiement < $2::timestamp + interval '1 day'
		GROUP BY p.source`, debut, fin)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	parSource := make(map[string]RevenuSource)
	for lignes.Next() {
		var rs RevenuSource
		if err := lignes.Scan(&rs.Source, &rs.NbPaiements, &rs.Revenu, &rs.MontantEncaisse); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		parSource[rs.Source] = rs
	}
	for _, source := range sourcesPaiement {
		rs, trouve := parSource[source]
		if !trouve {
			rs = RevenuSource{Source: source}
		}
		resume.RevenusParSource = append(resume.RevenusParSource, rs)
		resume.TotalRevenus += rs.Revenu
		if source == sourceVente {
			resume.VolumeVentes = rs.MontantEncaisse
		}
	}

	err = s.db.QueryRow(
		"SELECT COALESCE(sum(montant), 0) FROM charges WHERE date_charge BETWEEN $1::date AND $2::date",
		debut, fin,
	).Scan(&resume.TotalCharges)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	resume.TotalRevenus = arrondirCentimes(resume.TotalRevenus)
	resume.Resultat = arrondirCentimes(resume.TotalRevenus - resume.TotalCharges)

	err = s.db.QueryRow(`
		SELECT (SELECT count(*) FROM abonnements WHERE statut = 'actif'),
		       (SELECT count(*) FROM campagnes_publicitaires WHERE statut = 'active'),
		       (SELECT COALESCE(sum(montant_mensuel), 0) FROM abonnements WHERE statut = 'actif')
		     + (SELECT COALESCE(sum(montant_mensuel), 0) FROM campagnes_publicitaires WHERE statut = 'active')`,
	).Scan(&resume.AbonnesActifs, &resume.CampagnesActives, &resume.RevenuMensuelRecurrent)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	mois, err := s.db.Query(`
		SELECT to_char(m.debut, 'YYYY-MM'),
		       (SELECT COALESCE(sum(`+revenuPlateforme+`), 0) FROM paiements p
		        WHERE p.statut = 'reussi'
		          AND p.date_paiement >= greatest(m.debut, $1::timestamp)
		          AND p.date_paiement < least(m.debut + interval '1 month', $2::timestamp + interval '1 day')),
		       (SELECT COALESCE(sum(c.montant), 0) FROM charges c
		        WHERE c.date_charge >= greatest(m.debut, $1::timestamp)
		          AND c.date_charge < least(m.debut + interval '1 month', $2::timestamp + interval '1 day'))
		FROM generate_series(date_trunc('month', $1::timestamp), $2::timestamp, interval '1 month') AS m(debut)
		ORDER BY m.debut`, debut, fin)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer mois.Close()

	for mois.Next() {
		var b BilanMois
		if err := mois.Scan(&b.Mois, &b.Revenus, &b.Charges); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		b.Resultat = arrondirCentimes(b.Revenus - b.Charges)
		resume.ParMois = append(resume.ParMois, b)
	}

	envoyerJSON(w, http.StatusOK, resume)
}

func (s *Serveur) gererListePaiements(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	debut, fin, message := periodeDepuisRequete(r)
	if message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	args := []any{debut, fin}
	conditions := []string{"p.date_paiement >= $1::timestamp", "p.date_paiement < $2::timestamp + interval '1 day'"}
	if source := r.URL.Query().Get("source"); slices.Contains(sourcesPaiement, source) {
		args = append(args, source)
		conditions = append(conditions, fmt.Sprintf("p.source = $%d", len(args)))
	}
	if statut := r.URL.Query().Get("statut"); slices.Contains(statutsPaiement, statut) {
		args = append(args, statut)
		conditions = append(conditions, fmt.Sprintf("p.statut = $%d", len(args)))
	}
	where := strings.Join(conditions, " AND ")

	resultat := PagePaiements{Paiements: make([]Paiement, 0), Page: page, ParPage: paiementsParPage}
	if err := s.db.QueryRow("SELECT count(*) FROM paiements p WHERE "+where, args...).Scan(&resultat.Total); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	args = append(args, paiementsParPage, (page-1)*paiementsParPage)
	lignes, err := s.db.Query(fmt.Sprintf(`
		SELECT p.id, p.utilisateur_id, u.prenom || ' ' || u.nom, p.source, p.reference_id, p.montant,
		       p.taux_commission, p.montant_commission, p.statut, p.stripe_id, p.date_paiement
		FROM paiements p
		JOIN utilisateurs u ON u.id = p.utilisateur_id
		WHERE %s
		ORDER BY p.date_paiement DESC, p.id DESC
		LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args)), args...)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	for lignes.Next() {
		var p Paiement
		if err := lignes.Scan(&p.ID, &p.UtilisateurID, &p.Utilisateur, &p.Source, &p.ReferenceID, &p.Montant,
			&p.TauxCommission, &p.MontantCommission, &p.Statut, &p.StripeID, &p.DatePaiement); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		resultat.Paiements = append(resultat.Paiements, p)
	}

	envoyerJSON(w, http.StatusOK, resultat)
}

func (s *Serveur) lireParametresFinances() (ParametresFinances, error) {
	var p ParametresFinances
	err := s.db.QueryRow(
		"SELECT taux_commission, prix_premium_mensuel, date_modification FROM parametres_finances WHERE id = 1",
	).Scan(&p.TauxCommission, &p.PrixPremiumMensuel, &p.DateModification)
	return p, err
}

func (s *Serveur) gererParametresFinances(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	p, err := s.lireParametresFinances()
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	envoyerJSON(w, http.StatusOK, p)
}

func (s *Serveur) gererModificationParametresFinances(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	var entree ParametresFinances
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if entree.TauxCommission < tauxCommissionMin || entree.TauxCommission > tauxCommissionMax {
		envoyerErreur(w, http.StatusBadRequest, fmt.Sprintf("le taux de commission doit être compris entre %d et %d %%", tauxCommissionMin, tauxCommissionMax))
		return
	}
	if entree.PrixPremiumMensuel < prixPremiumMin || entree.PrixPremiumMensuel > prixPremiumMax {
		envoyerErreur(w, http.StatusBadRequest, fmt.Sprintf("le prix du premium doit être compris entre %d et %d €", prixPremiumMin, prixPremiumMax))
		return
	}

	_, err := s.db.Exec(
		"UPDATE parametres_finances SET taux_commission = $1, prix_premium_mensuel = $2, date_modification = now() WHERE id = 1",
		entree.TauxCommission, entree.PrixPremiumMensuel,
	)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	p, err := s.lireParametresFinances()
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	envoyerJSON(w, http.StatusOK, p)
}

func (s *Serveur) gererListeAbonnementsAdmin(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	lignes, err := s.db.Query(`
		SELECT a.id, u.prenom || ' ' || u.nom, u.email, a.montant_mensuel, a.statut, a.date_debut, a.date_fin,
		       (SELECT COALESCE(sum(p.montant), 0) FROM paiements p
		        WHERE p.source = 'abonnement' AND p.reference_id = a.id AND p.statut = 'reussi')
		FROM abonnements a
		JOIN utilisateurs u ON u.id = a.utilisateur_id
		ORDER BY a.statut = 'actif' DESC, a.date_debut DESC, a.id DESC`)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	abonnements := make([]AbonnementAdmin, 0)
	for lignes.Next() {
		var a AbonnementAdmin
		if err := lignes.Scan(&a.ID, &a.Utilisateur, &a.Email, &a.MontantMensuel, &a.Statut, &a.DateDebut, &a.DateFin, &a.TotalPaye); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		abonnements = append(abonnements, a)
	}

	envoyerJSON(w, http.StatusOK, abonnements)
}

func (s *Serveur) gererListeCampagnesAdmin(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	lignes, err := s.db.Query(`
		SELECT c.id, u.prenom || ' ' || u.nom, c.cible_type, c.cible_id,
		       CASE c.cible_type
		           WHEN 'projet' THEN (SELECT titre FROM projets_upcycling WHERE id = c.cible_id)
		           WHEN 'annonce' THEN (SELECT titre FROM annonces WHERE id = c.cible_id)
		       END,
		       c.montant_mensuel, c.statut, c.date_debut, c.date_fin,
		       (SELECT COALESCE(sum(p.montant), 0) FROM paiements p
		        WHERE p.source = 'publicite' AND p.reference_id = c.id AND p.statut = 'reussi')
		FROM campagnes_publicitaires c
		JOIN utilisateurs u ON u.id = c.utilisateur_id
		ORDER BY c.statut = 'active' DESC, c.date_debut DESC, c.id DESC`)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	campagnes := make([]CampagneAdmin, 0)
	for lignes.Next() {
		var c CampagneAdmin
		if err := lignes.Scan(&c.ID, &c.Annonceur, &c.CibleType, &c.CibleID, &c.CibleTitre, &c.MontantMensuel,
			&c.Statut, &c.DateDebut, &c.DateFin, &c.TotalPaye); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		campagnes = append(campagnes, c)
	}

	envoyerJSON(w, http.StatusOK, campagnes)
}

func (s *Serveur) gererListeCharges(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	debut, fin, message := periodeDepuisRequete(r)
	if message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}

	lignes, err := s.db.Query(`
		SELECT c.id, c.libelle, c.categorie, c.montant, c.date_charge, u.prenom || ' ' || u.nom
		FROM charges c
		LEFT JOIN utilisateurs u ON u.id = c.auteur_id
		WHERE c.date_charge BETWEEN $1::date AND $2::date
		ORDER BY c.date_charge DESC, c.id DESC`, debut, fin)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	charges := make([]Charge, 0)
	for lignes.Next() {
		var c Charge
		if err := lignes.Scan(&c.ID, &c.Libelle, &c.Categorie, &c.Montant, &c.DateCharge, &c.Auteur); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		charges = append(charges, c)
	}

	envoyerJSON(w, http.StatusOK, charges)
}

func validerCharge(entree ChargeEntree) string {
	if entree.Libelle == "" {
		return "libellé obligatoire"
	}
	if len([]rune(entree.Libelle)) > 150 {
		return "libellé trop long (150 caractères maximum)"
	}
	if !slices.Contains(categoriesCharge, entree.Categorie) {
		return "catégorie invalide"
	}
	if entree.Montant <= 0 {
		return "le montant doit être positif"
	}
	if _, err := time.Parse(formatDate, entree.DateCharge); err != nil {
		return "date invalide (AAAA-MM-JJ attendu)"
	}
	return ""
}

func (s *Serveur) gererCreationCharge(w http.ResponseWriter, r *http.Request) {
	adminID, ok := s.exigerRole(w, r, roleAdministrateur)
	if !ok {
		return
	}

	var entree ChargeEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	entree.Libelle = strings.TrimSpace(entree.Libelle)
	if message := validerCharge(entree); message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}

	var c Charge
	err := s.db.QueryRow(`
		INSERT INTO charges (libelle, categorie, montant, date_charge, auteur_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, libelle, categorie, montant, date_charge`,
		entree.Libelle, entree.Categorie, entree.Montant, entree.DateCharge, adminID,
	).Scan(&c.ID, &c.Libelle, &c.Categorie, &c.Montant, &c.DateCharge)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusCreated, c)
}

func (s *Serveur) gererSuppressionCharge(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var supprime int
	err = s.db.QueryRow("DELETE FROM charges WHERE id = $1 RETURNING id", id).Scan(&supprime)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "charge introuvable")
		return
	} else if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
