package api

import (
	"database/sql"
	"net/http"
	"time"
)

const (
	pointsDepotRecupere    = 10
	pointsAnnonceCedee     = 15
	pointsInscriptionOffre = 5
)

// RegleScore décrit une action qui rapporte des points (affichée sur la page du score).
type RegleScore struct {
	Cle    string `json:"cle"` // clé de traduction du libellé de l'action
	Points int    `json:"points"`
}

// gererBaremeScore renvoie le barème réellement appliqué par l'API : le front l'affiche
// sans recopier les valeurs, qui ne sont définies qu'à un seul endroit (les constantes).
func (s *Serveur) gererBaremeScore(w http.ResponseWriter, r *http.Request) {
	envoyerJSON(w, http.StatusOK, []RegleScore{
		{Cle: "score.regle.inscription", Points: pointsInscriptionOffre},
		{Cle: "score.regle.depot", Points: pointsDepotRecupere},
		{Cle: "score.regle.annonce", Points: pointsAnnonceCedee},
	})
}

// executeur est satisfait à la fois par *sql.DB et *sql.Tx : permet d'appeler
// ajouterPointsScore aussi bien hors transaction (depots, annonces) que dans une
// transaction existante (inscriptions).
type executeur interface {
	Exec(query string, args ...any) (sql.Result, error)
}

func ajouterPointsScore(ex executeur, utilisateurID int, delta int, motif string) error {
	if _, err := ex.Exec(
		"INSERT INTO score_historique (utilisateur_id, delta, motif) VALUES ($1, $2, $3)",
		utilisateurID, delta, motif,
	); err != nil {
		return err
	}
	_, err := ex.Exec(
		"UPDATE utilisateurs SET upcycling_score = upcycling_score + $1 WHERE id = $2",
		delta, utilisateurID,
	)
	return err
}

func (s *Serveur) gererMonScore(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	var score int
	err = s.db.QueryRow("SELECT upcycling_score FROM utilisateurs WHERE id = $1", utilisateurID).Scan(&score)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]int{"upcycling_score": score})
}

type LigneScoreHistorique struct {
	ID    int       `json:"id"`
	Delta int       `json:"delta"`
	Motif *string   `json:"motif"`
	Date  time.Time `json:"date"`
}

func (s *Serveur) gererMonScoreHistorique(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	lignes, err := s.db.Query(
		"SELECT id, delta, motif, date FROM score_historique WHERE utilisateur_id = $1 ORDER BY date DESC",
		utilisateurID,
	)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	historique := make([]LigneScoreHistorique, 0)
	for lignes.Next() {
		var l LigneScoreHistorique
		var motif sql.NullString
		if err := lignes.Scan(&l.ID, &l.Delta, &motif, &l.Date); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		if motif.Valid {
			l.Motif = &motif.String
		}
		historique = append(historique, l)
	}

	envoyerJSON(w, http.StatusOK, historique)
}
