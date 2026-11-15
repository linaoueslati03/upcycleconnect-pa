package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	notificationsParPage    = 20
	destinatairesParEnvoi   = 2000
	longueurMaxTitre        = 100
	longueurMaxMessage      = 500
	statutNotificationOK    = "envoyee"
	statutNotificationEchec = "echec"
)

var (
	urlOneSignal       = "https://api.onesignal.com/notifications"
	ciblesNotification = []string{"tous", "particuliers", "professionnels", "utilisateur", "prediction"}
)

type ContenuNotification struct {
	TitreFR   string
	MessageFR string
	TitreEN   string
	MessageEN string
}

type DemandeNotification struct {
	Cible          string `json:"cible"`
	DestinataireID *int   `json:"destinataire_id"`
	ClassePredite  string `json:"classe_predite"`
	TitreFR        string `json:"titre_fr"`
	MessageFR      string `json:"message_fr"`
	TitreEN        string `json:"titre_en"`
	MessageEN      string `json:"message_en"`
}

type Notification struct {
	ID              int       `json:"id"`
	TitreFR         string    `json:"titre_fr"`
	MessageFR       string    `json:"message_fr"`
	TitreEN         string    `json:"titre_en"`
	MessageEN       string    `json:"message_en"`
	Cible           string    `json:"cible"`
	DestinataireID  *int      `json:"destinataire_id"`
	Destinataire    *string   `json:"destinataire"`
	ClassePredite   *string   `json:"classe_predite"`
	NbDestinataires int       `json:"nb_destinataires"`
	Statut          string    `json:"statut"`
	Erreur          *string   `json:"erreur"`
	Auteur          *string   `json:"auteur"`
	DateEnvoi       time.Time `json:"date_envoi"`
}

type PageNotifications struct {
	Notifications []Notification `json:"notifications"`
	Total         int            `json:"total"`
	Page          int            `json:"page"`
	ParPage       int            `json:"par_page"`
}

type reponseOneSignal struct {
	ID     string `json:"id"`
	Errors any    `json:"errors"`
}

func envoyerOneSignal(destinataires []int, contenu ContenuNotification) (string, error) {
	appID := os.Getenv("ONESIGNAL_APP_ID")
	cleREST := os.Getenv("ONESIGNAL_REST_API_KEY")
	if appID == "" || cleREST == "" {
		return "", fmt.Errorf("configuration OneSignal manquante")
	}

	identifiants := make([]string, 0)
	for debut := 0; debut < len(destinataires); debut += destinatairesParEnvoi {
		paquet := destinataires[debut:min(debut+destinatairesParEnvoi, len(destinataires))]
		externes := make([]string, len(paquet))
		for i, id := range paquet {
			externes[i] = strconv.Itoa(id)
		}
		corps, err := json.Marshal(map[string]any{
			"app_id":          appID,
			"include_aliases": map[string][]string{"external_id": externes},
			"target_channel":  "push",
			"headings":        map[string]string{"fr": contenu.TitreFR, "en": contenu.TitreEN},
			"contents":        map[string]string{"fr": contenu.MessageFR, "en": contenu.MessageEN},
		})
		if err != nil {
			return "", fmt.Errorf("préparation de la notification : %w", err)
		}

		requete, err := http.NewRequest(http.MethodPost, urlOneSignal, bytes.NewBuffer(corps))
		if err != nil {
			return "", fmt.Errorf("création de la requête : %w", err)
		}
		requete.Header.Set("Content-Type", "application/json")
		requete.Header.Set("Authorization", "Key "+cleREST)

		reponse, err := (&http.Client{Timeout: 10 * time.Second}).Do(requete)
		if err != nil {
			return "", fmt.Errorf("OneSignal injoignable : %w", err)
		}
		details, _ := io.ReadAll(reponse.Body)
		reponse.Body.Close()
		if reponse.StatusCode < 200 || reponse.StatusCode >= 300 {
			return "", fmt.Errorf("OneSignal HTTP %d : %s", reponse.StatusCode, strings.TrimSpace(string(details)))
		}

		var resultat reponseOneSignal
		if err := json.Unmarshal(details, &resultat); err != nil {
			return "", fmt.Errorf("réponse OneSignal illisible : %w", err)
		}
		if resultat.ID == "" {
			return "", fmt.Errorf("aucun destinataire abonné aux notifications (%v)", resultat.Errors)
		}
		identifiants = append(identifiants, resultat.ID)
	}
	return strings.Join(identifiants, ","), nil
}

func (d *DemandeNotification) validerCible(classes []string) string {
	if !slices.Contains(ciblesNotification, d.Cible) {
		return "cible invalide"
	}
	if d.Cible == "utilisateur" && (d.DestinataireID == nil || *d.DestinataireID <= 0) {
		return "choisissez l'utilisateur destinataire"
	}
	if d.Cible == "prediction" && !slices.Contains(classes, d.ClassePredite) {
		return "choisissez la prochaine demande probable ciblée"
	}
	return ""
}

func (d *DemandeNotification) valider(classes []string) string {
	if message := d.validerCible(classes); message != "" {
		return message
	}
	d.TitreFR = strings.TrimSpace(d.TitreFR)
	d.MessageFR = strings.TrimSpace(d.MessageFR)
	d.TitreEN = strings.TrimSpace(d.TitreEN)
	d.MessageEN = strings.TrimSpace(d.MessageEN)
	if d.TitreEN == "" {
		d.TitreEN = d.TitreFR
	}
	if d.MessageEN == "" {
		d.MessageEN = d.MessageFR
	}

	if d.TitreFR == "" || d.MessageFR == "" {
		return "le titre et le message en français sont obligatoires"
	}
	if utf8.RuneCountInString(d.TitreFR) > longueurMaxTitre || utf8.RuneCountInString(d.TitreEN) > longueurMaxTitre {
		return fmt.Sprintf("le titre ne doit pas dépasser %d caractères", longueurMaxTitre)
	}
	if utf8.RuneCountInString(d.MessageFR) > longueurMaxMessage || utf8.RuneCountInString(d.MessageEN) > longueurMaxMessage {
		return fmt.Sprintf("le message ne doit pas dépasser %d caractères", longueurMaxMessage)
	}
	return ""
}

func (s *Serveur) destinatairesNotification(d DemandeNotification) ([]int, error) {
	if d.Cible == "prediction" {
		lignes, err := s.db.Query(requeteVariablesActeurs + " WHERE u.statut = 'actif'")
		if err != nil {
			return nil, err
		}
		defer lignes.Close()
		ids := make([]int, 0)
		for lignes.Next() {
			acteur, variables, err := lireActeur(lignes)
			if err != nil {
				return nil, err
			}
			if s.modele.Predire(variables).Classe == d.ClassePredite {
				ids = append(ids, acteur.ID)
			}
		}
		return ids, lignes.Err()
	}

	requete := `
		SELECT u.id
		FROM utilisateurs u
		JOIN roles r ON r.id = u.role_id
		WHERE u.statut = 'actif' AND r.code IN ('particulier', 'professionnel')`
	args := []any{}
	switch d.Cible {
	case "particuliers":
		requete += " AND r.code = 'particulier'"
	case "professionnels":
		requete += " AND r.code = 'professionnel'"
	case "utilisateur":
		requete += " AND u.id = $1"
		args = append(args, *d.DestinataireID)
	}

	lignes, err := s.db.Query(requete+" ORDER BY u.id", args...)
	if err != nil {
		return nil, err
	}
	defer lignes.Close()
	ids := make([]int, 0)
	for lignes.Next() {
		var id int
		if err := lignes.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, lignes.Err()
}

func demandeDepuisRequete(r *http.Request) DemandeNotification {
	d := DemandeNotification{Cible: r.URL.Query().Get("cible"), ClassePredite: r.URL.Query().Get("classe_predite")}
	if id, err := strconv.Atoi(r.URL.Query().Get("destinataire_id")); err == nil {
		d.DestinataireID = &id
	}
	return d
}

func (s *Serveur) gererApercuDestinataires(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}
	demande := demandeDepuisRequete(r)
	if message := demande.validerCible(s.modele.Classes); message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}
	ids, err := s.destinatairesNotification(demande)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	envoyerJSON(w, http.StatusOK, map[string]int{"nombre": len(ids)})
}

func (s *Serveur) gererEnvoiNotificationAdmin(w http.ResponseWriter, r *http.Request) {
	auteurID, ok := s.exigerRole(w, r, roleAdministrateur)
	if !ok {
		return
	}
	var demande DemandeNotification
	if err := json.NewDecoder(r.Body).Decode(&demande); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "données invalides")
		return
	}
	if message := demande.valider(s.modele.Classes); message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}

	ids, err := s.destinatairesNotification(demande)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if len(ids) == 0 {
		envoyerErreur(w, http.StatusBadRequest, "aucun destinataire actif pour cette cible")
		return
	}

	statut := statutNotificationOK
	var erreurEnvoi, onesignalID *string
	identifiant, err := envoyerOneSignal(ids, ContenuNotification{
		TitreFR: demande.TitreFR, MessageFR: demande.MessageFR, TitreEN: demande.TitreEN, MessageEN: demande.MessageEN,
	})
	if err != nil {
		statut = statutNotificationEchec
		message := err.Error()
		erreurEnvoi = &message
	} else {
		onesignalID = &identifiant
	}

	var destinataireID *int
	if demande.Cible == "utilisateur" {
		destinataireID = demande.DestinataireID
	}
	var classe *string
	if demande.Cible == "prediction" {
		classe = &demande.ClassePredite
	}

	var id int
	err = s.db.QueryRow(`
		INSERT INTO notifications (titre_fr, message_fr, titre_en, message_en, cible, destinataire_id, classe_predite,
		                           nb_destinataires, statut, erreur, onesignal_id, auteur_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id`,
		demande.TitreFR, demande.MessageFR, demande.TitreEN, demande.MessageEN, demande.Cible, destinataireID, classe,
		len(ids), statut, erreurEnvoi, onesignalID, auteurID,
	).Scan(&id)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	notification, err := s.lireNotification(id)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	envoyerJSON(w, http.StatusCreated, notification)
}

const requeteNotifications = `
	SELECT n.id, n.titre_fr, n.message_fr, n.titre_en, n.message_en, n.cible, n.destinataire_id,
	       CASE WHEN d.id IS NULL THEN NULL ELSE d.prenom || ' ' || d.nom END,
	       n.classe_predite, n.nb_destinataires, n.statut, n.erreur,
	       CASE WHEN a.id IS NULL THEN NULL ELSE a.prenom || ' ' || a.nom END,
	       n.date_envoi
	FROM notifications n
	LEFT JOIN utilisateurs d ON d.id = n.destinataire_id
	LEFT JOIN utilisateurs a ON a.id = n.auteur_id`

func lireLigneNotification(ligne scanneur) (Notification, error) {
	var n Notification
	err := ligne.Scan(&n.ID, &n.TitreFR, &n.MessageFR, &n.TitreEN, &n.MessageEN, &n.Cible, &n.DestinataireID,
		&n.Destinataire, &n.ClassePredite, &n.NbDestinataires, &n.Statut, &n.Erreur, &n.Auteur, &n.DateEnvoi)
	return n, err
}

func (s *Serveur) lireNotification(id int) (Notification, error) {
	return lireLigneNotification(s.db.QueryRow(requeteNotifications+" WHERE n.id = $1", id))
}

func (s *Serveur) gererHistoriqueNotifications(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	resultat := PageNotifications{Notifications: make([]Notification, 0), Page: page, ParPage: notificationsParPage}
	if err := s.db.QueryRow("SELECT count(*) FROM notifications").Scan(&resultat.Total); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	lignes, err := s.db.Query(requeteNotifications+" ORDER BY n.date_envoi DESC, n.id DESC LIMIT $1 OFFSET $2",
		notificationsParPage, (page-1)*notificationsParPage)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()
	for lignes.Next() {
		n, err := lireLigneNotification(lignes)
		if err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		resultat.Notifications = append(resultat.Notifications, n)
	}
	if err := lignes.Err(); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	envoyerJSON(w, http.StatusOK, resultat)
}
