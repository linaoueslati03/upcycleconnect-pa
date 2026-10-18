package api

import (
	"database/sql"
	"net/http"
	"time"
)

type EntreePlanning struct {
	InscriptionID int        `json:"inscription_id"`
	Type          string     `json:"type"`
	OffreID       int        `json:"offre_id"`
	Titre         string     `json:"titre"`
	Description   *string    `json:"description"`
	DateDebut     time.Time  `json:"date_debut"`
	DateFin       *time.Time `json:"date_fin"`
	Lieu          *string    `json:"lieu"`
}

func (s *Serveur) gererMonPlanning(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	requete := `
		SELECT i.id, 'formation', f.id, f.titre, f.description, f.date_debut, f.date_fin, f.lieu
			FROM inscriptions i JOIN formations f ON i.offre_id = f.id AND i.type_offre = 'formation'
			WHERE i.utilisateur_id = $1
		UNION ALL
		SELECT i.id, 'atelier', a.id, a.titre, a.description, a.date_debut, a.date_fin, a.lieu
			FROM inscriptions i JOIN ateliers a ON i.offre_id = a.id AND i.type_offre = 'atelier'
			WHERE i.utilisateur_id = $1
		UNION ALL
		SELECT i.id, 'evenement', e.id, e.titre, e.description, e.date_debut, e.date_fin, e.lieu
			FROM inscriptions i JOIN evenements e ON i.offre_id = e.id AND i.type_offre = 'evenement'
			WHERE i.utilisateur_id = $1
		ORDER BY date_debut ASC`

	lignes, err := s.db.Query(requete, utilisateurID)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	planning := make([]EntreePlanning, 0)
	for lignes.Next() {
		var e EntreePlanning
		var description, lieu sql.NullString
		var dateFin sql.NullTime

		if err := lignes.Scan(&e.InscriptionID, &e.Type, &e.OffreID, &e.Titre, &description, &e.DateDebut, &dateFin, &lieu); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		if description.Valid {
			e.Description = &description.String
		}
		if dateFin.Valid {
			e.DateFin = &dateFin.Time
		}
		if lieu.Valid {
			e.Lieu = &lieu.String
		}

		planning = append(planning, e)
	}

	envoyerJSON(w, http.StatusOK, planning)
}

type EntreePlanningSalarie struct {
	Type      string     `json:"type"`
	ID        int        `json:"id"`
	Titre     string     `json:"titre"`
	DateDebut time.Time  `json:"date_debut"`
	DateFin   *time.Time `json:"date_fin"`
	Lieu      *string    `json:"lieu"`
	Statut    string     `json:"statut"`
}

func (s *Serveur) gererPlanningSalarie(w http.ResponseWriter, r *http.Request) {
	salarieID, _, ok := s.exigerSalarie(w, r)
	if !ok {
		return
	}

	lignes, err := s.db.Query(`
		SELECT 'evenement', e.id, e.titre, e.date_debut, e.date_fin, e.lieu, e.statut
		FROM evenements e
		WHERE e.createur_id = $1
		   OR e.id IN (SELECT evenement_id FROM evenement_salarie WHERE salarie_id = $1)
		UNION ALL
		SELECT 'atelier', a.id, a.titre, a.date_debut, a.date_fin, a.lieu, a.statut
		FROM ateliers a
		WHERE a.createur_id = $1 OR a.responsable_id = $1
		   OR a.id IN (SELECT atelier_id FROM atelier_intervenant WHERE salarie_id = $1)
		ORDER BY 4`, salarieID)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	planning := make([]EntreePlanningSalarie, 0)
	for lignes.Next() {
		var e EntreePlanningSalarie
		if err := lignes.Scan(&e.Type, &e.ID, &e.Titre, &e.DateDebut, &e.DateFin, &e.Lieu, &e.Statut); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		planning = append(planning, e)
	}

	envoyerJSON(w, http.StatusOK, planning)
}

type EntreePlanningGlobal struct {
	UtilisateurID int       `json:"utilisateur_id"`
	Personne      string    `json:"personne"`
	RoleID        int       `json:"role_id"`
	Participation string    `json:"participation"`
	Type          string    `json:"type"`
	OffreID       int       `json:"offre_id"`
	Titre         string    `json:"titre"`
	DateDebut     time.Time `json:"date_debut"`
	Lieu          *string   `json:"lieu"`
}

func (s *Serveur) gererPlanningsAdmin(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	lignes, err := s.db.Query(`
		SELECT u.id, u.prenom || ' ' || u.nom, u.role_id, 'inscrit', o.type, o.id, o.titre, o.date_debut, o.lieu
		FROM inscriptions i
		JOIN utilisateurs u ON u.id = i.utilisateur_id
		JOIN (
			SELECT 'formation' AS type, id, titre, date_debut, lieu FROM formations
			UNION ALL SELECT 'atelier', id, titre, date_debut, lieu FROM ateliers
			UNION ALL SELECT 'evenement', id, titre, date_debut, lieu FROM evenements
		) o ON o.type = i.type_offre AND o.id = i.offre_id
		UNION ALL
		SELECT u.id, u.prenom || ' ' || u.nom, u.role_id, 'anime', a.type, a.id, a.titre, a.date_debut, a.lieu
		FROM (
			SELECT 'formation' AS type, id, titre, date_debut, lieu, createur_id AS animateur_id FROM formations
			UNION SELECT 'atelier', id, titre, date_debut, lieu, createur_id FROM ateliers
			UNION SELECT 'atelier', a.id, a.titre, a.date_debut, a.lieu, ai.salarie_id
				FROM ateliers a JOIN atelier_intervenant ai ON ai.atelier_id = a.id
			UNION SELECT 'evenement', id, titre, date_debut, lieu, createur_id FROM evenements
			UNION SELECT 'evenement', e.id, e.titre, e.date_debut, e.lieu, es.salarie_id
				FROM evenements e JOIN evenement_salarie es ON es.evenement_id = e.id
		) a
		JOIN utilisateurs u ON u.id = a.animateur_id
		ORDER BY 8, 2`)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	plannings := make([]EntreePlanningGlobal, 0)
	for lignes.Next() {
		var e EntreePlanningGlobal
		if err := lignes.Scan(&e.UtilisateurID, &e.Personne, &e.RoleID, &e.Participation, &e.Type, &e.OffreID,
			&e.Titre, &e.DateDebut, &e.Lieu); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		plannings = append(plannings, e)
	}

	envoyerJSON(w, http.StatusOK, plannings)
}
