package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type RequeteNotification struct {
	UtilisateurID int    `json:"utilisateur_id"`
	Titre         string `json:"titre"`
	Message       string `json:"message"`
}

func envoyerNotification(utilisateurID int, titre, message string) error {
	_, err := envoyerOneSignal([]int{utilisateurID}, ContenuNotification{
		TitreFR: titre, MessageFR: message, TitreEN: titre, MessageEN: message,
	})
	return err
}

func notifierUtilisateur(utilisateurID int, titre, message string) {
	notifierAvecLien(utilisateurID, titre, message, "")
}

// Notification enregistrée dans le site (cloche de l'espace), puis envoyée en push OneSignal.
// Le lien (facultatif) est la page de l'espace à ouvrir depuis la notification.
func notifierAvecLien(utilisateurID int, titre, message, lien string) {
	if utilisateurID <= 0 {
		return
	}
	enregistrerNotificationSite(utilisateurID, titre, message, lien)
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
