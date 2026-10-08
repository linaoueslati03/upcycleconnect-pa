// Package database ouvre la connexion à la base PostgreSQL.
package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // pilote PostgreSQL, enregistré sous le nom "postgres"
)

// Connecter ouvre la connexion à partir des variables d'environnement DB_*
// et vérifie que la base répond.
func Connecter() (*sql.DB, error) {
	chaineConnexion := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		variable("DB_HOST", "localhost"),
		variable("DB_PORT", "5433"),
		variable("DB_USER", "postgres"),
		variable("DB_PASSWORD", "postgres"),
		variable("DB_NAME", "upcycleconnect"),
	)

	db, err := sql.Open("postgres", chaineConnexion)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// variable renvoie la variable d'environnement cle, ou valeurParDefaut si elle est vide.
func variable(cle, valeurParDefaut string) string {
	if valeur := os.Getenv(cle); valeur != "" {
		return valeur
	}
	return valeurParDefaut
}
