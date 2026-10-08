package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Prestation struct {
	ID        int     `json:"id"`
	Titre     string  `json:"titre"`
	Categorie string  `json:"categorie"`
	Tarif     float64 `json:"tarif"`
	Statut    string  `json:"statut"`
}

func validerPrestation(p Prestation) string {
	if strings.TrimSpace(p.Titre) == "" || strings.TrimSpace(p.Categorie) == "" {
		return "titre et catégorie obligatoires"
	}
	if p.Tarif < 0 {
		return "le tarif ne peut pas être négatif"
	}
	if p.Statut != "Brouillon" && p.Statut != "Publiée" {
		return "statut invalide (Brouillon ou Publiée attendu)"
	}
	return ""
}

func (s *Serveur) gererListePrestations(w http.ResponseWriter, r *http.Request) {
	lignes, err := s.db.Query("SELECT id, titre, categorie, tarif, statut FROM prestations ORDER BY id")
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	prestations := make([]Prestation, 0)
	for lignes.Next() {
		var p Prestation
		if err := lignes.Scan(&p.ID, &p.Titre, &p.Categorie, &p.Tarif, &p.Statut); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		prestations = append(prestations, p)
	}

	envoyerJSON(w, http.StatusOK, prestations)
}

func (s *Serveur) gererCreationPrestation(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	var p Prestation
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if message := validerPrestation(p); message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}

	err := s.db.QueryRow(
		"INSERT INTO prestations (titre, categorie, tarif, statut) VALUES ($1, $2, $3, $4) RETURNING id",
		p.Titre, p.Categorie, p.Tarif, p.Statut,
	).Scan(&p.ID)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusCreated, p)
}

func (s *Serveur) gererModificationPrestation(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	var p Prestation
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if message := validerPrestation(p); message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}

	err = s.db.QueryRow(
		"UPDATE prestations SET titre = $1, categorie = $2, tarif = $3, statut = $4 WHERE id = $5 RETURNING id",
		p.Titre, p.Categorie, p.Tarif, p.Statut, id,
	).Scan(&p.ID)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "prestation introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, p)
}

func (s *Serveur) gererSuppressionPrestation(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	resultat, err := s.db.Exec("DELETE FROM prestations WHERE id = $1", id)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if lignes, _ := resultat.RowsAffected(); lignes == 0 {
		envoyerErreur(w, http.StatusNotFound, "prestation introuvable")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
