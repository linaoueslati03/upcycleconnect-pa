package api

import (
	"database/sql"
	"net/http"
	"time"
)

type DashboardReponse struct {
	AnnoncesActives     int        `json:"annonces_actives"`
	DernierDepotStatut  *string    `json:"dernier_depot_statut"`
	UpcyclingScore      int        `json:"upcycling_score"`
	ProchaineOffreTitre *string    `json:"prochaine_offre_titre"`
	ProchaineOffreDate  *time.Time `json:"prochaine_offre_date"`
}

func (s *Serveur) gererDashboard(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	var reponse DashboardReponse

	err = s.db.QueryRow(
		"SELECT COUNT(*) FROM annonces WHERE utilisateur_id = $1 AND statut = 'en_ligne'",
		utilisateurID,
	).Scan(&reponse.AnnoncesActives)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	var statutDepot sql.NullString
	err = s.db.QueryRow(
		"SELECT statut FROM depots WHERE utilisateur_id = $1 ORDER BY date_demande DESC LIMIT 1",
		utilisateurID,
	).Scan(&statutDepot)
	if err != nil && err != sql.ErrNoRows {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if statutDepot.Valid {
		reponse.DernierDepotStatut = &statutDepot.String
	}

	err = s.db.QueryRow(
		"SELECT upcycling_score FROM utilisateurs WHERE id = $1",
		utilisateurID,
	).Scan(&reponse.UpcyclingScore)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	var titre sql.NullString
	var date sql.NullTime
	err = s.db.QueryRow(`
		SELECT titre, date_debut FROM (
			SELECT f.titre, f.date_debut FROM inscriptions i
				JOIN formations f ON i.offre_id = f.id AND i.type_offre = 'formation'
				WHERE i.utilisateur_id = $1
			UNION ALL
			SELECT a.titre, a.date_debut FROM inscriptions i
				JOIN ateliers a ON i.offre_id = a.id AND i.type_offre = 'atelier'
				WHERE i.utilisateur_id = $1
			UNION ALL
			SELECT e.titre, e.date_debut FROM inscriptions i
				JOIN evenements e ON i.offre_id = e.id AND i.type_offre = 'evenement'
				WHERE i.utilisateur_id = $1
		) AS prochaines_offres
		WHERE date_debut > now()
		ORDER BY date_debut ASC
		LIMIT 1`,
		utilisateurID,
	).Scan(&titre, &date)
	if err != nil && err != sql.ErrNoRows {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if titre.Valid {
		reponse.ProchaineOffreTitre = &titre.String
		reponse.ProchaineOffreDate = &date.Time
	}

	envoyerJSON(w, http.StatusOK, reponse)
}
