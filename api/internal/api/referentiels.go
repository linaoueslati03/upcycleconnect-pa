package api

import "net/http"

type Categorie struct {
	ID      int    `json:"id"`
	Code    string `json:"code"`
	Libelle string `json:"libelle"`
	Actif   bool   `json:"actif"`
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
	lignes, err := s.db.Query("SELECT id, code, libelle, actif FROM categories_materiaux ORDER BY libelle")
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	categories := make([]Categorie, 0)
	for lignes.Next() {
		var c Categorie
		if err := lignes.Scan(&c.ID, &c.Code, &c.Libelle, &c.Actif); err != nil {
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
