package api

import "net/http"

func (s *Serveur) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/sante", s.gererSante)
	mux.HandleFunc("POST /api/comptes", s.gererCreationCompte)
	mux.HandleFunc("POST /api/login", s.gererLogin)
	mux.HandleFunc("POST /api/logout", s.gererLogout)
	mux.HandleFunc("GET /api/moi", s.gererMonCompte)
	mux.HandleFunc("PUT /api/moi", s.gererModificationMonCompte)
	mux.HandleFunc("PUT /api/moi/tutoriel", s.gererTutorielVu)
	mux.HandleFunc("PUT /api/moi/mot-de-passe", s.gererChangementMotDePasse)
	mux.HandleFunc("PUT /api/moi/langue", s.gererChangementLanguePreferee)
	mux.HandleFunc("GET /api/moi/dashboard", s.gererDashboard)

	mux.HandleFunc("GET /api/annonces", s.gererListeAnnonces)
	mux.HandleFunc("GET /api/annonces/{id}", s.gererDetailAnnonce)
	mux.HandleFunc("POST /api/annonces", s.gererCreationAnnonce)
	mux.HandleFunc("PUT /api/annonces/{id}", s.gererModificationAnnonce)
	mux.HandleFunc("DELETE /api/annonces/{id}", s.gererSuppressionAnnonce)
	mux.HandleFunc("GET /api/admin/annonces", s.gererModerationAnnonces)
	mux.HandleFunc("POST /api/annonces/{id}/valider", s.gererValidationAnnonce)
	mux.HandleFunc("POST /api/annonces/{id}/refuser", s.gererRefusAnnonce)

	mux.HandleFunc("POST /api/depots", s.gererCreationDepot)
	mux.HandleFunc("GET /api/moi/depots", s.gererMesDepots)
	mux.HandleFunc("GET /api/depots", s.gererListeDepots)
	mux.HandleFunc("GET /api/depots/{id}", s.gererDetailDepot)
	mux.HandleFunc("PUT /api/depots/{id}/statut", s.gererChangementStatutDepot)

	mux.HandleFunc("GET /api/categories", s.gererListeCategories)
	mux.HandleFunc("GET /api/conteneurs", s.gererListeConteneurs)
	mux.HandleFunc("GET /api/langues", s.gererListeLangues)
	mux.HandleFunc("GET /api/admin/langues", s.gererListeLanguesAdmin)
	mux.HandleFunc("POST /api/admin/langues", s.gererCreationLangue)
	mux.HandleFunc("PUT /api/admin/langues/{id}", s.gererModificationLangue)
	mux.HandleFunc("PUT /api/admin/langues/{id}/traduction", s.gererImportTraduction)

	mux.HandleFunc("GET /api/conseils", s.gererListeConseils)
	mux.HandleFunc("GET /api/conseils/{id}", s.gererDetailConseil)
	mux.HandleFunc("POST /api/conseils", s.gererCreationConseil)
	mux.HandleFunc("PUT /api/conseils/{id}", s.gererModificationConseil)

	mux.HandleFunc("GET /api/catalogue", s.gererCatalogue)
	mux.HandleFunc("POST /api/inscriptions", s.gererCreationInscription)

	mux.HandleFunc("GET /api/moi/score", s.gererMonScore)
	mux.HandleFunc("GET /api/moi/score/historique", s.gererMonScoreHistorique)
	mux.HandleFunc("GET /api/score/bareme", s.gererBaremeScore)
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
	mux.HandleFunc("PUT /api/forums/messages/{id}/statut", s.gererStatutMessage)
	mux.HandleFunc("PUT /api/forums/sujets/{id}/statut", s.gererStatutSujet)
	mux.HandleFunc("GET /api/moderation/signalements", s.gererListeSignalements)
	mux.HandleFunc("PUT /api/moderation/signalements/{id}", s.gererDecisionSignalement)
	mux.HandleFunc("GET /api/moderation/sujets", s.gererListeSujetsModeration)

	mux.HandleFunc("GET /api/utilisateurs", s.gererListeUtilisateurs)
	mux.HandleFunc("POST /api/utilisateurs", s.gererCreationUtilisateur)
	mux.HandleFunc("GET /api/utilisateurs/{id}", s.gererDetailUtilisateur)
	mux.HandleFunc("PUT /api/utilisateurs/{id}", s.gererModificationUtilisateur)
	mux.HandleFunc("DELETE /api/utilisateurs/{id}", s.gererSuppressionUtilisateur)
	mux.HandleFunc("PUT /api/utilisateurs/{id}/statut", s.gererChangementStatutUtilisateur)
	mux.HandleFunc("PUT /api/utilisateurs/{id}/responsable", s.gererChangementResponsable)

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
	mux.HandleFunc("GET /api/ateliers/{id}", s.gererDetailAtelier)
	mux.HandleFunc("POST /api/ateliers", s.gererCreationAtelier)
	mux.HandleFunc("PUT /api/ateliers/{id}", s.gererModificationAtelier)
	mux.HandleFunc("DELETE /api/ateliers/{id}", s.gererSuppressionAtelier)
	mux.HandleFunc("POST /api/ateliers/{id}/valider", s.gererValidationAtelier)

	mux.HandleFunc("GET /api/salaries/planning", s.gererPlanningSalarie)
	mux.HandleFunc("GET /api/admin/plannings", s.gererPlanningsAdmin)
	mux.HandleFunc("GET /api/salaries/conseils", s.gererMesConseils)

	return mux
}
