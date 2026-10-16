package api

import (
	"crypto/rand"
	"crypto/sha256"
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

const (
	roleParticulier    = "particulier"
	roleProfessionnel  = "professionnel"
	roleSalarie        = "salarie"
	roleAdministrateur = "administrateur"
)

const (
	statutActif    = "actif"
	statutSuspendu = "suspendu"
)

const abonnementPremium = "premium"

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

func empreinteJeton(token string) string {
	empreinte := sha256.Sum256([]byte(token))
	return hex.EncodeToString(empreinte[:])
}

func (s *Serveur) gererLogin(w http.ResponseWriter, r *http.Request) {
	var creds identifiants
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}

	var utilisateurID int
	var hash, role, langue, statut string
	err := s.db.QueryRow(`
		SELECT u.id, u.mot_de_passe_hash, r.code, COALESCE(l.code, ''), u.statut
		FROM utilisateurs u
		JOIN roles r ON r.id = u.role_id
		LEFT JOIN langues l ON l.id = u.langue_preferee_id
		WHERE u.email = $1`,
		creds.Email,
	).Scan(&utilisateurID, &hash, &role, &langue, &statut)

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

	if statut == statutSuspendu {
		envoyerErreur(w, http.StatusForbidden, "compte suspendu, contactez UpcycleConnect")
		return
	}

	token, err := genererToken()
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	_, err = s.db.Exec(
		"INSERT INTO sessions (empreinte_token, utilisateur_id, date_expiration) VALUES ($1, $2, $3)",
		empreinteJeton(token), utilisateurID, time.Now().Add(dureeSession),
	)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]string{"token": token, "role": role, "langue": langue})
}

func (s *Serveur) gererLogout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		envoyerErreur(w, http.StatusBadRequest, "token manquant")
		return
	}

	if _, err := s.db.Exec("DELETE FROM sessions WHERE empreinte_token = $1", empreinteJeton(token)); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	envoyerJSON(w, http.StatusOK, map[string]string{"message": "déconnecté"})
}

func (s *Serveur) utilisateurConnecte(r *http.Request) (int, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return 0, errors.New("token manquant")
	}

	var utilisateurID int
	err := s.db.QueryRow(`
		SELECT s.utilisateur_id
		FROM sessions s
		JOIN utilisateurs u ON u.id = s.utilisateur_id
		WHERE s.empreinte_token = $1 AND s.date_expiration > now() AND u.statut = $2`,
		empreinteJeton(token), statutActif,
	).Scan(&utilisateurID)
	if err != nil {
		return 0, errors.New("session invalide ou expirée")
	}

	return utilisateurID, nil
}

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

func (s *Serveur) exigerSalarie(w http.ResponseWriter, r *http.Request) (id int, estResponsable bool, ok bool) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return 0, false, false
	}

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

func (s *Serveur) exigerSalarieOuAdmin(w http.ResponseWriter, r *http.Request) (id int, peutValider bool, ok bool) {
	utilisateurID, role := s.roleConnecte(r)
	if role == roleAdministrateur {
		return utilisateurID, true, true
	}
	return s.exigerSalarie(w, r)
}

func (s *Serveur) exigerPremium(w http.ResponseWriter, r *http.Request) (int, bool) {
	utilisateurID, ok := s.exigerRole(w, r, roleProfessionnel, roleAdministrateur)
	if !ok {
		return 0, false
	}

	var role, abonnement string
	err := s.db.QueryRow(
		"SELECT r.code, u.abonnement FROM utilisateurs u JOIN roles r ON r.id = u.role_id WHERE u.id = $1",
		utilisateurID,
	).Scan(&role, &abonnement)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return 0, false
	}

	if role == roleProfessionnel && abonnement != abonnementPremium {
		envoyerErreur(w, http.StatusForbidden, "réservé aux abonnés Premium")
		return 0, false
	}
	return utilisateurID, true
}

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
