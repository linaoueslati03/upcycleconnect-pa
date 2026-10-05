package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	var err error
	db, err = connecterBaseDeDonnees()
	if err != nil {
		log.Fatalf("connexion base de données impossible : %v", err)
	}
	defer db.Close()
	log.Println("connexion à la base de données réussie")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/sante", gererSante)
	mux.HandleFunc("GET /api/utilisateurs", gererListeUtilisateurs)
	mux.HandleFunc("POST /api/utilisateurs", gererCreationUtilisateur)
	mux.HandleFunc("GET /api/utilisateurs/{id}", gererDetailUtilisateur)
	mux.HandleFunc("PUT /api/utilisateurs/{id}", gererModificationUtilisateur)
	mux.HandleFunc("DELETE /api/utilisateurs/{id}", gererSuppressionUtilisateur)

	mux.HandleFunc("GET /api/prestations", gererListePrestations)
	mux.HandleFunc("POST /api/prestations", gererCreationPrestation)
	mux.HandleFunc("PUT /api/prestations/{id}", gererModificationPrestation)
	mux.HandleFunc("DELETE /api/prestations/{id}", gererSuppressionPrestation)

	mux.HandleFunc("GET /evenements", gererListeEvenements)
	mux.HandleFunc("GET /evenements/{id}", gererDetailEvenement)
	mux.HandleFunc("POST /evenements", gererCreationEvenement)
	mux.HandleFunc("PUT /evenements/{id}", gererModificationEvenement)
	mux.HandleFunc("DELETE /evenements/{id}", gererSuppressionEvenement)
	mux.HandleFunc("POST /evenements/{id}/valider", gererValidationEvenement)

	log.Println("API démarrée sur http://localhost:8081")
	if err := http.ListenAndServe(":8081", autoriserCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func gererSante(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := db.Ping(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"statut": "erreur", "message": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"statut": "ok"})
}
