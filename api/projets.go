package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type EtapeProjet struct {
	ID          int       `json:"id"`
	ProjetID    int       `json:"projet_id"`
	Description string    `json:"description"`
	PhotoURL    *string   `json:"photo_url"`
	Ordre       int       `json:"ordre"`
	Date        time.Time `json:"date"`
}

type Projet struct {
	ID             int           `json:"id"`
	UtilisateurID  int           `json:"utilisateur_id"`
	Titre          string        `json:"titre"`
	Description    string        `json:"description"`
	PartagePublic  bool          `json:"partage_public"`
	Sponsorise     bool          `json:"sponsorise"`
	DateCreation   time.Time     `json:"date_creation"`
	Etapes         []EtapeProjet `json:"etapes,omitempty"`
}

type ProjetEntree struct {
	Titre         string `json:"titre"`
	Description   string `json:"description"`
	PartagePublic bool   `json:"partage_public"`
}

type EtapeEntree struct {
	Description string  `json:"description"`
	PhotoURL    *string `json:"photo_url"`
	Ordre       int     `json:"ordre"`
}

func gererCreationProjet(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	var entree ProjetEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if strings.TrimSpace(entree.Titre) == "" {
		envoyerErreur(w, http.StatusBadRequest, "titre obligatoire")
		return
	}

	var p Projet
	err = db.QueryRow(`
		INSERT INTO projets_upcycling (utilisateur_id, titre, description, partage_public)
		VALUES ($1, $2, $3, $4)
		RETURNING id, utilisateur_id, titre, description, partage_public, sponsorise, date_creation`,
		utilisateurID, entree.Titre, entree.Description, entree.PartagePublic,
	).Scan(&p.ID, &p.UtilisateurID, &p.Titre, &p.Description, &p.PartagePublic, &p.Sponsorise, &p.DateCreation)

	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusCreated, p)
}

func gererMesProjets(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	lignes, err := db.Query(`
		SELECT id, utilisateur_id, titre, description, partage_public, sponsorise, date_creation
		FROM projets_upcycling WHERE utilisateur_id = $1 ORDER BY date_creation DESC`, utilisateurID)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	projets := make([]Projet, 0)
	for lignes.Next() {
		var p Projet
		if err := lignes.Scan(&p.ID, &p.UtilisateurID, &p.Titre, &p.Description, &p.PartagePublic, &p.Sponsorise, &p.DateCreation); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		projets = append(projets, p)
	}

	envoyerJSON(w, http.StatusOK, projets)
}

func gererDetailProjet(w http.ResponseWriter, r *http.Request) {
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var p Projet
	err = db.QueryRow(`
		SELECT id, utilisateur_id, titre, description, partage_public, sponsorise, date_creation
		FROM projets_upcycling WHERE id = $1`, id,
	).Scan(&p.ID, &p.UtilisateurID, &p.Titre, &p.Description, &p.PartagePublic, &p.Sponsorise, &p.DateCreation)

	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "projet introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	if !p.PartagePublic {
		utilisateurID, err := utilisateurConnecte(r)
		if err != nil || utilisateurID != p.UtilisateurID {
			envoyerErreur(w, http.StatusForbidden, "ce projet n'est pas public")
			return
		}
	}

	lignes, err := db.Query(`
		SELECT id, projet_id, description, photo_url, ordre, date
		FROM etapes_projet WHERE projet_id = $1 ORDER BY ordre ASC`, id)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	p.Etapes = make([]EtapeProjet, 0)
	for lignes.Next() {
		var e EtapeProjet
		var photoURL sql.NullString
		if err := lignes.Scan(&e.ID, &e.ProjetID, &e.Description, &photoURL, &e.Ordre, &e.Date); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		if photoURL.Valid {
			e.PhotoURL = &photoURL.String
		}
		p.Etapes = append(p.Etapes, e)
	}

	envoyerJSON(w, http.StatusOK, p)
}

func verifierProprietaireProjet(w http.ResponseWriter, r *http.Request, projetID int) (int, bool) {
	utilisateurID, err := utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return 0, false
	}

	var proprietaireID int
	err = db.QueryRow("SELECT utilisateur_id FROM projets_upcycling WHERE id = $1", projetID).Scan(&proprietaireID)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "projet introuvable")
		return 0, false
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return 0, false
	}
	if proprietaireID != utilisateurID {
		envoyerErreur(w, http.StatusForbidden, "vous n'êtes pas propriétaire de ce projet")
		return 0, false
	}

	return utilisateurID, true
}

func gererModificationProjet(w http.ResponseWriter, r *http.Request) {
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	if _, ok := verifierProprietaireProjet(w, r, id); !ok {
		return
	}

	var entree ProjetEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if strings.TrimSpace(entree.Titre) == "" {
		envoyerErreur(w, http.StatusBadRequest, "titre obligatoire")
		return
	}

	var p Projet
	err = db.QueryRow(`
		UPDATE projets_upcycling SET titre = $1, description = $2, partage_public = $3
		WHERE id = $4
		RETURNING id, utilisateur_id, titre, description, partage_public, sponsorise, date_creation`,
		entree.Titre, entree.Description, entree.PartagePublic, id,
	).Scan(&p.ID, &p.UtilisateurID, &p.Titre, &p.Description, &p.PartagePublic, &p.Sponsorise, &p.DateCreation)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, p)
}

func gererSuppressionProjet(w http.ResponseWriter, r *http.Request) {
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	if _, ok := verifierProprietaireProjet(w, r, id); !ok {
		return
	}

	db.Exec("DELETE FROM projets_upcycling WHERE id = $1", id)
	envoyerJSON(w, http.StatusOK, map[string]string{"message": "projet supprimé"})
}

func gererCreationEtape(w http.ResponseWriter, r *http.Request) {
	projetID, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	if _, ok := verifierProprietaireProjet(w, r, projetID); !ok {
		return
	}

	var entree EtapeEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if strings.TrimSpace(entree.Description) == "" {
		envoyerErreur(w, http.StatusBadRequest, "description obligatoire")
		return
	}

	var e EtapeProjet
	var photoURL sql.NullString
	err = db.QueryRow(`
		INSERT INTO etapes_projet (projet_id, description, photo_url, ordre)
		VALUES ($1, $2, $3, $4)
		RETURNING id, projet_id, description, photo_url, ordre, date`,
		projetID, entree.Description, entree.PhotoURL, entree.Ordre,
	).Scan(&e.ID, &e.ProjetID, &e.Description, &photoURL, &e.Ordre, &e.Date)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if photoURL.Valid {
		e.PhotoURL = &photoURL.String
	}

	envoyerJSON(w, http.StatusCreated, e)
}
