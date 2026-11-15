package ml

import (
	"encoding/json"
	"math"
	"os"
	"strings"
	"testing"
)

type casDeReference struct {
	Variables    Variables `json:"variables"`
	Classe       string    `json:"classe"`
	Probabilites []float64 `json:"probabilites"`
}

func chargerModele(t *testing.T) *Modele {
	t.Helper()
	m, err := Charger()
	if err != nil {
		t.Fatalf("chargement du modèle : %v", err)
	}
	return m
}

func TestModeleValide(t *testing.T) {
	m := chargerModele(t)
	if len(m.Classes) != 6 {
		t.Errorf("%d classes, 6 attendues", len(m.Classes))
	}
	if len(m.Colonnes) != len(m.Encoder(Variables{})) {
		t.Errorf("l'encodage produit %d valeurs pour %d colonnes", len(m.Encoder(Variables{})), len(m.Colonnes))
	}
	if m.Version == "" {
		t.Error("version du modèle absente")
	}
}

func TestCasDeReference(t *testing.T) {
	m := chargerModele(t)
	donnees, err := os.ReadFile("cas_test.json")
	if err != nil {
		t.Fatalf("lecture de cas_test.json : %v", err)
	}
	var cas []casDeReference
	if err := json.Unmarshal(donnees, &cas); err != nil {
		t.Fatalf("cas_test.json illisible : %v", err)
	}
	if len(cas) == 0 {
		t.Fatal("aucun cas de référence")
	}

	for i, c := range cas {
		prediction := m.Predire(c.Variables)
		if prediction.Classe != c.Classe {
			t.Errorf("cas %d : classe %s, %s attendue", i, prediction.Classe, c.Classe)
		}
		obtenues := make(map[string]float64)
		for _, p := range prediction.Probabilites {
			obtenues[p.Classe] = p.Probabilite
		}
		for j, classe := range m.Classes {
			if math.Abs(obtenues[classe]-c.Probabilites[j]) > 1e-12 {
				t.Errorf("cas %d, classe %s : probabilité %v, %v attendue", i, classe, obtenues[classe], c.Probabilites[j])
			}
		}
	}
}

func TestProbabilitesTriees(t *testing.T) {
	m := chargerModele(t)
	prediction := m.Predire(Variables{Role: "professionnel", NbAchatsAnnonces: 6, DerniereDemande: "achat_annonce"})
	for i := 1; i < len(prediction.Probabilites); i++ {
		if prediction.Probabilites[i].Probabilite > prediction.Probabilites[i-1].Probabilite {
			t.Fatalf("probabilités non triées : %+v", prediction.Probabilites)
		}
	}
	if len(prediction.Probabilites) != len(m.Classes) {
		t.Errorf("%d probabilités pour %d classes", len(prediction.Probabilites), len(m.Classes))
	}
}

func TestValeursInconnues(t *testing.T) {
	m := chargerModele(t)
	prediction := m.Predire(Variables{
		Role: "particulier", SitePrincipal: "Lyon", CategoriePreferee: "nouvelle_categorie", DerniereDemande: "aucune",
		JoursDepuisDerniereDemande: 999,
	})
	if prediction.Classe == "" {
		t.Error("aucune classe prédite pour un acteur aux valeurs inconnues")
	}
}

func TestJamaisPremiumPourUnParticulier(t *testing.T) {
	m := chargerModele(t)
	for _, v := range []Variables{
		{Role: "particulier", NbAchatsAnnonces: 10, DerniereDemande: "achat_annonce"},
		{Role: "particulier", Premium: 1, NbAchatsAnnonces: 3},
		{Role: "particulier"},
	} {
		if classe := m.Predire(v).Classe; classe == classePremium {
			t.Errorf("premium prédit pour un particulier : %+v", v)
		}
	}
}

func TestChemin(t *testing.T) {
	m := chargerModele(t)
	prediction := m.Predire(Variables{Role: "professionnel", NbAchatsAnnonces: 2, NbFormations: 3})
	if len(prediction.Chemin) == 0 || len(prediction.Chemin) > m.Profondeur {
		t.Fatalf("chemin de %d étapes pour une profondeur de %d", len(prediction.Chemin), m.Profondeur)
	}
	premiere := prediction.Chemin[0]
	if premiere.Colonne != m.Colonnes[m.Noeuds[0].Variable] {
		t.Errorf("la première étape porte sur %s au lieu de la question de la racine", premiere.Colonne)
	}
	if premiere.InferieurOuEgal != (premiere.Valeur <= premiere.Seuil) {
		t.Errorf("direction incohérente à la première étape : %+v", premiere)
	}
}

func TestModeleInvalideRefuse(t *testing.T) {
	m := chargerModele(t)
	original, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	casInvalides := map[string]string{
		"enfant inexistant":   `"gauche":1,`,
		"colonne inexistante": `"variable":`,
	}
	for nom, motif := range casInvalides {
		texte := string(original)
		switch nom {
		case "enfant inexistant":
			texte = strings.Replace(texte, motif, `"gauche":9999,`, 1)
		case "colonne inexistante":
			texte = strings.Replace(texte, motif, `"variable":9999,"ancienne":`, 1)
		}
		if _, err := ChargerDepuis([]byte(texte)); err == nil {
			t.Errorf("%s : un modèle invalide a été accepté", nom)
		}
	}
	if _, err := ChargerDepuis([]byte("{")); err == nil {
		t.Error("un JSON illisible a été accepté")
	}
}
