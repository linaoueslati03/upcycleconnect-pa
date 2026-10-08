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
