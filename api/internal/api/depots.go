package api

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Depot struct {
	ID                          int        `json:"id"`
	UtilisateurID               int        `json:"utilisateur_id"`
	ConteneurID                 *int       `json:"conteneur_id"`
	DescriptionObjet            string     `json:"description_objet"`
	CodeOuverture               *string    `json:"code_ouverture"`
	CodeBarre                   *string    `json:"code_barre"`
	Statut                      string     `json:"statut"`
	ProfessionnelRecuperateurID *int       `json:"professionnel_recuperateur_id"`
	DateDemande                 time.Time  `json:"date_demande"`
	DateValidation              *time.Time `json:"date_validation"`
	DateRecuperation            *time.Time `json:"date_recuperation"`
}

type DepotEntree struct {
	ConteneurID      *int   `json:"conteneur_id"`
	DescriptionObjet string `json:"description_objet"`
}

type DepotChangementStatut struct {
	Statut                      string `json:"statut"`
	ProfessionnelRecuperateurID *int   `json:"professionnel_recuperateur_id"`
}

var ordreStatutsDepot = map[string]int{
	"demande":   0,
	"validee":   1,
	"deposee":   2,
	"recuperee": 3,
}

func scannerDepot(ligne interface{ Scan(...any) error }, d *Depot) error {
	var conteneurID, professionnelID sql.NullInt64
	var codeOuverture, codeBarre sql.NullString
	var dateValidation, dateRecuperation sql.NullTime

	err := ligne.Scan(
		&d.ID, &d.UtilisateurID, &conteneurID, &d.DescriptionObjet,
		&codeOuverture, &codeBarre, &d.Statut, &professionnelID,
		&d.DateDemande, &dateValidation, &dateRecuperation,
	)
	if err != nil {
		return err
	}
	if conteneurID.Valid {
		v := int(conteneurID.Int64)
		d.ConteneurID = &v
	}
	if professionnelID.Valid {
		v := int(professionnelID.Int64)
		d.ProfessionnelRecuperateurID = &v
	}
	if codeOuverture.Valid {
		d.CodeOuverture = &codeOuverture.String
	}
	if codeBarre.Valid {
		d.CodeBarre = &codeBarre.String
	}
	if dateValidation.Valid {
		d.DateValidation = &dateValidation.Time
	}
	if dateRecuperation.Valid {
		d.DateRecuperation = &dateRecuperation.Time
	}
	return nil
}

const colonnesDepot = `id, utilisateur_id, conteneur_id, description_objet, code_ouverture, code_barre,
	statut, professionnel_recuperateur_id, date_demande, date_validation, date_recuperation`

func (s *Serveur) gererCreationDepot(w http.ResponseWriter, r *http.Request) {
	utilisateurID, ok := s.exigerRole(w, r, roleParticulier, roleAdministrateur)
	if !ok {
		return
	}

	var entree DepotEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if entree.DescriptionObjet == "" {
		envoyerErreur(w, http.StatusBadRequest, "description_objet obligatoire")
		return
	}
	if entree.ConteneurID != nil {
		var statutConteneur string
		err := s.db.QueryRow("SELECT statut FROM conteneurs WHERE id = $1", *entree.ConteneurID).Scan(&statutConteneur)
		if errors.Is(err, sql.ErrNoRows) {
			envoyerErreur(w, http.StatusBadRequest, "conteneur introuvable")
			return
		} else if err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		if statutConteneur == "maintenance" {
			envoyerErreur(w, http.StatusConflict, "ce conteneur est en maintenance, choisissez-en un autre")
			return
		}
	}

	var d Depot
	ligne := s.db.QueryRow(`
		INSERT INTO depots (utilisateur_id, conteneur_id, description_objet)
		VALUES ($1, $2, $3)
		RETURNING `+colonnesDepot,
		utilisateurID, entree.ConteneurID, entree.DescriptionObjet,
	)
	if err := scannerDepot(ligne, &d); err != nil {
		gererErreurPostgres(w, err)
		return
	}

	envoyerJSON(w, http.StatusCreated, d)
}

func (s *Serveur) gererMesDepots(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	lignes, err := s.db.Query(`SELECT `+colonnesDepot+` FROM depots WHERE utilisateur_id = $1 ORDER BY date_demande DESC`, utilisateurID)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	depots := make([]Depot, 0)
	for lignes.Next() {
		var d Depot
		if err := scannerDepot(lignes, &d); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		depots = append(depots, d)
	}

	envoyerJSON(w, http.StatusOK, depots)
}

func (s *Serveur) gererDetailDepot(w http.ResponseWriter, r *http.Request) {
	utilisateurID, role := s.roleConnecte(r)
	if utilisateurID == 0 {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var d Depot
	ligne := s.db.QueryRow(`SELECT `+colonnesDepot+` FROM depots WHERE id = $1`, id)
	if err := scannerDepot(ligne, &d); errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "dépôt introuvable")
		return
	} else if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	estProprietaire := d.UtilisateurID == utilisateurID
	estRecuperateur := d.ProfessionnelRecuperateurID != nil && *d.ProfessionnelRecuperateurID == utilisateurID
	estEquipe := role == roleSalarie || role == roleAdministrateur
	if !estProprietaire && !estRecuperateur && !estEquipe {
		envoyerErreur(w, http.StatusForbidden, "accès non autorisé à ce dépôt")
		return
	}

	envoyerJSON(w, http.StatusOK, d)
}

func genererCodeDepot(longueurOctets int) (string, error) {
	octets := make([]byte, longueurOctets)
	if _, err := rand.Read(octets); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", octets), nil
}

func (s *Serveur) gererChangementStatutDepot(w http.ResponseWriter, r *http.Request) {

	appelantID, ok := s.exigerRole(w, r, roleSalarie, roleAdministrateur, roleProfessionnel)
	if !ok {
		return
	}
	_, roleAppelant := s.roleConnecte(r)

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var entree DepotChangementStatut
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}

	nouvelIndex, statutConnu := ordreStatutsDepot[entree.Statut]
	if !statutConnu {
		envoyerErreur(w, http.StatusBadRequest, "statut invalide")
		return
	}

	if roleAppelant == roleProfessionnel {
		if entree.Statut != "recuperee" {
			envoyerErreur(w, http.StatusForbidden, "un professionnel peut seulement enregistrer la récupération")
			return
		}

		entree.ProfessionnelRecuperateurID = &appelantID
	} else if entree.Statut == "recuperee" && entree.ProfessionnelRecuperateurID != nil {
		var roleRecuperateur string
		err := s.db.QueryRow("SELECT r.code FROM utilisateurs u JOIN roles r ON r.id = u.role_id WHERE u.id = $1",
			*entree.ProfessionnelRecuperateurID).Scan(&roleRecuperateur)
		if err != nil || roleRecuperateur != roleProfessionnel {
			envoyerErreur(w, http.StatusBadRequest, "professionnel_recuperateur_id doit désigner un compte professionnel")
			return
		}
	}

	var statutActuel string
	if err := s.db.QueryRow("SELECT statut FROM depots WHERE id = $1", id).Scan(&statutActuel); errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "dépôt introuvable")
		return
	} else if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	if nouvelIndex != ordreStatutsDepot[statutActuel]+1 {
		envoyerErreur(w, http.StatusBadRequest, "transition de statut invalide (l'ordre est demande -> validee -> deposee -> recuperee)")
		return
	}

	tx, err := s.db.Begin()
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer tx.Rollback()

	var d Depot
	var ligne *sql.Row

	switch entree.Statut {
	case "validee":
		if message := verifierConteneurDisponible(tx, id); message != "" {
			envoyerErreur(w, http.StatusConflict, message)
			return
		}
		codeOuverture, err1 := genererCodeDepot(3)
		codeBarre, err2 := genererCodeDepot(8)
		if err1 != nil || err2 != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		ligne = tx.QueryRow(`
			UPDATE depots SET statut = 'validee', code_ouverture = $1, code_barre = $2, date_validation = now()
			WHERE id = $3 AND statut = 'demande'
			RETURNING `+colonnesDepot,
			codeOuverture, codeBarre, id,
		)
	case "deposee":
		ligne = tx.QueryRow(`
			UPDATE depots SET statut = 'deposee' WHERE id = $1 AND statut = 'validee'
			RETURNING `+colonnesDepot, id)
	case "recuperee":
		if entree.ProfessionnelRecuperateurID == nil {
			envoyerErreur(w, http.StatusBadRequest, "professionnel_recuperateur_id obligatoire")
			return
		}
		ligne = tx.QueryRow(`
			UPDATE depots SET statut = 'recuperee', professionnel_recuperateur_id = $1, date_recuperation = now()
			WHERE id = $2 AND statut = 'deposee'
			RETURNING `+colonnesDepot,
			entree.ProfessionnelRecuperateurID, id,
		)
	}

	if err := scannerDepot(ligne, &d); errors.Is(err, sql.ErrNoRows) {

		envoyerErreur(w, http.StatusConflict, "le dépôt a déjà changé de statut")
		return
	} else if err != nil {
		gererErreurPostgres(w, err)
		return
	}

	if d.ConteneurID != nil {
		nouveauStatutConteneur := ""
		switch entree.Statut {
		case "deposee":
			nouveauStatutConteneur = "occupe"
		case "recuperee":
			nouveauStatutConteneur = "vide"
		}
		if nouveauStatutConteneur != "" {
			if err := changerStatutConteneur(tx, *d.ConteneurID, nouveauStatutConteneur, d.ID, appelantID); err != nil {
				envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
				return
			}
		}
	}

	if entree.Statut == "recuperee" {
		if err := ajouterPointsScore(tx, d.UtilisateurID, pointsDepotRecupere, "Dépôt récupéré"); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, d)
}

func verifierConteneurDisponible(tx *sql.Tx, depotID int) string {
	var conteneurID sql.NullInt64
	if err := tx.QueryRow("SELECT conteneur_id FROM depots WHERE id = $1", depotID).Scan(&conteneurID); err != nil {
		return "dépôt introuvable"
	}
	if !conteneurID.Valid {
		return "aucun conteneur n'est associé à ce dépôt"
	}

	var statut string
	if err := tx.QueryRow("SELECT statut FROM conteneurs WHERE id = $1 FOR UPDATE", conteneurID.Int64).Scan(&statut); err != nil {
		return "conteneur introuvable"
	}
	if statut == "maintenance" {
		return "ce conteneur est en maintenance"
	}

	var occupe bool
	err := tx.QueryRow(`
		SELECT EXISTS (SELECT 1 FROM depots WHERE conteneur_id = $1 AND statut IN ('validee', 'deposee'))`,
		conteneurID.Int64).Scan(&occupe)
	if err != nil {
		return "erreur serveur"
	}
	if occupe {
		return "ce conteneur est déjà réservé ou occupé par un autre objet"
	}
	return ""
}

func (s *Serveur) gererListeDepots(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleSalarie, roleAdministrateur); !ok {
		return
	}

	requete := `SELECT ` + colonnesDepot + ` FROM depots`
	args := []any{}
	if statut := r.URL.Query().Get("statut"); statut != "" {
		requete += ` WHERE statut = $1`
		args = append(args, statut)
	}
	requete += ` ORDER BY date_demande ASC`

	lignes, err := s.db.Query(requete, args...)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	depots := make([]Depot, 0)
	for lignes.Next() {
		var d Depot
		if err := scannerDepot(lignes, &d); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		depots = append(depots, d)
	}

	envoyerJSON(w, http.StatusOK, depots)
}
