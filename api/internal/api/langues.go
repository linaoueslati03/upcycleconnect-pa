package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/lib/pq"
)

const langueParDefaut = "fr"

var formatCodeLangue = regexp.MustCompile(`^[a-z]{2}(-[a-z]{2})?$`)

type LangueAdmin struct {
	ID      int    `json:"id"`
	Code    string `json:"code"`
	Libelle string `json:"libelle"`
	Actif   bool   `json:"actif"`
}

type LangueEntree struct {
	Code    string `json:"code"`
	Libelle string `json:"libelle"`
	Actif   *bool  `json:"actif"`
}

func (s *Serveur) gererListeLanguesAdmin(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	lignes, err := s.db.Query("SELECT id, code, libelle, actif FROM langues ORDER BY libelle")
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	langues := make([]LangueAdmin, 0)
	for lignes.Next() {
		var l LangueAdmin
		if err := lignes.Scan(&l.ID, &l.Code, &l.Libelle, &l.Actif); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		langues = append(langues, l)
	}

	envoyerJSON(w, http.StatusOK, langues)
}

func (s *Serveur) gererCreationLangue(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	var entree LangueEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	entree.Code = strings.ToLower(strings.TrimSpace(entree.Code))
	entree.Libelle = strings.TrimSpace(entree.Libelle)
	if !formatCodeLangue.MatchString(entree.Code) {
		envoyerErreur(w, http.StatusBadRequest, "code de langue invalide (deux lettres, par exemple es ou de)")
		return
	}
	if entree.Libelle == "" {
		envoyerErreur(w, http.StatusBadRequest, "nom de la langue obligatoire")
		return
	}

	l := LangueAdmin{Code: entree.Code, Libelle: entree.Libelle}
	err := s.db.QueryRow(
		"INSERT INTO langues (code, libelle, actif) VALUES ($1, $2, FALSE) RETURNING id, actif",
		entree.Code, entree.Libelle,
	).Scan(&l.ID, &l.Actif)
	var erreurPg *pq.Error
	if errors.As(err, &erreurPg) && erreurPg.Code == "23505" {
		envoyerErreur(w, http.StatusConflict, "cette langue existe déjà")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusCreated, l)
}

func (s *Serveur) gererModificationLangue(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var entree LangueEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	entree.Libelle = strings.TrimSpace(entree.Libelle)
	if entree.Libelle == "" || entree.Actif == nil {
		envoyerErreur(w, http.StatusBadRequest, "nom et statut actif obligatoires")
		return
	}

	var l LangueAdmin
	err = s.db.QueryRow("SELECT id, code FROM langues WHERE id = $1", id).Scan(&l.ID, &l.Code)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "langue introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if l.Code == langueParDefaut && !*entree.Actif {
		envoyerErreur(w, http.StatusBadRequest, "la langue par défaut ne peut pas être désactivée")
		return
	}

	_, err = s.db.Exec("UPDATE langues SET libelle = $1, actif = $2 WHERE id = $3", entree.Libelle, *entree.Actif, id)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	l.Libelle = entree.Libelle
	l.Actif = *entree.Actif
	envoyerJSON(w, http.StatusOK, l)
}

type LanguePrefereeEntree struct {
	Code string `json:"code"`
}

func (s *Serveur) gererChangementLanguePreferee(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	var entree LanguePrefereeEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}

	resultat, err := s.db.Exec(`
		UPDATE utilisateurs SET langue_preferee_id = l.id
		FROM langues l
		WHERE utilisateurs.id = $1 AND l.code = $2 AND l.actif`,
		utilisateurID, entree.Code,
	)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if lignes, _ := resultat.RowsAffected(); lignes == 0 {
		envoyerErreur(w, http.StatusBadRequest, "langue inconnue ou désactivée")
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]string{"langue": entree.Code})
}

const tailleMaximaleTraduction = 1 << 20

func dossierTraductions() string {
	if dossier := os.Getenv("DOSSIER_TRADUCTIONS"); dossier != "" {
		return dossier
	}
	return filepath.Join("..", "front", "assets", "locales")
}

func listerCles(prefixe string, valeur any, cles map[string]bool) bool {
	switch contenu := valeur.(type) {
	case map[string]any:
		for cle, enfant := range contenu {
			if !listerCles(prefixe+cle+".", enfant, cles) {
				return false
			}
		}
		return true
	case string:
		cles[strings.TrimSuffix(prefixe, ".")] = true
		return true
	default:
		return false
	}
}

func lireClesFichier(chemin string) (map[string]bool, error) {
	contenu, err := os.ReadFile(chemin)
	if err != nil {
		return nil, err
	}
	var traductions map[string]any
	if err := json.Unmarshal(contenu, &traductions); err != nil {
		return nil, err
	}
	cles := make(map[string]bool)
	listerCles("", traductions, cles)
	return cles, nil
}

func (s *Serveur) gererImportTraduction(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var code string
	err = s.db.QueryRow("SELECT code FROM langues WHERE id = $1", id).Scan(&code)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "langue introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if code == langueParDefaut {
		envoyerErreur(w, http.StatusBadRequest, "le fichier de référence français ne peut pas être remplacé depuis l'interface")
		return
	}

	var traductions map[string]any
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, tailleMaximaleTraduction)).Decode(&traductions); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "fichier invalide : ce n'est pas un fichier JSON de traduction")
		return
	}
	clesImportees := make(map[string]bool)
	if !listerCles("", traductions, clesImportees) {
		envoyerErreur(w, http.StatusBadRequest, "fichier invalide : chaque traduction doit être un texte")
		return
	}

	clesReference, err := lireClesFichier(filepath.Join(dossierTraductions(), langueParDefaut+".json"))
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "fichier de référence introuvable")
		return
	}
	manquantes := make([]string, 0)
	for cle := range clesReference {
		if !clesImportees[cle] {
			manquantes = append(manquantes, cle)
		}
	}
	if len(manquantes) > 0 {
		slices.Sort(manquantes)
		exemples := manquantes[:min(len(manquantes), 5)]
		envoyerErreur(w, http.StatusBadRequest, fmt.Sprintf("%d texte(s) non traduit(s), par exemple : %s",
			len(manquantes), strings.Join(exemples, ", ")))
		return
	}

	contenu, err := json.MarshalIndent(traductions, "", "    ")
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	chemin := filepath.Join(dossierTraductions(), code+".json")
	temporaire := chemin + ".tmp"
	if err := os.WriteFile(temporaire, append(contenu, '\n'), 0644); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "impossible d'enregistrer le fichier de traduction")
		return
	}
	if err := os.Rename(temporaire, chemin); err != nil {
		os.Remove(temporaire)
		envoyerErreur(w, http.StatusInternalServerError, "impossible d'enregistrer le fichier de traduction")
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]any{"langue": code, "textes": len(clesImportees)})
}
