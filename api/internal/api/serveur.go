package api

import (
	"database/sql"
	"log"
	"net/http"
)

type Serveur struct {
	db *sql.DB
}

func NouveauServeur(db *sql.DB) *Serveur {
	return &Serveur{db: db}
}

func (s *Serveur) gererSante(w http.ResponseWriter, r *http.Request) {
	if err := s.db.Ping(); err != nil {
		log.Println("base de données injoignable :", err)
		envoyerJSON(w, http.StatusServiceUnavailable, map[string]string{"statut": "erreur", "message": "base de données injoignable"})
		return
	}
	envoyerJSON(w, http.StatusOK, map[string]string{"statut": "ok"})
}
