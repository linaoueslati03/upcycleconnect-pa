package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

type MonCompteEntree struct {
	Nom              string `json:"nom"`
	Prenom           string `json:"prenom"`
	Email            string `json:"email"`
	LanguePrefereeID *int   `json:"langue_preferee_id"`
}

func gererMonCompte(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	var u Utilisateur
	var langueID sql.NullInt64
	err = db.QueryRow(`
		SELECT id, role_id, nom, prenom, email, statut, langue_preferee_id, upcycling_score, date_creation
		FROM utilisateurs WHERE id = $1`, utilisateurID,
	).Scan(&u.ID, &u.RoleID, &u.Nom, &u.Prenom, &u.Email, &u.Statut, &langueID, &u.UpcyclingScore, &u.DateCreation)

	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if langueID.Valid {
		v := int(langueID.Int64)
		u.LanguePrefereeID = &v
	}

	envoyerJSON(w, http.StatusOK, u)
}

func gererModificationMonCompte(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	var entree MonCompteEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}

	if strings.TrimSpace(entree.Nom) == "" || strings.TrimSpace(entree.Prenom) == "" || strings.TrimSpace(entree.Email) == "" {
		envoyerErreur(w, http.StatusBadRequest, "champs obligatoires manquants")
		return
	}
	if !strings.Contains(entree.Email, "@") {
		envoyerErreur(w, http.StatusBadRequest, "email invalide")
		return
	}

	var u Utilisateur
	err = db.QueryRow(`
		UPDATE utilisateurs SET nom = $1, prenom = $2, email = $3, langue_preferee_id = $4
		WHERE id = $5
		RETURNING id, role_id, statut, upcycling_score, date_creation`,
		entree.Nom, entree.Prenom, entree.Email, entree.LanguePrefereeID, utilisateurID,
	).Scan(&u.ID, &u.RoleID, &u.Statut, &u.UpcyclingScore, &u.DateCreation)

	if err != nil {
		gererErreurPostgres(w, err)
		return
	}

	u.Nom = entree.Nom
	u.Prenom = entree.Prenom
	u.Email = entree.Email
	u.LanguePrefereeID = entree.LanguePrefereeID

	envoyerJSON(w, http.StatusOK, u)
}

func gererTutorielVu(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	_, err = db.Exec(
		"UPDATE utilisateurs SET tutoriel_vu = TRUE, date_tutoriel_vu = now() WHERE id = $1",
		utilisateurID,
	)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]string{"message": "tutoriel marqué comme vu"})
}
