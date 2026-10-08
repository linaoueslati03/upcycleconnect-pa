// Commande api : démarre l'API HTTP d'UpcycleConnect sur le port 8081.
package main

import (
	"log"
	"net/http"

	"github.com/linaoueslati03/upcycleconnect-pa/api/internal/api"
	"github.com/linaoueslati03/upcycleconnect-pa/api/internal/database"
	"github.com/linaoueslati03/upcycleconnect-pa/api/internal/middleware"
)

func main() {
	db, err := database.Connecter()
	if err != nil {
		log.Fatalf("connexion base de données impossible : %v", err)
	}
	defer db.Close()
	log.Println("connexion à la base de données réussie")

	serveur := api.NouveauServeur(db)

	log.Println("API démarrée sur http://localhost:8081")
	if err := http.ListenAndServe(":8081", middleware.CORS(serveur.Routes())); err != nil {
		log.Fatal(err)
	}
}
