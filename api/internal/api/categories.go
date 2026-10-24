package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/lib/pq"
)

const (
	tableCategoriesPrestations = "categories_prestations"
	tableCategoriesMateriaux   = "categories_materiaux"
)

var accents = strings.NewReplacer(
	"à", "a", "â", "a", "ä", "a", "ç", "c", "é", "e", "è", "e", "ê", "e", "ë", "e",
	"î", "i", "ï", "i", "ô", "o", "ö", "o", "ù", "u", "û", "u", "ü", "u", "ÿ", "y", "œ", "oe", "æ", "ae",
)

var caracteresNonAutorises = regexp.MustCompile(`[^a-z0-9]+`)

type CategoriePrestation struct {
	ID      int    `json:"id"`
	Libelle string `json:"libelle"`
	Actif   bool   `json:"actif"`
}

type CategorieEntree struct {
	Libelle string `json:"libelle"`
	Actif   *bool  `json:"actif"`
}

func codeDepuisLibelle(libelle string) string {
	code := accents.Replace(strings.ToLower(libelle))
	return strings.Trim(caracteresNonAutorises.ReplaceAllString(code, "_"), "_")
}

func (s *Serveur) verifierCategorie(table string, id, idActuel int) (string, string, error) {
	var libelle string
	var actif bool
	err := s.db.QueryRow("SELECT libelle, actif FROM "+table+" WHERE id = $1", id).Scan(&libelle, &actif)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "catégorie inconnue", nil
	}
	if err != nil {
		return "", "", err
	}
	if !actif && id != idActuel {
		return "", "cette catégorie est désactivée", nil
	}
	return libelle, "", nil
}

func lireCategorie(w http.ResponseWriter, r *http.Request, actifObligatoire bool) (CategorieEntree, bool) {
	var entree CategorieEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return entree, false
	}
	entree.Libelle = strings.TrimSpace(entree.Libelle)
	if entree.Libelle == "" {
		envoyerErreur(w, http.StatusBadRequest, "nom de la catégorie obligatoire")
		return entree, false
	}
	if actifObligatoire && entree.Actif == nil {
		envoyerErreur(w, http.StatusBadRequest, "statut actif obligatoire")
		return entree, false
	}
	return entree, true
}

func estDoublon(err error) bool {
	var erreurPg *pq.Error
	return errors.As(err, &erreurPg) && erreurPg.Code == "23505"
}

func (s *Serveur) gererListeCategoriesPrestations(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	lignes, err := s.db.Query("SELECT id, libelle, actif FROM " + tableCategoriesPrestations + " ORDER BY libelle")
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	categories := make([]CategoriePrestation, 0)
	for lignes.Next() {
		var c CategoriePrestation
		if err := lignes.Scan(&c.ID, &c.Libelle, &c.Actif); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		categories = append(categories, c)
	}

	envoyerJSON(w, http.StatusOK, categories)
}

func (s *Serveur) gererCreationCategoriePrestation(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	entree, ok := lireCategorie(w, r, false)
	if !ok {
		return
	}

	c := CategoriePrestation{Libelle: entree.Libelle}
	err := s.db.QueryRow(
		"INSERT INTO "+tableCategoriesPrestations+" (libelle) VALUES ($1) RETURNING id, actif", entree.Libelle,
	).Scan(&c.ID, &c.Actif)
	if estDoublon(err) {
		envoyerErreur(w, http.StatusConflict, "cette catégorie existe déjà")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusCreated, c)
}

func (s *Serveur) gererCreationCategorieMateriau(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	entree, ok := lireCategorie(w, r, false)
	if !ok {
		return
	}
	code := codeDepuisLibelle(entree.Libelle)
	if code == "" {
		envoyerErreur(w, http.StatusBadRequest, "nom de la catégorie invalide")
		return
	}

	c := Categorie{Code: code, Libelle: entree.Libelle}
	err := s.db.QueryRow(
		"INSERT INTO "+tableCategoriesMateriaux+" (code, libelle) VALUES ($1, $2) RETURNING id, actif", code, entree.Libelle,
	).Scan(&c.ID, &c.Actif)
	if estDoublon(err) {
		envoyerErreur(w, http.StatusConflict, "cette catégorie existe déjà")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusCreated, c)
}

func (s *Serveur) modifierCategorie(w http.ResponseWriter, r *http.Request, table string) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	entree, ok := lireCategorie(w, r, true)
	if !ok {
		return
	}

	resultat, err := s.db.Exec("UPDATE "+table+" SET libelle = $1, actif = $2 WHERE id = $3", entree.Libelle, *entree.Actif, id)
	if estDoublon(err) {
		envoyerErreur(w, http.StatusConflict, "une autre catégorie porte déjà ce nom")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if lignes, _ := resultat.RowsAffected(); lignes == 0 {
		envoyerErreur(w, http.StatusNotFound, "catégorie introuvable")
		return
	}

	envoyerJSON(w, http.StatusOK, CategoriePrestation{ID: id, Libelle: entree.Libelle, Actif: *entree.Actif})
}

func (s *Serveur) gererModificationCategoriePrestation(w http.ResponseWriter, r *http.Request) {
	s.modifierCategorie(w, r, tableCategoriesPrestations)
}

func (s *Serveur) gererModificationCategorieMateriau(w http.ResponseWriter, r *http.Request) {
	s.modifierCategorie(w, r, tableCategoriesMateriaux)
}
