package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type RequeteNotification struct {
	UtilisateurID int    `json:"utilisateur_id"`
	Titre         string `json:"titre"`
	Message       string `json:"message"`
}

func envoyerNotification(utilisateurID int, titre, message string) error {
	appID := os.Getenv("ONESIGNAL_APP_ID")
	cleREST := os.Getenv("ONESIGNAL_REST_API_KEY")

	if appID == "" || cleREST == "" {
		return fmt.Errorf("configuration OneSignal manquante")
	}

	donnees := map[string]interface{}{
		"app_id": appID,
		"include_aliases": map[string][]string{
			"external_id": {fmt.Sprintf("%d", utilisateurID)},
		},
		"target_channel": "push",
		"headings": map[string]string{
			"fr": titre,
		},
		"contents": map[string]string{
			"fr": message,
		},
	}

	corps, err := json.Marshal(donnees)
	if err != nil {
		return fmt.Errorf("erreur préparation notification : %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.onesignal.com/notifications",
		bytes.NewBuffer(corps),
	)
	if err != nil {
		return fmt.Errorf("erreur création requête : %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Key "+cleREST)

	client := &http.Client{Timeout: 10 * time.Second}

	reponse, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erreur envoi OneSignal : %w", err)
	}
	defer reponse.Body.Close()

	if reponse.StatusCode < 200 || reponse.StatusCode >= 300 {
		details, _ := io.ReadAll(reponse.Body)
		return fmt.Errorf("OneSignal HTTP %d : %s", reponse.StatusCode, string(details))
	}

	return nil
}

func notifierUtilisateur(utilisateurID int, titre, message string) {
	if utilisateurID <= 0 {
		return
	}
	if err := envoyerNotification(utilisateurID, titre, message); err != nil {
		log.Printf("notification OneSignal impossible pour l'utilisateur %d : %v", utilisateurID, err)
	}
}

func (s *Serveur) gererEnvoiNotification(w http.ResponseWriter, r *http.Request) {
	var demande RequeteNotification

	if err := json.NewDecoder(r.Body).Decode(&demande); err != nil {
		envoyerJSON(w, http.StatusBadRequest, map[string]string{
			"erreur": "Requête invalide",
		})
		return
	}

	if demande.UtilisateurID <= 0 || demande.Titre == "" || demande.Message == "" {
		envoyerJSON(w, http.StatusBadRequest, map[string]string{
			"erreur": "Identifiant, titre et message obligatoires",
		})
		return
	}

	idConnecte, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerJSON(w, http.StatusUnauthorized, map[string]string{
			"erreur": "Authentification requise",
		})
		return
	}

	if idConnecte != demande.UtilisateurID {
		envoyerJSON(w, http.StatusForbidden, map[string]string{
			"erreur": "Vous ne pouvez pas envoyer une notification à cet utilisateur",
		})
		return
	}

	if err := envoyerNotification(
		demande.UtilisateurID,
		demande.Titre,
		demande.Message,
	); err != nil {
		envoyerJSON(w, http.StatusBadGateway, map[string]string{
			"erreur": "Impossible d'envoyer la notification",
		})
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]string{
		"message": "Notification envoyée",
	})
}
