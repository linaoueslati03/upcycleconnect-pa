package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
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

func (s *Serveur) gererListeAteliers(w http.ResponseWriter, r *http.Request) {
	lignes, err := s.db.Query(`SELECT ` + colonnesAtelier + ` FROM ateliers ORDER BY date_debut ASC`)
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

func (s *Serveur) gererCreationAtelier(w http.ResponseWriter, r *http.Request) {
	salarieID, _, ok := s.exigerSalarie(w, r)
	if !ok {
		return
	}

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
	ligne := s.db.QueryRow(`
		INSERT INTO ateliers (titre, description, date_debut, date_fin, lieu, statut, createur_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING `+colonnesAtelier,
		entree.Titre, entree.Description, entree.DateDebut, entree.DateFin, entree.Lieu, entree.Statut, salarieID,
	)
	if err := scannerAtelier(ligne, &a); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusCreated, a)
}

func (s *Serveur) gererDetailAtelier(w http.ResponseWriter, r *http.Request) {
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	var a Atelier
	err = scannerAtelier(s.db.QueryRow(`SELECT `+colonnesAtelier+` FROM ateliers WHERE id = $1`, id), &a)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "atelier introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, a)
}

func (s *Serveur) gererModificationAtelier(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := s.exigerSalarie(w, r); !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	var entree AtelierEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if strings.TrimSpace(entree.Titre) == "" || entree.DateDebut.IsZero() {
		envoyerErreur(w, http.StatusBadRequest, "titre et date_debut obligatoires")
		return
	}

	var a Atelier
	err = scannerAtelier(s.db.QueryRow(`
		UPDATE ateliers SET titre = $1, description = $2, date_debut = $3, date_fin = $4, lieu = $5,
		    statut = COALESCE(NULLIF($6, ''), statut), updated_at = now()
		WHERE id = $7
		RETURNING `+colonnesAtelier,
		entree.Titre, entree.Description, entree.DateDebut, entree.DateFin, entree.Lieu, statutModifiable(entree.Statut), id,
	), &a)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "atelier introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, a)
}

func (s *Serveur) gererSuppressionAtelier(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := s.exigerSalarie(w, r); !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	resultat, err := s.db.Exec("DELETE FROM ateliers WHERE id = $1", id)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if lignes, _ := resultat.RowsAffected(); lignes == 0 {
		envoyerErreur(w, http.StatusNotFound, "atelier introuvable")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Même règle que les événements : seul un responsable publie un atelier en attente,
// et il en devient le responsable.
func (s *Serveur) gererValidationAtelier(w http.ResponseWriter, r *http.Request) {
	responsableID, estResponsable, ok := s.exigerSalarie(w, r)
	if !ok {
		return
	}
	if !estResponsable {
		envoyerErreur(w, http.StatusForbidden, "seul un responsable peut valider un atelier")
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	var a Atelier
	err = scannerAtelier(s.db.QueryRow(`
		UPDATE ateliers SET statut = 'publie', responsable_id = $1, updated_at = now()
		WHERE id = $2 AND statut = 'en_attente'
		RETURNING `+colonnesAtelier,
		responsableID, id,
	), &a)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusConflict, "atelier introuvable ou pas en attente de validation")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, a)
}
