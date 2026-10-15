package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

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

func variable(cle, valeurParDefaut string) string {
	if valeur := os.Getenv(cle); valeur != "" {
		return valeur
	}
	return valeurParDefaut
}
