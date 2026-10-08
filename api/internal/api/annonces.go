package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Annonce struct {
	ID            int       `json:"id"`
	UtilisateurID int       `json:"utilisateur_id"`
	Titre         string    `json:"titre"`
	Description   string    `json:"description"`
	Type          string    `json:"type"`
	Prix          *float64  `json:"prix"`
	CategorieID   *int      `json:"categorie_id"`
	Localisation  string    `json:"localisation"`
	Statut        string    `json:"statut"`
	DateCreation  time.Time `json:"date_creation"`
}

type AnnonceEntree struct {
	Titre        string   `json:"titre"`
	Description  string   `json:"description"`
	Type         string   `json:"type"`
	Prix         *float64 `json:"prix"`
	CategorieID  *int     `json:"categorie_id"`
	Localisation string   `json:"localisation"`
	Statut       string   `json:"statut"`
}

func scannerAnnonce(lignes interface{ Scan(...any) error }, a *Annonce) error {
	var prix sql.NullFloat64
	var categorieID sql.NullInt64
	err := lignes.Scan(&a.ID, &a.UtilisateurID, &a.Titre, &a.Description, &a.Type, &prix, &categorieID, &a.Localisation, &a.Statut, &a.DateCreation)
	if err != nil {
		return err
	}
	if prix.Valid {
		a.Prix = &prix.Float64
	}
	if categorieID.Valid {
		v := int(categorieID.Int64)
		a.CategorieID = &v
	}
	return nil
}

func (s *Serveur) gererListeAnnonces(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("mine") == "true" {
		utilisateurID, err := s.utilisateurConnecte(r)
		if err != nil {
			envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
			return
		}
		s.listerAnnonces(w, "utilisateur_id = $1", utilisateurID)
		return
	}

	conditions := []string{"statut = 'en_ligne'"}
	args := []any{}
	if typeAnnonce := r.URL.Query().Get("type"); typeAnnonce != "" {
		args = append(args, typeAnnonce)
		conditions = append(conditions, "type = $"+strconv.Itoa(len(args)))
	}
	if localisation := r.URL.Query().Get("localisation"); localisation != "" {
		args = append(args, "%"+localisation+"%")
		conditions = append(conditions, "localisation ILIKE $"+strconv.Itoa(len(args)))
	}
	s.listerAnnonces(w, strings.Join(conditions, " AND "), args...)
}

func (s *Serveur) listerAnnonces(w http.ResponseWriter, whereClause string, args ...any) {
	requete := `SELECT id, utilisateur_id, titre, description, type, prix, categorie_id, localisation, statut, date_creation
		FROM annonces WHERE ` + whereClause + ` ORDER BY date_creation DESC`

	lignes, err := s.db.Query(requete, args...)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	annonces := make([]Annonce, 0)
	for lignes.Next() {
		var a Annonce
		if err := scannerAnnonce(lignes, &a); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		annonces = append(annonces, a)
	}

	envoyerJSON(w, http.StatusOK, annonces)
}

func (s *Serveur) gererDetailAnnonce(w http.ResponseWriter, r *http.Request) {
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var a Annonce
	ligne := s.db.QueryRow(`SELECT id, utilisateur_id, titre, description, type, prix, categorie_id, localisation, statut, date_creation
		FROM annonces WHERE id = $1`, id)
	if err := scannerAnnonce(ligne, &a); errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "annonce introuvable")
		return
	} else if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, a)
}

func validerAnnonce(entree AnnonceEntree) string {
	if strings.TrimSpace(entree.Titre) == "" {
		return "titre obligatoire"
	}
	if entree.Type != "don" && entree.Type != "vente" {
		return "type invalide (don ou vente attendu)"
	}
	if entree.Type == "vente" && (entree.Prix == nil || *entree.Prix <= 0) {
		return "prix obligatoire pour une vente"
	}
	return ""
}

func (s *Serveur) gererCreationAnnonce(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	var entree AnnonceEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if message := validerAnnonce(entree); message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}

	var a Annonce
	ligne := s.db.QueryRow(`
		INSERT INTO annonces (utilisateur_id, titre, description, type, prix, categorie_id, localisation)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, utilisateur_id, titre, description, type, prix, categorie_id, localisation, statut, date_creation`,
		utilisateurID, entree.Titre, entree.Description, entree.Type, entree.Prix, entree.CategorieID, entree.Localisation,
	)
	if err := scannerAnnonce(ligne, &a); err != nil {
		gererErreurPostgres(w, err)
		return
	}

	envoyerJSON(w, http.StatusCreated, a)
}

func (s *Serveur) gererModificationAnnonce(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var proprietaireID int
	var statutAvant string
	err = s.db.QueryRow("SELECT utilisateur_id, statut FROM annonces WHERE id = $1", id).Scan(&proprietaireID, &statutAvant)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "annonce introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if proprietaireID != utilisateurID {
		envoyerErreur(w, http.StatusForbidden, "vous n'êtes pas propriétaire de cette annonce")
		return
	}

	var entree AnnonceEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if message := validerAnnonce(entree); message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}
	if entree.Statut == "" {
		entree.Statut = "en_ligne"
	}

	var a Annonce
	ligne := s.db.QueryRow(`
		UPDATE annonces SET titre=$1, description=$2, type=$3, prix=$4, categorie_id=$5, localisation=$6, statut=$7
		WHERE id = $8
		RETURNING id, utilisateur_id, titre, description, type, prix, categorie_id, localisation, statut, date_creation`,
		entree.Titre, entree.Description, entree.Type, entree.Prix, entree.CategorieID, entree.Localisation, entree.Statut, id,
	)
	if err := scannerAnnonce(ligne, &a); err != nil {
		gererErreurPostgres(w, err)
		return
	}

	if entree.Statut == "cedee" && statutAvant != "cedee" {
		ajouterPointsScore(s.db, a.UtilisateurID, pointsAnnonceCedee, "Annonce cédée")
	}

	envoyerJSON(w, http.StatusOK, a)
}

func (s *Serveur) gererSuppressionAnnonce(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var proprietaireID int
	err = s.db.QueryRow("SELECT utilisateur_id FROM annonces WHERE id = $1", id).Scan(&proprietaireID)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "annonce introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if proprietaireID != utilisateurID {
		envoyerErreur(w, http.StatusForbidden, "vous n'êtes pas propriétaire de cette annonce")
		return
	}

	s.db.Exec("DELETE FROM annonces WHERE id = $1", id)
	envoyerJSON(w, http.StatusOK, map[string]string{"message": "annonce supprimée"})
}
