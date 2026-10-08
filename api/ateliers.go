package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/lib/pq"
)

type Atelier struct {
	ID            int        `json:"id"`
	Titre         string     `json:"titre"`
	Description   *string    `json:"description"`
	DateDebut     time.Time  `json:"date_debut"`
	DateFin       *time.Time `json:"date_fin"`
	Lieu          *string    `json:"lieu"`
	Statut        string     `json:"statut"`
	CreateurID    int        `json:"createur_id"`
	ResponsableID *int       `json:"responsable_id"`
}

type AtelierEntree struct {
	Titre       string     `json:"titre"`
	Description string     `json:"description"`
	DateDebut   time.Time  `json:"date_debut"`
	DateFin     *time.Time `json:"date_fin"`
	Lieu        string     `json:"lieu"`
	Statut      string     `json:"statut"`
	CreateurID  int        `json:"createur_id"`
}

const colonnesAtelier = `id, titre, description, date_debut, date_fin, lieu, statut, createur_id, responsable_id`

func scannerAtelier(ligne interface{ Scan(...any) error }, a *Atelier) error {
	var description, lieu sql.NullString
	var dateFin sql.NullTime
	var responsableID sql.NullInt64

	err := ligne.Scan(&a.ID, &a.Titre, &description, &a.DateDebut, &dateFin, &lieu, &a.Statut, &a.CreateurID, &responsableID)
	if err != nil {
		return err
	}
	if description.Valid {
		a.Description = &description.String
	}
	if dateFin.Valid {
		a.DateFin = &dateFin.Time
	}
	if lieu.Valid {
		a.Lieu = &lieu.String
	}
	if responsableID.Valid {
		v := int(responsableID.Int64)
		a.ResponsableID = &v
	}
	return nil
}

func gererListeAteliers(w http.ResponseWriter, r *http.Request) {
	lignes, err := db.Query(`SELECT ` + colonnesAtelier + ` FROM ateliers ORDER BY date_debut ASC`)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	ateliers := make([]Atelier, 0)
	for lignes.Next() {
		var a Atelier
		if err := scannerAtelier(lignes, &a); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		ateliers = append(ateliers, a)
	}

	envoyerJSON(w, http.StatusOK, ateliers)
}

func gererCreationAtelier(w http.ResponseWriter, r *http.Request) {
	var entree AtelierEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if strings.TrimSpace(entree.Titre) == "" {
		envoyerErreur(w, http.StatusBadRequest, "titre obligatoire")
		return
	}
	if entree.DateDebut.IsZero() {
		envoyerErreur(w, http.StatusBadRequest, "date_debut obligatoire")
		return
	}

	// Même règle que pour les événements : un salarié crée un atelier en brouillon ou
	// le soumet à validation ; la publication est faite ensuite par un responsable.
	if entree.Statut != "brouillon" && entree.Statut != "en_attente" {
		entree.Statut = "brouillon"
	}

	var a Atelier
	ligne := db.QueryRow(`
		INSERT INTO ateliers (titre, description, date_debut, date_fin, lieu, statut, createur_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING `+colonnesAtelier,
		entree.Titre, entree.Description, entree.DateDebut, entree.DateFin, entree.Lieu, entree.Statut, entree.CreateurID,
	)
	if err := scannerAtelier(ligne, &a); err != nil {
		// 23503 = clé étrangère : createur_id ne correspond à aucun salarié
		var erreurPg *pq.Error
		if errors.As(err, &erreurPg) && erreurPg.Code == "23503" {
			envoyerErreur(w, http.StatusBadRequest, "createur_id invalide")
			return
		}
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusCreated, a)
}
