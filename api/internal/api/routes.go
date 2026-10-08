package api

import "net/http"

// Routes associe chaque route de l'API à son handler, avec le routeur de la
// bibliothèque standard (méthode HTTP et paramètres {id} gérés depuis Go 1.22).
//
// Les droits sont vérifiés au début de chaque handler (exigerRole, exigerSalarie) :
//   - publiques : santé, création de compte, login, lectures du catalogue, annonces, conseils,
//     événements, ateliers, prestations et forums ;
//   - utilisateur connecté : /api/moi/*, annonces, dépôts, inscriptions, projets, forums ;
//   - salarié : écriture des événements et ateliers (validation : responsable) ;
//   - administrateur : gestion des utilisateurs et des prestations.
func (s *Serveur) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/sante", s.gererSante)
	mux.HandleFunc("POST /api/comptes", s.gererCreationCompte)
	mux.HandleFunc("POST /api/login", s.gererLogin)
	mux.HandleFunc("POST /api/logout", s.gererLogout)
	mux.HandleFunc("GET /api/moi", s.gererMonCompte)
	mux.HandleFunc("PUT /api/moi", s.gererModificationMonCompte)
	mux.HandleFunc("PUT /api/moi/tutoriel", s.gererTutorielVu)
	mux.HandleFunc("GET /api/moi/dashboard", s.gererDashboard)

	mux.HandleFunc("GET /api/annonces", s.gererListeAnnonces)
	mux.HandleFunc("GET /api/annonces/{id}", s.gererDetailAnnonce)
	mux.HandleFunc("POST /api/annonces", s.gererCreationAnnonce)
	mux.HandleFunc("PUT /api/annonces/{id}", s.gererModificationAnnonce)
	mux.HandleFunc("DELETE /api/annonces/{id}", s.gererSuppressionAnnonce)

	mux.HandleFunc("POST /api/depots", s.gererCreationDepot)
	mux.HandleFunc("GET /api/moi/depots", s.gererMesDepots)
	mux.HandleFunc("GET /api/depots/{id}", s.gererDetailDepot)
	mux.HandleFunc("PUT /api/depots/{id}/statut", s.gererChangementStatutDepot)

	mux.HandleFunc("GET /api/conseils", s.gererListeConseils)
	mux.HandleFunc("GET /api/conseils/{id}", s.gererDetailConseil)

	mux.HandleFunc("GET /api/catalogue", s.gererCatalogue)
	mux.HandleFunc("POST /api/inscriptions", s.gererCreationInscription)

	mux.HandleFunc("GET /api/moi/score", s.gererMonScore)
	mux.HandleFunc("GET /api/moi/score/historique", s.gererMonScoreHistorique)
	mux.HandleFunc("GET /api/moi/planning", s.gererMonPlanning)

	mux.HandleFunc("POST /api/projets", s.gererCreationProjet)
	mux.HandleFunc("GET /api/moi/projets", s.gererMesProjets)
	mux.HandleFunc("GET /api/projets/{id}", s.gererDetailProjet)
	mux.HandleFunc("PUT /api/projets/{id}", s.gererModificationProjet)
	mux.HandleFunc("DELETE /api/projets/{id}", s.gererSuppressionProjet)
	mux.HandleFunc("POST /api/projets/{id}/etapes", s.gererCreationEtape)

	mux.HandleFunc("GET /api/forums/sujets", s.gererListeSujetsForum)
	mux.HandleFunc("POST /api/forums/sujets", s.gererCreationSujetForum)
	mux.HandleFunc("GET /api/forums/sujets/{id}/messages", s.gererMessagesSujet)
	mux.HandleFunc("POST /api/forums/sujets/{id}/messages", s.gererReponseSujet)
	mux.HandleFunc("POST /api/forums/messages/{id}/signalement", s.gererSignalementMessage)

	mux.HandleFunc("GET /api/utilisateurs", s.gererListeUtilisateurs)
	mux.HandleFunc("POST /api/utilisateurs", s.gererCreationUtilisateur)
	mux.HandleFunc("GET /api/utilisateurs/{id}", s.gererDetailUtilisateur)
	mux.HandleFunc("PUT /api/utilisateurs/{id}", s.gererModificationUtilisateur)
	mux.HandleFunc("DELETE /api/utilisateurs/{id}", s.gererSuppressionUtilisateur)

	mux.HandleFunc("GET /api/prestations", s.gererListePrestations)
	mux.HandleFunc("POST /api/prestations", s.gererCreationPrestation)
	mux.HandleFunc("PUT /api/prestations/{id}", s.gererModificationPrestation)
	mux.HandleFunc("DELETE /api/prestations/{id}", s.gererSuppressionPrestation)
	mux.HandleFunc("GET /api/prestations/{id}/pdf", s.gererPDFPrestation)

	mux.HandleFunc("GET /api/evenements", s.gererListeEvenements)
	mux.HandleFunc("GET /api/evenements/{id}", s.gererDetailEvenement)
	mux.HandleFunc("POST /api/evenements", s.gererCreationEvenement)
	mux.HandleFunc("PUT /api/evenements/{id}", s.gererModificationEvenement)
	mux.HandleFunc("DELETE /api/evenements/{id}", s.gererSuppressionEvenement)
	mux.HandleFunc("POST /api/evenements/{id}/valider", s.gererValidationEvenement)

	mux.HandleFunc("GET /api/ateliers", s.gererListeAteliers)
	mux.HandleFunc("POST /api/ateliers", s.gererCreationAtelier)

	return mux
}
