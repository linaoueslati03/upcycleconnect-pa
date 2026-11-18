package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Génération automatique des documents PDF (« générés automatiquement… un double sera conservé »).
// Les PDF sont construits ici, puis archivés dans la table documents :
//   - la fiche d'une prestation quand elle est publiée ;
//   - l'attestation d'une inscription quand un particulier s'inscrit ;
//   - le bon de dépôt (code d'ouverture + code-barre) quand un dépôt est validé.
// La génération se fait après l'enregistrement de l'action : si le PDF échoue, l'action reste faite
// et l'erreur est notée dans les journaux.

var libellesTypeOffre = map[string]string{
	"atelier":   "Atelier",
	"evenement": "Événement",
	"formation": "Formation",
}

// Fiche d'une prestation : renvoie aussi son titre et son statut (la route vérifie les droits avec).
func (s *Serveur) construireFicheOffre(typeOffre string, id int) (titre, statut string, contenu []byte, err error) {
	table, ok := tableParTypeOffre[typeOffre]
	if !ok {
		return "", "", nil, fmt.Errorf("type d'offre invalide : %s", typeOffre)
	}

	var description string
	var dateDebut time.Time
	var dateFin sql.NullTime
	var lieu sql.NullString
	err = s.db.QueryRow(
		"SELECT titre, description, date_debut, date_fin, lieu, statut FROM "+table+" WHERE id = $1", id,
	).Scan(&titre, &description, &dateDebut, &dateFin, &lieu, &statut)
	if err != nil {
		return "", "", nil, err
	}

	libelleType := libellesTypeOffre[typeOffre]
	lignes := []string{
		"Type : " + libelleType,
		"Titre : " + titre,
		"Description : " + description,
		"Date de début : " + dateDebut.Format("02/01/2006 15:04"),
		"Statut : " + statut,
	}
	if dateFin.Valid {
		lignes = append(lignes, "Date de fin : "+dateFin.Time.Format("02/01/2006 15:04"))
	}
	if lieu.Valid {
		lignes = append(lignes, "Lieu : "+lieu.String)
	}

	contenu, err = genererPDF("Fiche "+libelleType, lignes)
	return titre, statut, contenu, err
}

func (s *Serveur) genererFicheOffre(typeOffre string, id, auteurID int) (int, error) {
	titre, _, contenu, err := s.construireFicheOffre(typeOffre, id)
	if err != nil {
		return 0, err
	}
	return s.enregistrerDocument("Fiche "+libellesTypeOffre[typeOffre]+" : "+titre, "fiche_"+typeOffre,
		nil, &typeOffre, &id, nil, auteurID, contenu)
}

// Attestation d'inscription : renvoie aussi l'inscription (la route vérifie le propriétaire avec).
func (s *Serveur) construireAttestation(inscriptionID int) (inscription Inscription, titreOffre string, contenu []byte, err error) {
	err = s.db.QueryRow(`
		SELECT id, utilisateur_id, type_offre, offre_id, date_inscription, statut_paiement, montant_paye
		FROM inscriptions WHERE id = $1`, inscriptionID,
	).Scan(&inscription.ID, &inscription.UtilisateurID, &inscription.TypeOffre, &inscription.OffreID,
		&inscription.DateInscription, &inscription.StatutPaiement, &inscription.MontantPaye)
	if err != nil {
		return inscription, "", nil, err
	}

	var nom, prenom string
	if err = s.db.QueryRow("SELECT nom, prenom FROM utilisateurs WHERE id = $1", inscription.UtilisateurID).
		Scan(&nom, &prenom); err != nil {
		return inscription, "", nil, err
	}

	table, ok := tableParTypeOffre[inscription.TypeOffre]
	if !ok {
		return inscription, "", nil, fmt.Errorf("type d'offre invalide en base : %s", inscription.TypeOffre)
	}
	if err = s.db.QueryRow("SELECT titre FROM "+table+" WHERE id = $1", inscription.OffreID).Scan(&titreOffre); err != nil {
		return inscription, "", nil, err
	}

	contenu, err = genererPDF("Attestation d'inscription", []string{
		"Attestation d'inscription",
		"Bénéficiaire : " + prenom + " " + nom,
		"Offre : " + titreOffre,
		"Type d'offre : " + inscription.TypeOffre,
		"Date d'inscription : " + inscription.DateInscription.Format("02/01/2006"),
		"Statut du paiement : " + inscription.StatutPaiement,
		fmt.Sprintf("Montant enregistré : %.2f €", inscription.MontantPaye),
	})
	return inscription, titreOffre, contenu, err
}

// Une seule attestation automatique par inscription : si elle existe déjà, on renvoie son numéro.
func (s *Serveur) genererAttestation(inscriptionID, auteurID int) (int, error) {
	var existant int
	err := s.db.QueryRow(`SELECT id FROM documents
		WHERE type_document = 'attestation_inscription' AND inscription_id = $1 LIMIT 1`, inscriptionID).Scan(&existant)
	if err == nil {
		return existant, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	inscription, titreOffre, contenu, err := s.construireAttestation(inscriptionID)
	if err != nil {
		return 0, err
	}
	return s.enregistrerDocument("Attestation d'inscription : "+titreOffre, "attestation_inscription",
		&inscription.UtilisateurID, &inscription.TypeOffre, &inscription.OffreID, &inscription.ID, auteurID, contenu)
}

// Bon de dépôt : les codes envoyés au particulier après validation de sa demande (un seul par dépôt).
func (s *Serveur) genererBonDepot(depotID, auteurID int) (int, error) {
	var existant int
	err := s.db.QueryRow("SELECT id FROM documents WHERE type_document = 'bon_depot' AND depot_id = $1 LIMIT 1", depotID).
		Scan(&existant)
	if err == nil {
		return existant, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	var deposantID int
	var prenom, nom, description string
	var codeOuverture, codeBarre, site sql.NullString
	var conteneurID sql.NullInt64
	var dateValidation sql.NullTime
	err = s.db.QueryRow(`
		SELECT d.utilisateur_id, u.prenom, u.nom, d.description_objet, d.code_ouverture, d.code_barre,
		       d.conteneur_id, c.site, d.date_validation
		FROM depots d
		JOIN utilisateurs u ON u.id = d.utilisateur_id
		LEFT JOIN conteneurs c ON c.id = d.conteneur_id
		WHERE d.id = $1`, depotID,
	).Scan(&deposantID, &prenom, &nom, &description, &codeOuverture, &codeBarre, &conteneurID, &site, &dateValidation)
	if err != nil {
		return 0, err
	}
	if !codeOuverture.Valid || !codeBarre.Valid {
		return 0, fmt.Errorf("le dépôt %d n'a pas encore de codes", depotID)
	}

	lignes := []string{
		"Déposant : " + prenom + " " + nom,
		"Objet : " + description,
	}
	if conteneurID.Valid {
		lignes = append(lignes, fmt.Sprintf("Conteneur : n°%d - %s", conteneurID.Int64, site.String))
	}
	if dateValidation.Valid {
		lignes = append(lignes, "Demande validée le : "+dateValidation.Time.Format("02/01/2006"))
	}
	lignes = append(lignes,
		"Code d'ouverture du conteneur : "+codeOuverture.String,
		"Code-barre de récupération : "+codeBarre.String,
		"Saisissez le code d'ouverture sur le conteneur pour y déposer l'objet.",
		"Le code-barre permet au professionnel qui récupère l'objet de le retirer du conteneur.",
	)
	contenu, err := genererPDF("Bon de dépôt", lignes)
	if err != nil {
		return 0, err
	}

	var id int
	err = s.db.QueryRow(`
		INSERT INTO documents (type_document, titre, utilisateur_id, depot_id, cree_par, contenu_pdf)
		VALUES ('bon_depot', $1, $2, $3, $4, $5)
		RETURNING id`,
		"Bon de dépôt : "+description, deposantID, depotID, auteurID, contenu,
	).Scan(&id)
	return id, err
}

// Lance une génération automatique sans jamais bloquer l'action qui l'a déclenchée.
func archiverAutomatiquement(description string, generer func() (int, error)) {
	if _, err := generer(); err != nil {
		log.Printf("archivage automatique impossible (%s) : %v", description, err)
	}
}

func (s *Serveur) archiverFicheOffre(typeOffre string, id, auteurID int) {
	archiverAutomatiquement(fmt.Sprintf("fiche %s %d", typeOffre, id), func() (int, error) {
		return s.genererFicheOffre(typeOffre, id, auteurID)
	})
}

// Rattrapage : génère une fois les PDF manquants pour ce qui existait avant la génération automatique.
type BilanRattrapage struct {
	Fiches       int `json:"fiches"`
	Attestations int `json:"attestations"`
	BonsDepot    int `json:"bons_depot"`
	Erreurs      int `json:"erreurs"`
}

func (s *Serveur) idsManquants(requete string) ([][2]any, error) {
	lignes, err := s.db.Query(requete)
	if err != nil {
		return nil, err
	}
	defer lignes.Close()
	var resultat [][2]any
	for lignes.Next() {
		var a, b any
		if err := lignes.Scan(&a, &b); err != nil {
			return nil, err
		}
		resultat = append(resultat, [2]any{a, b})
	}
	return resultat, lignes.Err()
}

func (s *Serveur) gererRattrapageDocuments(w http.ResponseWriter, r *http.Request) {
	adminID, ok := s.exigerRole(w, r, roleAdministrateur)
	if !ok {
		return
	}
	var bilan BilanRattrapage

	fiches, err := s.idsManquants(`
		SELECT o.type, o.id FROM (
			SELECT 'formation' AS type, id FROM formations WHERE statut IN ('publie', 'termine')
			UNION ALL SELECT 'atelier', id FROM ateliers WHERE statut IN ('publie', 'termine')
			UNION ALL SELECT 'evenement', id FROM evenements WHERE statut IN ('publie', 'termine')
		) o
		WHERE NOT EXISTS (SELECT 1 FROM documents d
		                  WHERE d.type_document = 'fiche_' || o.type AND d.type_offre = o.type AND d.offre_id = o.id)
		ORDER BY o.type, o.id`)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	for _, f := range fiches {
		if _, err := s.genererFicheOffre(f[0].(string), int(f[1].(int64)), adminID); err != nil {
			log.Printf("rattrapage fiche %v %v : %v", f[0], f[1], err)
			bilan.Erreurs++
		} else {
			bilan.Fiches++
		}
	}

	attestations, err := s.idsManquants(`
		SELECT i.id, i.utilisateur_id FROM inscriptions i
		WHERE NOT EXISTS (SELECT 1 FROM documents d
		                  WHERE d.type_document = 'attestation_inscription' AND d.inscription_id = i.id)
		ORDER BY i.id`)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	for _, a := range attestations {
		// L'auteur d'une attestation est la personne inscrite, comme pour une inscription faite aujourd'hui
		if _, err := s.genererAttestation(int(a[0].(int64)), int(a[1].(int64))); err != nil {
			log.Printf("rattrapage attestation %v : %v", a[0], err)
			bilan.Erreurs++
		} else {
			bilan.Attestations++
		}
	}

	bons, err := s.idsManquants(`
		SELECT dp.id, dp.utilisateur_id FROM depots dp
		WHERE dp.code_ouverture IS NOT NULL AND dp.code_barre IS NOT NULL
		  AND NOT EXISTS (SELECT 1 FROM documents d WHERE d.type_document = 'bon_depot' AND d.depot_id = dp.id)
		ORDER BY dp.id`)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	for _, b := range bons {
		if _, err := s.genererBonDepot(int(b[0].(int64)), adminID); err != nil {
			log.Printf("rattrapage bon de dépôt %v : %v", b[0], err)
			bilan.Erreurs++
		} else {
			bilan.BonsDepot++
		}
	}

	envoyerJSON(w, http.StatusOK, bilan)
}
