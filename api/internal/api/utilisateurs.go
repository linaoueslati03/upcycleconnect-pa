package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Utilisateur struct {
	ID               int       `json:"id"`
	RoleID           int       `json:"role_id"`
	Nom              string    `json:"nom"`
	Prenom           string    `json:"prenom"`
	Email            string    `json:"email"`
	Statut           string    `json:"statut"`
	LanguePrefereeID *int      `json:"langue_preferee_id"`
	UpcyclingScore   int       `json:"upcycling_score"`
	DateCreation     time.Time `json:"date_creation"`
}

type UtilisateurEntree struct {
	RoleID           int    `json:"role_id"`
	Nom              string `json:"nom"`
	Prenom           string `json:"prenom"`
	Email            string `json:"email"`
	MotDePasse       string `json:"mot_de_passe"`
	LanguePrefereeID *int   `json:"langue_preferee_id"`
}

func (s *Serveur) gererListeUtilisateurs(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	lignes, err := s.db.Query(`
		SELECT id, role_id, nom, prenom, email, statut, langue_preferee_id, upcycling_score, date_creation
		FROM utilisateurs ORDER BY id`)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	utilisateurs := make([]Utilisateur, 0)
	for lignes.Next() {
		var u Utilisateur
		var langueID sql.NullInt64
		if err := lignes.Scan(&u.ID, &u.RoleID, &u.Nom, &u.Prenom, &u.Email, &u.Statut, &langueID, &u.UpcyclingScore, &u.DateCreation); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		if langueID.Valid {
			v := int(langueID.Int64)
			u.LanguePrefereeID = &v
		}
		utilisateurs = append(utilisateurs, u)
	}

	envoyerJSON(w, http.StatusOK, utilisateurs)
}

func (s *Serveur) gererDetailUtilisateur(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var u Utilisateur
	var langueID sql.NullInt64
	err = s.db.QueryRow(`
		SELECT id, role_id, nom, prenom, email, statut, langue_preferee_id, upcycling_score, date_creation
		FROM utilisateurs WHERE id = $1`, id,
	).Scan(&u.ID, &u.RoleID, &u.Nom, &u.Prenom, &u.Email, &u.Statut, &langueID, &u.UpcyclingScore, &u.DateCreation)

	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "utilisateur introuvable")
		return
	}
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

func (s *Serveur) gererModificationUtilisateur(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var entree UtilisateurEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}

	if message := validerEntree(entree, false); message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}

	u := Utilisateur{
		ID:               id,
		RoleID:           entree.RoleID,
		Nom:              entree.Nom,
		Prenom:           entree.Prenom,
		Email:            entree.Email,
		LanguePrefereeID: entree.LanguePrefereeID,
	}

	err = s.db.QueryRow(`
		UPDATE utilisateurs SET role_id = $1, nom = $2, prenom = $3, email = $4, langue_preferee_id = $5
		WHERE id = $6
		RETURNING statut, upcycling_score, date_creation`,
		entree.RoleID, entree.Nom, entree.Prenom, entree.Email, entree.LanguePrefereeID, id,
	).Scan(&u.Statut, &u.UpcyclingScore, &u.DateCreation)

	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "utilisateur introuvable")
		return
	}
	if err != nil {
		gererErreurPostgres(w, err)
		return
	}

	if err := s.creerProfilSalarie(u.ID, u.RoleID); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, u)
}

func (s *Serveur) gererSuppressionUtilisateur(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	resultat, err := s.db.Exec("DELETE FROM utilisateurs WHERE id = $1", id)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	lignesAffectees, _ := resultat.RowsAffected()
	if lignesAffectees == 0 {
		envoyerErreur(w, http.StatusNotFound, "utilisateur introuvable")
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]string{"message": "utilisateur supprimé"})
}

func validerEntree(entree UtilisateurEntree, motDePasseObligatoire bool) string {
	if strings.TrimSpace(entree.Nom) == "" || strings.TrimSpace(entree.Prenom) == "" || strings.TrimSpace(entree.Email) == "" {
		return "champs obligatoires manquants"
	}
	if motDePasseObligatoire && strings.TrimSpace(entree.MotDePasse) == "" {
		return "champs obligatoires manquants"
	}
	if !strings.Contains(entree.Email, "@") {
		return "email invalide"
	}
	if entree.RoleID <= 0 {
		return "role_id invalide"
	}
	return ""
}

func gererErreurPostgres(w http.ResponseWriter, err error) {
	var erreurPg *pq.Error
	if errors.As(err, &erreurPg) {
		switch erreurPg.Code {
		case "23505":
			envoyerErreur(w, http.StatusConflict, "cet email est déjà utilisé")
			return
		case "23503":
			envoyerErreur(w, http.StatusBadRequest, "role_id ou langue_preferee_id invalide")
			return
		}
	}
	envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
}

func (s *Serveur) gererCreationUtilisateur(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	var entree UtilisateurEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}

	if message := validerEntree(entree, true); message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(entree.MotDePasse), bcrypt.DefaultCost)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	u := Utilisateur{
		RoleID:           entree.RoleID,
		Nom:              entree.Nom,
		Prenom:           entree.Prenom,
		Email:            entree.Email,
		LanguePrefereeID: entree.LanguePrefereeID,
	}

	err = s.db.QueryRow(`
		INSERT INTO utilisateurs (role_id, nom, prenom, email, mot_de_passe_hash, langue_preferee_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, statut, upcycling_score, date_creation`,
		entree.RoleID, entree.Nom, entree.Prenom, entree.Email, string(hash), entree.LanguePrefereeID,
	).Scan(&u.ID, &u.Statut, &u.UpcyclingScore, &u.DateCreation)

	if err != nil {
		gererErreurPostgres(w, err)
		return
	}

	if err := s.creerProfilSalarie(u.ID, u.RoleID); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusCreated, u)
}

// creerProfilSalarie ajoute la ligne salaries d'un compte qui a le rôle salarié
// (sans effet pour les autres rôles, ni si le profil existe déjà).
func (s *Serveur) creerProfilSalarie(utilisateurID, roleID int) error {
	_, err := s.db.Exec(`
		INSERT INTO salaries (utilisateur_id)
		SELECT $1 FROM roles WHERE id = $2 AND code = $3
		ON CONFLICT (utilisateur_id) DO NOTHING`,
		utilisateurID, roleID, roleSalarie,
	)
	return err
}

// CompteEntree est le corps de POST /api/comptes : comme UtilisateurEntree,
// mais le rôle est donné par son code et seuls particulier et professionnel sont permis.
type CompteEntree struct {
	Role             string `json:"role"`
	Nom              string `json:"nom"`
	Prenom           string `json:"prenom"`
	Email            string `json:"email"`
	MotDePasse       string `json:"mot_de_passe"`
	LanguePrefereeID *int   `json:"langue_preferee_id"`
}

// gererCreationCompte crée un compte depuis le site public. Contrairement à
// POST /api/utilisateurs (réservé aux administrateurs), on ne peut pas y choisir
// le rôle salarié ou administrateur.
func (s *Serveur) gererCreationCompte(w http.ResponseWriter, r *http.Request) {
	var entree CompteEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}

	if entree.Role == "" {
		entree.Role = roleParticulier
	}
	if entree.Role != roleParticulier && entree.Role != roleProfessionnel {
		envoyerErreur(w, http.StatusBadRequest, "rôle invalide (particulier ou professionnel attendu)")
		return
	}

	var roleID int
	if err := s.db.QueryRow("SELECT id FROM roles WHERE code = $1", entree.Role).Scan(&roleID); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	utilisateur := UtilisateurEntree{
		RoleID:           roleID,
		Nom:              entree.Nom,
		Prenom:           entree.Prenom,
		Email:            entree.Email,
		MotDePasse:       entree.MotDePasse,
		LanguePrefereeID: entree.LanguePrefereeID,
	}
	if message := validerEntree(utilisateur, true); message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(entree.MotDePasse), bcrypt.DefaultCost)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	u := Utilisateur{
		RoleID:           roleID,
		Nom:              entree.Nom,
		Prenom:           entree.Prenom,
		Email:            entree.Email,
		LanguePrefereeID: entree.LanguePrefereeID,
	}
	err = s.db.QueryRow(`
		INSERT INTO utilisateurs (role_id, nom, prenom, email, mot_de_passe_hash, langue_preferee_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, statut, upcycling_score, date_creation`,
		roleID, entree.Nom, entree.Prenom, entree.Email, string(hash), entree.LanguePrefereeID,
	).Scan(&u.ID, &u.Statut, &u.UpcyclingScore, &u.DateCreation)
	if err != nil {
		gererErreurPostgres(w, err)
		return
	}

	envoyerJSON(w, http.StatusCreated, u)
}
