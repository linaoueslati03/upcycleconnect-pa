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
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
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
	if !estProprietaire && !estRecuperateur {
		envoyerErreur(w, http.StatusForbidden, "accès non autorisé à ce dépôt")
		return
	}

	envoyerJSON(w, http.StatusOK, d)
}

// genererCodeDepot produit une chaîne hexadécimale aléatoire de la longueur demandée.
// Simplifié : pas un vrai format code-barre (EAN13...), juste une valeur unique et
// non réutilisable, suffisant pour la traçabilité demandée par le sujet.
func genererCodeDepot(longueurOctets int) (string, error) {
	octets := make([]byte, longueurOctets)
	if _, err := rand.Read(octets); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", octets), nil
}

func (s *Serveur) gererChangementStatutDepot(w http.ResponseWriter, r *http.Request) {
	// Le suivi d'un dépôt (validation, dépôt, récupération) est fait par le personnel
	// ou par le professionnel qui récupère l'objet, jamais par le particulier lui-même.
	if _, ok := s.exigerRole(w, r, roleSalarie, roleAdministrateur, roleProfessionnel); !ok {
		return
	}

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

	var d Depot
	var ligne *sql.Row

	switch entree.Statut {
	case "validee":
		codeOuverture, err1 := genererCodeDepot(3)
		codeBarre, err2 := genererCodeDepot(8)
		if err1 != nil || err2 != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		ligne = s.db.QueryRow(`
			UPDATE depots SET statut = 'validee', code_ouverture = $1, code_barre = $2, date_validation = now()
			WHERE id = $3
			RETURNING `+colonnesDepot,
			codeOuverture, codeBarre, id,
		)
	case "deposee":
		ligne = s.db.QueryRow(`
			UPDATE depots SET statut = 'deposee' WHERE id = $1
			RETURNING `+colonnesDepot, id)
	case "recuperee":
		if entree.ProfessionnelRecuperateurID == nil {
			envoyerErreur(w, http.StatusBadRequest, "professionnel_recuperateur_id obligatoire")
			return
		}
		ligne = s.db.QueryRow(`
			UPDATE depots SET statut = 'recuperee', professionnel_recuperateur_id = $1, date_recuperation = now()
			WHERE id = $2
			RETURNING `+colonnesDepot,
			entree.ProfessionnelRecuperateurID, id,
		)
	}

	if err := scannerDepot(ligne, &d); err != nil {
		gererErreurPostgres(w, err)
		return
	}

	if entree.Statut == "recuperee" {
		ajouterPointsScore(s.db, d.UtilisateurID, pointsDepotRecupere, "Dépôt récupéré")
	}

	envoyerJSON(w, http.StatusOK, d)
}
