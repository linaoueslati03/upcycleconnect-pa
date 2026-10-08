package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type Conseil struct {
	ID        int       `json:"id"`
	Titre     string    `json:"titre"`
	Contenu   string    `json:"contenu"`
	Categorie *string   `json:"categorie"`
	AuteurID  int       `json:"auteur_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Serveur) gererListeConseils(w http.ResponseWriter, r *http.Request) {
	requete := `SELECT id, titre, contenu, categorie, auteur_id, created_at, updated_at
		FROM conseils WHERE statut = 'publie' ORDER BY created_at DESC`
	args := []any{}

	if categorie := r.URL.Query().Get("categorie"); categorie != "" {
		requete = `SELECT id, titre, contenu, categorie, auteur_id, created_at, updated_at
			FROM conseils WHERE statut = 'publie' AND categorie = $1 ORDER BY created_at DESC`
		args = append(args, categorie)
	}

	lignes, err := s.db.Query(requete, args...)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	conseils := make([]Conseil, 0)
	for lignes.Next() {
		var c Conseil
		var categorie sql.NullString
		if err := lignes.Scan(&c.ID, &c.Titre, &c.Contenu, &categorie, &c.AuteurID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		if categorie.Valid {
			c.Categorie = &categorie.String
		}
		conseils = append(conseils, c)
	}

	envoyerJSON(w, http.StatusOK, conseils)
}

func (s *Serveur) gererDetailConseil(w http.ResponseWriter, r *http.Request) {
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var c Conseil
	var categorie sql.NullString
	err = s.db.QueryRow(`
		SELECT id, titre, contenu, categorie, auteur_id, created_at, updated_at
		FROM conseils WHERE id = $1 AND statut = 'publie'`, id,
	).Scan(&c.ID, &c.Titre, &c.Contenu, &categorie, &c.AuteurID, &c.CreatedAt, &c.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "conseil introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if categorie.Valid {
		c.Categorie = &categorie.String
	}

	envoyerJSON(w, http.StatusOK, c)
}

type ConseilEntree struct {
	Titre     string  `json:"titre"`
	Contenu   string  `json:"contenu"`
	Categorie *string `json:"categorie"`
	Statut    string  `json:"statut"`
}

// ConseilRedaction est un conseil vu par son auteur : avec son statut (brouillon ou publié).
type ConseilRedaction struct {
	Conseil
	Statut string `json:"statut"`
}

const colonnesConseilRedaction = `id, titre, contenu, categorie, auteur_id, created_at, updated_at, statut`

func scannerConseilRedaction(ligne interface{ Scan(...any) error }, c *ConseilRedaction) error {
	return ligne.Scan(&c.ID, &c.Titre, &c.Contenu, &c.Categorie, &c.AuteurID, &c.CreatedAt, &c.UpdatedAt, &c.Statut)
}

func validerConseil(entree ConseilEntree) string {
	if strings.TrimSpace(entree.Titre) == "" || strings.TrimSpace(entree.Contenu) == "" {
		return "titre et contenu obligatoires"
	}
	if entree.Statut != "brouillon" && entree.Statut != "publie" {
		return "statut invalide (brouillon ou publie attendu)"
	}
	return ""
}

// gererMesConseils liste les articles du salarié connecté, brouillons compris.
func (s *Serveur) gererMesConseils(w http.ResponseWriter, r *http.Request) {
	salarieID, _, ok := s.exigerSalarie(w, r)
	if !ok {
		return
	}

	lignes, err := s.db.Query(`SELECT `+colonnesConseilRedaction+` FROM conseils WHERE auteur_id = $1 ORDER BY updated_at DESC`, salarieID)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	conseils := make([]ConseilRedaction, 0)
	for lignes.Next() {
		var c ConseilRedaction
		if err := scannerConseilRedaction(lignes, &c); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		conseils = append(conseils, c)
	}

	envoyerJSON(w, http.StatusOK, conseils)
}

func (s *Serveur) gererCreationConseil(w http.ResponseWriter, r *http.Request) {
	salarieID, _, ok := s.exigerSalarie(w, r)
	if !ok {
		return
	}

	var entree ConseilEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if message := validerConseil(entree); message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}

	var c ConseilRedaction
	err := scannerConseilRedaction(s.db.QueryRow(`
		INSERT INTO conseils (titre, contenu, categorie, statut, auteur_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING `+colonnesConseilRedaction,
		entree.Titre, entree.Contenu, entree.Categorie, entree.Statut, salarieID,
	), &c)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusCreated, c)
}

// gererModificationConseil met à jour un article de son auteur et garde l'ancien contenu
// dans conseils_historique, dans la même transaction.
func (s *Serveur) gererModificationConseil(w http.ResponseWriter, r *http.Request) {
	salarieID, _, ok := s.exigerSalarie(w, r)
	if !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	var entree ConseilEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if message := validerConseil(entree); message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}

	tx, err := s.db.Begin()
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer tx.Rollback()

	resultat, err := tx.Exec(`
		INSERT INTO conseils_historique (conseil_id, contenu)
		SELECT id, contenu FROM conseils WHERE id = $1 AND auteur_id = $2`, id, salarieID)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if lignes, _ := resultat.RowsAffected(); lignes == 0 {
		envoyerErreur(w, http.StatusNotFound, "conseil introuvable parmi vos articles")
		return
	}

	var c ConseilRedaction
	err = scannerConseilRedaction(tx.QueryRow(`
		UPDATE conseils SET titre = $1, contenu = $2, categorie = $3, statut = $4, updated_at = now()
		WHERE id = $5
		RETURNING `+colonnesConseilRedaction,
		entree.Titre, entree.Contenu, entree.Categorie, entree.Statut, id,
	), &c)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	if err := tx.Commit(); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, c)
}
