package api

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"encoding/json"

	"github.com/go-pdf/fpdf"
)

//Génération et archivage communs

type DocumentArchive struct {
	ID             int       `json:"id"`
	TypeDocument   string    `json:"type_document"`
	Titre          string    `json:"titre"`
	UtilisateurID  *int      `json:"utilisateur_id,omitempty"`
	TypeOffre      *string   `json:"type_offre,omitempty"`
	OffreID        *int      `json:"offre_id,omitempty"`
	InscriptionID  *int      `json:"inscription_id,omitempty"`
	DateCreation   time.Time `json:"date_creation"`
}

func genererPDF(titre string, lignes []string) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	texte := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetTitle(titre, true)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 18)
	pdf.MultiCell(0, 10, texte(titre), "", "C", false)
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)

	for _, ligne := range lignes {
		pdf.MultiCell(0, 8, texte(ligne), "", "L", false)
		pdf.Ln(2)
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 9)
	pdf.CellFormat(
		0, 8,
		texte("Document généré automatiquement par UpcycleConnect"),
		"", 1, "C", false, 0, "",
	)

	var contenu bytes.Buffer
	if err := pdf.Output(&contenu); err != nil {
		return nil, err
	}

	return contenu.Bytes(), nil
}

func (s *Serveur) enregistrerDocument(
	titre string,
	typeDocument string,
	utilisateurID *int,
	typeOffre *string,
	offreID *int,
	inscriptionID *int,
	creePar int,
	contenu []byte,
) (int, error) {
	var id int

	err := s.db.QueryRow(`
		INSERT INTO documents (
			type_document, titre, utilisateur_id,
			type_offre, offre_id, inscription_id,
			cree_par, contenu_pdf
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`,
		typeDocument,
		titre,
		utilisateurID,
		typeOffre,
		offreID,
		inscriptionID,
		creePar,
		contenu,
	).Scan(&id)

	return id, err
}

func envoyerPDF(w http.ResponseWriter, contenu []byte, nomFichier string) {
	nomFichier = strings.ReplaceAll(nomFichier, `"`, "")
	nomFichier = strings.ReplaceAll(nomFichier, "\r", "")
	nomFichier = strings.ReplaceAll(nomFichier, "\n", "")

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, nomFichier),
	)
	w.Header().Set("Content-Length", fmt.Sprint(len(contenu)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(contenu)
}

// Générer une fiche d’atelier, d’événement ou de formation

func (s *Serveur) gererCreationFicheOffrePDF(
	w http.ResponseWriter,
	r *http.Request,
) {
	typeOffre := r.PathValue("type")

	tables := map[string]string{
		"atelier":   "ateliers",
		"evenement": "evenements",
		"formation": "formations",
	}

	table, ok := tables[typeOffre]
	if !ok {
		envoyerErreur(w, http.StatusBadRequest, "type d'offre invalide")
		return
	}

	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	utilisateurID, role := s.roleConnecte(r)
	if utilisateurID == 0 {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	var titre, description, statut string
	var dateDebut time.Time
	var dateFin sql.NullTime
	var lieu sql.NullString

	err = s.db.QueryRow(
		"SELECT titre, description, date_debut, date_fin, lieu, statut FROM "+table+" WHERE id = $1",
		id,
	).Scan(&titre, &description, &dateDebut, &dateFin, &lieu, &statut)

	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "offre introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	// Une offre non publiée est réservée au personnel.
	if statut != "publie" &&
		role != roleSalarie &&
		role != roleAdministrateur {
		envoyerErreur(w, http.StatusNotFound, "offre introuvable")
		return
	}

	libelleType := map[string]string{
		"atelier":   "Atelier",
		"evenement": "Événement",
		"formation": "Formation",
	}[typeOffre]

	lignes := []string{
		"Type : " + libelleType,
		"Titre : " + titre,
		"Description : " + description,
		"Date de début : " + dateDebut.Format("02/01/2006 15:04"),
		"Statut : " + statut,
	}

	if dateFin.Valid {
		lignes = append(
			lignes,
			"Date de fin : "+dateFin.Time.Format("02/01/2006 15:04"),
		)
	}
	if lieu.Valid {
		lignes = append(lignes, "Lieu : "+lieu.String)
	}

	contenu, err := genererPDF("Fiche "+libelleType, lignes)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur de génération du PDF")
		return
	}

	typeDocument := "fiche_" + typeOffre
	var typeOffrePtr = &typeOffre
	var offreID = &id

	documentID, err := s.enregistrerDocument(
		"Fiche "+libelleType+" : "+titre,
		typeDocument,
		nil,
		typeOffrePtr,
		offreID,
		nil,
		utilisateurID,
		contenu,
	)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur d'archivage du PDF")
		return
	}

	// La réponse indique l'identifiant du PDF archivé.
	envoyerJSON(w, http.StatusCreated, map[string]any{
		"id":       documentID,
		"titre":    titre,
		"pdf":      "/api/documents/" + fmt.Sprint(documentID) + "/pdf",
		"message":  "PDF généré et archivé",
	})
}

//Générer une attestation d’inscription

func (s *Serveur) gererCreationAttestationPDF(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	utilisateurID, role := s.roleConnecte(r)
	if utilisateurID == 0 {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	var inscription Inscription
	err = s.db.QueryRow(`
		SELECT id, utilisateur_id, type_offre, offre_id,
		       date_inscription, statut_paiement, montant_paye
		FROM inscriptions
		WHERE id = $1
	`, id).Scan(
		&inscription.ID,
		&inscription.UtilisateurID,
		&inscription.TypeOffre,
		&inscription.OffreID,
		&inscription.DateInscription,
		&inscription.StatutPaiement,
		&inscription.MontantPaye,
	)

	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "inscription introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	// Un particulier ne peut générer que sa propre attestation.
	if role == roleParticulier && inscription.UtilisateurID != utilisateurID {
		envoyerErreur(w, http.StatusForbidden, "accès refusé")
		return
	}
	if role != roleParticulier &&
		role != roleSalarie &&
		role != roleAdministrateur {
		envoyerErreur(w, http.StatusForbidden, "accès refusé")
		return
	}

	var nom, prenom string
	err = s.db.QueryRow(`
		SELECT nom, prenom
		FROM utilisateurs
		WHERE id = $1
	`, inscription.UtilisateurID).Scan(&nom, &prenom)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	var titreOffre string
	table, ok := tableParTypeOffre[inscription.TypeOffre]
	if !ok {
		envoyerErreur(w, http.StatusInternalServerError, "type d'offre invalide en base")
		return
	}

	err = s.db.QueryRow(
		"SELECT titre FROM "+table+" WHERE id = $1",
		inscription.OffreID,
	).Scan(&titreOffre)

	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "offre introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	lignes := []string{
		"Attestation d'inscription",
		"Bénéficiaire : " + prenom + " " + nom,
		"Offre : " + titreOffre,
		"Type d'offre : " + inscription.TypeOffre,
		"Date d'inscription : " + inscription.DateInscription.Format("02/01/2006"),
		"Statut du paiement : " + inscription.StatutPaiement,
		fmt.Sprintf("Montant enregistré : %.2f €", inscription.MontantPaye),
	}

	contenu, err := genererPDF("Attestation d'inscription", lignes)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur de génération du PDF")
		return
	}

	documentID, err := s.enregistrerDocument(
		"Attestation d'inscription : "+titreOffre,
		"attestation_inscription",
		&inscription.UtilisateurID,
		&inscription.TypeOffre,
		&inscription.OffreID,
		&inscription.ID,
		utilisateurID,
		contenu,
	)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur d'archivage du PDF")
		return
	}

	envoyerJSON(w, http.StatusCreated, map[string]any{
		"id":      documentID,
		"titre":   "Attestation d'inscription : " + titreOffre,
		"pdf":     "/api/documents/" + fmt.Sprint(documentID) + "/pdf",
		"message": "Attestation générée et archivée",
	})
}

//créer un contrat ou un document administratif

type DocumentAdministratifEntree struct {
	UtilisateurID int    `json:"utilisateur_id"`
	TypeDocument  string `json:"type_document"`
	Titre         string `json:"titre"`
	Contenu       string `json:"contenu"`
}

func (s *Serveur) gererCreationDocumentAdministratif(
	w http.ResponseWriter,
	r *http.Request,
) {
	createurID, ok := s.exigerRole(
		w, r, roleAdministrateur, roleSalarie,
	)
	if !ok {
		return
	}

	var entree DocumentAdministratifEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps JSON invalide")
		return
	}

	if strings.TrimSpace(entree.Titre) == "" ||
		strings.TrimSpace(entree.Contenu) == "" ||
		entree.UtilisateurID <= 0 {
		envoyerErreur(w, http.StatusBadRequest, "bénéficiaire, titre et contenu obligatoires")
		return
	}

	if entree.TypeDocument != "contrat" &&
		entree.TypeDocument != "document_administratif" {
		envoyerErreur(w, http.StatusBadRequest, "type de document invalide")
		return
	}

	var existe bool
	err := s.db.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM utilisateurs WHERE id = $1)",
		entree.UtilisateurID,
	).Scan(&existe)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	if !existe {
		envoyerErreur(w, http.StatusNotFound, "bénéficiaire introuvable")
		return
	}

	contenu, err := genererPDF(
		entree.Titre,
		strings.Split(entree.Contenu, "\n"),
	)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur de génération du PDF")
		return
	}

	documentID, err := s.enregistrerDocument(
		entree.Titre,
		entree.TypeDocument,
		&entree.UtilisateurID,
		nil,
		nil,
		nil,
		createurID,
		contenu,
	)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur d'archivage du PDF")
		return
	}

	envoyerJSON(w, http.StatusCreated, map[string]any{
		"id":      documentID,
		"titre":   entree.Titre,
		"pdf":     "/api/documents/" + fmt.Sprint(documentID) + "/pdf",
		"message": "Document généré et archivé",
	})
}

// lister les documents personnels

func (s *Serveur) gererMesDocuments(
	w http.ResponseWriter,
	r *http.Request,
) {
	utilisateurID, err := s.utilisateurConnecte(r)
	if err != nil {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	lignes, err := s.db.Query(`
		SELECT id, type_document, titre, utilisateur_id,
		       type_offre, offre_id, inscription_id, date_creation
		FROM documents
		WHERE utilisateur_id = $1
		ORDER BY date_creation DESC
	`, utilisateurID)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer lignes.Close()

	documents := make([]DocumentArchive, 0)

	for lignes.Next() {
		var d DocumentArchive
		var utilisateurIDSQL, offreIDSQL, inscriptionIDSQL sql.NullInt64
		var typeOffreSQL sql.NullString

		if err := lignes.Scan(
			&d.ID,
			&d.TypeDocument,
			&d.Titre,
			&utilisateurIDSQL,
			&typeOffreSQL,
			&offreIDSQL,
			&inscriptionIDSQL,
			&d.DateCreation,
		); err != nil {
			envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
			return
		}

		if utilisateurIDSQL.Valid {
			v := int(utilisateurIDSQL.Int64)
			d.UtilisateurID = &v
		}
		if typeOffreSQL.Valid {
			v := typeOffreSQL.String
			d.TypeOffre = &v
		}
		if offreIDSQL.Valid {
			v := int(offreIDSQL.Int64)
			d.OffreID = &v
		}
		if inscriptionIDSQL.Valid {
			v := int(inscriptionIDSQL.Int64)
			d.InscriptionID = &v
		}

		documents = append(documents, d)
	}

	if err := lignes.Err(); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	envoyerJSON(w, http.StatusOK, documents)
}

//consulter ou télécharger un PDF avec contrôle des droits

func (s *Serveur) gererTelechargementDocument(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return
	}

	utilisateurID, role := s.roleConnecte(r)
	if utilisateurID == 0 {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	var d DocumentArchive
	var contenu []byte
	var creePar sql.NullInt64

	err = s.db.QueryRow(`
		SELECT id, type_document, titre, utilisateur_id,
		       type_offre, offre_id, inscription_id, cree_par,
		       date_creation, contenu_pdf
		FROM documents
		WHERE id = $1
	`, id).Scan(
		&d.ID,
		&d.TypeDocument,
		&d.Titre,
		&d.UtilisateurID,
		&d.TypeOffre,
		&d.OffreID,
		&d.InscriptionID,
		&creePar,
		&d.DateCreation,
		&contenu,
	)

	if errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusNotFound, "document introuvable")
		return
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	autorise := role == roleAdministrateur

	if d.UtilisateurID != nil && *d.UtilisateurID == utilisateurID {
		autorise = true
	}
	if role == roleSalarie && creePar.Valid &&
		int(creePar.Int64) == utilisateurID {
		autorise = true
	}

	// Une fiche d'offre publiée peut être consultée par tout utilisateur connecté.
	if !autorise &&
		(d.TypeDocument == "fiche_atelier" ||
			d.TypeDocument == "fiche_evenement" ||
			d.TypeDocument == "fiche_formation") &&
		d.TypeOffre != nil &&
		d.OffreID != nil {
		tables := map[string]string{
			"atelier":   "ateliers",
			"evenement": "evenements",
			"formation": "formations",
		}

		if table, ok := tables[*d.TypeOffre]; ok {
			var statut string
			err := s.db.QueryRow(
				"SELECT statut FROM "+table+" WHERE id = $1",
				*d.OffreID,
			).Scan(&statut)

			if err == nil && statut == "publie" {
				autorise = true
			}
		}
	}

	if !autorise {
		// Ne pas révéler l'existence d'un document privé.
		envoyerErreur(w, http.StatusNotFound, "document introuvable")
		return
	}

	nomFichier := fmt.Sprintf("document_%d.pdf", d.ID)
	envoyerPDF(w, contenu, nomFichier)
}