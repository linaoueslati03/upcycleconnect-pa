package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

const nomComplet = `u.prenom || ' ' || u.nom`

type SignalementModeration struct {
	ID            int       `json:"id"`
	Motif         *string   `json:"motif"`
	Statut        string    `json:"statut"`
	CreatedAt     time.Time `json:"created_at"`
	SignalePar    string    `json:"signale_par"`
	MessageID     int       `json:"message_id"`
	Contenu       string    `json:"contenu"`
	StatutMessage string    `json:"statut_message"`
	AuteurMessage string    `json:"auteur_message"`
	SujetID       int       `json:"sujet_id"`
	TitreSujet    string    `json:"titre_sujet"`
}

type SujetModeration struct {
	ID        int       `json:"id"`
	Titre     string    `json:"titre"`
	Auteur    string    `json:"auteur"`
	Statut    string    `json:"statut"`
	CreatedAt time.Time `json:"created_at"`
}

type DecisionEntree struct {
	Statut string `json:"statut"`
}

func lireDecision(w http.ResponseWriter, r *http.Request, statutsAutorises ...string) (string, bool) {
	var entree DecisionEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return "", false
	}
	for _, statut := range statutsAutorises {
		if entree.Statut == statut {
			return entree.Statut, true
		}
	}
	envoyerErreur(w, http.StatusBadRequest, "statut invalide")
	return "", false
}

func journaliserModeration(tx *sql.Tx, moderateurID int, cibleType string, cibleID int, action string) error {
	_, err := tx.Exec(
		"INSERT INTO moderation_actions (moderateur_id, cible_type, cible_id, action) VALUES ($1, $2, $3, $4)",
		moderateurID, cibleType, cibleID, action,
	)
	return err
}

func (s *Serveur) gererListeSignalements(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := s.exigerSalarieOuAdmin(w, r); !ok {
		return
	}

	statut := r.URL.Query().Get("statut")
	if statut == "" {
		statut = "a_traiter"
	}

	lignes, err := s.db.Query(`
		SELECT sg.id, sg.motif, sg.statut, sg.created_at, `+nomComplet+`,
		       m.id, m.contenu, m.statut, auteur.prenom || ' ' || auteur.nom, su.id, su.titre
		FROM forum_signalements sg
		JOIN utilisateurs u ON u.id = sg.signale_par_utilisateur_id
		JOIN forum_messages m ON m.id = sg.message_id
		JOIN utilisateurs auteur ON auteur.id = m.auteur_utilisateur_id
		JOIN forum_sujets su ON su.id = m.sujet_id
		WHERE $1 = 'tous' OR sg.statut = $1
		ORDER BY sg.created_at DESC`, statut)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	signalements := make([]SignalementModeration, 0)
	for lignes.Next() {
		var sg SignalementModeration
		if err := lignes.Scan(&sg.ID, &sg.Motif, &sg.Statut, &sg.CreatedAt, &sg.SignalePar, &sg.MessageID, &sg.Contenu,
			&sg.StatutMessage, &sg.AuteurMessage, &sg.SujetID, &sg.TitreSujet); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		signalements = append(signalements, sg)
	}

	envoyerJSON(w, http.StatusOK, signalements)
}

func (s *Serveur) gererDecisionSignalement(w http.ResponseWriter, r *http.Request) {
	moderateurID, _, ok := s.exigerSalarieOuAdmin(w, r)
	if !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	statut, ok := lireDecision(w, r, "traite", "rejete")
	if !ok {
		return
	}

	tx, err := s.db.Begin()
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer tx.Rollback()

	resultat, err := tx.Exec("UPDATE forum_signalements SET statut = $1 WHERE id = $2", statut, id)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if lignes, _ := resultat.RowsAffected(); lignes == 0 {
		envoyerErreur(w, http.StatusNotFound, "signalement introuvable")
		return
	}

	if err := journaliserModeration(tx, moderateurID, "signalement", id, statut); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if err := tx.Commit(); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]string{"statut": statut})
}

func (s *Serveur) gererStatutMessage(w http.ResponseWriter, r *http.Request) {
	moderateurID, _, ok := s.exigerSalarieOuAdmin(w, r)
	if !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	statut, ok := lireDecision(w, r, "masque", "visible")
	if !ok {
		return
	}

	tx, err := s.db.Begin()
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer tx.Rollback()

	resultat, err := tx.Exec("UPDATE forum_messages SET statut = $1 WHERE id = $2", statut, id)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if lignes, _ := resultat.RowsAffected(); lignes == 0 {
		envoyerErreur(w, http.StatusNotFound, "message introuvable")
		return
	}

	if statut == "masque" {
		_, err := tx.Exec("UPDATE forum_signalements SET statut = 'traite' WHERE message_id = $1 AND statut = 'a_traiter'", id)
		if err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
	}

	if err := journaliserModeration(tx, moderateurID, "message", id, statut); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if err := tx.Commit(); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]string{"statut": statut})
}

func (s *Serveur) gererListeSujetsModeration(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := s.exigerSalarieOuAdmin(w, r); !ok {
		return
	}

	statut := r.URL.Query().Get("statut")
	if statut == "" {
		statut = "en_attente"
	}

	lignes, err := s.db.Query(`
		SELECT f.id, f.titre, `+nomComplet+`, f.statut, f.created_at
		FROM forum_sujets f JOIN utilisateurs u ON u.id = f.auteur_utilisateur_id
		WHERE $1 = 'tous' OR f.statut = $1
		ORDER BY f.created_at DESC`, statut)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	sujets := make([]SujetModeration, 0)
	for lignes.Next() {
		var su SujetModeration
		if err := lignes.Scan(&su.ID, &su.Titre, &su.Auteur, &su.Statut, &su.CreatedAt); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		sujets = append(sujets, su)
	}

	envoyerJSON(w, http.StatusOK, sujets)
}

func (s *Serveur) gererStatutSujet(w http.ResponseWriter, r *http.Request) {
	moderateurID, _, ok := s.exigerSalarieOuAdmin(w, r)
	if !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	statut, ok := lireDecision(w, r, "valide", "rejete")
	if !ok {
		return
	}

	tx, err := s.db.Begin()
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer tx.Rollback()

	resultat, err := tx.Exec("UPDATE forum_sujets SET statut = $1 WHERE id = $2", statut, id)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if lignes, _ := resultat.RowsAffected(); lignes == 0 {
		envoyerErreur(w, http.StatusNotFound, "sujet introuvable")
		return
	}

	if err := journaliserModeration(tx, moderateurID, "sujet", id, statut); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if err := tx.Commit(); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]string{"statut": statut})
}
