package api

import "net/http"

// Listes de référence lues en base plutôt que codées en dur dans le front : ajouter une
// catégorie, un conteneur ou une langue se fait par une simple insertion SQL.

type Categorie struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
}

type Conteneur struct {
	ID     int    `json:"id"`
	Site   string `json:"site"`
	Statut string `json:"statut"`
}

type Langue struct {
	ID      int    `json:"id"`
	Code    string `json:"code"`
	Libelle string `json:"libelle"`
}

func (s *Serveur) gererListeCategories(w http.ResponseWriter, r *http.Request) {
	lignes, err := s.db.Query("SELECT id, code FROM categories_materiaux ORDER BY code")
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	categories := make([]Categorie, 0)
	for lignes.Next() {
		var c Categorie
		if err := lignes.Scan(&c.ID, &c.Code); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		categories = append(categories, c)
	}

	envoyerJSON(w, http.StatusOK, categories)
}

func (s *Serveur) gererListeConteneurs(w http.ResponseWriter, r *http.Request) {
	lignes, err := s.db.Query("SELECT id, site, statut FROM conteneurs ORDER BY site")
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	conteneurs := make([]Conteneur, 0)
	for lignes.Next() {
		var c Conteneur
		if err := lignes.Scan(&c.ID, &c.Site, &c.Statut); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		conteneurs = append(conteneurs, c)
	}

	envoyerJSON(w, http.StatusOK, conteneurs)
}

func (s *Serveur) gererListeLangues(w http.ResponseWriter, r *http.Request) {
	lignes, err := s.db.Query("SELECT id, code, libelle FROM langues WHERE actif ORDER BY libelle")
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	langues := make([]Langue, 0)
	for lignes.Next() {
		var l Langue
		if err := lignes.Scan(&l.ID, &l.Code, &l.Libelle); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		langues = append(langues, l)
	}

	envoyerJSON(w, http.StatusOK, langues)
}
