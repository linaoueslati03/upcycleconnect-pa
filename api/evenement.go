package main

import (
	"encoding/json"
	"log"
	"net/http"
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

func gererListeEvenements(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id, titre, description, date_debut, date_fin, lieu, site, statut, createur_id, valide_par_id
		FROM evenements
		ORDER BY date_debut ASC
	`)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var evenements []Evenement
	for rows.Next() {
		var e Evenement
		if err := rows.Scan(&e.ID, &e.Titre, &e.Description, &e.DateDebut, &e.DateFin, &e.Lieu, &e.Site, &e.Statut, &e.CreateurID, &e.ValideParID); err != nil {
			log.Println("Erreur de scan :", err)
			continue
		}
		evenements = append(evenements, e)
	}

	envoyerJSON(w, http.StatusOK, evenements)
}

func gererDetailEvenement(w http.ResponseWriter, r *http.Request) {
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	var e Evenement
	err = db.QueryRow(`
		SELECT id, titre, description, date_debut, date_fin, lieu, site, statut, createur_id, valide_par_id
		FROM evenements WHERE id=$1
	`, id).Scan(&e.ID, &e.Titre, &e.Description, &e.DateDebut, &e.DateFin, &e.Lieu, &e.Site, &e.Statut, &e.CreateurID, &e.ValideParID)
	if err != nil {
		envoyerErreur(w, http.StatusNotFound, "événement introuvable")
		return
	}

	envoyerJSON(w, http.StatusOK, e)
}

func gererCreationEvenement(w http.ResponseWriter, r *http.Request) {
	var e Evenement
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête invalide")
		return
	}

	// Un salarié ne peut créer un événement qu'en brouillon ou en attente de validation
	// (les boutons "Enregistrer en brouillon" / "Soumettre à validation" du formulaire
	// envoient cette valeur ; toute autre valeur est ramenée à "brouillon" par sécurité)
	if e.Statut != "brouillon" && e.Statut != "en_attente" {
		e.Statut = "brouillon"
	}

	err := db.QueryRow(`
		INSERT INTO evenements (titre, description, date_debut, date_fin, lieu, site, statut, createur_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`, e.Titre, e.Description, e.DateDebut, e.DateFin, e.Lieu, e.Site, e.Statut, e.CreateurID).Scan(&e.ID)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, err.Error())
		return
	}

	envoyerJSON(w, http.StatusCreated, e)
}

func gererModificationEvenement(w http.ResponseWriter, r *http.Request) {
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

	_, err = db.Exec(`
		UPDATE evenements
		SET titre=$1, description=$2, date_debut=$3, date_fin=$4, lieu=$5, site=$6
		WHERE id=$7
	`, e.Titre, e.Description, e.DateDebut, e.DateFin, e.Lieu, e.Site, id)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, err.Error())
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]string{"statut": "modifié"})
}

func gererSuppressionEvenement(w http.ResponseWriter, r *http.Request) {
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	_, err = db.Exec("DELETE FROM evenements WHERE id=$1", id)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Validation d'un événement par un responsable : passe le statut à "publie"
// et enregistre quel salarié a validé (règle de gestion : seul un responsable
// peut valider, ce contrôle sera fait via le middleware d'authentification)
func gererValidationEvenement(w http.ResponseWriter, r *http.Request) {
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	var corps struct {
		ValidateurID int `json:"validateur_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&corps); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête invalide")
		return
	}

	_, err = db.Exec(
		"UPDATE evenements SET statut='publie', valide_par_id=$1 WHERE id=$2",
		corps.ValidateurID, id,
	)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, err.Error())
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]string{"statut": "publié"})
}