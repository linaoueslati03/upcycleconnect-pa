package api

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dureeSession = 24 * time.Hour

// Codes des rôles, identiques à la colonne code de la table roles.
const (
	roleParticulier    = "particulier"
	roleProfessionnel  = "professionnel"
	roleSalarie        = "salarie"
	roleAdministrateur = "administrateur"
)

type identifiants struct {
	Email      string `json:"email"`
	MotDePasse string `json:"mot_de_passe"`
}

func genererToken() (string, error) {
	octets := make([]byte, 32)
	if _, err := rand.Read(octets); err != nil {
		return "", err
	}
	return hex.EncodeToString(octets), nil
}

func (s *Serveur) gererLogin(w http.ResponseWriter, r *http.Request) {
	var creds identifiants
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}

	var utilisateurID int
	var hash, role, langue string
	err := s.db.QueryRow(`
		SELECT u.id, u.mot_de_passe_hash, r.code, COALESCE(l.code, '')
		FROM utilisateurs u
		JOIN roles r ON r.id = u.role_id
		LEFT JOIN langues l ON l.id = u.langue_preferee_id
		WHERE u.email = $1`,
		creds.Email,
	).Scan(&utilisateurID, &hash, &role, &langue)

	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusUnauthorized, "email ou mot de passe incorrect")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(creds.MotDePasse)) != nil {
		envoyerErreur(w, http.StatusUnauthorized, "email ou mot de passe incorrect")
		return
	}

	token, err := genererToken()
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	_, err = s.db.Exec(
		"INSERT INTO sessions (token, utilisateur_id, date_expiration) VALUES ($1, $2, $3)",
		token, utilisateurID, time.Now().Add(dureeSession),
	)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	// Le rôle sert seulement au front à choisir l'espace à afficher : chaque route
	// protégée revérifie le rôle côté API, le front ne décide jamais des droits.
	// langue vaut "" si l'utilisateur n'a pas encore choisi : le front garde alors la langue affichée
	envoyerJSON(w, http.StatusOK, map[string]string{"token": token, "role": role, "langue": langue})
}

func (s *Serveur) gererLogout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		envoyerErreur(w, http.StatusBadRequest, "token manquant")
		return
	}

	if _, err := s.db.Exec("DELETE FROM sessions WHERE token = $1", token); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	envoyerJSON(w, http.StatusOK, map[string]string{"message": "déconnecté"})
}

// utilisateurConnecte retourne l'id de l'utilisateur associé au token dans le header
// Authorization, ou une erreur si le token est absent, inconnu ou expiré.
func (s *Serveur) utilisateurConnecte(r *http.Request) (int, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return 0, errors.New("token manquant")
	}

	var utilisateurID int
	err := s.db.QueryRow(
		"SELECT utilisateur_id FROM sessions WHERE token = $1 AND date_expiration > now()",
		token,
	).Scan(&utilisateurID)
	if err != nil {
		return 0, errors.New("session invalide ou expirée")
	}

	return utilisateurID, nil
}

// exigerRole vérifie que la requête vient d'un utilisateur connecté qui a l'un des rôles
// autorisés. Elle renvoie son id, ou écrit l'erreur (401 non connecté, 403 rôle interdit)
// et renvoie false : le handler doit alors s'arrêter.
func (s *Serveur) exigerRole(w http.ResponseWriter, r *http.Request, rolesAutorises ...string) (int, bool) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return 0, false
	}

	var role string
	err = s.db.QueryRow(
		"SELECT r.code FROM utilisateurs u JOIN roles r ON r.id = u.role_id WHERE u.id = $1",
		utilisateurID,
	).Scan(&role)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return 0, false
	}

	if !slices.Contains(rolesAutorises, role) {
		envoyerErreur(w, http.StatusForbidden, "accès refusé")
		return 0, false
	}

	return utilisateurID, true
}

// exigerSalarie vérifie que l'utilisateur connecté a un profil salarié (table salaries)
// et renvoie son id et s'il est responsable. Même principe d'arrêt que exigerRole.
func (s *Serveur) exigerSalarie(w http.ResponseWriter, r *http.Request) (id int, estResponsable bool, ok bool) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return 0, false, false
	}

	// On vérifie aussi le rôle du compte : si un administrateur a retiré le rôle salarié,
	// l'ancienne ligne de la table salaries ne donne plus aucun droit.
	err = s.db.QueryRow(`
		SELECT s.est_responsable
		FROM salaries s
		JOIN utilisateurs u ON u.id = s.utilisateur_id
		JOIN roles r ON r.id = u.role_id
		WHERE s.utilisateur_id = $1 AND r.code = $2`,
		utilisateurID, roleSalarie,
	).Scan(&estResponsable)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusForbidden, "réservé aux salariés")
		return 0, false, false
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return 0, false, false
	}

	return utilisateurID, estResponsable, true
}

// roleConnecte renvoie l'id et le rôle de l'utilisateur connecté, ou (0, "") pour un
// visiteur. Contrairement à exigerRole, elle n'écrit pas d'erreur : elle sert aux routes
// publiques dont la réponse change selon le visiteur (ex. brouillons visibles du personnel).
func (s *Serveur) roleConnecte(r *http.Request) (int, string) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		return 0, ""
	}
	var role string
	err = s.db.QueryRow(
		"SELECT r.code FROM utilisateurs u JOIN roles r ON r.id = u.role_id WHERE u.id = $1", utilisateurID,
	).Scan(&role)
	if err != nil {
		return 0, ""
	}
	return utilisateurID, role
}
