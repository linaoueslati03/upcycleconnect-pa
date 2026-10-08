package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type Evenement struct {
	ID          int        `json:"id"`
	Titre       string     `json:"titre"`
	Description string     `json:"description"`
	DateDebut   time.Time  `json:"date_debut"`
	DateFin     *time.Time `json:"date_fin,omitempty"`
	Lieu        string     `json:"lieu"`
	Site        string     `json:"site"`
	Statut      string     `json:"statut"`
	CreateurID  int        `json:"createur_id"`
	ValideParID *int       `json:"valide_par_id,omitempty"`
}

const colonnesEvenement = `id, titre, description, date_debut, date_fin, lieu, site, statut, createur_id, valide_par_id`

func scannerEvenement(ligne interface{ Scan(...any) error }, e *Evenement) error {
	return ligne.Scan(&e.ID, &e.Titre, &e.Description, &e.DateDebut, &e.DateFin, &e.Lieu, &e.Site, &e.Statut, &e.CreateurID, &e.ValideParID)
}

// verifierDroitOffre vérifie qu'un salarié peut modifier ou supprimer une offre : il doit
// en être le créateur, ou être responsable. Renvoie le statut actuel de l'offre ; en cas
// de refus, l'erreur HTTP est déjà écrite. table vaut "evenements" ou "ateliers" (valeurs
// fixées dans le code, jamais venues de la requête).
func (s *Serveur) verifierDroitOffre(w http.ResponseWriter, table string, id, salarieID int, estResponsable bool) (string, bool) {
	var createurID int
	var statut string
	err := s.db.QueryRow("SELECT createur_id, statut FROM "+table+" WHERE id = $1", id).Scan(&createurID, &statut)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "introuvable")
		return "", false
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return "", false
	}
	if createurID != salarieID && !estResponsable {
		envoyerErreur(w, http.StatusForbidden, "seuls le créateur et les responsables peuvent modifier cette offre")
		return "", false
	}
	return statut, true
}

// statutApresModificationOffre : une offre déjà publiée et modifiée par un salarié non
// responsable repart en validation ; sinon le salarié choisit brouillon ou en_attente,
// et un statut absent laisse le statut actuel.
func statutApresModificationOffre(statutActuel, statutDemande string, estResponsable bool) string {
	if statutActuel == "publie" {
		if estResponsable {
			return "publie"
		}
		return "en_attente"
	}
	if choisi := statutModifiable(statutDemande); choisi != "" {
		return choisi
	}
	return statutActuel
}

// statutModifiable garde seulement les statuts qu'un salarié peut choisir lui-même
// (« publie » est réservé à la validation par un responsable) ; "" = statut inchangé.
func statutModifiable(statut string) string {
	if statut == "brouillon" || statut == "en_attente" {
		return statut
	}
	return ""
}

func validerEvenement(e Evenement) string {
	if strings.TrimSpace(e.Titre) == "" {
		return "titre obligatoire"
	}
	if e.DateDebut.IsZero() {
		return "date_debut obligatoire"
	}
	return ""
}

func (s *Serveur) gererListeEvenements(w http.ResponseWriter, r *http.Request) {
	// Brouillons et offres en attente ne sont visibles que du personnel
	requete := `SELECT ` + colonnesEvenement + ` FROM evenements WHERE statut = 'publie' ORDER BY date_debut ASC`
	if _, role := s.roleConnecte(r); role == roleSalarie || role == roleAdministrateur {
		requete = `SELECT ` + colonnesEvenement + ` FROM evenements ORDER BY date_debut ASC`
	}
	lignes, err := s.db.Query(requete)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	evenements := make([]Evenement, 0)
	for lignes.Next() {
		var e Evenement
		if err := scannerEvenement(lignes, &e); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		evenements = append(evenements, e)
	}

	envoyerJSON(w, http.StatusOK, evenements)
}

func (s *Serveur) gererDetailEvenement(w http.ResponseWriter, r *http.Request) {
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	var e Evenement
	err = scannerEvenement(s.db.QueryRow(`SELECT `+colonnesEvenement+` FROM evenements WHERE id = $1`, id), &e)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "événement introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	// Un brouillon ou une offre en attente n'est visible que du personnel
	if e.Statut != "publie" {
		if _, role := s.roleConnecte(r); role != roleSalarie && role != roleAdministrateur {
			envoyerErreur(w, http.StatusNotFound, "introuvable")
			return
		}
	}

	envoyerJSON(w, http.StatusOK, e)
}

func (s *Serveur) gererCreationEvenement(w http.ResponseWriter, r *http.Request) {
	salarieID, _, ok := s.exigerSalarie(w, r)
	if !ok {
		return
	}

	var e Evenement
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête invalide")
		return
	}
	if message := validerEvenement(e); message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}

	// Un salarié ne peut créer un événement qu'en brouillon ou en attente de validation
	// (les boutons "Enregistrer en brouillon" / "Soumettre à validation" du formulaire
	// envoient cette valeur ; toute autre valeur est ramenée à "brouillon" par sécurité)
	if e.Statut != "brouillon" && e.Statut != "en_attente" {
		e.Statut = "brouillon"
	}

	// Le créateur est le salarié connecté : on ignore un éventuel createur_id envoyé par
	// le front, sinon n'importe qui pourrait créer un événement au nom d'un autre.
	err := scannerEvenement(s.db.QueryRow(`
		INSERT INTO evenements (titre, description, date_debut, date_fin, lieu, site, statut, createur_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING `+colonnesEvenement,
		e.Titre, e.Description, e.DateDebut, e.DateFin, e.Lieu, e.Site, e.Statut, salarieID,
	), &e)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusCreated, e)
}

func (s *Serveur) gererModificationEvenement(w http.ResponseWriter, r *http.Request) {
	salarieID, estResponsable, ok := s.exigerSalarie(w, r)
	if !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	var e Evenement
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête invalide")
		return
	}
	if message := validerEvenement(e); message != "" {
		envoyerErreur(w, http.StatusBadRequest, message)
		return
	}

	statutActuel, autorise := s.verifierDroitOffre(w, "evenements", id, salarieID, estResponsable)
	if !autorise {
		return
	}

	err = scannerEvenement(s.db.QueryRow(`
		UPDATE evenements
		SET titre = $1, description = $2, date_debut = $3, date_fin = $4, lieu = $5, site = $6,
		    statut = $7, updated_at = now()
		WHERE id = $8
		RETURNING `+colonnesEvenement,
		e.Titre, e.Description, e.DateDebut, e.DateFin, e.Lieu, e.Site, statutApresModificationOffre(statutActuel, e.Statut, estResponsable), id,
	), &e)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "événement introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, e)
}

func (s *Serveur) gererSuppressionEvenement(w http.ResponseWriter, r *http.Request) {
	salarieID, estResponsable, ok := s.exigerSalarie(w, r)
	if !ok {
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	if _, autorise := s.verifierDroitOffre(w, "evenements", id, salarieID, estResponsable); !autorise {
		return
	}

	resultat, err := s.db.Exec("DELETE FROM evenements WHERE id = $1", id)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if lignes, _ := resultat.RowsAffected(); lignes == 0 {
		envoyerErreur(w, http.StatusNotFound, "événement introuvable")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Validation d'un événement par un responsable : passe le statut de "en_attente" à
// "publie" et enregistre quel salarié a validé (le responsable connecté).
func (s *Serveur) gererValidationEvenement(w http.ResponseWriter, r *http.Request) {
	responsableID, estResponsable, ok := s.exigerSalarie(w, r)
	if !ok {
		return
	}
	if !estResponsable {
		envoyerErreur(w, http.StatusForbidden, "seul un responsable peut valider un événement")
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	var e Evenement
	err = scannerEvenement(s.db.QueryRow(`
		UPDATE evenements SET statut = 'publie', valide_par_id = $1, updated_at = now()
		WHERE id = $2 AND statut = 'en_attente'
		RETURNING `+colonnesEvenement,
		responsableID, id,
	), &e)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusConflict, "événement introuvable ou pas en attente de validation")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, e)
}
