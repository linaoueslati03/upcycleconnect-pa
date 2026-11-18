package api

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Back-office : accès aux PDF archivés (table documents, alimentée par documents_pdf.go)
// et aux codes de dépôt (code d'ouverture du conteneur et code-barre de récupération).
// Les PDF eux-mêmes se téléchargent avec GET /api/documents/{id}/pdf, qui contrôle déjà les droits.

const elementsParPageDocuments = 20

var typesDocument = []string{
	"fiche_atelier", "fiche_evenement", "fiche_formation",
	"attestation_inscription", "contrat", "document_administratif", "bon_depot", "facture",
}

var statutsDepotCodes = []string{"validee", "deposee", "recuperee"}

type PersonneLiee struct {
	ID  int    `json:"id"`
	Nom string `json:"nom"`
}

type DocumentAdmin struct {
	ID           int           `json:"id"`
	TypeDocument string        `json:"type_document"`
	Titre        string        `json:"titre"`
	Beneficiaire *PersonneLiee `json:"beneficiaire"`
	Auteur       *PersonneLiee `json:"auteur"`
	TailleOctets int           `json:"taille_octets"`
	DateCreation time.Time     `json:"date_creation"`
}

type PageDocumentsAdmin struct {
	Documents []DocumentAdmin `json:"documents"`
	ParType   []Comptage      `json:"par_type"`
	Total     int             `json:"total"`
	Page      int             `json:"page"`
	ParPage   int             `json:"par_page"`
}

type CodeDepot struct {
	DepotID          int           `json:"depot_id"`
	DescriptionObjet string        `json:"description_objet"`
	Statut           string        `json:"statut"`
	CodeOuverture    *string       `json:"code_ouverture"`
	CodeBarre        *string       `json:"code_barre"`
	Deposant         PersonneLiee  `json:"deposant"`
	Recuperateur     *PersonneLiee `json:"recuperateur"`
	Site             *string       `json:"site"`
	DateValidation   *time.Time    `json:"date_validation"`
	DateRecuperation *time.Time    `json:"date_recuperation"`
}

type PageCodesDepot struct {
	Codes   []CodeDepot `json:"codes"`
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	ParPage int         `json:"par_page"`
}

func pageDepuisRequete(r *http.Request) int {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func motifRecherche(r *http.Request) string {
	recherche := strings.TrimSpace(r.URL.Query().Get("recherche"))
	if recherche == "" {
		return ""
	}
	return "%" + echappementRecherche.Replace(strings.ToLower(recherche)) + "%"
}

func personne(id *int, prenom, nom *string) *PersonneLiee {
	if id == nil {
		return nil
	}
	return &PersonneLiee{ID: *id, Nom: strings.TrimSpace(*prenom + " " + *nom)}
}

func (s *Serveur) gererListeDocumentsAdmin(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	conditions := []string{"TRUE"}
	args := []any{}
	if motif := motifRecherche(r); motif != "" {
		args = append(args, motif)
		conditions = append(conditions, fmt.Sprintf(
			"lower(d.titre || ' ' || COALESCE(b.prenom || ' ' || b.nom || ' ' || b.email, '')) LIKE $%d", len(args)))
	}
	// Le comptage par type ignore le filtre de type, pour que les onglets gardent leurs nombres
	conditionsSansType := strings.Join(conditions, " AND ")
	argsSansType := slices.Clone(args)
	if typeDocument := r.URL.Query().Get("type"); typeDocument != "" {
		if !slices.Contains(typesDocument, typeDocument) {
			envoyerErreur(w, http.StatusBadRequest, "type de document invalide")
			return
		}
		args = append(args, typeDocument)
		conditions = append(conditions, fmt.Sprintf("d.type_document = $%d", len(args)))
	}
	where := strings.Join(conditions, " AND ")
	jointures := ` FROM documents d LEFT JOIN utilisateurs b ON b.id = d.utilisateur_id `

	page := pageDepuisRequete(r)
	resultat := PageDocumentsAdmin{Documents: make([]DocumentAdmin, 0), Page: page, ParPage: elementsParPageDocuments}

	var err error
	resultat.ParType, err = s.compterParLibelle(typesDocument,
		`SELECT d.type_document, count(*)`+jointures+`WHERE `+conditionsSansType+` GROUP BY 1`, argsSansType...)
	if err == nil {
		err = s.db.QueryRow(`SELECT count(*)`+jointures+`WHERE `+where, args...).Scan(&resultat.Total)
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	args = append(args, elementsParPageDocuments, (page-1)*elementsParPageDocuments)
	lignes, err := s.db.Query(fmt.Sprintf(`
		SELECT d.id, d.type_document, d.titre, octet_length(d.contenu_pdf), d.date_creation,
		       b.id, b.prenom, b.nom, a.id, a.prenom, a.nom
		%s LEFT JOIN utilisateurs a ON a.id = d.cree_par
		WHERE %s
		ORDER BY d.date_creation DESC, d.id DESC
		LIMIT $%d OFFSET $%d`, jointures, where, len(args)-1, len(args)), args...)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	for lignes.Next() {
		var d DocumentAdmin
		var beneficiaireID, auteurID *int
		var beneficiairePrenom, beneficiaireNom, auteurPrenom, auteurNom *string
		if err := lignes.Scan(&d.ID, &d.TypeDocument, &d.Titre, &d.TailleOctets, &d.DateCreation,
			&beneficiaireID, &beneficiairePrenom, &beneficiaireNom, &auteurID, &auteurPrenom, &auteurNom); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		d.Beneficiaire = personne(beneficiaireID, beneficiairePrenom, beneficiaireNom)
		d.Auteur = personne(auteurID, auteurPrenom, auteurNom)
		resultat.Documents = append(resultat.Documents, d)
	}
	if err := lignes.Err(); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, resultat)
}

func (s *Serveur) gererListeCodesDepot(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.exigerRole(w, r, roleAdministrateur); !ok {
		return
	}

	conditions := []string{"(dp.code_ouverture IS NOT NULL OR dp.code_barre IS NOT NULL)"}
	args := []any{}
	if motif := motifRecherche(r); motif != "" {
		args = append(args, motif)
		conditions = append(conditions, fmt.Sprintf(`lower(COALESCE(dp.code_ouverture, '') || ' ' || COALESCE(dp.code_barre, '')
			|| ' ' || dp.description_objet || ' ' || u.prenom || ' ' || u.nom || ' ' || u.email
			|| ' ' || COALESCE(c.site, '')) LIKE $%d`, len(args)))
	}
	if statut := r.URL.Query().Get("statut"); statut != "" {
		if !slices.Contains(statutsDepotCodes, statut) {
			envoyerErreur(w, http.StatusBadRequest, "statut invalide")
			return
		}
		args = append(args, statut)
		conditions = append(conditions, fmt.Sprintf("dp.statut = $%d", len(args)))
	}
	where := strings.Join(conditions, " AND ")
	jointures := ` FROM depots dp
		JOIN utilisateurs u ON u.id = dp.utilisateur_id
		LEFT JOIN conteneurs c ON c.id = dp.conteneur_id
		LEFT JOIN utilisateurs p ON p.id = dp.professionnel_recuperateur_id `

	page := pageDepuisRequete(r)
	resultat := PageCodesDepot{Codes: make([]CodeDepot, 0), Page: page, ParPage: elementsParPageDocuments}
	if err := s.db.QueryRow(`SELECT count(*)`+jointures+`WHERE `+where, args...).Scan(&resultat.Total); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	args = append(args, elementsParPageDocuments, (page-1)*elementsParPageDocuments)
	lignes, err := s.db.Query(fmt.Sprintf(`
		SELECT dp.id, dp.description_objet, dp.statut, dp.code_ouverture, dp.code_barre,
		       u.id, u.prenom, u.nom, p.id, p.prenom, p.nom, c.site, dp.date_validation, dp.date_recuperation
		%s
		WHERE %s
		ORDER BY COALESCE(dp.date_validation, dp.date_demande) DESC, dp.id DESC
		LIMIT $%d OFFSET $%d`, jointures, where, len(args)-1, len(args)), args...)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	for lignes.Next() {
		var c CodeDepot
		var recuperateurID *int
		var deposantPrenom, deposantNom string
		var recuperateurPrenom, recuperateurNom *string
		if err := lignes.Scan(&c.DepotID, &c.DescriptionObjet, &c.Statut, &c.CodeOuverture, &c.CodeBarre,
			&c.Deposant.ID, &deposantPrenom, &deposantNom, &recuperateurID, &recuperateurPrenom, &recuperateurNom,
			&c.Site, &c.DateValidation, &c.DateRecuperation); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		c.Deposant.Nom = strings.TrimSpace(deposantPrenom + " " + deposantNom)
		c.Recuperateur = personne(recuperateurID, recuperateurPrenom, recuperateurNom)
		resultat.Codes = append(resultat.Codes, c)
	}
	if err := lignes.Err(); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, resultat)
}
