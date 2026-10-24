package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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
	EstResponsable   *bool     `json:"est_responsable,omitempty"`
}

type UtilisateurEntree struct {
	RoleID           int    `json:"role_id"`
	Nom              string `json:"nom"`
	Prenom           string `json:"prenom"`
	Email            string `json:"email"`
	MotDePasse       string `json:"mot_de_passe"`
	LanguePrefereeID *int   `json:"langue_preferee_id"`
}

const utilisateursParPage = 20

var echappementRecherche = strings.NewReplacer(`\`, `\\`, "%", `\%`, "_", `\_`)

type PageUtilisateurs struct {
	Utilisateurs []Utilisateur `json:"utilisateurs"`
	Total        int           `json:"total"`
	Page         int           `json:"page"`
	ParPage      int           `json:"par_page"`
}

func filtresUtilisateurs(r *http.Request) (string, []any) {
	conditions := []string{"TRUE"}
	args := []any{}

	if recherche := strings.TrimSpace(r.URL.Query().Get("recherche")); recherche != "" {
		args = append(args, "%"+echappementRecherche.Replace(strings.ToLower(recherche))+"%")
		conditions = append(conditions, fmt.Sprintf("lower(u.nom || ' ' || u.prenom || ' ' || u.email) LIKE $%d", len(args)))
	}
	if roleID, err := strconv.Atoi(r.URL.Query().Get("role_id")); err == nil {
		args = append(args, roleID)
		conditions = append(conditions, fmt.Sprintf("u.role_id = $%d", len(args)))
	}
	if statut := r.URL.Query().Get("statut"); statut == statutActif || statut == statutSuspendu {
		args = append(args, statut)
		conditions = append(conditions, fmt.Sprintf("u.statut = $%d", len(args)))
	}

	return strings.Join(conditions, " AND "), args
}

func (s *Serveur) gererListeUtilisateurs(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	where, args := filtresUtilisateurs(r)

	resultat := PageUtilisateurs{Utilisateurs: make([]Utilisateur, 0), Page: page, ParPage: utilisateursParPage}
	if err := s.db.QueryRow("SELECT count(*) FROM utilisateurs u WHERE "+where, args...).Scan(&resultat.Total); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	args = append(args, utilisateursParPage, (page-1)*utilisateursParPage)
	lignes, err := s.db.Query(fmt.Sprintf(`
		SELECT u.id, u.role_id, u.nom, u.prenom, u.email, u.statut, u.langue_preferee_id, u.upcycling_score,
		       u.date_creation, s.est_responsable
		FROM utilisateurs u LEFT JOIN salaries s ON s.utilisateur_id = u.id
		WHERE %s
		ORDER BY u.nom, u.prenom, u.id
		LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args)), args...)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	for lignes.Next() {
		var u Utilisateur
		var langueID sql.NullInt64
		var estResponsable sql.NullBool
		if err := lignes.Scan(&u.ID, &u.RoleID, &u.Nom, &u.Prenom, &u.Email, &u.Statut, &langueID, &u.UpcyclingScore,
			&u.DateCreation, &estResponsable); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		if langueID.Valid {
			v := int(langueID.Int64)
			u.LanguePrefereeID = &v
		}
		if estResponsable.Valid {
			u.EstResponsable = &estResponsable.Bool
		}
		resultat.Utilisateurs = append(resultat.Utilisateurs, u)
	}

	envoyerJSON(w, http.StatusOK, resultat)
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

type StatutCompteEntree struct {
	Statut string `json:"statut"`
}

func (s *Serveur) gererChangementStatutUtilisateur(w http.ResponseWriter, r *http.Request) {
	adminID, ok := s.exigerRole(w, r, roleAdministrateur)
	if !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var entree StatutCompteEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if entree.Statut != statutActif && entree.Statut != statutSuspendu {
		envoyerErreur(w, http.StatusBadRequest, "statut invalide (actif ou suspendu)")
		return
	}
	if id == adminID && entree.Statut == statutSuspendu {
		envoyerErreur(w, http.StatusBadRequest, "vous ne pouvez pas suspendre votre propre compte")
		return
	}

	resultat, err := s.db.Exec("UPDATE utilisateurs SET statut = $1 WHERE id = $2", entree.Statut, id)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if lignes, _ := resultat.RowsAffected(); lignes == 0 {
		envoyerErreur(w, http.StatusNotFound, "utilisateur introuvable")
		return
	}

	if entree.Statut == statutSuspendu {
		if _, err := s.db.Exec("DELETE FROM sessions WHERE utilisateur_id = $1", id); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
	}

	envoyerJSON(w, http.StatusOK, map[string]string{"statut": entree.Statut})
}

type ResponsableEntree struct {
	EstResponsable bool `json:"est_responsable"`
}

func (s *Serveur) gererChangementResponsable(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var entree ResponsableEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}

	resultat, err := s.db.Exec("UPDATE salaries SET est_responsable = $1 WHERE utilisateur_id = $2", entree.EstResponsable, id)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if lignes, _ := resultat.RowsAffected(); lignes == 0 {
		envoyerErreur(w, http.StatusNotFound, "salarié introuvable")
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]bool{"est_responsable": entree.EstResponsable})
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
	var erreurPg *pq.Error
	if errors.As(err, &erreurPg) && erreurPg.Code == "23503" {

		envoyerErreur(w, http.StatusConflict, "ce compte est lié à des données (annonces, inscriptions, offres…) et ne peut pas être supprimé")
		return
	}
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
	if motDePasseObligatoire {
		if message := validerMotDePasse(entree.MotDePasse); message != "" {
			return message
		}
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

			envoyerErreur(w, http.StatusBadRequest, "un identifiant fourni ne correspond à aucun élément existant")
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

func (s *Serveur) creerProfilSalarie(utilisateurID, roleID int) error {
	_, err := s.db.Exec(`
		INSERT INTO salaries (utilisateur_id)
		SELECT $1 FROM roles WHERE id = $2 AND code = $3
		ON CONFLICT (utilisateur_id) DO NOTHING`,
		utilisateurID, roleID, roleSalarie,
	)
	return err
}

type CompteEntree struct {
	Role             string `json:"role"`
	Nom              string `json:"nom"`
	Prenom           string `json:"prenom"`
	Email            string `json:"email"`
	MotDePasse       string `json:"mot_de_passe"`
	LanguePrefereeID *int   `json:"langue_preferee_id"`
}

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
