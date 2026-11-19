package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Achats : payer une formation, acheter un objet mis en vente.
//
// Parcours :
//  1. POST /api/commandes crée une commande « en attente » (l'objet est réservé, la place de formation est
//     retenue pendant 30 minutes) et renvoie l'adresse de la page de paiement ;
//  2. le paiement est confirmé → confirmerCommande() : inscription payée ou objet cédé, paiement enregistré
//     avec enregistrerPaiement(), facture PDF archivée, notifications ;
//  3. ou il est abandonné → annulerCommande() : l'objet redevient disponible.
//
// Tant que Stripe n'est pas branché, le paiement est SIMULÉ (variable PAIEMENT_FOURNISSEUR absente ou
// « simulation ») : la page de paiement est /espace/paiement, avec un bouton « Payer » de démonstration.
// Pour brancher Stripe (voir README, « Brancher Stripe ») :
//   - urlPaiement() crée une session Stripe Checkout et renvoie son adresse (reference_externe = id de session) ;
//   - un webhook « checkout.session.completed » appelle confirmerCommande(id, id_du_paiement_stripe) ;
//   - l'annulation (cancel_url) et l'expiration de la session appellent annulerCommande().
// Rien d'autre n'est à changer.

const dureeReservation = 30 * time.Minute

var designationsCommande = map[string]string{
	"formation":  "Formation",
	"vente":      "Achat de l'objet",
	"abonnement": "Abonnement",
	"publicite":  "Campagne publicitaire (1 mois) pour",
}

// Préfixe chaque colonne d'une liste « a, b, c » (pour les requêtes avec jointure).
func prefixer(colonnes, prefixe string) string {
	morceaux := strings.Split(colonnes, ",")
	for i, m := range morceaux {
		morceaux[i] = prefixe + strings.TrimSpace(m)
	}
	return strings.Join(morceaux, ", ")
}

var (
	errCommandeIntrouvable = errors.New("commande introuvable")
	errCommandeTerminee    = errors.New("cette commande n'est plus en attente de paiement")
	errCommandeExpiree     = errors.New("le délai de paiement est dépassé, la réservation a été annulée")
)

type Commande struct {
	ID             int        `json:"id"`
	UtilisateurID  int        `json:"utilisateur_id"`
	Type           string     `json:"type"`
	ReferenceID    int        `json:"reference_id"`
	Libelle        string     `json:"libelle"`
	Montant        float64    `json:"montant"`
	Statut         string     `json:"statut"`
	Fournisseur    string     `json:"fournisseur"`
	DateCreation   time.Time  `json:"date_creation"`
	DateExpiration time.Time  `json:"date_expiration"`
	DatePaiement   *time.Time `json:"date_paiement"`
	FactureID      *int       `json:"facture_id"`
}

type CommandeEntree struct {
	Type        string `json:"type"`
	ReferenceID int    `json:"reference_id"`
	// Campagne publicitaire : élément mis en avant (projet ou annonce) et budget mensuel
	CibleType string  `json:"cible_type"`
	CibleID   int     `json:"cible_id"`
	Montant   float64 `json:"montant"`
}

const colonnesCommande = `id, utilisateur_id, type, reference_id, libelle, montant, statut, fournisseur,
	date_creation, date_expiration, date_paiement`

func scannerCommande(ligne interface{ Scan(...any) error }, c *Commande) error {
	return ligne.Scan(&c.ID, &c.UtilisateurID, &c.Type, &c.ReferenceID, &c.Libelle, &c.Montant, &c.Statut,
		&c.Fournisseur, &c.DateCreation, &c.DateExpiration, &c.DatePaiement)
}

func fournisseurPaiement() string {
	if f := os.Getenv("PAIEMENT_FOURNISSEUR"); f != "" {
		return f
	}
	return "simulation"
}

// Adresse où le navigateur est envoyé pour payer.
func urlPaiement(c Commande) (string, error) {
	switch c.Fournisseur {
	case "simulation":
		return fmt.Sprintf("/espace/paiement?commande=%d", c.ID), nil
	default:
		return "", fmt.Errorf("le paiement par %s n'est pas encore branché", c.Fournisseur)
	}
}

// Libère les réservations dont le délai de paiement est dépassé (appelé avant de lister ou de commander).
func libererCommandesExpirees(ex executeur) error {
	if _, err := ex.Exec(`
		UPDATE annonces a SET statut = 'en_ligne'
		FROM commandes c
		WHERE c.type = 'vente' AND c.reference_id = a.id AND c.statut = 'en_attente'
		  AND c.date_expiration < now() AND a.statut = 'reservee' AND a.acheteur_id IS NULL`); err != nil {
		return err
	}
	if _, err := ex.Exec(`
		UPDATE campagnes_publicitaires cp SET statut = 'annulee'
		FROM commandes c
		WHERE c.type = 'publicite' AND c.reference_id = cp.id AND c.statut = 'en_attente'
		  AND c.date_expiration < now() AND cp.statut = 'en_attente'`); err != nil {
		return err
	}
	_, err := ex.Exec("UPDATE commandes SET statut = 'expiree' WHERE statut = 'en_attente' AND date_expiration < now()")
	return err
}

func (s *Serveur) gererCreationCommande(w http.ResponseWriter, r *http.Request) {
	utilisateurID, role := s.roleConnecte(r)
	if utilisateurID == 0 {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return
	}

	var entree CommandeEntree
	if err := json.NewDecoder(r.Body).Decode(&entree); err != nil {
		envoyerErreur(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}
	switch {
	case entree.Type == "formation" && role != roleParticulier:
		envoyerErreur(w, http.StatusForbidden, "seuls les particuliers peuvent s'inscrire à une formation")
		return
	case entree.Type == "vente" && role != roleParticulier && role != roleProfessionnel:
		envoyerErreur(w, http.StatusForbidden, "seuls les particuliers et les professionnels peuvent acheter un objet")
		return
	case (entree.Type == "abonnement" || entree.Type == "publicite") && role != roleProfessionnel:
		envoyerErreur(w, http.StatusForbidden, "réservé aux professionnels")
		return
	case entree.Type != "formation" && entree.Type != "vente" && entree.Type != "abonnement" && entree.Type != "publicite":
		envoyerErreur(w, http.StatusBadRequest, "type invalide (formation, vente, abonnement ou publicite attendu)")
		return
	}
	switch entree.Type {
	case "abonnement":
		// Un abonnement par professionnel : la référence est le professionnel lui-même
		entree.ReferenceID = utilisateurID
	case "publicite":
		entree.ReferenceID = 0
	}

	tx, err := s.db.Begin()
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	defer tx.Rollback()

	if err := libererCommandesExpirees(tx); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	// Une commande déjà en cours pour le même achat est reprise telle quelle (double clic, retour arrière…)
	var c Commande
	if entree.Type == "publicite" {
		err = scannerCommande(tx.QueryRow(`SELECT `+prefixer(colonnesCommande, "c.")+` FROM commandes c
			JOIN campagnes_publicitaires cp ON cp.id = c.reference_id
			WHERE c.utilisateur_id = $1 AND c.type = 'publicite' AND c.statut = 'en_attente'
			  AND cp.cible_type = $2 AND cp.cible_id = $3 AND cp.montant_mensuel = $4`,
			utilisateurID, entree.CibleType, entree.CibleID, entree.Montant), &c)
	} else {
		err = scannerCommande(tx.QueryRow(`SELECT `+colonnesCommande+` FROM commandes
			WHERE utilisateur_id = $1 AND type = $2 AND reference_id = $3 AND statut = 'en_attente'`,
			utilisateurID, entree.Type, entree.ReferenceID), &c)
	}
	if err == nil {
		s.repondreCommande(w, http.StatusOK, c)
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}

	var libelle string
	var montant float64
	var code int
	var message string
	switch entree.Type {
	case "formation":
		libelle, montant, code, message = verifierAchatFormation(tx, utilisateurID, entree.ReferenceID)
	case "vente":
		libelle, montant, code, message = verifierAchatObjet(tx, utilisateurID, entree.ReferenceID)
	case "abonnement":
		libelle, montant, code, message = verifierAchatAbonnement(tx, utilisateurID)
	case "publicite":
		entree.ReferenceID, libelle, montant, code, message = verifierAchatPublicite(tx, utilisateurID, entree)
	}
	if code != 0 {
		envoyerErreur(w, code, message)
		return
	}

	err = scannerCommande(tx.QueryRow(`
		INSERT INTO commandes (utilisateur_id, type, reference_id, libelle, montant, fournisseur, date_expiration)
		VALUES ($1, $2, $3, $4, $5, $6, now() + $7::interval)
		RETURNING `+colonnesCommande,
		utilisateurID, entree.Type, entree.ReferenceID, libelle, montant, fournisseurPaiement(),
		fmt.Sprintf("%d seconds", int(dureeReservation.Seconds())),
	), &c)
	if err == nil {
		err = tx.Commit()
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	s.repondreCommande(w, http.StatusCreated, c)
}

func (s *Serveur) repondreCommande(w http.ResponseWriter, code int, c Commande) {
	url, err := urlPaiement(c)
	if err != nil {
		envoyerErreur(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	envoyerJSON(w, code, map[string]any{"commande": c, "url_paiement": url})
}

// Vérifie qu'une formation payante peut être achetée ; renvoie un code HTTP d'erreur sinon.
func verifierAchatFormation(tx *sql.Tx, utilisateurID, formationID int) (string, float64, int, string) {
	var titre, statut string
	var tarif sql.NullFloat64
	var nbPlaces sql.NullInt64
	var dejaCommencee bool
	err := tx.QueryRow(`SELECT titre, statut, tarif, nb_places, date_debut <= now() FROM formations WHERE id = $1 FOR UPDATE`,
		formationID).Scan(&titre, &statut, &tarif, &nbPlaces, &dejaCommencee)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return "", 0, http.StatusNotFound, "formation introuvable"
	case err != nil:
		return "", 0, http.StatusInternalServerError, "erreur serveur"
	case statut != "publie":
		return "", 0, http.StatusBadRequest, "cette formation n'est pas publiée"
	case dejaCommencee:
		return "", 0, http.StatusBadRequest, "cette formation a déjà commencé"
	case !tarif.Valid || tarif.Float64 <= 0:
		return "", 0, http.StatusBadRequest, "cette formation est gratuite : inscrivez-vous directement"
	}

	var dejaInscrit bool
	var occupees int
	err = tx.QueryRow(`
		SELECT EXISTS (SELECT 1 FROM inscriptions WHERE utilisateur_id = $1 AND type_offre = 'formation' AND offre_id = $2),
		       (SELECT count(*) FROM inscriptions WHERE type_offre = 'formation' AND offre_id = $2)
		     + (SELECT count(*) FROM commandes WHERE type = 'formation' AND reference_id = $2 AND statut = 'en_attente')`,
		utilisateurID, formationID).Scan(&dejaInscrit, &occupees)
	switch {
	case err != nil:
		return "", 0, http.StatusInternalServerError, "erreur serveur"
	case dejaInscrit:
		return "", 0, http.StatusConflict, "vous êtes déjà inscrit à cette formation"
	case nbPlaces.Valid && occupees >= int(nbPlaces.Int64):
		return "", 0, http.StatusConflict, "plus de places disponibles"
	}
	return titre, tarif.Float64, 0, ""
}

// Vérifie qu'un objet peut être acheté et le réserve pendant le paiement.
func verifierAchatObjet(tx *sql.Tx, utilisateurID, annonceID int) (string, float64, int, string) {
	var vendeurID int
	var titre, typeAnnonce, statut string
	var prix sql.NullFloat64
	err := tx.QueryRow(`SELECT utilisateur_id, titre, type, prix, statut FROM annonces WHERE id = $1 FOR UPDATE`, annonceID).
		Scan(&vendeurID, &titre, &typeAnnonce, &prix, &statut)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return "", 0, http.StatusNotFound, "annonce introuvable"
	case err != nil:
		return "", 0, http.StatusInternalServerError, "erreur serveur"
	case typeAnnonce != "vente" || !prix.Valid || prix.Float64 <= 0:
		return "", 0, http.StatusBadRequest, "cet objet n'est pas en vente"
	case vendeurID == utilisateurID:
		return "", 0, http.StatusBadRequest, "vous ne pouvez pas acheter votre propre objet"
	case statut != "en_ligne":
		return "", 0, http.StatusConflict, "cet objet n'est plus disponible"
	}
	if _, err := tx.Exec("UPDATE annonces SET statut = 'reservee' WHERE id = $1", annonceID); err != nil {
		return "", 0, http.StatusInternalServerError, "erreur serveur"
	}
	return titre, prix.Float64, 0, ""
}

// Lit une commande en vérifiant qu'elle appartient à la personne connectée (ou que c'est un admin).
func (s *Serveur) commandeDeLaRequete(w http.ResponseWriter, r *http.Request) (Commande, bool) {
	var c Commande
	utilisateurID, role := s.roleConnecte(r)
	if utilisateurID == 0 {
		envoyerErreur(w, http.StatusUnauthorized, "non authentifié")
		return c, false
	}
	id, err := idDepuisChemin(r)
	if err != nil {
		envoyerErreur(w, http.StatusBadRequest, "id invalide")
		return c, false
	}
	if err := libererCommandesExpirees(s.db); err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return c, false
	}
	err = scannerCommande(s.db.QueryRow(`SELECT `+colonnesCommande+` FROM commandes WHERE id = $1`, id), &c)
	if errors.Is(err, sql.ErrNoRows) || (err == nil && c.UtilisateurID != utilisateurID && role != roleAdministrateur) {
		envoyerErreur(w, http.StatusNotFound, "commande introuvable")
		return c, false
	}
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return c, false
	}
	var facture sql.NullInt64
	if err := s.db.QueryRow("SELECT id FROM documents WHERE type_document = 'facture' AND commande_id = $1 LIMIT 1", c.ID).
		Scan(&facture); err == nil {
		v := int(facture.Int64)
		c.FactureID = &v
	}
	return c, true
}

func (s *Serveur) gererDetailCommande(w http.ResponseWriter, r *http.Request) {
	if c, ok := s.commandeDeLaRequete(w, r); ok {
		envoyerJSON(w, http.StatusOK, c)
	}
}

// Paiement de démonstration : n'existe que tant que Stripe n'est pas branché.
func (s *Serveur) gererPaiementSimule(w http.ResponseWriter, r *http.Request) {
	c, ok := s.commandeDeLaRequete(w, r)
	if !ok {
		return
	}
	utilisateurID, _ := s.roleConnecte(r)
	if c.Fournisseur != "simulation" || c.UtilisateurID != utilisateurID {
		envoyerErreur(w, http.StatusForbidden, "paiement simulé indisponible pour cette commande")
		return
	}
	reference, err := genererCodeDepot(8)
	if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	s.repondreConfirmation(w, r, c.ID, "sim_"+reference)
}

func (s *Serveur) repondreConfirmation(w http.ResponseWriter, r *http.Request, id int, reference string) {
	_, err := s.confirmerCommande(id, reference)
	switch {
	case errors.Is(err, errCommandeIntrouvable):
		envoyerErreur(w, http.StatusNotFound, err.Error())
		return
	case errors.Is(err, errCommandeTerminee), errors.Is(err, errCommandeExpiree):
		envoyerErreur(w, http.StatusConflict, err.Error())
		return
	case err != nil:
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	s.gererDetailCommande(w, r)
}

func (s *Serveur) gererAnnulationCommande(w http.ResponseWriter, r *http.Request) {
	c, ok := s.commandeDeLaRequete(w, r)
	if !ok {
		return
	}
	if err := s.annulerCommande(c.ID, "annulee"); errors.Is(err, errCommandeTerminee) {
		envoyerErreur(w, http.StatusConflict, err.Error())
		return
	} else if err != nil {
		envoyerErreur(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	s.gererDetailCommande(w, r)
}

// Confirme une commande payée. Rappeler deux fois avec la même commande ne fait rien de plus
// (Stripe peut envoyer le même webhook plusieurs fois).
func (s *Serveur) confirmerCommande(id int, referenceExterne string) (Commande, error) {
	var c Commande
	tx, err := s.db.Begin()
	if err != nil {
		return c, err
	}
	defer tx.Rollback()

	err = scannerCommande(tx.QueryRow(`SELECT `+colonnesCommande+` FROM commandes WHERE id = $1 FOR UPDATE`, id), &c)
	if errors.Is(err, sql.ErrNoRows) {
		return c, errCommandeIntrouvable
	}
	if err != nil {
		return c, err
	}
	if c.Statut == "payee" {
		return c, nil
	}
	if c.Statut != "en_attente" {
		return c, errCommandeTerminee
	}
	if time.Now().After(c.DateExpiration) {
		tx.Rollback()
		_ = s.annulerCommande(c.ID, "expiree")
		return c, errCommandeExpiree
	}

	var inscriptionID int
	var destinataireID int
	referencePaiement := c.ReferenceID
	var titreNotification, messageNotification, lienNotification string
	if c.Type == "formation" {
		err = tx.QueryRow(`
			INSERT INTO inscriptions (utilisateur_id, type_offre, offre_id, statut_paiement, montant_paye)
			VALUES ($1, 'formation', $2, 'paye', $3)
			RETURNING id`, c.UtilisateurID, c.ReferenceID, c.Montant).Scan(&inscriptionID)
		if err == nil {
			err = ajouterPointsScore(tx, c.UtilisateurID, pointsInscriptionOffre, "Inscription à une offre")
		}
		if err == nil {
			err = tx.QueryRow("SELECT createur_id FROM formations WHERE id = $1", c.ReferenceID).Scan(&destinataireID)
		}
		titreNotification = "Nouvelle inscription"
		messageNotification = "Une personne vient de s'inscrire et de payer ta formation « " + c.Libelle + " »."
		lienNotification = "formations"
	} else if c.Type == "abonnement" {
		var abonnementID int
		err = tx.QueryRow(`
			INSERT INTO abonnements (utilisateur_id, montant_mensuel, statut, stripe_abonnement_id)
			VALUES ($1, $2, 'actif', $3) RETURNING id`, c.UtilisateurID, c.Montant, referenceExterne).Scan(&abonnementID)
		if err == nil {
			_, err = tx.Exec("UPDATE utilisateurs SET abonnement = 'premium' WHERE id = $1", c.UtilisateurID)
		}
		// Le paiement d'abonnement référence l'abonnement (même convention que les données existantes)
		referencePaiement = abonnementID
	} else if c.Type == "publicite" {
		var resultat sql.Result
		resultat, err = tx.Exec(`UPDATE campagnes_publicitaires SET statut = 'active', date_debut = CURRENT_DATE
			WHERE id = $1 AND statut = 'en_attente'`, c.ReferenceID)
		if err == nil {
			if n, _ := resultat.RowsAffected(); n == 0 {
				return c, errCommandeTerminee
			}
		}
	} else {
		err = tx.QueryRow(`
			UPDATE annonces SET statut = 'cedee', acheteur_id = $1
			WHERE id = $2 AND statut = 'reservee' AND acheteur_id IS NULL
			RETURNING utilisateur_id`, c.UtilisateurID, c.ReferenceID).Scan(&destinataireID)
		if errors.Is(err, sql.ErrNoRows) {
			return c, errCommandeTerminee
		}
		if err == nil {
			err = ajouterPointsScore(tx, destinataireID, pointsAnnonceCedee, "Annonce cédée")
		}
		titreNotification = "Objet vendu !"
		messageNotification = "Ton objet « " + c.Libelle + " » a été acheté et payé."
		lienNotification = "mes-annonces"
	}
	if err != nil {
		return c, err
	}

	paiementID, err := enregistrerPaiement(tx, c.UtilisateurID, c.Type, &referencePaiement, c.Montant, referenceExterne)
	if err != nil {
		return c, err
	}
	err = scannerCommande(tx.QueryRow(`
		UPDATE commandes SET statut = 'payee', paiement_id = $1, reference_externe = $2, date_paiement = now()
		WHERE id = $3
		RETURNING `+colonnesCommande, paiementID, referenceExterne, c.ID), &c)
	if err == nil {
		err = tx.Commit()
	}
	if err != nil {
		return c, err
	}

	// Après l'enregistrement : documents et notifications ne peuvent plus faire échouer l'achat
	if inscriptionID != 0 {
		archiverAutomatiquement(fmt.Sprintf("attestation de l'inscription %d", inscriptionID), func() (int, error) {
			return s.genererAttestation(inscriptionID, c.UtilisateurID)
		})
	}
	archiverAutomatiquement(fmt.Sprintf("facture de la commande %d", c.ID), func() (int, error) {
		return s.genererFacture(c.ID)
	})
	if c.Type == "abonnement" || c.Type == "publicite" {
		archiverAutomatiquement(fmt.Sprintf("contrat de la commande %d", c.ID), func() (int, error) {
			return s.genererContratCommande(c.ID)
		})
	}
	if destinataireID != 0 && destinataireID != c.UtilisateurID {
		notifierAvecLien(destinataireID, titreNotification, messageNotification, lienNotification)
	}
	return c, nil
}

// Annule (ou fait expirer) une commande en attente et remet l'objet en vente.
func (s *Serveur) annulerCommande(id int, statutFinal string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var c Commande
	err = scannerCommande(tx.QueryRow(`SELECT `+colonnesCommande+` FROM commandes WHERE id = $1 FOR UPDATE`, id), &c)
	if errors.Is(err, sql.ErrNoRows) {
		return errCommandeIntrouvable
	}
	if err != nil {
		return err
	}
	if c.Statut != "en_attente" {
		if c.Statut == statutFinal {
			return nil
		}
		return errCommandeTerminee
	}
	if _, err := tx.Exec("UPDATE commandes SET statut = $1 WHERE id = $2", statutFinal, id); err != nil {
		return err
	}
	if c.Type == "vente" {
		if _, err := tx.Exec(`UPDATE annonces SET statut = 'en_ligne'
			WHERE id = $1 AND statut = 'reservee' AND acheteur_id IS NULL`, c.ReferenceID); err != nil {
			return err
		}
	}
	if c.Type == "publicite" {
		if _, err := tx.Exec(`UPDATE campagnes_publicitaires SET statut = 'annulee'
			WHERE id = $1 AND statut = 'en_attente'`, c.ReferenceID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Facture PDF de la commande, archivée au nom de l'acheteur (une seule par commande).
func (s *Serveur) genererFacture(commandeID int) (int, error) {
	var existant int
	err := s.db.QueryRow("SELECT id FROM documents WHERE type_document = 'facture' AND commande_id = $1 LIMIT 1", commandeID).
		Scan(&existant)
	if err == nil {
		return existant, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	var c Commande
	if err := scannerCommande(s.db.QueryRow(`SELECT `+colonnesCommande+` FROM commandes WHERE id = $1`, commandeID), &c); err != nil {
		return 0, err
	}
	if c.Statut != "payee" || c.DatePaiement == nil {
		return 0, fmt.Errorf("la commande %d n'est pas payée", commandeID)
	}
	var prenom, nom, email string
	var reference sql.NullString
	var taux, commission sql.NullFloat64
	err = s.db.QueryRow(`
		SELECT u.prenom, u.nom, u.email, p.stripe_id, p.taux_commission, p.montant_commission
		FROM commandes c
		JOIN utilisateurs u ON u.id = c.utilisateur_id
		LEFT JOIN paiements p ON p.id = c.paiement_id
		WHERE c.id = $1`, commandeID).Scan(&prenom, &nom, &email, &reference, &taux, &commission)
	if err != nil {
		return 0, err
	}

	numero := fmt.Sprintf("F-%d-%06d", c.DatePaiement.Year(), c.ID)
	moyen := "Carte bancaire (Stripe)"
	if c.Fournisseur == "simulation" {
		moyen = "Paiement simulé (mode démonstration, aucun débit réel)"
	}
	lignes := []string{
		"Facture n° " + numero,
		"Date : " + c.DatePaiement.Format("02/01/2006 15:04"),
		"Client : " + prenom + " " + nom + " (" + email + ")",
		"Désignation : " + designationsCommande[c.Type] + " « " + c.Libelle + " »",
	}
	if c.Type == "vente" {
		var vendeurPrenom, vendeurNom string
		if err := s.db.QueryRow(`SELECT u.prenom, u.nom FROM annonces a JOIN utilisateurs u ON u.id = a.utilisateur_id
			WHERE a.id = $1`, c.ReferenceID).Scan(&vendeurPrenom, &vendeurNom); err == nil {
			lignes = append(lignes, "Vendeur : "+vendeurPrenom+" "+vendeurNom)
		}
	}
	lignes = append(lignes, fmt.Sprintf("Montant payé : %.2f € TTC", c.Montant))
	if commission.Valid && taux.Valid {
		lignes = append(lignes, fmt.Sprintf("dont commission UpcycleConnect : %.2f € (%.2f %%)", commission.Float64, taux.Float64))
	}
	lignes = append(lignes, "Moyen de paiement : "+moyen)
	if reference.Valid {
		lignes = append(lignes, "Référence du paiement : "+reference.String)
	}

	contenu, err := genererPDF("Facture", lignes)
	if err != nil {
		return 0, err
	}
	var id int
	err = s.db.QueryRow(`
		INSERT INTO documents (type_document, titre, utilisateur_id, commande_id, cree_par, contenu_pdf)
		VALUES ('facture', $1, $2, $3, $2, $4)
		RETURNING id`, "Facture "+numero+" : "+c.Libelle, c.UtilisateurID, c.ID, contenu).Scan(&id)
	return id, err
}
