package api

import (
	"database/sql"
	"net/http"
	"time"
)

// Page d'accueil publique (front/index.php) : quelques chiffres, les dernières annonces en ligne
// et les prochains rendez-vous publiés. Route ouverte à tous : aucune donnée personnelle
// (ni nom, ni e-mail, ni identifiant d'utilisateur) n'est renvoyée.

const (
	nombreAnnoncesAccueil   = 6
	nombreRendezVousAccueil = 3
)

type ChiffresAccueil struct {
	Membres          int `json:"membres"`
	ObjetsRecuperes  int `json:"objets_recuperes"`
	AnnoncesEnLigne  int `json:"annonces_en_ligne"`
	RendezVousAVenir int `json:"rendez_vous_a_venir"`
}

type AnnonceAccueil struct {
	ID               int       `json:"id"`
	Titre            string    `json:"titre"`
	Type             string    `json:"type"`
	Prix             *float64  `json:"prix"`
	CategorieCode    *string   `json:"categorie_code"`
	CategorieLibelle *string   `json:"categorie_libelle"`
	Localisation     *string   `json:"localisation"`
	DateCreation     time.Time `json:"date_creation"`
}

type RendezVousAccueil struct {
	Type      string    `json:"type"`
	ID        int       `json:"id"`
	Titre     string    `json:"titre"`
	DateDebut time.Time `json:"date_debut"`
	Lieu      *string   `json:"lieu"`
}

type Accueil struct {
	Chiffres   ChiffresAccueil     `json:"chiffres"`
	Annonces   []AnnonceAccueil    `json:"annonces"`
	RendezVous []RendezVousAccueil `json:"rendez_vous"`
}

func (s *Serveur) gererAccueil(w http.ResponseWriter, r *http.Request) {
	accueil := Accueil{Annonces: make([]AnnonceAccueil, 0), RendezVous: make([]RendezVousAccueil, 0)}

	// 1. Chiffres : membres actifs (particuliers et professionnels), objets récupérés en conteneur,
	//    annonces en ligne et rendez-vous publiés à venir
	err := s.db.QueryRow(`
		SELECT
			(SELECT COUNT(*) FROM utilisateurs u JOIN roles r ON r.id = u.role_id
				WHERE r.code IN ('particulier', 'professionnel') AND u.statut = 'actif'),
			(SELECT COUNT(*) FROM depots WHERE statut = 'recuperee'),
			(SELECT COUNT(*) FROM annonces WHERE statut = 'en_ligne'),
			(SELECT COUNT(*) FROM (
				SELECT date_debut FROM formations WHERE statut = 'publie'
				UNION ALL SELECT date_debut FROM ateliers WHERE statut = 'publie'
				UNION ALL SELECT date_debut FROM evenements WHERE statut = 'publie'
			) offres WHERE date_debut > now())`,
	).Scan(&accueil.Chiffres.Membres, &accueil.Chiffres.ObjetsRecuperes,
		&accueil.Chiffres.AnnoncesEnLigne, &accueil.Chiffres.RendezVousAVenir)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	// 2. Dernières annonces en ligne
	lignes, err := s.db.Query(`
		SELECT a.id, a.titre, a.type, a.prix, c.code, c.libelle, a.localisation, a.date_creation
		FROM annonces a
		LEFT JOIN categories_materiaux c ON c.id = a.categorie_id
		WHERE a.statut = 'en_ligne'
		ORDER BY a.date_creation DESC
		LIMIT $1`, nombreAnnoncesAccueil)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()
	for lignes.Next() {
		var a AnnonceAccueil
		var prix sql.NullFloat64
		var code, libelle, localisation sql.NullString
		if err := lignes.Scan(&a.ID, &a.Titre, &a.Type, &prix, &code, &libelle, &localisation, &a.DateCreation); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		if prix.Valid {
			a.Prix = &prix.Float64
		}
		if code.Valid {
			a.CategorieCode = &code.String
		}
		if libelle.Valid {
			a.CategorieLibelle = &libelle.String
		}
		if localisation.Valid && localisation.String != "" {
			a.Localisation = &localisation.String
		}
		accueil.Annonces = append(accueil.Annonces, a)
	}

	// 3. Prochains rendez-vous publiés (formations, ateliers, événements), du plus proche au plus lointain
	lignesRendezVous, err := s.db.Query(`
		SELECT type, id, titre, date_debut, lieu FROM (
			SELECT 'formation' AS type, id, titre, date_debut, lieu FROM formations WHERE statut = 'publie'
			UNION ALL SELECT 'atelier', id, titre, date_debut, lieu FROM ateliers WHERE statut = 'publie'
			UNION ALL SELECT 'evenement', id, titre, date_debut, lieu FROM evenements WHERE statut = 'publie'
		) offres
		WHERE date_debut > now()
		ORDER BY date_debut ASC
		LIMIT $1`, nombreRendezVousAccueil)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignesRendezVous.Close()
	for lignesRendezVous.Next() {
		var rdv RendezVousAccueil
		var lieu sql.NullString
		if err := lignesRendezVous.Scan(&rdv.Type, &rdv.ID, &rdv.Titre, &rdv.DateDebut, &lieu); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		if lieu.Valid && lieu.String != "" {
			rdv.Lieu = &lieu.String
		}
		accueil.RendezVous = append(accueil.RendezVous, rdv)
	}

	envoyerJSON(w, http.StatusOK, accueil)
}
