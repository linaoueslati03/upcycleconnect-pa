package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

const delaiAntiSpamForum = 10 * time.Second

type SujetForum struct {
	ID                  int       `json:"id"`
	Titre               string    `json:"titre"`
	AuteurUtilisateurID int       `json:"auteur_utilisateur_id"`
	Statut              string    `json:"statut"`
	CreatedAt           time.Time `json:"created_at"`
}

type MessageForum struct {
	ID                  int       `json:"id"`
	SujetID             int       `json:"sujet_id"`
	AuteurUtilisateurID int       `json:"auteur_utilisateur_id"`
	Contenu             string    `json:"contenu"`
	Statut              string    `json:"statut"`
	CreatedAt           time.Time `json:"created_at"`
}

type SujetEntree struct {
	Titre   string `json:"titre"`
	Contenu string `json:"contenu"`
}

type MessageEntree struct {
	Contenu string `json:"contenu"`
}

type SignalementEntree struct {
	Motif string `json:"motif"`
}

func verifierAntiSpam(utilisateurID int) bool {
	var dernierMessage time.Time
	err := db.QueryRow(
		"SELECT created_at FROM forum_messages WHERE auteur_utilisateur_id = $1 ORDER BY created_at DESC LIMIT 1",
		utilisateurID,
	).Scan(&dernierMessage)
	if errors.Is(err, sql.ErrNoRows) {
		return true
	}
	if err != nil {
		return true
	}
	return time.Since(dernierMessage) >= delaiAntiSpamForum
}

func gererListeSujetsForum(w http.ResponseWriter, r *http.Request) {
	lignes, err := db.Query(`
		SELECT id, titre, auteur_utilisateur_id, statut, created_at
		FROM forum_sujets WHERE statut != 'rejete' ORDER BY created_at DESC`)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	sujets := make([]SujetForum, 0)
	for lignes.Next() {
		var s SujetForum
		if err := lignes.Scan(&s.ID, &s.Titre, &s.AuteurUtilisateurID, &s.Statut, &s.CreatedAt); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		sujets = append(sujets, s)
	}

	envoyerJSON(w, http.StatusOK, sujets)
}

func gererCreationSujetForum(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	var entree SujetEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if strings.TrimSpace(entree.Titre) == "" || strings.TrimSpace(entree.Contenu) == "" {
		envoyerErreur(w, http.StatusBadRequest, "titre et contenu obligatoires")
		return
	}
	if !verifierAntiSpam(utilisateurID) {
		envoyerErreur(w, http.StatusTooManyRequests, "veuillez patienter avant de publier à nouveau")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer tx.Rollback()

	var s SujetForum
	err = tx.QueryRow(
		"INSERT INTO forum_sujets (titre, auteur_utilisateur_id) VALUES ($1, $2) RETURNING id, titre, auteur_utilisateur_id, statut, created_at",
		entree.Titre, utilisateurID,
	).Scan(&s.ID, &s.Titre, &s.AuteurUtilisateurID, &s.Statut, &s.CreatedAt)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	_, err = tx.Exec(
		"INSERT INTO forum_messages (sujet_id, auteur_utilisateur_id, contenu) VALUES ($1, $2, $3)",
		s.ID, utilisateurID, entree.Contenu,
	)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	if err := tx.Commit(); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusCreated, s)
}

func gererMessagesSujet(w http.ResponseWriter, r *http.Request) {
	sujetID, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	lignes, err := db.Query(`
		SELECT id, sujet_id, auteur_utilisateur_id, contenu, statut, created_at
		FROM forum_messages WHERE sujet_id = $1 AND statut = 'visible' ORDER BY created_at ASC`, sujetID)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	messages := make([]MessageForum, 0)
	for lignes.Next() {
		var m MessageForum
		if err := lignes.Scan(&m.ID, &m.SujetID, &m.AuteurUtilisateurID, &m.Contenu, &m.Statut, &m.CreatedAt); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		messages = append(messages, m)
	}

	envoyerJSON(w, http.StatusOK, messages)
}

func gererReponseSujet(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	sujetID, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var entree MessageEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	if strings.TrimSpace(entree.Contenu) == "" {
		envoyerErreur(w, http.StatusBadRequest, "contenu obligatoire")
		return
	}
	if !verifierAntiSpam(utilisateurID) {
		envoyerErreur(w, http.StatusTooManyRequests, "veuillez patienter avant de publier à nouveau")
		return
	}

	var m MessageForum
	err = db.QueryRow(`
		INSERT INTO forum_messages (sujet_id, auteur_utilisateur_id, contenu)
		VALUES ($1, $2, $3)
		RETURNING id, sujet_id, auteur_utilisateur_id, contenu, statut, created_at`,
		sujetID, utilisateurID, entree.Contenu,
	).Scan(&m.ID, &m.SujetID, &m.AuteurUtilisateurID, &m.Contenu, &m.Statut, &m.CreatedAt)
	if err != nil {
		gererErreurPostgres(w, err)
		return
	}

	envoyerJSON(w, http.StatusCreated, m)
}

func gererSignalementMessage(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	messageID, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var entree SignalementEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}

	_, err = db.Exec(
		"INSERT INTO forum_signalements (message_id, signale_par_utilisateur_id, motif) VALUES ($1, $2, $3)",
		messageID, utilisateurID, entree.Motif,
	)
	if err != nil {
		gererErreurPostgres(w, err)
		return
	}

	envoyerJSON(w, http.StatusCreated, map[string]string{"message": "signalement enregistré"})
}
