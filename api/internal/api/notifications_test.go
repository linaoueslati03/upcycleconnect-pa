package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

type requeteRecue struct {
	Autorisation string
	Corps        struct {
		AppID    string              `json:"app_id"`
		Alias    map[string][]string `json:"include_aliases"`
		Canal    string              `json:"target_channel"`
		Titres   map[string]string   `json:"headings"`
		Contenus map[string]string   `json:"contents"`
	}
}

func fauxOneSignal(t *testing.T, statut int, reponse func(numero int) string) (*[]requeteRecue, func()) {
	t.Helper()
	var verrou sync.Mutex
	recues := make([]requeteRecue, 0)
	serveur := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var recue requeteRecue
		recue.Autorisation = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&recue.Corps); err != nil {
			t.Errorf("corps illisible : %v", err)
		}
		verrou.Lock()
		recues = append(recues, recue)
		numero := len(recues)
		verrou.Unlock()
		w.WriteHeader(statut)
		fmt.Fprint(w, reponse(numero))
	}))
	ancienne := urlOneSignal
	urlOneSignal = serveur.URL
	t.Setenv("ONESIGNAL_APP_ID", "app-test")
	t.Setenv("ONESIGNAL_REST_API_KEY", "cle-test")
	return &recues, func() {
		urlOneSignal = ancienne
		serveur.Close()
	}
}

var contenuTest = ContenuNotification{TitreFR: "Bonjour", MessageFR: "Message", TitreEN: "Hello", MessageEN: "Message EN"}

func TestEnvoiContenuEtDestinataires(t *testing.T) {
	recues, fermer := fauxOneSignal(t, http.StatusOK, func(int) string { return `{"id":"notif-1"}` })
	defer fermer()

	identifiant, err := envoyerOneSignal([]int{4, 12}, contenuTest)
	if err != nil {
		t.Fatalf("envoi : %v", err)
	}
	if identifiant != "notif-1" {
		t.Errorf("identifiant %q, notif-1 attendu", identifiant)
	}
	if len(*recues) != 1 {
		t.Fatalf("%d requêtes envoyées, 1 attendue", len(*recues))
	}
	r := (*recues)[0]
	if r.Autorisation != "Key cle-test" || r.Corps.AppID != "app-test" || r.Corps.Canal != "push" {
		t.Errorf("en-tête ou paramètres incorrects : %+v", r)
	}
	if got := strings.Join(r.Corps.Alias["external_id"], ","); got != "4,12" {
		t.Errorf("destinataires %q, 4,12 attendu", got)
	}
	if r.Corps.Titres["fr"] != "Bonjour" || r.Corps.Titres["en"] != "Hello" ||
		r.Corps.Contenus["fr"] != "Message" || r.Corps.Contenus["en"] != "Message EN" {
		t.Errorf("titres ou contenus fr / en incorrects : %+v %+v", r.Corps.Titres, r.Corps.Contenus)
	}
}

func TestEnvoiDecoupeParPaquets(t *testing.T) {
	recues, fermer := fauxOneSignal(t, http.StatusOK, func(n int) string { return fmt.Sprintf(`{"id":"notif-%d"}`, n) })
	defer fermer()

	destinataires := make([]int, 2*destinatairesParEnvoi+5)
	for i := range destinataires {
		destinataires[i] = i + 1
	}
	identifiant, err := envoyerOneSignal(destinataires, contenuTest)
	if err != nil {
		t.Fatalf("envoi : %v", err)
	}
	if len(*recues) != 3 {
		t.Fatalf("%d requêtes, 3 attendues", len(*recues))
	}
	total := 0
	for _, r := range *recues {
		total += len(r.Corps.Alias["external_id"])
	}
	if total != len(destinataires) || identifiant != "notif-1,notif-2,notif-3" {
		t.Errorf("%d destinataires envoyés sur %d, identifiants %q", total, len(destinataires), identifiant)
	}
}

func TestEnvoiErreurs(t *testing.T) {
	cas := map[string]struct {
		statut  int
		reponse string
		attendu string
	}{
		"refus HTTP":        {http.StatusBadRequest, `{"errors":["bad request"]}`, "HTTP 400"},
		"aucun abonné":      {http.StatusOK, `{"id":"","errors":["All included players are not subscribed"]}`, "aucun destinataire abonné"},
		"réponse illisible": {http.StatusOK, `pas du json`, "illisible"},
	}
	for nom, c := range cas {
		t.Run(nom, func(t *testing.T) {
			_, fermer := fauxOneSignal(t, c.statut, func(int) string { return c.reponse })
			defer fermer()
			_, err := envoyerOneSignal([]int{1}, contenuTest)
			if err == nil || !strings.Contains(err.Error(), c.attendu) {
				t.Errorf("erreur %v, %q attendu", err, c.attendu)
			}
		})
	}
}

func TestEnvoiSansConfiguration(t *testing.T) {
	t.Setenv("ONESIGNAL_APP_ID", "")
	t.Setenv("ONESIGNAL_REST_API_KEY", "")
	if _, err := envoyerOneSignal([]int{1}, contenuTest); err == nil || !strings.Contains(err.Error(), "configuration") {
		t.Errorf("erreur %v, configuration manquante attendue", err)
	}
}

func TestValidationDemande(t *testing.T) {
	classes := []string{"formation", "depot"}
	id := 5
	cas := []struct {
		nom     string
		demande DemandeNotification
		valide  bool
	}{
		{"tous", DemandeNotification{Cible: "tous", TitreFR: "T", MessageFR: "M"}, true},
		{"cible inconnue", DemandeNotification{Cible: "salaries", TitreFR: "T", MessageFR: "M"}, false},
		{"sans titre", DemandeNotification{Cible: "tous", MessageFR: "M"}, false},
		{"utilisateur sans id", DemandeNotification{Cible: "utilisateur", TitreFR: "T", MessageFR: "M"}, false},
		{"utilisateur", DemandeNotification{Cible: "utilisateur", DestinataireID: &id, TitreFR: "T", MessageFR: "M"}, true},
		{"prédiction inconnue", DemandeNotification{Cible: "prediction", ClassePredite: "cours", TitreFR: "T", MessageFR: "M"}, false},
		{"prédiction", DemandeNotification{Cible: "prediction", ClassePredite: "depot", TitreFR: "T", MessageFR: "M"}, true},
		{"titre trop long", DemandeNotification{Cible: "tous", TitreFR: strings.Repeat("é", longueurMaxTitre+1), MessageFR: "M"}, false},
	}
	for _, c := range cas {
		d := c.demande
		if valide := d.valider(classes) == ""; valide != c.valide {
			t.Errorf("%s : valide=%v, %v attendu", c.nom, valide, c.valide)
		}
	}

	d := DemandeNotification{Cible: "tous", TitreFR: " Titre ", MessageFR: "Message"}
	d.valider(classes)
	if d.TitreEN != "Titre" || d.MessageEN != "Message" {
		t.Errorf("l'anglais vide doit reprendre le français : %+v", d)
	}
}

func TestNotificationAutomatiqueBilingue(t *testing.T) {
	recues, fermer := fauxOneSignal(t, http.StatusOK, func(int) string { return `{"id":"notif-auto"}` })
	defer fermer()

	if err := envoyerNotification(7, "Dépôt validé", "Ta demande de dépôt a été validée."); err != nil {
		t.Fatalf("envoi : %v", err)
	}
	r := (*recues)[0]
	if r.Corps.Titres["en"] != "Dépôt validé" || r.Corps.Contenus["en"] != "Ta demande de dépôt a été validée." {
		t.Errorf("le texte anglais doit reprendre le français : %+v %+v", r.Corps.Titres, r.Corps.Contenus)
	}
	if got := strings.Join(r.Corps.Alias["external_id"], ","); got != "7" {
		t.Errorf("destinataire %q, 7 attendu", got)
	}
}
