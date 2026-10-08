package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
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

	if !slices.Contains(statutsAnnonceValidee, a.Statut) && !s.peutVoirAnnonceNonValidee(r, a.UtilisateurID) {
		envoyerErreur(w, http.StatusNotFound, "annonce introuvable")
		return
	}

	envoyerJSON(w, http.StatusOK, a)
}

// Cycle de vie d'une annonce : créée « en_attente », elle est validée (« en_ligne ») ou
// refusée (« refusee ») par un administrateur. Une fois en ligne, son auteur peut la passer
// en « reservee » ou « cedee ». Le sujet impose cette validation par le service administratif.
var statutsAnnonceValidee = []string{"en_ligne", "reservee", "cedee"}

// statutApresModification calcule le statut d'une annonce modifiée par son auteur :
//   - une annonce cédée est définitive (sinon on pourrait regagner les points en boucle) ;
//   - une annonce pas encore validée, refusée, ou dont le contenu change, repart en
//     validation : l'auteur ne peut jamais publier lui-même un contenu non vérifié ;
//   - sinon l'auteur peut seulement la passer en ligne, réservée ou cédée.
func statutApresModification(statutAvant, statutDemande string, contenuModifie bool) string {
	if statutAvant == "cedee" {
		return "cedee"
	}
	if !slices.Contains(statutsAnnonceValidee, statutAvant) || contenuModifie {
		return "en_attente"
	}
	if slices.Contains(statutsAnnonceValidee, statutDemande) {
		return statutDemande
	}
	return statutAvant
}

// memeContenu compare le contenu d'une annonce (hors statut) avec les nouvelles valeurs.
func memeContenu(a Annonce, e AnnonceEntree) bool {
	memePrix := (a.Prix == nil && e.Prix == nil) || (a.Prix != nil && e.Prix != nil && *a.Prix == *e.Prix)
	memeCategorie := (a.CategorieID == nil && e.CategorieID == nil) ||
		(a.CategorieID != nil && e.CategorieID != nil && *a.CategorieID == *e.CategorieID)
	return a.Titre == e.Titre && a.Description == e.Description && a.Type == e.Type &&
		a.Localisation == e.Localisation && memePrix && memeCategorie
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
		INSERT INTO annonces (utilisateur_id, titre, description, type, prix, categorie_id, localisation, statut)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'en_attente')
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

	var avant Annonce
	err = scannerAnnonce(s.db.QueryRow(`SELECT id, utilisateur_id, titre, description, type, prix, categorie_id, localisation, statut, date_creation
		FROM annonces WHERE id = $1`, id), &avant)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "annonce introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if avant.UtilisateurID != utilisateurID {
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
	if avant.Statut == "cedee" && !memeContenu(avant, entree) {
		envoyerErreur(w, http.StatusConflict, "une annonce cédée ne peut plus être modifiée")
		return
	}
	entree.Statut = statutApresModification(avant.Statut, entree.Statut, !memeContenu(avant, entree))

	// La mise à jour et l'éventuel ajout de points se font dans la même transaction
	tx, err := s.db.Begin()
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer tx.Rollback()

	var a Annonce
	ligne := tx.QueryRow(`
		UPDATE annonces SET titre=$1, description=$2, type=$3, prix=$4, categorie_id=$5, localisation=$6, statut=$7
		WHERE id = $8
		RETURNING id, utilisateur_id, titre, description, type, prix, categorie_id, localisation, statut, date_creation`,
		entree.Titre, entree.Description, entree.Type, entree.Prix, entree.CategorieID, entree.Localisation, entree.Statut, id,
	)
	if err := scannerAnnonce(ligne, &a); err != nil {
		gererErreurPostgres(w, err)
		return
	}

	if entree.Statut == "cedee" && avant.Statut != "cedee" {
		if err := ajouterPointsScore(tx, a.UtilisateurID, pointsAnnonceCedee, "Annonce cédée"); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
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

	if _, err := s.db.Exec("DELETE FROM annonces WHERE id = $1", id); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	envoyerJSON(w, http.StatusOK, map[string]string{"message": "annonce supprimée"})
}

// peutVoirAnnonceNonValidee : l'auteur de l'annonce ou un administrateur.
func (s *Serveur) peutVoirAnnonceNonValidee(r *http.Request, auteurID int) bool {
	utilisateurID, role := s.roleConnecte(r)
	return utilisateurID != 0 && (utilisateurID == auteurID || role == roleAdministrateur)
}

// AnnonceModeration est une annonce vue par l'administrateur, avec son auteur.
type AnnonceModeration struct {
	Annonce
	Auteur string `json:"auteur"`
}

// gererModerationAnnonces liste les annonces pour le service administratif
// (par défaut celles en attente ; ?statut= pour un autre statut, ?statut=tous pour tout).
func (s *Serveur) gererModerationAnnonces(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	statut := r.URL.Query().Get("statut")
	if statut == "" {
		statut = "en_attente"
	}

	requete := `SELECT a.id, a.utilisateur_id, a.titre, a.description, a.type, a.prix, a.categorie_id,
			a.localisation, a.statut, a.date_creation, u.prenom || ' ' || u.nom || ' (' || u.email || ')'
		FROM annonces a JOIN utilisateurs u ON u.id = a.utilisateur_id`
	args := []any{}
	if statut != "tous" {
		requete += ` WHERE a.statut = $1`
		args = append(args, statut)
	}
	requete += ` ORDER BY a.date_creation ASC`

	lignes, err := s.db.Query(requete, args...)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	annonces := make([]AnnonceModeration, 0)
	for lignes.Next() {
		var a AnnonceModeration
		var prix sql.NullFloat64
		var categorieID sql.NullInt64
		err := lignes.Scan(&a.ID, &a.UtilisateurID, &a.Titre, &a.Description, &a.Type, &prix, &categorieID,
			&a.Localisation, &a.Statut, &a.DateCreation, &a.Auteur)
		if err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		if prix.Valid {
			a.Prix = &prix.Float64
		}
		if categorieID.Valid {
			v := int(categorieID.Int64)
			a.CategorieID = &v
		}
		annonces = append(annonces, a)
	}

	envoyerJSON(w, http.StatusOK, annonces)
}

func (s *Serveur) gererValidationAnnonce(w http.ResponseWriter, r *http.Request) {
	s.deciderAnnonce(w, r, "en_ligne")
}

func (s *Serveur) gererRefusAnnonce(w http.ResponseWriter, r *http.Request) {
	s.deciderAnnonce(w, r, "refusee")
}

// deciderAnnonce applique la décision de l'administrateur à une annonce en attente.
func (s *Serveur) deciderAnnonce(w http.ResponseWriter, r *http.Request, nouveauStatut string) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var a Annonce
	ligne := s.db.QueryRow(`
		UPDATE annonces SET statut = $1
		WHERE id = $2 AND statut = 'en_attente'
		RETURNING id, utilisateur_id, titre, description, type, prix, categorie_id, localisation, statut, date_creation`,
		nouveauStatut, id,
	)
	if err := scannerAnnonce(ligne, &a); errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusConflict, "annonce introuvable ou déjà traitée")
		return
	} else if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, a)
}
