package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type Formation struct {
	ID                 int        `json:"id"`
	Titre              string     `json:"titre"`
	Description        string     `json:"description"`
	DateDebut          time.Time  `json:"date_debut"`
	DateFin            *time.Time `json:"date_fin,omitempty"`
	Lieu               string     `json:"lieu"`
	NbPlaces           *int       `json:"nb_places,omitempty"`
	Tarif              float64    `json:"tarif"`
	MaterielNecessaire string     `json:"materiel_necessaire"`
	Statut             string     `json:"statut"`
	CreateurID         int        `json:"createur_id"`
	ValideParID        *int       `json:"valide_par_id,omitempty"`
	MotifRefus         *string    `json:"motif_refus,omitempty"`
}

type FormationEntree struct {
	Titre              string     `json:"titre"`
	Description        string     `json:"description"`
	DateDebut          time.Time  `json:"date_debut"`
	DateFin            *time.Time `json:"date_fin"`
	Lieu               string     `json:"lieu"`
	NbPlaces           *int       `json:"nb_places"`
	Tarif              float64    `json:"tarif"`
	MaterielNecessaire string     `json:"materiel_necessaire"`
	Statut             string     `json:"statut"`
}

const colonnesFormation = `
	id, titre, description, date_debut, date_fin, lieu,
	nb_places, tarif, materiel_necessaire, statut,
	createur_id, valide_par_id, motif_refus
`

func scannerFormation(ligne interface{ Scan(...any) error }, f *Formation) error {
	var description, lieu, materiel sql.NullString
	var dateFin sql.NullTime
	var nbPlaces sql.NullInt64
	var valideParID sql.NullInt64
	var motifRefus sql.NullString

	err := ligne.Scan(
		&f.ID,
		&f.Titre,
		&description,
		&f.DateDebut,
		&dateFin,
		&lieu,
		&nbPlaces,
		&f.Tarif,
		&materiel,
		&f.Statut,
		&f.CreateurID,
		&valideParID,
		&motifRefus,
	)
	if err != nil {
		return err
	}

	f.Description = description.String
	f.Lieu = lieu.String
	f.MaterielNecessaire = materiel.String

	if dateFin.Valid {
		f.DateFin = &dateFin.Time
	}
	if nbPlaces.Valid {
		v := int(nbPlaces.Int64)
		f.NbPlaces = &v
	}
	if valideParID.Valid {
		v := int(valideParID.Int64)
		f.ValideParID = &v
	}
	if motifRefus.Valid {
		f.MotifRefus = &motifRefus.String
	}

	return nil
}

func validerFormation(f FormationEntree) string {
	if strings.TrimSpace(f.Titre) == "" {
		return "titre obligatoire"
	}
	if f.DateDebut.IsZero() {
		return "date_debut obligatoire"
	}
	if f.DateFin != nil && f.DateFin.Before(f.DateDebut) {
		return "date_fin doit être postérieure à date_debut"
	}
	if f.NbPlaces != nil && *f.NbPlaces < 1 {
		return "nb_places doit être supérieur à zéro"
	}
	if f.Tarif < 0 {
		return "tarif ne peut pas être négatif"
	}
	return ""
}

func (s *Serveur) gererListeFormations(w http.ResponseWriter, r *http.Request) {
	requete := `SELECT ` + colonnesFormation + `
		FROM formations WHERE statut = 'publie'
		ORDER BY date_debut ASC`

	if _, role := s.roleConnecte(r); role == roleSalarie || role == roleAdministrateur {
		requete = `SELECT ` + colonnesFormation + `
			FROM formations ORDER BY date_debut ASC`
	}

	lignes, err := s.db.Query(requete)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	formations := make([]Formation, 0)
	for lignes.Next() {
		var f Formation
		if err := scannerFormation(lignes, &f); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		formations = append(formations, f)
	}

	if err := lignes.Err(); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, formations)
}

func (s *Serveur) gererDetailFormation(w http.ResponseWriter, r *http.Request) {
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	var f Formation
	err = scannerFormation(
		s.db.QueryRow(
			`SELECT `+colonnesFormation+` FROM formations WHERE id = $1`,
			id,
		),
		&f,
	)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "formation introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	if f.Statut != "publie" {
		_, role := s.roleConnecte(r)
		if role != roleSalarie && role != roleAdministrateur {
			envoyerErreur(w, http.StatusNotFound, "introuvable")
			return
		}
	}

	envoyerJSON(w, http.StatusOK, f)
}

func (s *Serveur) gererCreationFormation(w http.ResponseWriter, r *http.Request) {
	salarieID, _, ok := s.exigerSalarieOuAdmin(w, r)
	if !ok {
		return
	}

	var entree FormationEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}

	if message := validerFormation(entree); message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}

	if entree.Statut != "brouillon" && entree.Statut != "en_attente" {
		entree.Statut = "brouillon"
	}

	var f Formation
	err := scannerFormation(s.db.QueryRow(`
		INSERT INTO formations (
			titre, description, date_debut, date_fin, lieu,
			nb_places, tarif, materiel_necessaire, statut, createur_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING `+colonnesFormation,
		entree.Titre,
		entree.Description,
		entree.DateDebut,
		entree.DateFin,
		entree.Lieu,
		entree.NbPlaces,
		entree.Tarif,
		entree.MaterielNecessaire,
		entree.Statut,
		salarieID,
	), &f)

	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusCreated, f)
}

func (s *Serveur) gererModificationFormation(w http.ResponseWriter, r *http.Request) {
	salarieID, estResponsable, ok := s.exigerSalarieOuAdmin(w, r)
	if !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	var entree FormationEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}

	if message := validerFormation(entree); message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}

	statutActuel, autorise := s.verifierDroitOffre(
		w, "formations", id, salarieID, estResponsable,
	)
	if !autorise {
		return
	}

	statut := statutApresModificationOffre(
		statutActuel, entree.Statut, estResponsable,
	)

	var f Formation
	err = scannerFormation(s.db.QueryRow(`
		UPDATE formations
		SET titre = $1,
			description = $2,
			date_debut = $3,
			date_fin = $4,
			lieu = $5,
			nb_places = $6,
			tarif = $7,
			materiel_necessaire = $8,
			statut = $9,
			updated_at = now()
		WHERE id = $10
		RETURNING `+colonnesFormation,
		entree.Titre,
		entree.Description,
		entree.DateDebut,
		entree.DateFin,
		entree.Lieu,
		entree.NbPlaces,
		entree.Tarif,
		entree.MaterielNecessaire,
		statut,
		id,
	), &f)

	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "formation introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, f)
}

func (s *Serveur) gererSuppressionFormation(w http.ResponseWriter, r *http.Request) {
	salarieID, estResponsable, ok := s.exigerSalarieOuAdmin(w, r)
	if !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	if _, autorise := s.verifierDroitOffre(
		w, "formations", id, salarieID, estResponsable,
	); !autorise {
		return
	}

	resultat, err := s.db.Exec(`DELETE FROM formations WHERE id = $1`, id)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	lignes, err := resultat.RowsAffected()
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if lignes == 0 {
		envoyerErreur(w, http.StatusNotFound, "formation introuvable")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Serveur) gererValidationFormation(w http.ResponseWriter, r *http.Request) {
	responsableID, estResponsable, ok := s.exigerSalarieOuAdmin(w, r)
	if !ok {
		return
	}
	if !estResponsable {
		envoyerErreur(w, http.StatusForbidden, "seul un responsable peut valider une formation")
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	var f Formation
	err = scannerFormation(s.db.QueryRow(`
		UPDATE formations
		SET statut = 'publie',
			valide_par_id = $1,
			motif_refus = NULL,
			updated_at = now()
		WHERE id = $2 AND statut = 'en_attente'
		RETURNING `+colonnesFormation,
		responsableID, id,
	), &f)

	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusConflict, "formation introuvable ou pas en attente de validation")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, f)
}