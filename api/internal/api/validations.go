package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type tableOffre struct {
	table      string
	validateur string
}

var tablesOffres = map[string]tableOffre{
	"formation": {table: "formations", validateur: "valide_par_id"},
	"atelier":   {table: "ateliers", validateur: "responsable_id"},
	"evenement": {table: "evenements", validateur: "valide_par_id"},
}

type OffreAValider struct {
	Type        string     `json:"type"`
	ID          int        `json:"id"`
	Titre       string     `json:"titre"`
	Description string     `json:"description"`
	DateDebut   time.Time  `json:"date_debut"`
	DateFin     *time.Time `json:"date_fin"`
	Lieu        string     `json:"lieu"`
	Site        *string    `json:"site"`
	NbPlaces    *int       `json:"nb_places"`
	Tarif       *float64   `json:"tarif"`
	Createur    string     `json:"createur"`
	CreatedAt   time.Time  `json:"created_at"`
}

type DecisionOffreEntree struct {
	Decision string `json:"decision"`
	Motif    string `json:"motif"`
}

func (s *Serveur) gererListeValidations(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	lignes, err := s.db.Query(`
		SELECT 'formation', f.id, f.titre, COALESCE(f.description, ''), f.date_debut, f.date_fin, COALESCE(f.lieu, ''),
		       NULL::varchar, f.nb_places, f.tarif::float8, ` + nomComplet + `, f.created_at
		FROM formations f JOIN utilisateurs u ON u.id = f.createur_id
		WHERE f.statut = 'en_attente'
		UNION ALL
		SELECT 'atelier', a.id, a.titre, COALESCE(a.description, ''), a.date_debut, a.date_fin, COALESCE(a.lieu, ''),
		       NULL::varchar, NULL::int, NULL::float8, ` + nomComplet + `, a.created_at
		FROM ateliers a JOIN utilisateurs u ON u.id = a.createur_id
		WHERE a.statut = 'en_attente'
		UNION ALL
		SELECT 'evenement', e.id, e.titre, e.description, e.date_debut, e.date_fin, e.lieu,
		       e.site, NULL::int, NULL::float8, ` + nomComplet + `, e.created_at
		FROM evenements e JOIN utilisateurs u ON u.id = e.createur_id
		WHERE e.statut = 'en_attente'
		ORDER BY 12 ASC`)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	offres := make([]OffreAValider, 0)
	for lignes.Next() {
		var o OffreAValider
		if err := lignes.Scan(&o.Type, &o.ID, &o.Titre, &o.Description, &o.DateDebut, &o.DateFin, &o.Lieu,
			&o.Site, &o.NbPlaces, &o.Tarif, &o.Createur, &o.CreatedAt); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		offres = append(offres, o)
	}

	envoyerJSON(w, http.StatusOK, offres)
}

func (s *Serveur) gererDecisionOffre(w http.ResponseWriter, r *http.Request) {
	validateurID, peutValider, ok := s.exigerSalarieOuAdmin(w, r)
	if !ok {
		return
	}
	if !peutValider {
		envoyerErreur(w, http.StatusForbidden, "seul un responsable ou un administrateur peut valider une offre")
		return
	}

	offre, typeConnu := tablesOffres[r.PathValue("type")]
	if !typeConnu {
		envoyerErreur(w, http.StatusBadRequest, "type invalide (formation, atelier ou evenement attendu)")
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "identifiant invalide")
		return
	}

	var entree DecisionOffreEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	entree.Motif = strings.TrimSpace(entree.Motif)

	var nouveauStatut string
	var motif *string
	switch entree.Decision {
	case "valider":
		nouveauStatut = "publie"
	case "refuser":
		if entree.Motif == "" {
			envoyerErreur(w, http.StatusBadRequest, "motif obligatoire pour un refus")
			return
		}
		nouveauStatut = "refuse"
		motif = &entree.Motif
	default:
		envoyerErreur(w, http.StatusBadRequest, "décision invalide (valider ou refuser attendu)")
		return
	}

	var statutActuel string
	err = s.db.QueryRow("SELECT statut FROM "+offre.table+" WHERE id = $1", id).Scan(&statutActuel)
	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "offre introuvable")
		return
	} else if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	resultat, err := s.db.Exec(`
		UPDATE `+offre.table+` SET statut = $1, `+offre.validateur+` = $2, motif_refus = $3, updated_at = now()
		WHERE id = $4 AND statut = 'en_attente'`,
		nouveauStatut, validateurID, motif, id,
	)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if lignes, _ := resultat.RowsAffected(); lignes == 0 {
		envoyerErreur(w, http.StatusConflict, "cette offre n'est pas en attente de validation")
		return
	}

	envoyerJSON(w, http.StatusOK, map[string]string{"statut": nouveauStatut})
}
