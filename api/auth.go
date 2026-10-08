package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dureeSession = 24 * time.Hour

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

func gererLogin(w http.ResponseWriter, r *http.Request) {
	var creds identifiants
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}

	var utilisateurID int
	var hash string
	err := db.QueryRow(
		"SELECT id, mot_de_passe_hash FROM utilisateurs WHERE email = $1",
		creds.Email,
	).Scan(&utilisateurID, &hash)

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

	_, err = db.Exec(
		"INSERT INTO sessions (token, utilisateur_id, date_expiration) VALUES ($1, $2, $3)",
		token, utilisateurID, time.Now().Add(dureeSession),
	)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]string{"token": token})
}

func gererLogout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		envoyerErreur(w, http.StatusBadRequest, "token manquant")
		return
	}

	db.Exec("DELETE FROM sessions WHERE token = $1", token)
	envoyerJSON(w, http.StatusOK, map[string]string{"message": "déconnecté"})
}

// utilisateurConnecte retourne l'id de l'utilisateur associé au token dans le header
// Authorization, ou une erreur si le token est absent, inconnu ou expiré.
func utilisateurConnecte(r *http.Request) (int, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return 0, errors.New("token manquant")
	}

	var utilisateurID int
	err := db.QueryRow(
		"SELECT utilisateur_id FROM sessions WHERE token = $1 AND date_expiration > now()",
		token,
	).Scan(&utilisateurID)
	if err != nil {
		return 0, errors.New("session invalide ou expirée")
	}

	return utilisateurID, nil
}
