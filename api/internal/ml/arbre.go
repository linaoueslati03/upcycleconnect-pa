package ml

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"strings"
)

//go:embed arbre.json
var arbreJSON []byte

const (
	classePremium   = "premium"
	roleParticulier = "particulier"
)

type Variables struct {
	Role                       string  `json:"role"`
	Premium                    int     `json:"premium"`
	AncienneteJours            int     `json:"anciennete_jours"`
	UpcyclingScore             int     `json:"upcycling_score"`
	NbAnnonces                 int     `json:"nb_annonces"`
	NbAnnoncesVente            int     `json:"nb_annonces_vente"`
	NbFormations               int     `json:"nb_formations"`
	NbAteliers                 int     `json:"nb_ateliers"`
	NbEvenements               int     `json:"nb_evenements"`
	NbDepots                   int     `json:"nb_depots"`
	NbAchatsAnnonces           int     `json:"nb_achats_annonces"`
	MontantDepense             float64 `json:"montant_depense"`
	JoursDepuisDerniereDemande int     `json:"jours_depuis_derniere_demande"`
	DerniereDemande            string  `json:"derniere_demande"`
	CategoriePreferee          string  `json:"categorie_preferee"`
	SitePrincipal              string  `json:"site_principal"`
}

func (v Variables) numerique(nom string) (float64, error) {
	switch nom {
	case "premium":
		return float64(v.Premium), nil
	case "anciennete_jours":
		return float64(v.AncienneteJours), nil
	case "upcycling_score":
		return float64(v.UpcyclingScore), nil
	case "nb_annonces":
		return float64(v.NbAnnonces), nil
	case "nb_annonces_vente":
		return float64(v.NbAnnoncesVente), nil
	case "nb_formations":
		return float64(v.NbFormations), nil
	case "nb_ateliers":
		return float64(v.NbAteliers), nil
	case "nb_evenements":
		return float64(v.NbEvenements), nil
	case "nb_depots":
		return float64(v.NbDepots), nil
	case "nb_achats_annonces":
		return float64(v.NbAchatsAnnonces), nil
	case "montant_depense":
		return v.MontantDepense, nil
	case "jours_depuis_derniere_demande":
		return float64(v.JoursDepuisDerniereDemande), nil
	}
	return 0, fmt.Errorf("variable numérique inconnue : %s", nom)
}

func (v Variables) texte(nom string) (string, error) {
	switch nom {
	case "role":
		return v.Role, nil
	case "derniere_demande":
		return v.DerniereDemande, nil
	case "categorie_preferee":
		return v.CategoriePreferee, nil
	case "site_principal":
		return v.SitePrincipal, nil
	}
	return "", fmt.Errorf("variable texte inconnue : %s", nom)
}

type Encodage struct {
	Numeriques    []string            `json:"numeriques"`
	Parts         []string            `json:"parts"`
	Categorielles map[string][]string `json:"categorielles"`
}

type Noeud struct {
	Variable     int       `json:"variable"`
	Seuil        float64   `json:"seuil"`
	Gauche       int       `json:"gauche"`
	Droite       int       `json:"droite"`
	Probabilites []float64 `json:"probabilites"`
	Effectif     int       `json:"effectif"`
}

func (n Noeud) estFeuille() bool {
	return n.Probabilites != nil
}

type genreColonne int

const (
	colonneNumerique genreColonne = iota
	colonnePart
	colonneIndicateur
)

type colonne struct {
	genre    genreColonne
	variable string
	valeur   string
}

type Modele struct {
	Version        string   `json:"version"`
	ExactitudeTest float64  `json:"exactitude_test"`
	F1MacroTest    float64  `json:"f1_macro_test"`
	Profondeur     int      `json:"profondeur"`
	Classes        []string `json:"classes"`
	Encodage       Encodage `json:"encodage"`
	Colonnes       []string `json:"colonnes"`
	Noeuds         []Noeud  `json:"noeuds"`
	plan           []colonne
}

type ProbabiliteClasse struct {
	Classe      string  `json:"classe"`
	Probabilite float64 `json:"probabilite"`
}

type Etape struct {
	Colonne         string  `json:"colonne"`
	Seuil           float64 `json:"seuil"`
	Valeur          float64 `json:"valeur"`
	InferieurOuEgal bool    `json:"inferieur_ou_egal"`
}

type Prediction struct {
	Classe       string              `json:"classe"`
	Probabilites []ProbabiliteClasse `json:"probabilites"`
	Effectif     int                 `json:"effectif"`
	Chemin       []Etape             `json:"chemin"`
}

func Charger() (*Modele, error) {
	return ChargerDepuis(arbreJSON)
}

func ChargerDepuis(donnees []byte) (*Modele, error) {
	var m Modele
	if err := json.Unmarshal(donnees, &m); err != nil {
		return nil, fmt.Errorf("modèle illisible : %w", err)
	}
	if err := m.preparer(); err != nil {
		return nil, fmt.Errorf("modèle invalide : %w", err)
	}
	return &m, nil
}

func (m *Modele) preparer() error {
	if len(m.Classes) == 0 || len(m.Noeuds) == 0 {
		return fmt.Errorf("classes ou nœuds absents")
	}

	var vide Variables
	attendues := len(m.Encodage.Numeriques) + len(m.Encodage.Parts)
	for _, valeurs := range m.Encodage.Categorielles {
		attendues += len(valeurs)
	}
	if attendues != len(m.Colonnes) {
		return fmt.Errorf("%d colonnes annoncées par l'encodage, %d dans la liste", attendues, len(m.Colonnes))
	}

	m.plan = make([]colonne, len(m.Colonnes))
	for i, nom := range m.Colonnes {
		switch {
		case strings.Contains(nom, "="):
			variable, valeur, _ := strings.Cut(nom, "=")
			if _, err := vide.texte(variable); err != nil {
				return err
			}
			m.plan[i] = colonne{genre: colonneIndicateur, variable: variable, valeur: valeur}
		case strings.HasPrefix(nom, "part_"):
			compteur := "nb_" + strings.TrimPrefix(nom, "part_")
			if !slices.Contains(m.Encodage.Parts, compteur) {
				return fmt.Errorf("colonne %s : compteur %s absent de l'encodage", nom, compteur)
			}
			m.plan[i] = colonne{genre: colonnePart, variable: compteur}
		default:
			if _, err := vide.numerique(nom); err != nil {
				return err
			}
			m.plan[i] = colonne{genre: colonneNumerique, variable: nom}
		}
	}
	for _, compteur := range m.Encodage.Parts {
		if _, err := vide.numerique(compteur); err != nil {
			return err
		}
	}

	for i, n := range m.Noeuds {
		if n.estFeuille() {
			if len(n.Probabilites) != len(m.Classes) {
				return fmt.Errorf("feuille %d : %d probabilités pour %d classes", i, len(n.Probabilites), len(m.Classes))
			}
			continue
		}
		if n.Variable < 0 || n.Variable >= len(m.Colonnes) {
			return fmt.Errorf("nœud %d : colonne %d inexistante", i, n.Variable)
		}
		if n.Gauche <= i || n.Gauche >= len(m.Noeuds) || n.Droite <= i || n.Droite >= len(m.Noeuds) {
			return fmt.Errorf("nœud %d : enfants %d et %d invalides", i, n.Gauche, n.Droite)
		}
	}
	return nil
}

func (m *Modele) Encoder(v Variables) []float64 {
	total := 0.0
	for _, compteur := range m.Encodage.Parts {
		valeur, _ := v.numerique(compteur)
		total += valeur
	}

	vecteur := make([]float64, len(m.plan))
	for i, c := range m.plan {
		switch c.genre {
		case colonneNumerique:
			vecteur[i], _ = v.numerique(c.variable)
		case colonnePart:
			if total > 0 {
				valeur, _ := v.numerique(c.variable)
				vecteur[i] = valeur / total
			}
		case colonneIndicateur:
			if texte, _ := v.texte(c.variable); texte == c.valeur {
				vecteur[i] = 1
			}
		}
	}
	return vecteur
}

func (m *Modele) Predire(v Variables) Prediction {
	vecteur := m.Encoder(v)
	chemin := make([]Etape, 0, m.Profondeur)
	noeud := m.Noeuds[0]
	for !noeud.estFeuille() {
		valeur := vecteur[noeud.Variable]
		aGauche := valeur <= noeud.Seuil
		chemin = append(chemin, Etape{
			Colonne:         m.Colonnes[noeud.Variable],
			Seuil:           noeud.Seuil,
			Valeur:          valeur,
			InferieurOuEgal: aGauche,
		})
		if aGauche {
			noeud = m.Noeuds[noeud.Gauche]
		} else {
			noeud = m.Noeuds[noeud.Droite]
		}
	}

	probabilites := make([]ProbabiliteClasse, len(m.Classes))
	meilleure := -1
	for i, classe := range m.Classes {
		probabilites[i] = ProbabiliteClasse{Classe: classe, Probabilite: noeud.Probabilites[i]}
		interdite := classe == classePremium && v.Role == roleParticulier
		if !interdite && (meilleure < 0 || noeud.Probabilites[i] > noeud.Probabilites[meilleure]) {
			meilleure = i
		}
	}
	sort.SliceStable(probabilites, func(a, b int) bool {
		return probabilites[a].Probabilite > probabilites[b].Probabilite
	})

	return Prediction{
		Classe:       m.Classes[meilleure],
		Probabilites: probabilites,
		Effectif:     noeud.Effectif,
		Chemin:       chemin,
	}
}
