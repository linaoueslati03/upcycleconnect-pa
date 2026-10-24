package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

var statutsConteneur = []string{"vide", "occupe", "maintenance"}

type ConteneurAdmin struct {
	ID               int     `json:"id"`
	Site             string  `json:"site"`
	Statut           string  `json:"statut"`
	DepotID          *int    `json:"depot_id"`
	DepotStatut      *string `json:"depot_statut"`
	DescriptionObjet *string `json:"description_objet"`
	CodeOuverture    *string `json:"code_ouverture"`
	CodeBarre        *string `json:"code_barre"`
	Deposant         *string `json:"deposant"`
}

type HistoriqueConteneur struct {
	ID             int       `json:"id"`
	AncienStatut   *string   `json:"ancien_statut"`
	NouveauStatut  string    `json:"nouveau_statut"`
	DepotID        *int      `json:"depot_id"`
	Auteur         *string   `json:"auteur"`
	Motif          *string   `json:"motif"`
	DateChangement time.Time `json:"date_changement"`
}

type ConteneurEntree struct {
	Site string `json:"site"`
}

type ConteneurStatutEntree struct {
	Statut string `json:"statut"`
	Motif  string `json:"motif"`
}

func statutConteneurValide(statut string) bool {
	for _, s := range statutsConteneur {
		if s == statut {
			return true
		}
	}
	return false
}

func journaliserConteneur(tx *sql.Tx, conteneurID int, ancienStatut *string, nouveauStatut string, depotID *int, auteurID int, motif string) error {
	var motifEnBase *string
	if motif != "" {
		motifEnBase = &motif
	}
	_, err := tx.Exec(`
		INSERT INTO historique_conteneurs (conteneur_id, ancien_statut, nouveau_statut, depot_id, auteur_id, motif)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		conteneurID, ancienStatut, nouveauStatut, depotID, auteurID, motifEnBase,
	)
	return err
}

func changerStatutConteneur(tx *sql.Tx, conteneurID int, nouveauStatut string, depotID int, auteurID int) error {
	var ancienStatut string
	if err := tx.QueryRow("SELECT statut FROM conteneurs WHERE id = $1 FOR UPDATE", conteneurID).Scan(&ancienStatut); err != nil {
		return err
	}
	if ancienStatut == nouveauStatut {
		return nil
	}
	if _, err := tx.Exec("UPDATE conteneurs SET statut = $1 WHERE id = $2", nouveauStatut, conteneurID); err != nil {
		return err
	}
	return journaliserConteneur(tx, conteneurID, &ancienStatut, nouveauStatut, &depotID, auteurID, "")
}

func (s *Serveur) gererListeConteneursAdmin(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	lignes, err := s.db.Query(`
		SELECT c.id, c.site, c.statut, d.id, d.statut, d.description_objet, d.code_ouverture, d.code_barre,
		       u.prenom || ' ' || u.nom
		FROM conteneurs c
		LEFT JOIN depots d ON d.conteneur_id = c.id AND d.statut IN ('validee', 'deposee')
		LEFT JOIN utilisateurs u ON u.id = d.utilisateur_id
		ORDER BY c.site, c.id`)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	conteneurs := make([]ConteneurAdmin, 0)
	for lignes.Next() {
		var c ConteneurAdmin
		if err := lignes.Scan(&c.ID, &c.Site, &c.Statut, &c.DepotID, &c.DepotStatut, &c.DescriptionObjet,
			&c.CodeOuverture, &c.CodeBarre, &c.Deposant); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		conteneurs = append(conteneurs, c)
	}

	envoyerJSON(w, http.StatusOK, conteneurs)
}

func (s *Serveur) gererCreationConteneur(w http.ResponseWriter, r *http.Request) {
	adminID, ok := s.exigerRole(w, r, roleAdministrateur)
	if !ok {
		return
	}

	var entree ConteneurEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	entree.Site = strings.TrimSpace(entree.Site)
	if entree.Site == "" {
		envoyerErreur(w, http.StatusBadRequest, "site obligatoire")
		return
	}

	tx, err := s.db.Begin()
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer tx.Rollback()

	var c Conteneur
	err = tx.QueryRow("INSERT INTO conteneurs (site) VALUES ($1) RETURNING id, site, statut", entree.Site).
		Scan(&c.ID, &c.Site, &c.Statut)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	if err := journaliserConteneur(tx, c.ID, nil, c.Statut, nil, adminID, "Mise en service"); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if err := tx.Commit(); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusCreated, c)
}

func (s *Serveur) gererChangementStatutConteneur(w http.ResponseWriter, r *http.Request) {
	adminID, ok := s.exigerRole(w, r, roleAdministrateur)
	if !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var entree ConteneurStatutEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if !statutConteneurValide(entree.Statut) {
		envoyerErreur(w, http.StatusBadRequest, "statut invalide")
		return
	}
	entree.Motif = strings.TrimSpace(entree.Motif)
	if entree.Motif == "" {
		envoyerErreur(w, http.StatusBadRequest, "motif obligatoire")
		return
	}

	tx, err := s.db.Begin()
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer tx.Rollback()

	var ancienStatut string
	err = tx.QueryRow("SELECT statut FROM conteneurs WHERE id = $1 FOR UPDATE", id).Scan(&ancienStatut)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "conteneur introuvable")
		return
	} else if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if ancienStatut == entree.Statut {
		envoyerErreur(w, http.StatusBadRequest, "le conteneur a déjà ce statut")
		return
	}

	if _, err := tx.Exec("UPDATE conteneurs SET statut = $1 WHERE id = $2", entree.Statut, id); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if err := journaliserConteneur(tx, id, &ancienStatut, entree.Statut, nil, adminID, entree.Motif); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if err := tx.Commit(); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]string{"statut": entree.Statut})
}

func (s *Serveur) gererHistoriqueConteneur(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	lignes, err := s.db.Query(`
		SELECT h.id, h.ancien_statut, h.nouveau_statut, h.depot_id, u.prenom || ' ' || u.nom, h.motif, h.date_changement
		FROM historique_conteneurs h
		LEFT JOIN utilisateurs u ON u.id = h.auteur_id
		WHERE h.conteneur_id = $1
		ORDER BY h.date_changement DESC, h.id DESC`, id)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	historique := make([]HistoriqueConteneur, 0)
	for lignes.Next() {
		var h HistoriqueConteneur
		if err := lignes.Scan(&h.ID, &h.AncienStatut, &h.NouveauStatut, &h.DepotID, &h.Auteur, &h.Motif, &h.DateChangement); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		historique = append(historique, h)
	}

	envoyerJSON(w, http.StatusOK, historique)
}
