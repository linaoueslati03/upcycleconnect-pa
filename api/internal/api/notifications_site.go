package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

// Notifications affichées dans le site (cloche de l'espace). Chaque appel à notifierUtilisateur()
// les enregistre ici, ce qui marche sans HTTPS ni OneSignal ; le push reste un bonus.

// Base utilisée par notifierUtilisateur(), qui est appelée depuis tout le code sans le Serveur.
var baseNotifications *sql.DB

const notificationsAffichees = 30

type NotificationSite struct {
	ID           int       `json:"id"`
	Titre        string    `json:"titre"`
	Message      string    `json:"message"`
	Lien         *string   `json:"lien"`
	Lue          bool      `json:"lue"`
	DateCreation time.Time `json:"date_creation"`
}

func enregistrerNotificationSite(utilisateurID int, titre, message, lien string) {
	if baseNotifications == nil {
		return
	}
	var lienSQL any
	if lien != "" {
		lienSQL = lien
	}
	if _, err := baseNotifications.Exec(
		"INSERT INTO notifications_utilisateurs (utilisateur_id, titre, message, lien) VALUES ($1, $2, $3, $4)",
		utilisateurID, titre, message, lienSQL,
	); err != nil {
		log.Printf("notification du site impossible pour l'utilisateur %d : %v", utilisateurID, err)
	}
}

func (s *Serveur) gererMesNotifications(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	resultat := struct {
		Notifications []NotificationSite `json:"notifications"`
		NonLues       int                `json:"non_lues"`
	}{Notifications: make([]NotificationSite, 0)}

	if err := s.db.QueryRow("SELECT count(*) FROM notifications_utilisateurs WHERE utilisateur_id = $1 AND NOT lue",
		utilisateurID).Scan(&resultat.NonLues); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	lignes, err := s.db.Query(`
		SELECT id, titre, message, lien, lue, date_creation FROM notifications_utilisateurs
		WHERE utilisateur_id = $1 ORDER BY date_creation DESC, id DESC LIMIT $2`, utilisateurID, notificationsAffichees)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()
	for lignes.Next() {
		var n NotificationSite
		if err := lignes.Scan(&n.ID, &n.Titre, &n.Message, &n.Lien, &n.Lue, &n.DateCreation); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		resultat.Notifications = append(resultat.Notifications, n)
	}
	envoyerJSON(w, http.StatusOK, resultat)
}

func (s *Serveur) gererNotificationLue(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}
	resultat, err := s.db.Exec("UPDATE notifications_utilisateurs SET lue = true WHERE id = $1 AND utilisateur_id = $2",
		id, utilisateurID)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if n, _ := resultat.RowsAffected(); n == 0 {
		envoyerErreur(w, http.StatusNotFound, "notification introuvable")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Serveur) gererToutesNotificationsLues(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}
	if _, err := s.db.Exec("UPDATE notifications_utilisateurs SET lue = true WHERE utilisateur_id = $1 AND NOT lue",
		utilisateurID); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
