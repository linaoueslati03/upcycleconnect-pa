// Package api contient les handlers HTTP de l'API UpcycleConnect : une méthode de
// Serveur par route, regroupées par domaine (un fichier par domaine).
package api

import (
	"database/sql"
	"net/http"
)

// Serveur regroupe les dépendances partagées par tous les handlers.
// La connexion à la base est passée à la création au lieu d'être une variable globale.
type Serveur struct {
	db *sql.DB
}

// NouveauServeur crée un serveur qui utilise la connexion db.
func NouveauServeur(db *sql.DB) *Serveur {
	return &Serveur{db: db}
}

func (s *Serveur) gererSante(w http.ResponseWriter, r *http.Request) {
	if err := s.db.Ping(); err != nil {
		envoyerJSON(w, http.StatusServiceUnavailable, map[string]string{"statut": "erreur", "message": err.Error()})
		return
	}
	envoyerJSON(w, http.StatusOK, map[string]string{"statut": "ok"})
}
