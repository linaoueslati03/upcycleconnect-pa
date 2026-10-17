package api

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"math/big"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

const longueurMotDePasseTemporaire = 16

var famillesCaracteres = []string{
	"ABCDEFGHJKLMNPQRSTUVWXYZ",
	"abcdefghijkmnopqrstuvwxyz",
	"23456789",
	"!@#$%&*?-+",
}

var errUtilisateurIntrouvable = errors.New("utilisateur introuvable")

func caractereAuHasard(alphabet string) (byte, error) {
	position, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
	if err != nil {
		return 0, err
	}
	return alphabet[position.Int64()], nil
}

func genererMotDePasseTemporaire() (string, error) {
	tousLesCaracteres := ""
	motDePasse := make([]byte, 0, longueurMotDePasseTemporaire)
	for _, famille := range famillesCaracteres {
		tousLesCaracteres += famille
		caractere, err := caractereAuHasard(famille)
		if err != nil {
			return "", err
		}
		motDePasse = append(motDePasse, caractere)
	}
	for len(motDePasse) < longueurMotDePasseTemporaire {
		caractere, err := caractereAuHasard(tousLesCaracteres)
		if err != nil {
			return "", err
		}
		motDePasse = append(motDePasse, caractere)
	}

	for i := len(motDePasse) - 1; i > 0; i-- {
		position, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		j := position.Int64()
		motDePasse[i], motDePasse[j] = motDePasse[j], motDePasse[i]
	}
	return string(motDePasse), nil
}

func (s *Serveur) reinitialiserMotDePasse(utilisateurID int) (string, error) {
	motDePasse, err := genererMotDePasseTemporaire()
	if err != nil {
		return "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(motDePasse), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	resultat, err := tx.Exec(
		"UPDATE utilisateurs SET mot_de_passe_hash = $1, doit_changer_mot_de_passe = TRUE WHERE id = $2",
		string(hash), utilisateurID,
	)
	if err != nil {
		return "", err
	}
	if lignes, _ := resultat.RowsAffected(); lignes == 0 {
		return "", errUtilisateurIntrouvable
	}
	if _, err := tx.Exec("DELETE FROM sessions WHERE utilisateur_id = $1", utilisateurID); err != nil {
		return "", err
	}
	if err := tx.Commit(); err != nil {
		return "", err
	}
	return motDePasse, nil
}

func (s *Serveur) ReinitialiserMotDePasseParEmail(email string) (string, error) {
	var utilisateurID int
	err := s.db.QueryRow("SELECT id FROM utilisateurs WHERE email = $1", email).Scan(&utilisateurID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errUtilisateurIntrouvable
	}
	if err != nil {
		return "", err
	}
	return s.reinitialiserMotDePasse(utilisateurID)
}

func (s *Serveur) gererReinitialisationMotDePasse(w http.ResponseWriter, r *http.Request) {
	adminID, ok := s.exigerRole(w, r, roleAdministrateur)
	if !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}
	if id == adminID {
		envoyerErreur(w, http.StatusBadRequest, "changez votre propre mot de passe depuis « Mon compte »")
		return
	}

	motDePasse, err := s.reinitialiserMotDePasse(id)
	if errors.Is(err, errUtilisateurIntrouvable) {
		envoyerErreur(w, http.StatusNotFound, "utilisateur introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]string{"mot_de_passe_temporaire": motDePasse})
}
