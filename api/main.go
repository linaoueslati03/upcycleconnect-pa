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
	mux.HandleFunc("POST /api/login", gererLogin)
	mux.HandleFunc("POST /api/logout", gererLogout)
	mux.HandleFunc("GET /api/moi", gererMonCompte)
	mux.HandleFunc("PUT /api/moi", gererModificationMonCompte)
	mux.HandleFunc("PUT /api/moi/tutoriel", gererTutorielVu)
	mux.HandleFunc("GET /api/moi/dashboard", gererDashboard)

	mux.HandleFunc("GET /api/annonces", gererListeAnnonces)
	mux.HandleFunc("GET /api/annonces/{id}", gererDetailAnnonce)
	mux.HandleFunc("POST /api/annonces", gererCreationAnnonce)
	mux.HandleFunc("PUT /api/annonces/{id}", gererModificationAnnonce)
	mux.HandleFunc("DELETE /api/annonces/{id}", gererSuppressionAnnonce)

	mux.HandleFunc("POST /api/depots", gererCreationDepot)
	mux.HandleFunc("GET /api/moi/depots", gererMesDepots)
	mux.HandleFunc("GET /api/depots/{id}", gererDetailDepot)
	mux.HandleFunc("PUT /api/depots/{id}/statut", gererChangementStatutDepot)

	mux.HandleFunc("GET /api/conseils", gererListeConseils)
	mux.HandleFunc("GET /api/conseils/{id}", gererDetailConseil)

	mux.HandleFunc("GET /api/catalogue", gererCatalogue)
	mux.HandleFunc("POST /api/inscriptions", gererCreationInscription)

	mux.HandleFunc("GET /api/moi/score", gererMonScore)
	mux.HandleFunc("GET /api/moi/score/historique", gererMonScoreHistorique)
	mux.HandleFunc("GET /api/moi/planning", gererMonPlanning)

	mux.HandleFunc("POST /api/projets", gererCreationProjet)
	mux.HandleFunc("GET /api/moi/projets", gererMesProjets)
	mux.HandleFunc("GET /api/projets/{id}", gererDetailProjet)
	mux.HandleFunc("PUT /api/projets/{id}", gererModificationProjet)
	mux.HandleFunc("DELETE /api/projets/{id}", gererSuppressionProjet)
	mux.HandleFunc("POST /api/projets/{id}/etapes", gererCreationEtape)

	mux.HandleFunc("GET /api/forums/sujets", gererListeSujetsForum)
	mux.HandleFunc("POST /api/forums/sujets", gererCreationSujetForum)
	mux.HandleFunc("GET /api/forums/sujets/{id}/messages", gererMessagesSujet)
	mux.HandleFunc("POST /api/forums/sujets/{id}/messages", gererReponseSujet)
	mux.HandleFunc("POST /api/forums/messages/{id}/signalement", gererSignalementMessage)
	mux.HandleFunc("GET /api/utilisateurs", gererListeUtilisateurs)
	mux.HandleFunc("POST /api/utilisateurs", gererCreationUtilisateur)
	mux.HandleFunc("GET /api/utilisateurs/{id}", gererDetailUtilisateur)
	mux.HandleFunc("PUT /api/utilisateurs/{id}", gererModificationUtilisateur)
	mux.HandleFunc("DELETE /api/utilisateurs/{id}", gererSuppressionUtilisateur)

	mux.HandleFunc("GET /api/prestations", gererListePrestations)
	mux.HandleFunc("POST /api/prestations", gererCreationPrestation)
	mux.HandleFunc("PUT /api/prestations/{id}", gererModificationPrestation)
	mux.HandleFunc("DELETE /api/prestations/{id}", gererSuppressionPrestation)

	mux.HandleFunc("GET /api/evenements", gererListeEvenements)
	mux.HandleFunc("GET /api/evenements/{id}", gererDetailEvenement)
	mux.HandleFunc("POST /api/evenements", gererCreationEvenement)
	mux.HandleFunc("PUT /api/evenements/{id}", gererModificationEvenement)
	mux.HandleFunc("DELETE /api/evenements/{id}", gererSuppressionEvenement)
	mux.HandleFunc("POST /api/evenements/{id}/valider", gererValidationEvenement)

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
