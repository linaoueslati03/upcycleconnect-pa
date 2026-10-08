package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type OffreCatalogue struct {
	Type        string     `json:"type"`
	ID          int        `json:"id"`
	Titre       string     `json:"titre"`
	Description *string    `json:"description"`
	DateDebut   time.Time  `json:"date_debut"`
	DateFin     *time.Time `json:"date_fin"`
	Lieu        *string    `json:"lieu"`
	Tarif       *float64   `json:"tarif"`
	NbPlaces    *int       `json:"nb_places"`
	// Places encore disponibles (capacité moins inscrits) ; null si l'offre n'a pas de limite
	PlacesRestantes *int `json:"places_restantes"`
}

func (s *Serveur) gererCatalogue(w http.ResponseWriter, r *http.Request) {
	typeFiltre := r.URL.Query().Get("type")

	requete := `
		SELECT 'formation' AS type, f.id, f.titre, f.description, f.date_debut, f.date_fin, f.lieu, f.tarif, f.nb_places,
				f.nb_places - (SELECT COUNT(*) FROM inscriptions i WHERE i.type_offre = 'formation' AND i.offre_id = f.id)
			FROM formations f WHERE f.statut = 'publie'
		UNION ALL
		SELECT 'atelier' AS type, id, titre, description, date_debut, date_fin, lieu, NULL, NULL, NULL
			FROM ateliers WHERE statut = 'publie'
		UNION ALL
		SELECT 'evenement' AS type, id, titre, description, date_debut, date_fin, lieu, NULL, NULL, NULL
			FROM evenements WHERE statut = 'publie'
		ORDER BY date_debut ASC`

	lignes, err := s.db.Query(requete)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	offres := make([]OffreCatalogue, 0)
	for lignes.Next() {
		var o OffreCatalogue
		var description, lieu sql.NullString
		var dateFin sql.NullTime
		var tarif sql.NullFloat64
		var nbPlaces, placesRestantes sql.NullInt64

		if err := lignes.Scan(&o.Type, &o.ID, &o.Titre, &description, &o.DateDebut, &dateFin, &lieu, &tarif, &nbPlaces, &placesRestantes); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		if description.Valid {
			o.Description = &description.String
		}
		if dateFin.Valid {
			o.DateFin = &dateFin.Time
		}
		if lieu.Valid {
			o.Lieu = &lieu.String
		}
		if tarif.Valid {
			o.Tarif = &tarif.Float64
		}
		if nbPlaces.Valid {
			v := int(nbPlaces.Int64)
			o.NbPlaces = &v
		}
		if placesRestantes.Valid {
			v := int(placesRestantes.Int64)
			o.PlacesRestantes = &v
		}

		if typeFiltre == "" || typeFiltre == o.Type {
			offres = append(offres, o)
		}
	}

	envoyerJSON(w, http.StatusOK, offres)
}

type InscriptionEntree struct {
	TypeOffre string `json:"type_offre"`
	OffreID   int    `json:"offre_id"`
}

type Inscription struct {
	ID              int       `json:"id"`
	UtilisateurID   int       `json:"utilisateur_id"`
	TypeOffre       string    `json:"type_offre"`
	OffreID         int       `json:"offre_id"`
	DateInscription time.Time `json:"date_inscription"`
	StatutPaiement  string    `json:"statut_paiement"`
	MontantPaye     float64   `json:"montant_paye"`
}

var tableParTypeOffre = map[string]string{
	"formation": "formations",
	"atelier":   "ateliers",
	"evenement": "evenements",
}

func (s *Serveur) gererCreationInscription(w http.ResponseWriter, r *http.Request) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	var entree InscriptionEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}

	table, typeValide := tableParTypeOffre[entree.TypeOffre]
	if !typeValide {
		envoyerErreur(w, http.StatusBadRequest, "type_offre invalide (formation, atelier ou evenement attendu)")
		return
	}

	tx, err := s.db.Begin()
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer tx.Rollback()

	var statutOffre string
	var nbPlaces sql.NullInt64
	var dejaCommencee bool
	err = tx.QueryRow(
		"SELECT statut, "+placesColonne(entree.TypeOffre)+", date_debut <= now() FROM "+table+" WHERE id = $1 FOR UPDATE",
		entree.OffreID,
	).Scan(&statutOffre, &nbPlaces, &dejaCommencee)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "offre introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if statutOffre != "publie" {
		envoyerErreur(w, http.StatusBadRequest, "cette offre n'est pas publiée")
		return
	}
	if dejaCommencee {
		envoyerErreur(w, http.StatusBadRequest, "cette offre a déjà commencé")
		return
	}

	// L'offre est verrouillée (FOR UPDATE) : deux requêtes simultanées du même utilisateur
	// ne peuvent pas passer toutes les deux cette vérification.
	var dejaInscrit bool
	err = tx.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM inscriptions WHERE utilisateur_id = $1 AND type_offre = $2 AND offre_id = $3)",
		utilisateurID, entree.TypeOffre, entree.OffreID,
	).Scan(&dejaInscrit)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if dejaInscrit {
		envoyerErreur(w, http.StatusConflict, "vous êtes déjà inscrit à cette offre")
		return
	}

	if nbPlaces.Valid {
		var nbInscrits int
		err = tx.QueryRow(
			"SELECT COUNT(*) FROM inscriptions WHERE type_offre = $1 AND offre_id = $2",
			entree.TypeOffre, entree.OffreID,
		).Scan(&nbInscrits)
		if err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		if nbInscrits >= int(nbPlaces.Int64) {
			envoyerErreur(w, http.StatusConflict, "plus de places disponibles")
			return
		}
	}

	var ins Inscription
	err = tx.QueryRow(`
		INSERT INTO inscriptions (utilisateur_id, type_offre, offre_id)
		VALUES ($1, $2, $3)
		RETURNING id, utilisateur_id, type_offre, offre_id, date_inscription, statut_paiement, montant_paye`,
		utilisateurID, entree.TypeOffre, entree.OffreID,
	).Scan(&ins.ID, &ins.UtilisateurID, &ins.TypeOffre, &ins.OffreID, &ins.DateInscription, &ins.StatutPaiement, &ins.MontantPaye)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	if err := ajouterPointsScore(tx, utilisateurID, pointsInscriptionOffre, "Inscription à une offre"); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	if err := tx.Commit(); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusCreated, ins)
}

// placesColonne renvoie le nom de la colonne de places max pour un type d'offre, ou "NULL"
// si ce type n'a pas de limite de places (seules les formations en ont une).
func placesColonne(typeOffre string) string {
	if typeOffre == "formation" {
		return "nb_places"
	}
	return "NULL"
}
