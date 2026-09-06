package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func connecterBaseDeDonnees() (*sql.DB, error) {
	hote := getenvDefault("DB_HOST", "localhost")
	port := getenvDefault("DB_PORT", "5433")
	utilisateur := getenvDefault("DB_USER", "postgres")
	motDePasse := getenvDefault("DB_PASSWORD", "postgres")
	nomBase := getenvDefault("DB_NAME", "upcycleconnect")

	chaineConnexion := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		hote, port, utilisateur, motDePasse, nomBase,
	)

	db, err := sql.Open("postgres", chaineConnexion)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func getenvDefault(cle, valeurParDefaut string) string {
	valeur := os.Getenv(cle)
	if valeur == "" {
		return valeurParDefaut
	}
	return valeur
}
