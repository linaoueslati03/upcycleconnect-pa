package api

import (
	"database/sql"
	"errors"
	"net/http"
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
