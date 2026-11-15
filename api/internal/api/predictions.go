package api

import (
	"database/sql"
	"errors"
	"net/http"
	"slices"
	"sort"
	"strconv"

	"github.com/linaoueslati03/upcycleconnect-pa/api/internal/ml"
)

const (
	predictionsParPage      = 20
	requeteVariablesActeurs = `
		SELECT u.id, u.prenom, u.nom, u.email,
		       v.role, v.premium, v.anciennete_jours, v.upcycling_score,
		       v.nb_annonces, v.nb_annonces_vente, v.nb_formations, v.nb_ateliers, v.nb_evenements,
		       v.nb_depots, v.nb_achats_annonces, v.montant_depense, v.jours_depuis_derniere_demande,
		       v.derniere_demande, COALESCE(cm.code, 'aucune'), v.site_principal
		FROM vue_variables_acteurs v
		JOIN utilisateurs u ON u.id = v.utilisateur_id
		LEFT JOIN categories_materiaux cm ON cm.id = v.categorie_preferee`
)

type ActeurPredit struct {
	ID     int    `json:"id"`
	Prenom string `json:"prenom"`
	Nom    string `json:"nom"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

type InfosModele struct {
	Version        string  `json:"version"`
	ExactitudeTest float64 `json:"exactitude_test"`
	F1MacroTest    float64 `json:"f1_macro_test"`
}

type DetailPrediction struct {
	Utilisateur ActeurPredit  `json:"utilisateur"`
	Variables   ml.Variables  `json:"variables"`
	Prediction  ml.Prediction `json:"prediction"`
	Modele      InfosModele   `json:"modele"`
}

type LignePrediction struct {
	ActeurPredit
	Classe                     string  `json:"classe"`
	Probabilite                float64 `json:"probabilite"`
	JoursDepuisDerniereDemande int     `json:"jours_depuis_derniere_demande"`
}

type PagePredictions struct {
	Modele      InfosModele       `json:"modele"`
	Repartition []Comptage        `json:"repartition"`
	Predictions []LignePrediction `json:"predictions"`
	Total       int               `json:"total"`
	Page        int               `json:"page"`
	ParPage     int               `json:"par_page"`
}

type scanneur interface {
	Scan(dest ...any) error
}

func lireActeur(ligne scanneur) (ActeurPredit, ml.Variables, error) {
	var a ActeurPredit
	var v ml.Variables
	err := ligne.Scan(&a.ID, &a.Prenom, &a.Nom, &a.Email,
		&v.Role, &v.Premium, &v.AncienneteJours, &v.UpcyclingScore,
		&v.NbAnnonces, &v.NbAnnoncesVente, &v.NbFormations, &v.NbAteliers, &v.NbEvenements,
		&v.NbDepots, &v.NbAchatsAnnonces, &v.MontantDepense, &v.JoursDepuisDerniereDemande,
		&v.DerniereDemande, &v.CategoriePreferee, &v.SitePrincipal)
	a.Role = v.Role
	return a, v, err
}

func (s *Serveur) infosModele() InfosModele {
	return InfosModele{Version: s.modele.Version, ExactitudeTest: s.modele.ExactitudeTest, F1MacroTest: s.modele.F1MacroTest}
}

func (s *Serveur) gererPredictionActeur(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	acteur, variables, err := lireActeur(s.db.QueryRow(requeteVariablesActeurs+" WHERE v.utilisateur_id = $1", id))
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "acteur introuvable (seuls les particuliers et les professionnels ont une prédiction)")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, DetailPrediction{
		Utilisateur: acteur,
		Variables:   variables,
		Prediction:  s.modele.Predire(variables),
		Modele:      s.infosModele(),
	})
}

func (s *Serveur) gererListePredictions(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	filtreClasse := r.URL.Query().Get("classe")
	if filtreClasse != "" && !slices.Contains(s.modele.Classes, filtreClasse) {
		envoyerErreur(w, http.StatusBadRequest, "classe invalide")
		return
	}
	filtreRole := r.URL.Query().Get("role")
	if filtreRole != "" && !slices.Contains(rolesActeurs, filtreRole) {
		envoyerErreur(w, http.StatusBadRequest, "rôle invalide (particulier ou professionnel attendu)")
		return
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	lignes, err := s.db.Query(requeteVariablesActeurs)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	parClasse := make(map[string]int)
	retenues := make([]LignePrediction, 0)
	for lignes.Next() {
		acteur, variables, err := lireActeur(lignes)
		if err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		if filtreRole != "" && acteur.Role != filtreRole {
			continue
		}
		prediction := s.modele.Predire(variables)
		parClasse[prediction.Classe]++
		if filtreClasse != "" && prediction.Classe != filtreClasse {
			continue
		}
		ligne := LignePrediction{
			ActeurPredit:               acteur,
			Classe:                     prediction.Classe,
			JoursDepuisDerniereDemande: variables.JoursDepuisDerniereDemande,
		}
		for _, p := range prediction.Probabilites {
			if p.Classe == prediction.Classe {
				ligne.Probabilite = p.Probabilite
			}
		}
		retenues = append(retenues, ligne)
	}
	if err := lignes.Err(); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	sort.SliceStable(retenues, func(a, b int) bool {
		if retenues[a].Probabilite != retenues[b].Probabilite {
			return retenues[a].Probabilite > retenues[b].Probabilite
		}
		return retenues[a].ID < retenues[b].ID
	})

	resultat := PagePredictions{
		Modele:      s.infosModele(),
		Repartition: make([]Comptage, 0, len(s.modele.Classes)),
		Total:       len(retenues),
		Page:        page,
		ParPage:     predictionsParPage,
		Predictions: make([]LignePrediction, 0, predictionsParPage),
	}
	for _, classe := range s.modele.Classes {
		resultat.Repartition = append(resultat.Repartition, Comptage{Libelle: classe, Nombre: parClasse[classe]})
	}
	debut := (page - 1) * predictionsParPage
	if debut < len(retenues) {
		resultat.Predictions = append(resultat.Predictions, retenues[debut:min(debut+predictionsParPage, len(retenues))]...)
	}

	envoyerJSON(w, http.StatusOK, resultat)
}
