package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/linaoueslati03/upcycleconnect-pa/api/internal/ml"
)

type Serveur struct {
	db     *sql.DB
	modele *ml.Modele
}

func NouveauServeur(db *sql.DB, modele *ml.Modele) *Serveur {
	baseNotifications = db
	return &Serveur{db: db, modele: modele}
}

func (s *Serveur) gererSante(w http.ResponseWriter, r *http.Request) {
	if err := s.db.Ping(); err != nil {
		log.Println("base de données injoignable :", err)
		envoyerJSON(w, http.StatusServiceUnavailable, map[string]string{"statut": "erreur", "message": "base de données injoignable"})
		return
	}
	envoyerJSON(w, http.StatusOK, map[string]string{"statut": "ok"})
}
