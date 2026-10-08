package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// MonCompte est la réponse de GET /api/moi : le compte, plus l'état du tutoriel
// de première connexion (qui ne concerne que l'utilisateur lui-même).
type MonCompte struct {
	Utilisateur
	TutorielVu bool `json:"tutoriel_vu"`
}

type MotDePasseEntree struct {
	AncienMotDePasse  string `json:"ancien_mot_de_passe"`
	NouveauMotDePasse string `json:"nouveau_mot_de_passe"`
}

type MonCompteEntree struct {
	Nom              string `json:"nom"`
	Prenom           string `json:"prenom"`
	Email            string `json:"email"`
	LanguePrefereeID *int   `json:"langue_preferee_id"`
}

func (s *Serveur) gererMonCompte(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	var u MonCompte
	var langueID sql.NullInt64
	err = s.db.QueryRow(`
		SELECT id, role_id, nom, prenom, email, statut, langue_preferee_id, upcycling_score, date_creation, tutoriel_vu
		FROM utilisateurs WHERE id = $1`, utilisateurID,
	).Scan(&u.ID, &u.RoleID, &u.Nom, &u.Prenom, &u.Email, &u.Statut, &langueID, &u.UpcyclingScore, &u.DateCreation, &u.TutorielVu)

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

func (s *Serveur) gererModificationMonCompte(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := s.utilisateurConnecte(r)
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
	err = s.db.QueryRow(`
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

func (s *Serveur) gererTutorielVu(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	_, err = s.db.Exec(
		"UPDATE utilisateurs SET tutoriel_vu = TRUE, date_tutoriel_vu = now() WHERE id = $1",
		utilisateurID,
	)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]string{"message": "tutoriel marqué comme vu"})
}

// gererChangementMotDePasse vérifie l'ancien mot de passe avant d'enregistrer le
// nouveau (haché avec bcrypt, jamais stocké en clair).
func (s *Serveur) gererChangementMotDePasse(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	var entree MotDePasseEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if len(entree.NouveauMotDePasse) < 8 {
		envoyerErreur(w, http.StatusBadRequest, "le nouveau mot de passe doit faire au moins 8 caractères")
		return
	}

	var hashActuel string
	if err := s.db.QueryRow("SELECT mot_de_passe_hash FROM utilisateurs WHERE id = $1", utilisateurID).Scan(&hashActuel); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hashActuel), []byte(entree.AncienMotDePasse)) != nil {
		envoyerErreur(w, http.StatusBadRequest, "ancien mot de passe incorrect")
		return
	}

	nouveauHash, err := bcrypt.GenerateFromPassword([]byte(entree.NouveauMotDePasse), bcrypt.DefaultCost)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if _, err := s.db.Exec("UPDATE utilisateurs SET mot_de_passe_hash = $1 WHERE id = $2", string(nouveauHash), utilisateurID); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]string{"message": "mot de passe modifié"})
}
