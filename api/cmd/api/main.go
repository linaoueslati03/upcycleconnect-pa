package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

	if len(os.Args) == 3 && os.Args[1] == "reinitialiser-mot-de-passe" {
		motDePasse, err := serveur.ReinitialiserMotDePasseParEmail(os.Args[2])
		if err != nil {
			log.Fatalf("réinitialisation impossible : %v", err)
		}
		fmt.Printf("Mot de passe temporaire pour %s : %s\n", os.Args[2], motDePasse)
		fmt.Println("Il devra être changé à la prochaine connexion.")
		return
	}

	log.Println("API démarrée sur http://localhost:8081")
	if err := http.ListenAndServe(":8081", middleware.CORS(middleware.Erreurs(serveur.Routes()))); err != nil {
		log.Fatal(err)
	}
}
